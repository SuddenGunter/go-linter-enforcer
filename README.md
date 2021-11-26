# go-linter-enforcer
This simple app checks if provided git repos are following your golangci-lint config and if not - updates repo config and pushes update to a remote branch

## Current status
Very early version, absolutely not production ready. If it deletes all your production repos - that's on you.

## Modes
Architecture of the app allows writing integrations with different Git-based hosting providers independently. 
Currently, the only supported mode is `bitbucket`.

## Demo

Create some default user and repo in your bitbucket cloud. Your demo repository must have it's `Repository details > Advanced > Language` set to  `Go`. 

Set your env variables:
```env
GIT_SSH_PRIVATE_KEY_PATH=/home/gunter/.ssh/PRIVATE_KEY
SSH_PRIVATE_KEY_PASSWORD=YOUR_PASSWORD # optional
GIT_USERNAME=bitbucketuser42
GIT_EMAIL=user@bitbucket
BITBUCKET_LOGIN=bitbucketuser42 # same as your GIT_USERNAME, temporary inconvenience that could be fixed later
BITBUCKET_APP_PASSWORD=YOUR_BITBUCKET_APP_PASSWORD
BITBUCKET_ORGANIZATION=your_org
MODE=bitbucket
```

Replace `example.golangci.yaml` with your preferred config.

Build & run:
```cgo
go build -o app .
./app
```