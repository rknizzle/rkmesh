test:
	go test -v ./model/... ./auth/... ./user/...
start-local:
	docker-compose up
