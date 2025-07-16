test:
	go test ./... --cover

check-sec:
	gosec -exclude=G101,G104 ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

checks:
	make test
	make check-sec
	make fmt
	make lint

build-docker:
	mkdir data
	chmod +x ./scripts/migrate_up.sh && ./scripts/migrate_up.sh
	docker build --tag mrramonster/pantry_pal:latest .

push-docker:
	docker push mrramonster/pantry_pal:latest

build-push:
	make build-docker && make push-docker

