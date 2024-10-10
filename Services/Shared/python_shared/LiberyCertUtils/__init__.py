import os

SSL_CA_PATH = os.getenv('SSL_CA_PATH', "")
assert SSL_CA_PATH != "", "SSL_CA_PATH is not set"

def getSSLcaFile():
    with open(SSL_CA_PATH, 'rb') as f:
        return f.read()