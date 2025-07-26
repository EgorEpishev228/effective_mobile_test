FROM golang:1.22 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

COPY .env/config.yaml ./config.yaml

COPY docs/swagger.json ./docs/swagger.json
COPY docs/swagger.yaml ./docs/swagger.yaml

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder app/config.yaml .env/config.yaml
COPY --from=builder app/docs/swagger.json docs/swagger.json
COPY --from=builder app/docs/swagger.yaml docs/swagger.yaml

EXPOSE 8080

CMD ["./app"]