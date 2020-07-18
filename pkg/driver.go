package pkg

import (
	"errors"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"k8s.io/klog"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

type MinioDriver struct {
	Driver            *csicommon.CSIDriver
	Endpoint          string
	maxVolumesPerNode int64
	IDServer          csi.IdentityServer
	NodeServer        csi.NodeServer
	ControllerServer  csi.ControllerServer
}

func NewMinioDriver(driverName, nodeID, endpoint string, maxVolumesPerNode int64, version string) (*MinioDriver, error) {

	if driverName == "" {
		return nil, errors.New("no driver name provided")
	}

	if nodeID == "" {
		return nil, errors.New("no node id provided")
	}

	if endpoint == "" {
		return nil, errors.New("no driver endpoint provided")
	}

	minioDriver := &MinioDriver{Endpoint: endpoint}
	minioDriver.Driver = csicommon.NewCSIDriver(driverName, version, nodeID)

	minioDriver.Driver.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
	})

	minioDriver.Driver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER})

	// minio client
	minioClient := NewMinioClient()

	// identity server
	minioDriver.IDServer = NewIdentityServer(minioDriver.Driver)
	// node server
	minioDriver.NodeServer = NewNodeServer(minioDriver.Driver, minioClient)
	// controller server
	minioDriver.ControllerServer = NewControllerServer(minioDriver.Driver, minioClient)

	return minioDriver, nil

}

func (d *MinioDriver) Run() {

	klog.Infof("Start UDisk CSI Driver, version: %v", CSI_VERSION)
	s := csicommon.NewNonBlockingGRPCServer()
	s.Start(d.Endpoint, d.IDServer, d.ControllerServer, d.NodeServer)
	go cleanUpOnTermination(s)

	s.Wait()
}


func cleanUpOnTermination(s csicommon.NonBlockingGRPCServer) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	klog.Infof("Recieve signal %+v, will stop myself gracefully", sig)
	s.Stop()
	cmd := exec.Command("umount", PluginCSISockDir)
	err := cmd.Run()
	if err != nil {
		klog.Errorf("Umount %s failed, %v", PluginCSISockDir, err)
	}
	os.Exit(0)
}
