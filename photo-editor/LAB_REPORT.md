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
