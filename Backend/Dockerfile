FROM golang:1.19 AS build

WORKDIR /lb-demo-server
ADD . .
RUN  CGO_ENABLED=0 GOOS=linux go build -o /server main.go

CMD [ "/server" ]
