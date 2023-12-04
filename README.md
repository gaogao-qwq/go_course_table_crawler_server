# A course table crawler service implementation with go and chromedp
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fgaogao-qwq%2Fgo_course_table_crawler_server.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fgaogao-qwq%2Fgo_course_table_crawler_server?ref=badge_shield)


## Build from docker (Suggest)

**Change `ENV` in Dockerfile first.**

1. Clone repo
```bash
git clone https://github.com/gaogao-qwq/go_course_table_crawler_server.git
cd go_course_table_crawler_server
mv .env.template .env # Configure environment variables
```

2. Build image and run with .env file
```bash
docker build -t course_table_crawler_image .
docker run --name course_table_crawler_container --env-file .env -v /run/dbus:/run/dbus -p <host port>:56789 course_table_crawler_image
```

## Build from native

**Make sure that Chrome or Chromium is available in current environment**

```bash
go mod download
go run main.go -address <server listen address> -port <server listen port>
               -loginurl <login page url> -homeurl <home page url>
```

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fgaogao-qwq%2Fgo_course_table_crawler_server.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fgaogao-qwq%2Fgo_course_table_crawler_server?ref=badge_large)