FROM golang:alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o order-app ./cmd/order-service/main.go
CMD ["./order-app"]