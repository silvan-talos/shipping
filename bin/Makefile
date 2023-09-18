build:
	go build -o shipping ../cmd/main.go

fmt:
	go fmt ../...

coverage:
	go test -race -coverprofile=cov.out ../...
	go tool cover -html=cov.out

docs:
	swag init -g server.go -d ../http --parseDependency -o ../docs --instanceName ShippingAPI
	swag fmt -d ../http
