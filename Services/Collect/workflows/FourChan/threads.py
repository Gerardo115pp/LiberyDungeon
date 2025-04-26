from bs4 import BeautifulSoup, Tag
import helpers
import models
import httpx
from workflows.Communication.CommunicationUtils import impersonateTLSFingerPrint
import bleach

def getThreadImages(board_name: str, thread_uuid: str) -> tuple[list[str], Exception]:
    thread_url = helpers.getThreadUrl(board_name, thread_uuid)
    
    response = impersonateTLSFingerPrint(thread_url)
    if response.status_code != 200:
        return [], Exception(f"Error on 4chan communication for '{thread_url}': {response.status_code}")
    
    thread_content = response.text
    
    images: list[str] = []
    
    try:
        soup = BeautifulSoup(thread_content, 'html.parser')
        
        images_tags = soup.findAll('a', attrs={'class': 'fileThumb'})
        
        for image_tag in images_tags:
            image_url = image_tag['href']
            if not image_url.startswith("http"):
                image_url = f"https:{image_url}"
            
            images.append(image_url)
    except Exception as e:
        return [], e    
    
    return images, None

def getThreadContent(board_name: str, thread_uuid: str) -> tuple[models.FourChanThread, Exception|None]:
    thread_url = helpers.getThreadUrl(board_name, thread_uuid)
    
    response: httpx.Response = httpx.get(thread_url)
    if response.status_code != 200:
        return None, Exception(f"Error on 4chan communication: {response.status_code}")
    
    thread_content = response.text
    
    soup = BeautifulSoup(thread_content, 'html.parser')
    
    thread_model: models.FourChanThread = parseFourChanThread(thread_uuid, soup)
    
    replies_models: list[models.FourChanThreadReply] = parseFourChanThreadReplies(soup)
    
    thread_model.replies = replies_models
    
    return thread_model, None
    
def parseFourChanThread(thread_uuid: str, soup: BeautifulSoup) -> models.FourChanThread:
    date_element = soup.select_one('.postInfo.desktop > .dateTime')
    date_value = date_element.attrs.get('data-utc', "no-date")
    
    file_element = soup.select_one('.opContainer .file a')
    file_value = file_element.attrs.get('href', "no-file")
    if not file_value.startswith("http"):
        file_value = f"https:{file_value}"
    
    title_element = soup.select_one('.opContainer .subject')
    title_value = title_element.text if title_element else "no-title"
    
    description_element = soup.select_one('.opContainer .postMessage')
    description_value = __getContentWithoutXSS(description_element) 
    
    cover_image_element = soup.select_one('.opContainer .fileThumb img')
    cover_image_value = cover_image_element.attrs.get('src', "no-image-url")
    if not cover_image_value.startswith("http"):
        cover_image_value = f"https:{cover_image_value}"
    
    return models.FourChanThread( 
        uuid=thread_uuid,
        date=date_value,
        title=title_value,
        file=file_value,
        description=description_value,
        cover_image_url=cover_image_value,
        replies=[]
    )

def parseFourChanThreadReplies(soup: BeautifulSoup) -> list[models.FourChanThreadReply]:
    replies: list[models.FourChanThreadReply] = []
    
    reply_containers = soup.select('.thread .replyContainer')
    
    for reply_container_element in reply_containers:
        reply_model = parseFourChanThreadReply(reply_container_element)
        replies.append(reply_model)
        
    return replies
    
def parseFourChanThreadReply(reply_container_element: BeautifulSoup) -> models.FourChanThreadReply:
    uuid_value = reply_container_element.attrs.get('id', "no-uuid")
    
    message_element = reply_container_element.select_one('.postMessage')
    message_value = __getContentWithoutXSS(message_element)
    
    thumbnail_element = reply_container_element.select_one('.fileThumb img')
    thumbnail_value = "no-thumbnail-url"    
    if thumbnail_element:
        thumbnail_value = thumbnail_element.attrs.get('src', "no-thumbnail-url")
        if not thumbnail_value.startswith("http"):
            thumbnail_value = f"https:{thumbnail_value}"
    
    file_element = reply_container_element.select_one('.file a')
    file_value = "no-file"
    if file_element:
        file_value = file_element.attrs.get('href', "no-file")
        if not file_value.startswith("http"):
            file_value = f"https:{file_value}"
    
    date_element = reply_container_element.select_one('.postInfo.desktop > .dateTime')
    date_value = date_element.attrs.get('data-utc', "no-date")
    
    return models.FourChanThreadReply(
        uuid=uuid_value,
        message=message_value,
        thumbnail_url=thumbnail_value,
        file=file_value,
        date=date_value
    )
    
def __getContentWithoutXSS(content: Tag) -> str:
    if not content:
        return ""
    
    content_string = str(content.decode_contents())
    
    allowed_tags = list(bleach.ALLOWED_TAGS) + ['br', 'span', 'wbr']
    allowed_attributes = bleach.ALLOWED_ATTRIBUTES
    
    clean_html = bleach.clean(content_string, tags=allowed_tags, attributes=allowed_attributes) 
    
    return clean_html
    
    