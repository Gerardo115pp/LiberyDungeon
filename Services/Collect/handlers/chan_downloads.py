from fastapi import APIRouter, HTTPException, Response, Request
from helpers.jwt_helpers import signCluster
from fastapi.responses import JSONResponse
import Config as app_config
import workflows
import models

four_chan_downloads_handler = APIRouter()

@four_chan_downloads_handler.post("/thread/images", status_code=202)
async def getThreadImages(request_body: models.Requests.ChanDownloads.DownloadThreadImagesRequest):
    """
    This handler requests the thread content and scraps the images from it. Then it sends a message the categories server to create the 
    the target categorie. After that it sends the download batch request to the download service through grpc.
    """
    
    thread_images_urls: list[str]
    
    thread_images_urls, err = workflows.FourChan.Threads.getThreadImages(request_body.board_name, request_body.thread_uuid)
    if err:
        print(f"Error while getting thread images urls: {err}")
        raise HTTPException(status_code=502, detail=f"Error while getting thread images urls: {err}")
    
    cluster, err = workflows.CategoriesService.sendGetCategoriesClusterRequest(request_body.cluster_uuid)
    if err:
        print(f"In getThreadImages: Error while getting cluster because '{err}'")
        raise HTTPException(status_code=404, detail="Error while retrieving cluster data, it likely does not exist")
    
    category_uuid: str = workflows.CategoriesService.sendCreateCategoryRequest(request_body.target_category_name, request_body.parent_uuid, cluster.uuid)
    
    cluster_token: str = signCluster(cluster, app_config.DOMAIN_SECRET)
    
    download_uuid: str = workflows.DownloadsService.sendThreadImagesDownloadRequest(thread_images_urls, category_uuid, cluster_token, download_id=request_body.thread_uuid)
    
    return JSONResponse(
        content={"download_uuid": download_uuid}, 
        headers={"Cache-Control": "no-cache"}
    )
    