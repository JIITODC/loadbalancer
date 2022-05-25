FROM golang:1.18 AS build
WORKDIR /lb-demo-server
ADD . .
RUN cd Backend; CGO_ENABLED=0 GOOS=linux go build -o /server main.go

FROM scratch
COPY --from=build /server /
CMD [ "/server" ]
