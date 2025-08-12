FROM golang:1.24-alpine AS builder
WORKDIR /src
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app ./cmd/web/server.go

FROM alpine:3.18
RUN apk add --no-cache ca-certificates

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=builder /app /home/appuser/app
USER appuser

EXPOSE ${PORT}
ENV LOG_DIR=/home/appuser/logs
ENV VIDEO_DIR=/home/appuser/videos
ENV CONFIG_DIR=/home/appuser/.config
ENV PORT=${PORT}
ENTRYPOINT [ "/home/appuser/app" ]
