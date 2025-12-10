# Pro Photo Editor (Python + Web)

A web-based photo editor with a FastAPI (Python) backend and a modern HTML/JS frontend.

## Features
- **Adjustments**: Brightness, Contrast, Saturation, Sharpness.
- **Filters**: Grayscale, Sepia, Negative, Gaussian Blur, Sobel Edge Detection.
- **Transforms**: Rotate (90/180/270), Flip (H/V), Crop.
- **Tools**: Undo/Redo (Cmd+Z/Cmd+Shift+Z), Zoom (Mouse Wheel), Download.

## Architecture
- **Backend**: FastAPI serving static files and processing images with OpenCV.
- **Frontend**: Vanilla JS + HTML5 Canvas.
- **Communication**: Frontend sends images (Base64) + parameters to Backend APIs (`/process`, `/transform`).

## Setup & Run

1. **Install Dependencies**:
   ```bash
   cd backend
   python3 -m venv venv
   source venv/bin/activate
   pip install fastapi uvicorn python-multipart opencv-python-headless numpy
   ```

2. **Run Server**:
   ```bash
   ./venv/bin/python -m uvicorn main:app --reload
   ```

3. **Open App**:
   Navigate to [http://127.0.0.1:8000](http://127.0.0.1:8000)

## Usage
1. Upload an image.
2. Use the **Adjust** tab for sliders (real-time preview).
3. Use **Filters** for instant effects.
4. Use **Transform** to Rotate, Flip, or Crop.
   - For Crop: Click "Crop Mode", drag on canvas, "Apply Crop" confirmation.
5. **Undo**: Cmd+Z / Ctrl+Z. **Redo**: Cmd+Shift+Z.
6. **Zoom**: Scroll wheel on canvas.
7. **Download**: Click "Download" to save as PNG.
