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

