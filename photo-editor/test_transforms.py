import requests
import cv2
import numpy as np
import base64

# Create a 100x50 image (width=100, height=50). Use asymmetric colors.
# Left half Red, Right half Blue.
img = np.zeros((50, 100, 3), dtype=np.uint8)
img[:, :50] = (0, 0, 255)   # Red
img[:, 50:] = (255, 0, 0) # Blue
_, img_encoded = cv2.imencode('.png', img)
b64_img = "data:image/png;base64," + base64.b64encode(img_encoded.tobytes()).decode('utf-8')

url = 'http://127.0.0.1:8000/transform'

# Test 1: Rotate 90
payload_rot = {"image": b64_img, "rotate": 90}
try:
    print("Testing Rotate 90...")
    res = requests.post(url, json=payload_rot)
    if res.status_code == 200:
        nparr = np.frombuffer(res.content, np.uint8)
        out = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
        # Original 100x50. Rotated 90 -> 50x100.
        if out.shape[0] == 100 and out.shape[1] == 50:
            print("Rotate 90 Dimensions PASSED")
        else:
            print(f"Rotate 90 Dimensions FAILED: {out.shape}")
    else:
        print(f"Rotate Failed: {res.status_code}")
except Exception as e:
    print(f"Rotate Error: {e}")

# Test 2: Crop
# Crop center 20x20. x=40, y=15, w=20, h=20.
# Should contain Red and Blue border?
# Actually input is 100 width. 0-49 Red, 50-99 Blue.
# Crop x=40, width=20 -> x=40 to 60.
# 40-49 Red, 50-60 Blue.
# So left half of crop is Red, right half is Blue.
payload_crop = {"image": b64_img, "crop": {"x": 40, "y": 15, "w": 20, "h": 20}}
try:
    print("Testing Crop...")
    res = requests.post(url, json=payload_crop)
    if res.status_code == 200:
        nparr = np.frombuffer(res.content, np.uint8)
        out = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
        if out.shape[0] == 20 and out.shape[1] == 20:
            print("Crop Dimensions PASSED")
            # Check pixels
            # out[10, 0] should be Red (original x=40)
            # out[10, 19] should be Blue (original x=59)
            left_px = out[10, 0]
            right_px = out[10, 19]
            if left_px[0] == 0 and left_px[2] == 255: # Red
                print("Crop Content Left PASSED")
            else:
                print(f"Crop Content Left FAILED: {left_px}")
                
            if right_px[0] == 255 and right_px[2] == 0: # Blue
                print("Crop Content Right PASSED")
            else:
                print("Crop Content Right FAILED")
        else:
            print(f"Crop Dimensions FAILED: {out.shape}")
    else:
        print(f"Crop Failed: {res.status_code}")
except Exception as e:
    print(f"Crop Error: {e}")
