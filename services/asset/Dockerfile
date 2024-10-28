# Build Stage
FROM golang:1.22.4-alpine AS build

WORKDIR /app

# Zaman dilimini ayarla
RUN apk add --no-cache tzdata
ENV TZ=Europe/Istanbul

# Kök dizindeki go.mod ve go.sum dosyalarını kopyala
COPY ../../go.mod ../../go.sum ./
RUN go mod download

# Tüm proje dosyalarını kopyala (pkg ve asset servis dosyaları dahil)
COPY ../../ .

# Uygulamayı derle
RUN go build -o asset-service ./services/asset/main.go

# Run Stage
FROM alpine:latest

# Zaman dilimini ayarla
RUN apk add --no-cache tzdata
ENV TZ=Europe/Istanbul

WORKDIR /app

# Derlenmiş ikili dosyayı kopyala
COPY --from=build /app/asset-service .

# Portu tanımla
EXPOSE 8000

# Uygulamayı başlat
CMD ["./asset-service"]
