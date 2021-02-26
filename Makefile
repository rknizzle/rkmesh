test:
	go test -v ./model/... ./auth/...
start-local:
	docker-compose up
