# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY interverse-candidate/go.mod interverse-candidate/go.sum ./
COPY interverse-contracts /interverse-contracts
RUN go mod edit -replace=github.com/LimeOnTop/interverse-contracts=/interverse-contracts
RUN go mod download

COPY interverse-candidate/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o candidate-service .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
RUN adduser -D -s /bin/sh appuser
USER appuser
WORKDIR /app
COPY --from=builder /app/candidate-service .
EXPOSE 50053
CMD ["./candidate-service"]
