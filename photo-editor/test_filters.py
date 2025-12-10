import requests
import cv2
import numpy as np
import base64

# Create a colored image (100x100), red square
img = np.zeros((100, 100, 3), dtype=np.uint8)
img[:] = (0, 0, 255) # Red in BGR
_, img_encoded = cv2.imencode('.png', img)
b64_img = "data:image/png;base64," + base64.b64encode(img_encoded.tobytes()).decode('utf-8')

# Send to backend for grayscale
url = 'http://127.0.0.1:8000/process'
payload = {
    "image": b64_img,
    "grayscale": True
}

try:
    response = requests.post(url, json=payload)
    if response.status_code == 200:
        nparr = np.frombuffer(response.content, np.uint8)
        res_img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
        
        # Check center pixel. Red (0,0,255) -> Gray.
        # Gray = 0.299*R + 0.587*G + 0.114*B = 0.299*255 = ~76
        b, g, r = res_img[50, 50]
        print(f"Original: Red. Grayscale BGR: {b}, {g}, {r}")
        
        if 70 <= b <= 80 and 70 <= g <= 80 and 70 <= r <= 80:
            print("Grayscale test PASSED")
        else:
            print("Grayscale test FAILED")
    else:
        print(f"Process failed: {response.status_code} - {response.text}")

except Exception as e:
    print(f"Error: {e}")
