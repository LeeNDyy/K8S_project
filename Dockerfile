FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

RUN --mount=type=cache,target="/root/.cache/go-build" \
    go build -o /app/applinux ./cmd/

FROM alpine:3.20

WORKDIR /myapp

COPY --from=builder /app/applinux ./
COPY --from=builder /app/web/page /myapp/web/page
COPY --from=builder /app/web/static /myapp/web/static
# COPY --from=builder /app/web ./

RUN chmod +x /myapp/applinux

EXPOSE 7080

CMD [ "./applinux" ] 