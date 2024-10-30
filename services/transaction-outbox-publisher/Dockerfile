# Build Stage
FROM golang:1.22.4-alpine AS build

WORKDIR /app

RUN apk add --no-cache tzdata
ENV TZ=Europe/Istanbul

COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY ../../ .

RUN go build -o transaction-outbox-publisher ./services/transaction-outbox-publisher

FROM alpine:latest

RUN apk add --no-cache tzdata
ENV TZ=Europe/Istanbul

WORKDIR /app

COPY --from=build /app/transaction-outbox-publisher .

EXPOSE 8000

CMD ["./transaction-outbox-publisher"]
