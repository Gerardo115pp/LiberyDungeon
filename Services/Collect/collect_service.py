from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from handlers.chan_boards import four_chan_boards_handler
from handlers.chan_downloads import four_chan_downloads_handler
from handlers.chan_threads import four_chan_threads_handler
import uvicorn

def createApp():
    app = FastAPI(
        title="Libery collect service",
        description="Handles the scrapping of medias that will be downloaded, the download happens in the download service"
    )
    
    app.add_middleware(CORSMiddleware, allow_origins=["*"], allow_methods=["*"], allow_headers=["*"])
    
    app.include_router(four_chan_boards_handler, prefix="/4chan-boards", tags=["4chan-boards"])
    app.include_router(four_chan_downloads_handler, prefix="/4chan-downloads", tags=["4chan-downloads"])
    app.include_router(four_chan_threads_handler, prefix="/4chan-threads", tags=["4chan-threads"])
    
    return app

app = createApp()

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=6972)