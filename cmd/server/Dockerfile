FROM golang:1.21-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/server ./cmd/server

FROM alpine:3.18
COPY --from=build /usr/local/bin/server /usr/local/bin/server
EXPOSE 8080
CMD ["/usr/local/bin/server"]
