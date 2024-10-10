from inspect import getfullargspec
from typing import Callable, Any
from . import chan as __chan_module
from . import requests as __requests_models
from . import categories as __categories_models

def _allArgsPresent(obj: object, args: list[str]) -> bool:
    all_args = getfullargspec(obj).args
    are_all_present = True
    for arg in all_args:
        if arg not in args and arg != "self":
            print(f"Missing argument: {arg} in {obj.__name__}")
            are_all_present = False
            break
        
    return are_all_present

__chan_module.allArgsPresent = _allArgsPresent
__categories_models.allArgsPresent = _allArgsPresent

CatalogThread = __chan_module.CatalogThread
FourChanThread = __chan_module.FourChanThread
FourChanThreadReply = __chan_module.FourChanThreadReply
DungeonsCategory = __categories_models.Category
DungeonsCategoriesCluster = __categories_models.CategoriesCluster

class ChanFunctions:
    unpackFromCatalogData: Callable[[str, dict[str, Any]], CatalogThread]
    
ChanFunctions.unpackFromCatalogData = __chan_module.unpackFromCatalogData

Requests = __requests_models