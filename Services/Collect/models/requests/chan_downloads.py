# Purpose: Requests bodies for chan_downloads handler
from pydantic import BaseModel
import Config as app_config

class DownloadThreadImagesRequest(BaseModel):
    board_name: str
    thread_uuid: str
    target_category_name: str
    cluster_uuid: str
    parent_uuid: str



