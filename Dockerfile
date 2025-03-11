FROM golang:1.24-alpine AS builder

ENV GOPROXY=https://proxy.golang.org,direct
ENV PYTHONUNBUFFERED=1

WORKDIR /app

COPY cmd cmd
COPY internal internal
COPY go.mod go.mod

RUN go mod tidy
RUN go build -o app ./cmd


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 2112

CMD ["./app"]