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

makeRequests = True
page = 1

while makeRequests:
    response = requests.get('https://api.bitbucket.org/2.0/repositories/{0}?page={1}'.format(organization, page), headers=headers, auth=(username, password))
    data = response.json()
    if not data.next:
        makeRequests = False
         
    count += 1
    makeRequests = False
