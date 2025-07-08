FROM golang:tip-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o pantry_pal

# Final image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/pantry_pal .

COPY .env ./
COPY static ./static
COPY css ./css


EXPOSE 8080
CMD ["./pantry_pal"]