test:
	sh scripts/migrate_up.sh dev
	go test ./... --cover -v

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
	mkdir -p data
	chmod +x ./scripts/migrate_up.sh prod
	docker build --tag mrramonster/pantry_pal:latest .

push-docker:
	docker push mrramonster/pantry_pal:latest

build-push:
	make build-docker && make push-docker

