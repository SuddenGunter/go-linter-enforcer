import requests
import os
import sys


def getRequiredEnv(key):
    value = os.getenv(key)
    if not value: 
        sys.exit(key+" is empty")  
    else:
        return value



username = getRequiredEnv("BB_USERNAME")
password = getRequiredEnv("BB_PASSWORD")
organization = getRequiredEnv("BB_ORGANIZATION")

headers = {
    'Accept': 'application/json',
}

response = requests.get('https://api.bitbucket.org/2.0/repositories/'+organization, headers=headers, auth=(username, password))
print(response.json)

