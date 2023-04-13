.PHONY: build up upd start down destroy stop

build:
	docker compose -f build/dev/docker-compose.yml build

up:
	docker compose -f build/dev/docker-compose.yml up

upd:
	docker compose -f build/dev/docker-compose.yml up -d

start:
	docker compose -f build/dev/docker-compose.yml start

stop:
	docker compose -f build/dev/docker-compose.yml stop

down:
	docker compose -f build/dev/docker-compose.yml down

destroy:
	docker compose -f build/dev/docker-compose.yml down -v
