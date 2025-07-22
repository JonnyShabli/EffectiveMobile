FROM golang:1.23-alpine as go

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app_bin ./cmd/main.go

FROM alpine:3.14

WORKDIR /app

COPY --from=go /app_bin ./
COPY /config/local/* ./config/local/

COPY .env ./

EXPOSE 8080

CMD ["/app/app_bin"]