FROM golang:1.20

RUN apt-get update && apt-get install -y chromium

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/bin/course_table_crawler

ENV SERVER_ADDRESS="0.0.0.0" \
    SERVER_PORT="56789" \
    CRAWLER_LOGIN_URL="http://targeturl/login" \
    CRAWLER_HOME_URL="http://targeturl/home.action"

EXPOSE 56789

CMD ["/bin/sh", "-c", "/usr/bin/course_table_crawler", "-address", "$SERVER_ADDRESS", "-port", "$SERVER_PORT", "-loginurl", "$CRAWLER_LOGIN_URL", "-homeurl", "$CRAWLER_HOME_URL"]