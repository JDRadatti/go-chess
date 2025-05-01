# syntax=docker/dockerfile:1

# Stage 1: Build frontend
FROM node:20 AS frontend
WORKDIR /app
COPY app/package*.json ./app/
RUN cd app && npm install
COPY app ./app
RUN cd app && npm run build

# Stage 2: Build Go backend
FROM golang:latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
COPY --from=frontend /app/app/dist ./app/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o ./chess cmd/main.go

EXPOSE 8080
CMD ["./chess", "--addr", ":8080"]
