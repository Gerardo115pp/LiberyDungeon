from helpers import getBoardCatalogUrl
from bs4.element import ResultSet, Tag
from bs4 import BeautifulSoup
import Config as app_config
import models
import httpx
import json
import re

def boardExists(board: str) -> bool:
    board_url = getBoardCatalogUrl(board)
    
    response: httpx.Response = httpx.get(board_url)
    
    return response.status_code == 200

def findThreadsTag(lookup_pool: ResultSet) -> tuple[Tag, Exception]:
    threads_tag: Tag = None
    err: Exception = None
    
    for tag in lookup_pool:
        if tag.name == "script":
            if app_config.CATALOG_THREADS_LOOKUP_STRING in str(tag):
                threads_tag = tag
                break
    
    if not threads_tag:
        err = Exception("No threads tag found")
    
    return threads_tag, err

def getBoardCatalog(board: str) -> tuple[str, Exception]:
    board_url = getBoardCatalogUrl(board)
    
    response: httpx.Response = httpx.get(board_url)
    if response.status_code != 200:
        return "", Exception(f"Error on 4chan communication: {response.status_code}")
    
    return response.text, None

def getBoardCatalogContent(board: str) -> tuple[list[models.CatalogThread], Exception]:
    board_catalog, err = getBoardCatalog(board)
    if err:
        return [], err

    json_data = {}

    try:
        soup = BeautifulSoup(board_catalog, 'html.parser')

        threads_tags = soup.findAll('script')
        threads_tag, err = findThreadsTag(threads_tags)
        if err:
            return [], err

        raw_tables = "{" + re.search(r"var\scatalog\s=\s\{.+\};", str(threads_tag))[0].split("{",1)[1]
        
        json_data = json.loads(raw_tables[:-1])
    except Exception as e:
        return [], e
    
    if "threads" not in json_data:
        return [], Exception("No threads found in catalog")
    
    board_threads: list[models.CatalogThread] = []  
    
    for thread_uuid, thread_data in json_data['threads'].items():
        thread_data['board_name'] = board
        
        catalog_thread = models.ChanFunctions.unpackFromCatalogData(thread_uuid, thread_data)
        board_threads.append(catalog_thread)

    return board_threads, None