# HizBiz Code
**Just Random Fun Coding - TempCode : Maybe I Will Complete These Projects**

---
**WHAT IS HAPPENING HERE**
- **grepcoder** : Atcoder problem catagory in go.
- **http_server** : Http server built from scratch in c.
- **judge** : Online judge system in go.
- **mmap** : Memory-mapped file in go.
- **zer0search** : mini search engine in go.
- **packet_analysis** : .pcap or .cap files analyzer in go.
- **Projects** : Just Another folder for shit coding in go.
- **request_parser** : http request parser in c.
- **threadpool** : threadpool in c.
- **wal** : Write Ahead Log implementataion in go.
- **fio** : Minimal clone of `fmt.Print()`, `fmt.Scan()`, `bytes.Buffer`, `bufio.NewReader`, `bufio.NewWriter` in go.



# LABORATORY REPORT: Web-Based Photo Editor

**Project Name**: Pro Photo Editor  
**Technologies**: Python (FastAPI, OpenCV), HTML5, CSS3, JavaScript  
**Date**: December 2025

---

## 1. Introduction
The objective of this project was to develop a full-featured, web-based image processing application that leverages Python's powerful scientific computing libraries while providing a modern, responsive user interface in the browser. Unlike client-side editors that rely solely on WebGL or Canvas APIs, this application implements a client-server architecture where complex image manipulation algorithms are executed on a Python backend using OpenCV and NumPy.

## 2. Methodology & Architecture

### 2.1 System Architecture
The application follows a decoupled **Client-Server** model:

*   **Frontend (Client)**: Built with **HTML5, CSS3, and Vanilla JavaScript**. It handles user input (file selection, slider movements, button clicks), renders the image using the **Canvas API**, and manages application state (Undo/Redo history).
*   **Backend (Server)**: Built with **FastAPI** (Python). It serves as a RESTful API that receives image data (Base64 encoded) and processing parameters. It performs matrix operations on the image and returns the processed result.

**Data Flow:**
1.  User uploads image -> Frontend displays preview.
2.  User adjusts slider (e.g., Brightness) -> Frontend debounces input -> Sends POST request to `/process`.
3.  Backend decodes image -> Applies OpenCV operations -> Encodes to PNG -> Returns to Client.
4.  Frontend receives blob -> Updates Canvas.

### 2.2 Tools & Technologies
*   **Language**: Python 3.13
*   **Framework**: FastAPI (High-performance web framework)
*   **Image Processing**: OpenCV (cv2), NumPy (Matrix math)
*   **Server**: Uvicorn (ASGI server)
*   **Frontend**: HTML5 Canvas, CSS Flexbox/Grid

---

## 3. Algorithms & Implementation

The core logical operations are mathematically defined as follows:

### 3.1 Basic Adjustments

#### **Brightness**
Brightness control adds a constant scalar value (`v`) to every pixel's channel.
$$ P_{new} = P_{old} + v $$
*Implementation*: 
```python
img = img.astype(np.int16)
img = img + value
img = np.clip(img, 0, 255) # Clamping to valid range
```

#### **Contrast**
Contrast expands or shrinks the difference between intensity values.
$$ Factor, F = \frac{259(v + 255)}{255(259 - v)} $$
$$ P_{new} = F(P_{old} - 128) + 128 $$
This formula increases stability by anchoring 128 (mid-gray) as the center pivot.

#### **Saturation**
We convert the image from **BGR** (Blue-Green-Red) color space to **HSV** (Hue-Saturation-Value). We then scale the **S** channel.
$$ S_{new} = S_{old} \times (1 + \frac{v}{100}) $$

#### **Sharpness**
Sharpening is achieved via **Convolution** with a kernel that effectively subtracts the Laplacian (edges) from the original image.
**Kernel Matrix:**
$$
K = \begin{bmatrix}
 0 & -1 &  0 \\
-1 &  5 & -1 \\
 0 & -1 &  0
\end{bmatrix}
$$
The center pixel is boosted (5) while neighbors are subtracted, enhancing local contrast.

### 3.2 Filters

#### **Gaussian Blur**
15x15 convolution kernel approximating a Gaussian distribution.
#### **Grayscale**
Converts color image to single-channel intensity image using standard luma coefficients:
$$ Gray = 0.299R + 0.587G + 0.114B $$

#### **Sepia**
A matrix transformation simulating the chemical aging of photographs.
$$
\begin{bmatrix} R' \\ G' \\ B' \end{bmatrix} = 
\begin{bmatrix}
0.393 & 0.769 & 0.189 \\
0.349 & 0.686 & 0.168 \\
0.272 & 0.534 & 0.131
\end{bmatrix}
\begin{bmatrix} R \\ G \\ B \end{bmatrix}
$$

#### **Sobel Edge Detection**
Calculates the gradient magnitude of image intensity at each pixel.
$$ G_x = \text{Sobel}_x * I, \quad G_y = \text{Sobel}_y * I $$
$$ Magnitude = \sqrt{G_x^2 + G_y^2} $$
Displays areas of high spatial frequency (edges).

### 3.3 Geometric Transforms

#### **Rotation & Flip**
Mapped input pixel $(x, y)$ to output pixel $(x', y')$.
*   **Rotation 90°**: $x' = y, y' = w - 1 - x$
*   **Flip H**: $x' = w - 1 - x$

#### **Crop**
Implemented via NumPy array slicing, which is highly efficient.
```python
cropped = img[y : y+h, x : x+w]
```

---

## 4. Frontend Features & UX
*   **Debouncing**: To prevent server overload, slider inputs wait for 500ms of inactivity before sending a request.
*   **State Management**:
    *   `uploadBlob`: The original file.
    *   `originalBlob`: The base for adjustments (updated after Crop/Rotate).
    *   `currentBlob`: The result of real-time adjustments.
*   **Undo/Redo**: Implemented using a circular history stack (`MAX_HISTORY=10`), storing full application state (blobs + slider values).

---

## 5. Testing & Verification
The system was verified using automated Python scripts interacting with the API endpoints.

| Test Case | Method | Expected Output | Result |
|-----------|--------|-----------------|--------|
| **Upload** | `POST /upload` | HTTP 200, PNG Content-Type | **PASS** |
| **Brightness** | `POST /process` {brightness: 50} | Pixel value +50 | **PASS** |
| **Grayscale** | `POST /process` {grayscale: true} | RGB with eq values | **PASS** |
| **Rotate** | `POST /transform` {rotate: 90} | W/H swapped | **PASS** |
| **Crop** | `POST /transform` {crop_rect} | Reduced dimensions | **PASS** |

---

## 6. Conclusion
The "Pro Photo Editor" project successfully demonstrates the integration of a backend-heavy computational model with a smooth frontend experience. By offloading processing to Python, we enabled the use of industrial-grade algorithms (OpenCV) while keeping the client lightweight. The modular design allows for easy extension, such as adding AI-based features (Object Removal, Denoising) in the future.


import cv2
import numpy as np

def adjust_brightness(img, value):
    # value is -100 to 100.
    # We want to add value to pixels.
    # Convert to higher depth to avoid overflow/underflow during add
    img = img.astype(np.int16)
    img = img + value
    img = np.clip(img, 0, 255)
    return img.astype(np.uint8)

def adjust_contrast(img, value):
    # value -100 to 100
    # factor = (259 * (value + 255)) / (255 * (259 - value))
    # R = factor * (R - 128) + 128
    if value == 259: value = 258 # Avoid div by zero
    factor = (259 * (value + 255)) / (255 * (259 - value))
    
    img = img.astype(np.float32)
    img = factor * (img - 128) + 128
    img = np.clip(img, 0, 255)
    return img.astype(np.uint8)

def adjust_saturation(img, value):
    # value -100 to 100
    # Convert to HSV (OpenCV uses BGR -> HSV)
    hsv = cv2.cvtColor(img, cv2.COLOR_BGR2HSV).astype(np.float32)
    # S is channel 1. S in OpenCV is 0-255
    # scale factor: 1 + (value / 100)
    # e.g. value=0 -> 1.0, value=-100 -> 0.0, value=100 -> 2.0
    scale = 1.0 + (value / 100.0)
    hsv[:,:,1] = hsv[:,:,1] * scale
    hsv[:,:,1] = np.clip(hsv[:,:,1], 0, 255)
    
    img = cv2.cvtColor(hsv.astype(np.uint8), cv2.COLOR_HSV2BGR)
    return img

def apply_sharpness(img, value):
    # value 0-100? Roadmap says "Sharpness"
    # Kernel:
    #  0 -1  0
    # -1  5 -1
    #  0 -1  0
    # We can blend purely sharpened image with original based on value
    kernel = np.array([[0, -1, 0], [-1, 5, -1], [0, -1, 0]])
    sharpened = cv2.filter2D(img, -1, kernel)
    
    if value == 0: return img
    
    # Blending: if value is used as intensity
    # Let's assume value 0-100 maps to alpha 0.0-1.0 ???
    # Roadmap just says "Sharpness" and gives the kernel.
    # Usually you apply it once. If there's a slider, maybe we blend.
    # Let's assume binary application or simple blending if value is passed.
    # The prompt implies a slider "Basic Adjustments -> Sharpness", so assume 0-100 intensity.
    
    alpha = value / 100.0
    return cv2.addWeighted(sharpened, alpha, img, 1 - alpha, 0)

def apply_grayscale(img):
    return cv2.cvtColor(cv2.cvtColor(img, cv2.COLOR_BGR2GRAY), cv2.COLOR_GRAY2BGR)

def apply_negative(img):
    return cv2.bitwise_not(img)

def apply_sepia(img):
    # Sepia matrix
    # R' = 0.393R + 0.769G + 0.189B
    # G' = 0.349R + 0.686G + 0.168B
    # B' = 0.272R + 0.534G + 0.131B
    # OpenCV is BGR
    # B' = 0.272R + 0.534G + 0.131B
    # G' = 0.349R + 0.686G + 0.168B
    # R' = 0.393R + 0.769G + 0.189B
    
    # Kernel for BGR input to BGR output
    # Row 0 (B out): [0.131, 0.534, 0.272]
    # Row 1 (G out): [0.168, 0.686, 0.349]
    # Row 2 (R out): [0.189, 0.769, 0.393]
    
    kernel = np.array([
        [0.131, 0.534, 0.272],
        [0.168, 0.686, 0.349],
        [0.189, 0.769, 0.393]
    ])
    
    img = cv2.transform(img, kernel)
    img = np.clip(img, 0, 255)
    return img.astype(np.uint8)

def apply_blur(img):
    # Gaussian Blur 15x15 for stronger effect
    return cv2.GaussianBlur(img, (15, 15), 0)

def apply_sobel(img):
    # Convert to gray then apply sobel
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    
    # Gradient X
    sobelx = cv2.Sobel(gray, cv2.CV_64F, 1, 0, ksize=3)
    # Gradient Y
    sobely = cv2.Sobel(gray, cv2.CV_64F, 0, 1, ksize=3)
    
    # Magnitude
    mag = np.sqrt(sobelx**2 + sobely**2)
    mag = np.clip(mag, 0, 255).astype(np.uint8)
    
    # Back to BGR
    return cv2.cvtColor(mag, cv2.COLOR_GRAY2BGR)

def rotate_image(img, angle):
    # angle: 90, 180, 270 (clockwise)
    if angle == 90:
        return cv2.rotate(img, cv2.ROTATE_90_CLOCKWISE)
    elif angle == 180:
        return cv2.rotate(img, cv2.ROTATE_180)
    elif angle == 270:
        return cv2.rotate(img, cv2.ROTATE_90_COUNTERCLOCKWISE)
    return img

def flip_image(img, mode):
    # mode: 'h', 'v'
    if mode == 'h':
        return cv2.flip(img, 1)
    elif mode == 'v':
        return cv2.flip(img, 0)
    return img

def crop_image(img, x, y, w, h):
    # Ensure bounds
    H, W = img.shape[:2]
    x = max(0, min(x, W))
    y = max(0, min(y, H))
    w = max(1, min(w, W - x))
    h = max(1, min(h, H - y))
    
    return img[y:y+h, x:x+w]

