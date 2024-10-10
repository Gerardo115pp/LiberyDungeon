from CategoriesServicePB import categories_requests_pb2_grpc as CategoriesServicePB
import Config as app_config
import models
import grpc

if app_config.DEVELOPMENT_MODE:
    from LiberyCertUtils import getSSLcaFile
    
def sendCreateCategoryRequest(category_name: str, parent_uuid: str, cluster_uuid: str) -> str:
    response: CategoriesServicePB.categories__requests__pb2.CreateCategoryResponse = None
    
    credentials = grpc.ssl_channel_credentials()
    
    if app_config.DEVELOPMENT_MODE:
        credentials = grpc.ssl_channel_credentials(root_certificates=getSSLcaFile())
        
    with grpc.secure_channel(app_config.GRPC_SERVER, credentials) as channel:
        stub = CategoriesServicePB.CategoriesServiceStub(channel)
        
        request = CategoriesServicePB.categories__requests__pb2.CreateCategoryRequest(name=category_name, parent=parent_uuid, cluster=cluster_uuid)
        
        response = stub.CreateCategory(request)
        
    return response.uuid
        
def sendGetCategoriesClusterRequest(cluster_uuid: str) -> tuple[models.DungeonsCategoriesCluster, Exception]:
    response: CategoriesServicePB.categories__requests__pb2.GetCategoriesClusterResponse = None
    
    credentials = grpc.ssl_channel_credentials()
    
    if app_config.DEVELOPMENT_MODE:
        credentials = grpc.ssl_channel_credentials(root_certificates=getSSLcaFile())
        
    with grpc.secure_channel(app_config.GRPC_SERVER, credentials) as channel:
        stub = CategoriesServicePB.CategoriesServiceStub(channel)
        
        request = CategoriesServicePB.categories__requests__pb2.GetCategoriesClusterRequest(uuid=cluster_uuid)
        
        response = stub.GetCategoriesCluster(request)
        
    try: 
        cluster = models.DungeonsCategoriesCluster(
            response.cluster.uuid, 
            response.cluster.name, 
            response.cluster.fs_path, 
            response.cluster.filter_category, 
            response.cluster.root_category,
        )
    except Exception as e:
        return None, e
    
    return cluster, None