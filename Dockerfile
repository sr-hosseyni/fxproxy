FROM golang:latest

WORKDIR /app

RUN go get github.com/githubnemo/CompileDaemon

EXPOSE 8888

ENTRYPOINT CompileDaemon --build="go build -mod=mod fxproxy" --command=./fxproxy
