from fastapi import APIRouter, HTTPException, Response
from fastapi.encoders import jsonable_encoder
from fastapi.responses import JSONResponse
import Config as app_config
import workflows
import httpx
import models


four_chan_boards_handler = APIRouter()

@four_chan_boards_handler.get("/boards", status_code=200)
async def getBoards():
    boards: list[models.Board]
    err: Exception = None
    
    boards = app_config.TRACKED_BOARDS
    
    return boards

@four_chan_boards_handler.get("/board/catalog", status_code=202, response_class=JSONResponse)
async def getBoardCatalog(board_name: str, not_filter: bool = False):
    boards_content: dict
    err: Exception

    print(f"Getting board catalog for {board_name}")
    
    boards_content: list[models.CatalogThread] = []
    
    boards_content, err = workflows.FourChan.Boards.getBoardCatalogContent(board_name)
    if err:
        print(f"Error while getting boards content: {err}")
        raise HTTPException(status_code=502, detail=f"Error while getting boards content: {err}")
    
    if not_filter:
        return boards_content
    
    filtered_boards_content: list[models.CatalogThread] = []
    
    for board_thread in boards_content:
        if board_thread.hasImages and board_thread.hasCoverImage:
            filtered_boards_content.append(board_thread)

    json_content = jsonable_encoder(filtered_boards_content)       

    response = JSONResponse(
        content=json_content, 
        headers={
            "Cache-Control": "max-age=120, must-revalidate, private",
            "Age": "10"
        }
    )

    return response
    
@four_chan_boards_handler.get("/boards/proxy-media", status_code=200)
async def proxyGetBoardCoverImage(media_url: str):
    """
    This method requests the image from the url and returns it as a response. avoiding cors issues
    """

    fake_headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
    }
    
    chan_response = httpx.get(media_url, headers=fake_headers, follow_redirects=True)
    if chan_response.status_code != 200:
        raise HTTPException(status_code=502, detail=f"Error while getting requested image: {chan_response.status_code}")

    mime_type = chan_response.headers.get("Content-Type", "image/jpeg")

    response = Response(
        content=chan_response.content, 
        media_type=mime_type,
        headers={
            "Cache-Control": "max-age=13600, must-revalidate, private",
            "Age": "10"
        }
    )
    
    return response
