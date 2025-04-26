import httpx
from httpx import _types as http_types
from curl_cffi import requests as cc_requests, Response

def _getCommonBrowserHeaders() -> dict[str, str]:
    common_headers = {
            'accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8',
            'accept-language': 'en-US,en;q=0.9',
            'cache-control': 'no-cache',
            'pragma': 'no-cache',
            'priority': 'u=0, i',
            'referer': 'https://boards.4chan.org/b/catalog',
            'sec-ch-ua': '"Brave";v="135", "Not-A.Brand";v="8", "Chromium";v="135"',
            'sec-ch-ua-mobile': '?0',
            'sec-ch-ua-platform': '"Linux"',
            'sec-fetch-dest': 'document',
            'sec-fetch-mode': 'navigate',
            'sec-fetch-site': 'same-origin',
            'sec-fetch-user': '?1',
            'sec-gpc': '1',
            'upgrade-insecure-requests': '1',
            'user-agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36'
        }

    return common_headers

def disguiseHttpRequest(request_url: str, with_user_agent:bool = True, follow_redirects: bool = True) -> httpx.Response:
    """Abstracts protocols to disguise GET request and thereby avoid basic request filtering an blocking protocols.
    

    Args:
        request_url (str): the url the GET request will be sent to.
        with_user_agent (bool, optional): Whether to use a fake user agent. Defaults to True.
        follow_redirects (bool, optional): Whether to follow 3xx responses to the redirected url. Defaults to True.
    """

    fake_headers: None | http_types.HeaderTypes = None
    
    print(f"Requesting to: {request_url}")
    
    if with_user_agent: 
        fake_headers = _getCommonBrowserHeaders()

    print(f"headers: {fake_headers}")

    response = httpx.get(request_url, headers=fake_headers, follow_redirects=follow_redirects)
    
    response.request.headers
    
    return response

def impersonateTLSFingerPrint(request_url: str, with_user_agent:bool = True, follow_redirects: bool = True, impersonate: cc_requests.BrowserTypeLiteral = "chrome110") -> Response:
    """Abstracts protocols to disguise GET request and thereby avoid TLS fingerprinting detection.
    
    Args:
        request_url (str): the url the GET request will be sent to.
        with_user_agent (bool, optional): Whether to use a fake user agent. Defaults to True.
        follow_redirects (bool, optional): Whether to follow 3xx responses to the redirected url. Defaults to True.
        impersonate (cc_requests.BrowserTypeLiteral, optional): The browser type to impersonate. Defaults to "chrome110".
    """

    fake_headers = None
    
    if with_user_agent: 
        fake_headers = _getCommonBrowserHeaders()
    
    response = cc_requests.get(request_url, headers=fake_headers, allow_redirects=follow_redirects, impersonate=impersonate)
    
    return response