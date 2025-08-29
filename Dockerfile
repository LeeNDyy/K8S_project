FROM golang:1.24-alpine AS builder

WORKDIR /app


COPY ./go.mod ./go.sum ./

RUN go mod download
COPY ./cmd/ ./cmd/

RUN --mount=type=cache,target="/root/.cache/go-build" \
    go build -o monitor ./cmd/

FROM alpine:3.20

WORKDIR /myapp

COPY --from=builder /app/monitor /app/monitor
COPY ./templates ./templates

EXPOSE 7080

CMD [ "/app/monitor" ] 