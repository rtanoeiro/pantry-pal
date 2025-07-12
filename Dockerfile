FROM golang:tip-alpine3.22 AS builder

WORKDIR /app

RUN apk update && apk add --no-cache build-base && apk add sqlite

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o pantry_pal

# Final image
FROM alpine:3.22.0
WORKDIR /app
COPY --from=builder /app/pantry_pal .

COPY data ./data
COPY static ./static
COPY css ./css

RUN apk update && apk add sqlite

EXPOSE 8080
CMD ["./pantry_pal"]