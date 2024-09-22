# Start from an official Golang image to build the application
FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
# If you need to support arm64, uncomment the following 3 lines
#ENV CGO_ENABLED=0
#ENV GOOS=linux
#ENV GOARCH=arm64
RUN go build -o go-clickup-pagerduty-webhook ./cmd/main.go

# Use a lightweight image for the final stage
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates
WORKDIR /app
COPY --from=builder /app/go-clickup-pagerduty-webhook .
COPY --from=builder /app/config/rules.yaml ./config/rules.yaml
COPY --from=builder /app/config/groups.yaml ./config/groups.yaml
CMD ["./go-clickup-pagerduty-webhook"]
