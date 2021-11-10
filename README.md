# go-linter-enforcer
This simple app checks if provided git repos are following your golangci-lint config and if not - updates repo config and pushes update to a remote branch

## Demo
Run gitea with
```bash
docker-compose up -d
```
Create some default user and repo.

Create config file (democonfig.json in this example):
```json
{
  "repositories": [
    {
      "url": "http://localhost:3000/suddengunter/linterdemo.git",
      "name": "linterdemo",
      "mainBranch": "main"
    }
  ]
}
```

Set your env variables:
```env
CONFIG_FILE=democonfig.json
DEMO_PASSWORD=password
DEMO_USERNAME=user
```
Build & run:
```cgo
go build -o app .
./app
```