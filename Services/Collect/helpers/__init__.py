import Config as app_config

def getBoardCatalogUrl(board: str) -> str:
    board_template = app_config.BOARD_CATALOG_TEMPLATE
    
    board_url = board_template.replace('{board_name}', board)
    
    return board_url

def getThreadCoverImageUrl(board: str, cover_image_id: str) -> str:
    thread_cover_image_template = app_config.THREAD_COVER_IMAGE_TEMPLATE
    
    thread_cover_image_url = thread_cover_image_template.replace('{board_name}', board)
    thread_cover_image_url = thread_cover_image_url.replace('{thread_id}', cover_image_id)
    
    return thread_cover_image_url

def getThreadUrl(board_name: str, thread_uuid: str) -> str:
    thread_template = app_config.THREAD_TEMPLATE
    
    thread_url = thread_template.replace('{board_name}', board_name)
    thread_url = thread_url.replace('{thread_uuid}', thread_uuid)
    
    return thread_url