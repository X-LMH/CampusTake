# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder

WORKDIR /src

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags="-s -w" -o /out/campusservice ./campusservice.go

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S app -G app

WORKDIR /app

ENV TZ=Asia/Shanghai

COPY --from=builder /out/campusservice /app/campusservice
COPY etc /app/etc
COPY deploy/docker-config.yaml /app/deploy/docker-config.yaml

RUN mkdir -p /app/upload/avatar /app/upload/appeal /app/upload/campus_card \
    && chown -R app:app /app

USER app

EXPOSE 8888

VOLUME ["/app/upload"]

ENTRYPOINT ["/app/campusservice"]
CMD ["-f", "/app/deploy/docker-config.yaml"]
