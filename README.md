# Just a course table crawler api server for now.

You can just clone the repo and run the commands below to run the API server
```bash
go mod download
go run main.go -address <server listen address> -port <server listen port>
               -loginurl <login page url> -homeurl <home page url>
```