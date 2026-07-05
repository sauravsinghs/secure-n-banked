FROM golang:1.25rc2-alpine3.22 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/main .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

# Make shell scripts executable
RUN chmod +x start.sh wait-for.sh

EXPOSE 8080 9090
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]
