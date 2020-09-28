FROM golang:latest

WORKDIR /app

RUN go get github.com/githubnemo/CompileDaemon

RUN touch /var/log/fxproxy.access.log
RUN touch /var/log/fxproxy.error.log
ADD conf.yaml /etc/fxproxy.yml

EXPOSE 8888

ENTRYPOINT CompileDaemon --build="go build -mod=mod fxproxy" --command=./fxproxy
