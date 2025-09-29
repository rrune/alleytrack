FROM golang:latest
RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app/
RUN go build -o main ./cmd/alleytrack
CMD ["/app/main"]