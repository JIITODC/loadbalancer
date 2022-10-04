FROM golang:1.18 AS build

WORKDIR /lb-demo-server
ADD . .
RUN  CGO_ENABLED=0 GOOS=linux go build -o /server main.go

CMD [ "/server" ]
