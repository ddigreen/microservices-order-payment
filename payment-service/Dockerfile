FROM golang:alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o payment-app ./cmd/payment-service/main.go
CMD ["./payment-app"]