FROM golang:tip-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pantry_pal

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