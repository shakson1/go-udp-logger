FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY *.go ./
COPY *.mod ./
RUN go build -o /logger

FROM alpine

COPY --from=builder /logger /logger

CMD ["/logger"]