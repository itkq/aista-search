FROM golang:1.6

RUN go get github.com/codegangsta/gin
WORKDIR /go/src/app
EXPOSE 8080
EXPOSE 3000

CMD ["/go/bin/gin", "-a", "8080"]
