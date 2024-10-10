from DownloadServicePB import DownloadsPB
import Config as app_config
import grpc

if app_config.DEVELOPMENT_MODE:
    from LiberyCertUtils import getSSLcaFile

def sendThreadImagesDownloadRequest(thread_images_urls: list[str], category_uuid: str, cluster_token:str, download_id: str="") -> str:
    response: DownloadsPB.download__requests__pb2.DownloadBatchResponse = None

    credentials = grpc.ssl_channel_credentials()

    if app_config.DEVELOPMENT_MODE:
        credentials = grpc.ssl_channel_credentials(root_certificates=getSSLcaFile())
    
    with grpc.secure_channel(app_config.GRPC_SERVER, credentials) as channel:
        stub = DownloadsPB.DownloadServiceStub(channel)
        
        # If download_id is not provided, meaning it's an empty string, then the server will generate a new one
        request = DownloadsPB.download__requests__pb2.DownloadImagesBatchRequest(image_urls=thread_images_urls, category_uuid=category_uuid, cluster_token=cluster_token, download_uuid=download_id)
        
        response = stub.DownloadImagesBatch(request)
    
    return response.download_uuid
        
        
        