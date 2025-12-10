const API_URL = 'http://127.0.0.1:8000';

const canvas = document.getElementById('editor-canvas');
const ctx = canvas.getContext('2d');
const fileInput = document.getElementById('file-upload');
const loadingSpinner = document.getElementById('loading-spinner');

// State
let uploadBlob = null; // The absolute original upload
let originalBlob = null; // The base for current adjustments (may be cropped/rotated)
let currentBlob = null; // The displayed image (adjustments applied)

// History
const historyStack = [];
let historyPointer = -1;
const MAX_HISTORY = 10;

function pushHistory(state) {
    // state = { blob: Blob, adjustments: {...}, filters: {...} }
    // Remove future states if we are in middle
    if (historyPointer < historyStack.length - 1) {
        historyStack.splice(historyPointer + 1);
    }
    historyStack.push(state);
    if (historyStack.length > MAX_HISTORY) historyStack.shift();
    else historyPointer++;

    updateUndoRedoUI();
}

function getCurrentState() {
    return {
        blob: originalBlob, // We store the BASE image
        adjustments: {
            brightness: document.getElementById('brightness').value,
            contrast: document.getElementById('contrast').value,
            saturation: document.getElementById('saturation').value,
            sharpness: document.getElementById('sharpness').value
        },
        filters: { ...activeFilters }
    };
}

function restoreState(state) {
    originalBlob = state.blob;
    // Restore UI values
    document.getElementById('brightness').value = state.adjustments.brightness;
    document.getElementById('contrast').value = state.adjustments.contrast;
    document.getElementById('saturation').value = state.adjustments.saturation;
    document.getElementById('sharpness').value = state.adjustments.sharpness;
    Object.assign(activeFilters, state.filters);

    // Update Filter UI
    document.querySelectorAll('.filter-btn').forEach(btn => {
        const f = btn.dataset.filter;
        if (activeFilters[f]) {
            btn.style.borderColor = '#3b82f6';
            btn.style.color = '#3b82f6';
        } else {
            btn.style.borderColor = 'transparent';
            btn.style.color = 'white';
        }
    });

    // Trigger update
    applyAdjustments();
}

function updateUndoRedoUI() {
    // Enable/Disable buttons if they existed (add them to HTML first?)
    // Assuming buttons will be added or we just have keyboard shortcuts?
    // Project spec says "Add Undo". Let's add buttons in next step or assume shortcuts.
    // For now, let's just log.
}

// Download
document.getElementById('save-btn').addEventListener('click', () => {
    if (currentBlob) {
        const url = URL.createObjectURL(currentBlob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'edited-image.png';
        a.click();
        URL.revokeObjectURL(url);
    }
});

// Zoom
let zoomLevel = 1.0;
canvas.addEventListener('wheel', (e) => {
    e.preventDefault();
    const delta = e.deltaY > 0 ? -0.1 : 0.1;
    zoomLevel = Math.max(0.1, Math.min(5.0, zoomLevel + delta));
    canvas.style.transform = `scale(${zoomLevel})`;
});

// Add listeners to sliders
['brightness', 'contrast', 'saturation', 'sharpness'].forEach(id => {
    document.getElementById(id).addEventListener('input', () => {
        applyAdjustments();
        // Debounce handles the call, but we should probably push history ONLY on mouse up (change)
        // rather than every input? 
        // For now, history is only pushed on Transforms.
        // If we want history for slider changes, we need to listen for 'change' event (mouse up)
    });

    document.getElementById(id).addEventListener('change', () => {
        // Push history on release
        pushHistory(getCurrentState());
    });
});

// Global Undo/Redo shortcuts
window.addEventListener('keydown', (e) => {
    if ((e.metaKey || e.ctrlKey) && e.key === 'z') {
        e.preventDefault();
        if (e.shiftKey) {
            // Redo
            if (historyPointer < historyStack.length - 1) {
                historyPointer++;
                restoreState(historyStack[historyPointer]);
            }
        } else {
            // Undo
            if (historyPointer > 0) {
                historyPointer--;
                restoreState(historyStack[historyPointer]);
            }
        }
    }
});

// Crop State
let isCropMode = false;
let cropStart = null;
let cropRect = null; // {x, y, w, h}

// Tab switching
document.querySelectorAll('.tool-btn').forEach(btn => {
    btn.addEventListener('click', () => {
        document.querySelectorAll('.tool-btn').forEach(b => b.classList.remove('active'));
        document.querySelectorAll('.panel').forEach(p => p.classList.remove('active'));

        btn.classList.add('active');
        document.getElementById(`${btn.dataset.tab}-panel`).classList.add('active');
    });
});

// File Upload
fileInput.addEventListener('change', async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    const formData = new FormData();
    formData.append('file', file);

    showLoading(true);
    try {
        const response = await fetch(`${API_URL}/upload`, {
            method: 'POST',
            body: formData
        });

        if (response.ok) {
            const blob = await response.blob();
            uploadBlob = blob;
            originalBlob = blob;
            currentBlob = blob;

            // Initial History
            historyStack.length = 0;
            historyPointer = -1;
            pushHistory(getCurrentState());

            await loadImageToCanvas(blob);
        }
    } catch (err) {
        console.error("Upload failed", err);
    } finally {
        showLoading(false);
    }
});

function debounce(func, wait) {
    let timeout;
    return function (...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// Filter State
let activeFilters = {
    grayscale: false,
    sepia: false,
    negative: false,
    blur: false,
    sobel: false
};

const applyAdjustments = debounce(async () => {
    if (!originalBlob) return;

    const brightness = parseInt(document.getElementById('brightness').value);
    const contrast = parseInt(document.getElementById('contrast').value);
    const saturation = parseInt(document.getElementById('saturation').value);
    const sharpness = parseInt(document.getElementById('sharpness').value);

    // We always start processing from the ORIGINAL image to avoid destructive compounding
    // Convert blob to base64
    const reader = new FileReader();
    reader.readAsDataURL(originalBlob);
    reader.onloadend = async () => {
        const base64data = reader.result;

        showLoading(true);
        try {
            const body = {
                image: base64data,
                brightness,
                contrast,
                saturation,
                sharpness,
                ...activeFilters
            };

            const response = await fetch(`${API_URL}/process`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(body)
            });

            if (response.ok) {
                const blob = await response.blob();
                currentBlob = blob;
                await loadImageToCanvas(blob);
            }
        } catch (err) {
            console.error("Processing failed", err);
        } finally {
            showLoading(false);
        }
    };
}, 500); // 500ms debounce

// Filter Buttons
document.querySelectorAll('.filter-btn').forEach(btn => {
    btn.addEventListener('click', () => {
        const filter = btn.dataset.filter;
        // Toggle filter
        activeFilters[filter] = !activeFilters[filter];

        // Update visual state
        if (activeFilters[filter]) {
            btn.style.borderColor = '#3b82f6';
            btn.style.color = '#3b82f6';
        } else {
            btn.style.borderColor = 'transparent';
            btn.style.color = 'white';
        }

        applyAdjustments();
    });
});

// Transform Helper
async function applyTransform(params) {
    if (!currentBlob) return;

    // We transform the CURRENT displayed image (destructive)
    const reader = new FileReader();
    reader.readAsDataURL(currentBlob);
    reader.onloadend = async () => {
        const base64data = reader.result;
        showLoading(true);
        try {
            const response = await fetch(`${API_URL}/transform`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    image: base64data,
                    ...params
                })
            });

            if (response.ok) {
                const blob = await response.blob();
                originalBlob = blob;
                currentBlob = blob;

                // Reset Adjustments UI (for destructive transform)
                document.querySelectorAll('input[type="range"]').forEach(input => {
                    input.value = input.getAttribute('value');
                });
                // Reset Filters UI
                Object.keys(activeFilters).forEach(k => activeFilters[k] = false);
                document.querySelectorAll('.filter-btn').forEach(btn => {
                    btn.style.borderColor = 'transparent';
                    btn.style.color = 'white';
                });

                // Push new state
                pushHistory(getCurrentState());

                await loadImageToCanvas(blob);
            }
        } catch (err) {
            console.error(err);
        } finally {
            showLoading(false);
        }
    };
}

document.getElementById('rotate-90').addEventListener('click', () => applyTransform({ rotate: 90 }));
document.getElementById('flip-h').addEventListener('click', () => applyTransform({ flip: 'h' }));
document.getElementById('flip-v').addEventListener('click', () => applyTransform({ flip: 'v' }));

// Crop Mode (Simple Toggle)
const cropBtn = document.getElementById('crop-btn');
let cropOverlay = null;

cropBtn.addEventListener('click', () => {
    isCropMode = !isCropMode;
    cropBtn.classList.toggle('active');

    if (isCropMode) {
        // Init Overlay
        if (!cropOverlay) {
            cropOverlay = document.createElement('canvas');
            cropOverlay.style.position = 'absolute';
            cropOverlay.style.pointerEvents = 'all';
            cropOverlay.style.cursor = 'crosshair';
            canvas.parentNode.appendChild(cropOverlay);

            // Sync size
            const rect = canvas.getBoundingClientRect();
            cropOverlay.style.top = canvas.offsetTop + 'px';
            cropOverlay.style.left = canvas.offsetLeft + 'px';
            cropOverlay.width = canvas.width;
            cropOverlay.height = canvas.height;
            cropOverlay.style.width = canvas.style.width || canvas.width + 'px';
            cropOverlay.style.height = canvas.style.height || canvas.height + 'px';

            // Mouse Events
            cropOverlay.addEventListener('mousedown', startCrop);
            window.addEventListener('mousemove', drawCrop);
            window.addEventListener('mouseup', endCrop);
        }
        cropOverlay.classList.remove('hidden');
    } else {
        if (cropOverlay) cropOverlay.classList.add('hidden');
    }
});

function startCrop(e) {
    const rect = cropOverlay.getBoundingClientRect();
    const x = (e.clientX - rect.left) * (cropOverlay.width / rect.width);
    const y = (e.clientY - rect.top) * (cropOverlay.height / rect.height);
    cropStart = { x, y };
    cropRect = { x, y, w: 0, h: 0 };
}

function drawCrop(e) {
    if (!isCropMode || !cropStart) return;
    const rect = cropOverlay.getBoundingClientRect();
    const x = (e.clientX - rect.left) * (cropOverlay.width / rect.width);
    const y = (e.clientY - rect.top) * (cropOverlay.height / rect.height);

    cropRect.w = x - cropStart.x;
    cropRect.h = y - cropStart.y;

    const ctxOverlay = cropOverlay.getContext('2d');
    ctxOverlay.clearRect(0, 0, cropOverlay.width, cropOverlay.height);
    ctxOverlay.strokeStyle = 'red';
    ctxOverlay.lineWidth = 2;
    ctxOverlay.strokeRect(cropStart.x, cropStart.y, cropRect.w, cropRect.h);
}

async function endCrop() {
    if (!isCropMode || !cropStart) return;
    // Normalize rect
    let { x, y, w, h } = cropRect;
    if (w < 0) { x += w; w = -w; }
    if (h < 0) { y += h; h = -h; }

    if (confirm("Apply Crop?")) {
        await applyTransform({ crop: { x, y, w, h } });
        isCropMode = false;
        cropBtn.classList.remove('active');
        cropOverlay.classList.add('hidden');
    }
    cropStart = null;
}

async function loadImageToCanvas(blob) {
    const url = URL.createObjectURL(blob);
    const img = new Image();
    img.onload = () => {
        canvas.width = img.width;
        canvas.height = img.height;
        ctx.drawImage(img, 0, 0);

        // Fit canvas to container roughly if needed, CSS handles visual scaling
        URL.revokeObjectURL(url);
    };
    img.src = url;
}

function showLoading(show) {
    if (show) loadingSpinner.classList.remove('hidden');
    else loadingSpinner.classList.add('hidden');
}

// Reset
document.getElementById('reset-btn').addEventListener('click', async () => {
    if (uploadBlob) {
        originalBlob = uploadBlob;
        currentBlob = uploadBlob;
        await loadImageToCanvas(uploadBlob);

        // Reset all sliders
        document.querySelectorAll('input[type="range"]').forEach(input => {
            input.value = input.getAttribute('value');
        });

        // Reset filters
        Object.keys(activeFilters).forEach(k => activeFilters[k] = false);
        document.querySelectorAll('.filter-btn').forEach(btn => {
            btn.style.borderColor = 'transparent';
            btn.style.color = 'white';
        });
    }
});
