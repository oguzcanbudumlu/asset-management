FROM golang:1.22.4-alpine AS build

WORKDIR /app

RUN apk add --no-cache tzdata
ENV TZ=Europe/Istanbul

COPY ../../go.mod ../../go.sum ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3

COPY ../../ .

RUN swag init --dir ./services/wallet -o services/wallet/docs
RUN go build -o wallet-service ./services/wallet/

# Run Stage
FROM alpine:latest

RUN apk add --no-cache tzdata
ENV TZ=Europe/Istanbul

WORKDIR /app

COPY --from=build /app/wallet-service .

EXPOSE 8000

# Uygulamayı başlat
CMD ["./wallet-service"]
