FROM golang:1.19.2-alpine3.15 AS builder

LABEL build="builder-go-eventsam"

RUN apk add git

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -v -o /app/application /app/cmd/mysql

FROM alpine:3.15

RUN apk --no-cache add tzdata curl

COPY --from=builder /app/ /app/

WORKDIR /app

CMD ["/app/application"]

