build-docker-and-run: build-docker run

build-docker:
	docker build . -t pages-ms:latest

run:
	docker-compose up -d

stop:
	docker-compose down

setup-dummy-db-data:
	docker exec pages-ms-mongo mongoimport -u user -p pass \
	  --drop --authenticationDatabase admin --jsonArray --db test --collection seos \
	  --file /sample-data/sample-seos.json
	docker exec pages-ms-mongo mongoimport -u user -p pass \
	  --drop --authenticationDatabase admin --jsonArray --db test --collection products \
	  --file /sample-data/sample-products.json

test:
	go test ./src/...

build:
	go build -o ./target/pages-ms ./src/main.go