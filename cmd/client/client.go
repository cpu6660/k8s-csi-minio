package main

import (
	"context"
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("unix:///tmp/csi.sock", grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	// identity client
	identityClient := csi.NewIdentityClient(conn)
	request := csi.GetPluginCapabilitiesRequest{}
	res, err := identityClient.GetPluginCapabilities(context.TODO(), &request)
	fmt.Println(res, err)

	// controller client
	controllerClient := csi.NewControllerClient(conn)
	createVolumeRequest := csi.CreateVolumeRequest{}
	createVolumeRes, err := controllerClient.CreateVolume(context.TODO(), &createVolumeRequest)
	fmt.Println(createVolumeRes, err)

	// node client

}
