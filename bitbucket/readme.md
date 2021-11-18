# Bitbucket integration

## gen-repos.py

This script allows you to fetch all your bitbucket repositories, filter repos with Go set as primary language and create `repos.json` file for `go-linter-enforcer`.

## make-prs.py

Creates pull requests for all `lintenforcer/*` branches in all your repositories from `repos.json`.
