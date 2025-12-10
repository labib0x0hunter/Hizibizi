# Photo Editor Technical Report

## 1. Architecture
The application follows a client-server architecture:
- **Client (Frontend)**: Handles user interaction, state management (Undo/Redo history), and rendering (Canvas). It sends processing requests to the server.
- **Server (Backend)**: Built with FastAPI. It receives raw images and parameters, performs matrix operations using OpenCV/NumPy, and returns the processed image.

This decoupling ensures that computationally intensive operations (convolutions, geometric transforms) are handled by optimized Python libraries (OpenCV) rather than the browser, adhering to the requirement for Python-based image processing.

## 2. Algorithms & Math

### 2.1 Adjustments
- **Brightness**: 
  $$ P_{new} = P_{old} + \text{value} $$
  Implemented by adding a scalar to the pixel matrix and clipping to [0, 255].

- **Contrast**: 
  $$ F = \frac{259(\text{value} + 255)}{255(259 - \text{value})} $$
  $$ P_{new} = F(P_{old} - 128) + 128 $$
  Expands or contracts the dynamic range around the mid-tone (128).

- **Saturation**:
  Converted RGB to HSV (Hue, Saturation, Value).
  $$ S_{new} = S_{old} \times (1 + \frac{\text{value}}{100}) $$
  Then converted back to RGB.

- **Sharpness**:
  Used a kernel convolution to sharpen:
  $$ K = \begin{bmatrix} 0 & -1 & 0 \\ -1 & 5 & -1 \\ 0 & -1 & 0 \end{bmatrix} $$
  The result is blended with the original based on the slider intensity.

### 2.2 Filters
- **Grayscale**: Weighted sum of RGB channels: $0.299R + 0.587G + 0.114B$.
- **Sepia**: Matrix transformation mapping RGB to specific warm tones.
- **Gaussian Blur**: 5x5 convolution kernel approximating a Gaussian distribution.
- **Sobel Edge Detection**:
  Compute gradients $G_x$ and $G_y$ using Sobel kernels.
  $$ Magnitude = \sqrt{G_x^2 + G_y^2} $$
  
### 2.3 Transforms
- **Rotation**: Geometric transformation of the pixel matrix coordinates.
- **Crop**: Slicing the NumPy array: `img[y:y+h, x:x+w]`.

## 3. Implementation Details
- **Non-Destructive Adjustments**: Adjustments are applied to the base image on every change, preventing compound degradation.
- **Destructive Transforms**: Geometric changes (Crop/Rotate) update the base image to simplify the coordinate system for subsequent operations.
