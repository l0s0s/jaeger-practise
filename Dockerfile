FROM golang:1.17-alpine

WORKDIR /app

ARG PORT

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE ${PORT}

CMD ["./main"]