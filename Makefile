build:
	GOOS=linux GOARCH=amd64 go build  -o csi-minio cmd/main.go
	docker build -t uhub.service.ucloud.cn/safehouse/csi-minio:1.10  .
	docker push uhub.service.ucloud.cn/safehouse/csi-minio:1.10