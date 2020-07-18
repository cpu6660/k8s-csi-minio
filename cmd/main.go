package main

import (
	"flag"
	"fmt"
	"github.com/cpu6660/k8s-csi-minio/pkg"
	"k8s.io/klog"
	"os"
	"runtime"
)

var (
	CSIDriverName = flag.String("drivername", "", "name of the driver")
	CSIEndpoint  = flag.String("endpoint", "unix://var/run/csisocket", "CSI endpoint")
	CSIHostName  = flag.String("hostname", os.Getenv("KUBE_NODE_NAME"), "node id")
	CSIMaxVolume = flag.Int64("maxvolume", 100, "max volumes per node")
)

func showVersion() {
	fmt.Println("CSI Version:\t\t\t" + pkg.CSI_VERSION)
	fmt.Println("Go Version:\t\t\t" + runtime.Version())
	fmt.Printf("Go OS/Arch: \t\t\t%s/%s\n", runtime.GOOS, runtime.GOARCH)
}

func main() {

	flag.Parse()

	showVersion()

	klog.Infof("UCloud Kubernetes CSI Driver(%v) start with args %+v, Endpoint %v, HostName %v, MaxVolume %v",
		pkg.CSI_DRIVER_NAME, os.Args, *CSIEndpoint, *CSIHostName, *CSIMaxVolume)

	driver, err := pkg.NewMinioDriver(pkg.CSI_DRIVER_NAME, *CSIHostName, *CSIEndpoint, *CSIMaxVolume, pkg.CSI_VERSION)
	if err != nil {
		panic(err)
	}
	driver.Run()

}
