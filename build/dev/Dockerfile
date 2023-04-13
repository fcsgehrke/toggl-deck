FROM golang:1.20-alpine as builder

WORKDIR /app

COPY ../../go.* ./
RUN go mod download

COPY ./ ./
RUN go build -o server /app/cmd/deck/main.go

#--------------------------------------------

FROM alpine:latest

COPY --from=builder /app/server /app/server

CMD ["/app/server"]
