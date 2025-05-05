from fastapi import APIRouter, HTTPException, Response
from fastapi.encoders import jsonable_encoder
from fastapi.responses import JSONResponse
import Config as app_config
import workflows
import models
import httpx

four_chan_threads_handler = APIRouter()

@four_chan_threads_handler.get("/thread", status_code=200, response_class=JSONResponse)
async def getThread(board_name: str, thread_id: str):
    thread: models.FourChanThread
    err: Exception|None = None
    
    thread, err = workflows.FourChan.Threads.getThreadContent(board_name, thread_id)
    if err:
        raise HTTPException(status_code=502, detail=f"Error while getting thread content: {err}")
    
    json_content = jsonable_encoder(thread)
    
    response = JSONResponse(
        content=json_content, 
        headers = {
            "Cache-Control": "max-age=600, must-revalidate, private",
            "Age": "10"
        }
    )    
    
    return response 