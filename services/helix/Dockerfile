FROM golang:1.20-alpine AS builder

WORKDIR /build
COPY . .
RUN go build -o app .

FROM scratch
COPY --from=builder /build/app /app
CMD ["/app"]
