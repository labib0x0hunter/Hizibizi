import requests
import cv2
import numpy as np
import base64

# Create a dummy image (100x100 gray square, value 100)
img = np.zeros((100, 100, 3), dtype=np.uint8)
img[:] = (100, 100, 100)
_, img_encoded = cv2.imencode('.png', img)
b64_img = "data:image/png;base64," + base64.b64encode(img_encoded.tobytes()).decode('utf-8')

# Send to backend for brightness adjustment (+50)
url = 'http://127.0.0.1:8000/process'
payload = {
    "image": b64_img,
    "brightness": 50,
    "contrast": 0,
    "saturation": 0,
    "sharpness": 0
}

try:
    response = requests.post(url, json=payload)
    if response.status_code == 200:
        # Decode response
        nparr = np.frombuffer(response.content, np.uint8)
        res_img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
        
        # Check center pixel value. Should be approx 100 + 50 = 150.
        pixel_val = res_img[50, 50, 0]
        print(f"Original pixel: 100. Adjusted (+50): {pixel_val}")
        
        if 148 <= pixel_val <= 152:
            print("Brightness test PASSED")
        else:
            print("Brightness test FAILED")
    else:
        print(f"Process failed: {response.status_code} - {response.text}")

except Exception as e:
    print(f"Error: {e}")
