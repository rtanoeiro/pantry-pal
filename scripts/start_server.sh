# /bin/bash

# Use this when the applications is dockerized
# CGO_ENABLED=0 GOOS=linux GOARCH=amd64  

go build -o pantry-pal

./pantry-pal