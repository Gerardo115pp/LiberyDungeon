from dataclasses import dataclass, asdict
from typing import Any
import json

def allArgsPresent(obj: object, args: list[str]) -> bool:
    """ 
        this method should be reassigned from the __init__.py file
    """
    raise NotImplementedError("this shouldn't happen")

@dataclass
class Category:
    uuid: str
    name: str
    fullpath: str
    parent: str
    cluster: str

    def toJson(self) -> str:
        self_dict = asdict(self)
        
        json_str = json.dumps(self_dict)
        
    def toDict(self) -> dict[str, Any]:
        return asdict(self)
    
@dataclass
class CategoriesCluster:
    uuid: str
    name: str
    fs_path: str
    filter_category: str
    root_category: str
    
    def toJson(self) -> str:
        self_dict = asdict(self)
        
        json_str = json.dumps(self_dict)
        
    def toDict(self) -> dict[str, Any]:
        return asdict(self)