FROM golang:alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o poetry

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/poetry .
EXPOSE 8080
CMD ["./poetry"]
