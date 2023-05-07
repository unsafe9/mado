FROM golang:1.20-alpine AS builder

COPY . .
RUN go build -o app .

FROM scratch
COPY --from=builder app .
ENTRYPOINT ["app"]