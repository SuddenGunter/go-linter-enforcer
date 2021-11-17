# go-linter-enforcer
This simple app checks if provided git repos are following your golangci-lint config and if not - updates repo config and pushes update to a remote branch

## Current status
Very early version, absolutely not production ready. If it deletes all your production repos - that's on you.

## Demo
Run gitea with
```bash
docker-compose up -d
```
Create some default user and repo.

Create repositories list file (repos.json in this example):
```json
{
  "repositories": [
    {
      "url": "http://localhost:3000/gunter/linterdemo.git",
      "name": "linterdemo",
      "mainBranch": "main"
    }
  ]
}
```

Set your env variables:
```env
GIT_PASSWORD=password
GIT_USERNAME=gunter
GIT_EMAIL=email@example.com
```
Build & run:
```cgo
go build -o app .
./app
```