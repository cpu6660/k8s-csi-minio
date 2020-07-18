package pkg

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	csicommon "github.com/kubernetes-csi/drivers/pkg/csi-common"
)

type identityServer struct {
	*csicommon.DefaultIdentityServer
}

func NewIdentityServer(driver *csicommon.CSIDriver) csi.IdentityServer {
	return &identityServer{
		csicommon.NewDefaultIdentityServer(driver),
	}
}
