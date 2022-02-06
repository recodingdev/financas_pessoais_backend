FROM golang:1.16-alpine3.13

WORKDIR /go/src/app
ADD . .
RUN go mod init
RUN apk add --no-cache ca-certificates && update-ca-certificates
RUN go get github.com/silenceper/gowatch

RUN go get github.com/georgysavva/scany/sqlscan
RUN go get github.com/gorilla/mux
RUN go get github.com/lib/pq

RUN apk add --update curl && apk add --update tar && apk add --no-cache tzdata && rm -rf /var/cache/apk/*
ENV CGO_ENABLED=0
ENV GOOS=linux
EXPOSE 8000
EXPOSE 5432
RUN rm -f /etc/localtime; ln -s /usr/share/zoneinfo/America/Sao_Paulo /etc/localtime;

ENTRYPOINT gowatch