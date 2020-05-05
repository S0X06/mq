FROM golang:1.13
RUN mkdir /cmd
ADD . /cmd/
WORKDIR /cmd
RUN go mod download
RUN go build -o main ./...
CMD ["/cmd/main"]
