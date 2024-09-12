FROM golang:1.23

WORKDIR /usr/src/app

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/air-verse/air@latest


COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o main .

EXPOSE 50052

# CMD ["./main"]
CMD ["air", "-c", ".air.toml"]