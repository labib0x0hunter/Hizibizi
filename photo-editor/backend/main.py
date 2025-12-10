from fastapi import FastAPI, UploadFile, File, Response
from fastapi.staticfiles import StaticFiles
from fastapi.middleware.cors import CORSMiddleware
import uvicorn
import cv2
import numpy as np
import base64
from pydantic import BaseModel
from typing import Optional
import processor

app = FastAPI()

# CORS for development
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)



@app.post("/upload")
async def upload_image(file: UploadFile = File(...)):
    contents = await file.read()
    nparr = np.fromstring(contents, np.uint8)
    img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
    
    # Encode back to PNG for display
    _, encoded_img = cv2.imencode('.png', img)
    return Response(content=encoded_img.tobytes(), media_type="image/png")

class ProcessRequest(BaseModel):
    image: str # Base64 encoded image
    brightness: int = 0
    contrast: int = 0
    saturation: int = 0
    sharpness: int = 0
    # Filters
    grayscale: bool = False
    sepia: bool = False
    negative: bool = False
    blur: bool = False
    sobel: bool = False
    # Add other params later

@app.post("/process")
async def process_image(req: ProcessRequest):
    # Decode base64 image
    # Expected format: "data:image/png;base64,....." or just raw base64
    if "," in req.image:
        encoded_data = req.image.split(",")[1]
    else:
        encoded_data = req.image
        
    nparr = np.frombuffer(base64.b64decode(encoded_data), np.uint8)
    img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
    
    # Apply adjustments in order
    # Note: Order matters. Usually Brightness/Contrast/Saturation/Sharpness
    # Roadmap: B -> C -> S -> Sharpness
    
    if req.brightness != 0:
        img = processor.adjust_brightness(img, req.brightness)
        
    if req.contrast != 0:
        img = processor.adjust_contrast(img, req.contrast)
        
    if req.saturation != 0:
        img = processor.adjust_saturation(img, req.saturation)
        
    if req.sharpness > 0:
        img = processor.apply_sharpness(img, req.sharpness)
        
    # Apply filters
    if req.grayscale:
        img = processor.apply_grayscale(img)
    if req.sepia:
        img = processor.apply_sepia(img)
    if req.negative:
        img = processor.apply_negative(img)
    if req.blur:
        img = processor.apply_blur(img)
    if req.sobel:
        img = processor.apply_sobel(img)
        
    # Helper to encode response
    _, encoded_img = cv2.imencode('.png', img)
    return Response(content=encoded_img.tobytes(), media_type="image/png")

class TransformRequest(BaseModel):
    image: str # Base64
    rotate: int = 0 # 0, 90, 180, 270
    flip: Optional[str] = None # 'h' or 'v'
    crop: Optional[dict] = None # {x, y, w, h}

@app.post("/transform")
async def transform_image(req: TransformRequest):
    if "," in req.image:
        encoded_data = req.image.split(",")[1]
    else:
        encoded_data = req.image
        
    nparr = np.frombuffer(base64.b64decode(encoded_data), np.uint8)
    img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
    
    if req.rotate:
        img = processor.rotate_image(img, req.rotate)
        
    if req.flip:
        img = processor.flip_image(img, req.flip)
        
    if req.crop:
        x = int(req.crop['x'])
        y = int(req.crop['y'])
        w = int(req.crop['w'])
        h = int(req.crop['h'])
        img = processor.crop_image(img, x, y, w, h)
        
    _, encoded_img = cv2.imencode('.png', img)
    return Response(content=encoded_img.tobytes(), media_type="image/png")



# Serve frontend static files (Mount at root MUST be last)
app.mount("/", StaticFiles(directory="../frontend", html=True), name="frontend")

if __name__ == "__main__":
    uvicorn.run("main:app", host="127.0.0.1", port=8000, reload=True)
