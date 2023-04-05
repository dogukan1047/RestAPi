# İmajımızın alacağı temel imaj
FROM golang:1.2.0


WORKDIR /app

RUN apk update && \
    apk add libc-dev && \
    apk add gcc && \
    apk add make

COPY go.mod ./
COPY go.sum ./

# Go mod dosyalarını kullanarak bağımlılıkları yükle
RUN go mod download && go mod verify
RUN go get github.com/githubnemo/CompileDaemon

# Tüm uygulama dosyalarını kopyala
COPY . .
COPY ./entrypoint.sh /entrypoint.sh

ADD https://raw.githubusercontent.com/eficode/wait-for/v2.1.0/wait-for /usr/local/bin/wait-for
RUN chmod +rx /usr/local/bin/wait-for /entrypoint.sh


ENTRYPOINT [ "sh", "/entrypoint.sh" ]