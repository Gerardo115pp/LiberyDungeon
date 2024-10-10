from dataclasses import dataclass, asdict
from typing import Any
import helpers

def allArgsPresent(obj: object, args: list[str]) -> bool:
    """ 
        this method should be reassigned from the __init__.py file
    """
    raise NotImplementedError("this shouldn't happen")

@dataclass
class CatalogThread:
    uuid: str
    date: str
    file: str
    responses: int
    images: int
    teaser: str
    image_url: str
    teaser_thumb_width: int
    teaser_thumb_height: int
    board_name: str

    @property
    def hasImages(self) -> bool:
        return self.images > 0
    
    @property
    def hasCoverImage(self) -> bool:
        return self.image_url != "no-image-url"
    
def unpackFromCatalogData(uuid: str, catalog_data: dict[str, Any]) -> CatalogThread:
    catalog_thread = CatalogThread(
        uuid=uuid,
        date=catalog_data.get('date', "no-date"),
        file=catalog_data.get('file', "no-file"),
        image_url=catalog_data.get('imgurl', "no-image-url"),
        teaser_thumb_width=catalog_data.get('tn_w', 0),
        teaser_thumb_height=catalog_data.get('tn_h', 0),
        responses=catalog_data.get('r', 0),
        images=catalog_data.get('i', 0),
        teaser=catalog_data.get('teaser', "no-teaser"),
        board_name=catalog_data.get('board_name', "no-board-name")
    )

    if catalog_thread.image_url != "no-image-url":
        catalog_thread.image_url = helpers.getThreadCoverImageUrl(catalog_thread.board_name, catalog_thread.image_url)
    
    if "sub" in catalog_data:
        catalog_thread.teaser = f"{catalog_thread.teaser} - {catalog_data['sub']}"

    return catalog_thread

@dataclass 
class FourChanThreadReply:
    uuid: str
    message: str
    thumbnail_url: str
    file: str
    date: str

@dataclass
class FourChanThread:
    uuid: str
    date: str
    file: str
    title: str        
    description: str
    cover_image_url: str
    replies: list[FourChanThreadReply]
    