FROM uhub.service.ucloud.cn/wxyz/csi-basic:1.0

COPY csi-minio /usr/local/bin
RUN chmod +x /usr/local/bin/csi-minio
ENTRYPOINT ["/usr/local/bin/csi-minio"]