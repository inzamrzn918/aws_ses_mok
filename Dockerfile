FROM golang:latest

WORKDIR /app
COPY . .

ENV CGO_ENABLED=1 
ENV GOOS=linux

ENV GOARCH=amd64

RUN go mod tidy
RUN go build -o server main.go


CMD ["./server"]
