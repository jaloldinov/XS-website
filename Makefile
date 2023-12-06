run:
	go run ./cmd/main.go

swag-gen:
	swag init -g cmd/main.go -o internal/controller/docs

deploy:
	GOOS=linux GOARCH=amd64 go build cmd/main.go && scp -P 13197 main root@185.74.7.202:/var/www/senat-back