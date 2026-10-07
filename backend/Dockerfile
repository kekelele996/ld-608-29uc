FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache wget ca-certificates
COPY --from=build /app/server ./server
EXPOSE 3000
CMD ["./server"]
