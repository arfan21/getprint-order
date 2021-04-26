build-dev:
	docker build -f dev.Dockerfile -t getprint-service-order-dev .

build-prod:
	docker build -f prod.Dockerfile -t getprint-service-order-prod .