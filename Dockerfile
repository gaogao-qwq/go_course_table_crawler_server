FROM golang:1.20

ENV GO111MODULE="on" \
    GOPROXY="https://goproxy.cn,direct" \
    GIN_MODE="release" \
    SERVER_ADDRESS="0.0.0.0" \
    SERVER_PORT="56789" \
    CRAWLER_LOGIN_URL="http://jw.gzgs.edu.cn/eams/login.action" \
    CRAWLER_HOME_URL="http://jw.gzgs.edu.cn/eams/home.action"

# 更换国内源
COPY ./sources.list /etc/apt/sources.list

RUN apt-get update && apt-get install -y chromium

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/bin/course_table_crawler

EXPOSE 56789

CMD ["/bin/sh", "-c", "/usr/bin/course_table_crawler", "-address", "$SERVER_ADDRESS", "-port", "$SERVER_PORT", "-loginurl", "$CRAWLER_LOGIN_URL", "-homeurl", "$CRAWLER_HOME_URL"]