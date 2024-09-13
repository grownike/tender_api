FROM golang:1.22.1



WORKDIR /app

RUN docker-compose down -v
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o main main.go


ENTRYPOINT ["./main"]