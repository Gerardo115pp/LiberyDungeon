from time import time
import jwt
import models


def signCluster(cluster: models.DungeonsCategoriesCluster, secret: str) -> str:
    cluster_dict = cluster.toDict()
    expiration_time = time() + 30 # This is a system token so it should expire really fast
    
    cluster_dict['exp'] = int(expiration_time) # exp shouldn't include milliseconds or any fractional seconds. Python includes the milliseconds since epoch in time()
    
    jwt_token = jwt.encode(
        cluster_dict,
        secret,
        algorithm='HS256'
    )
    
    return jwt_token