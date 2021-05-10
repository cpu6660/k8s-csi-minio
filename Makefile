IMAGE=csi-minio:1.11
build:
	GOOS=linux GOARCH=amd64 go build  -o csi-minio cmd/main.go
	docker build -t uhub.service.ucloud.cn/safehouse/${IMAGE}  .
	docker push uhub.service.ucloud.cn/safehouse/${IMAGE}