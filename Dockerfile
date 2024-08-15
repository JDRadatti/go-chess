# syntax=docker/dockerfile:1

FROM golang:latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
ADD cmd ./cmd
ADD app ./app
ADD internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -o ./chess cmd/main.go
EXPOSE 8080
CMD ["./chess", "--addr", ":8080"]
