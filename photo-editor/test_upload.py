import requests
import cv2
import numpy as np

# Create a dummy image (100x100 red square)
img = np.zeros((100, 100, 3), dtype=np.uint8)
img[:] = (0, 0, 255)
_, img_encoded = cv2.imencode('.png', img)

# Send to backend
url = 'http://127.0.0.1:8000/upload'
files = {'file': ('test.png', img_encoded.tobytes(), 'image/png')}

try:
    response = requests.post(url, files=files)
    if response.status_code == 200:
        print("Upload successful!")
        # Verify content type
        if response.headers['content-type'] == 'image/png':
            print("Received valid PNG image.")
        else:
            print(f"Unexpected content type: {response.headers['content-type']}")
    else:
        print(f"Upload failed: {response.status_code} - {response.text}")
except Exception as e:
    print(f"Error: {e}")
