package pkg

const (
	CSI_DRIVER_NAME  = "minio.csi.ucloud.cn"
	CSI_VERSION      = "1.0.0"
	PluginCSISockDir = "/var/lib/kubelet/plugins/minio.csi.ucloud.cn"
)

const (
	KiB int64 = 1024
	MiB int64 = KiB * 1024
	GiB int64 = MiB * 1024
	TiB int64 = GiB * 1024
)
