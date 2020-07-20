package pkg

import (
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"github.com/minio/minio-go"
	"golang.org/x/net/context"
	"k8s.io/klog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type controllerServer struct {
	*csicommon.DefaultControllerServer
	minioClient *minio.Client
}

func NewControllerServer(driver *csicommon.CSIDriver, minioClient *minio.Client) *controllerServer {
	return &controllerServer{
		DefaultControllerServer: csicommon.NewDefaultControllerServer(driver),
		minioClient:             minioClient,
	}
}

// create minio bucket
func (cs *controllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {

	fmt.Println("<<<<<<<-req-params->>>>>>>>", req.GetParameters())
	fmt.Println("<<<<<<<-ctx->>>>>>>>>>>>>>", ctx)
	volumeParams := req.GetParameters()
	pvcName := volumeParams[CSI_PVC_NAME]
	pvcNameSpace := volumeParams[CSI_PVC_NAMESPACE]

	fmt.Println(GetPvcYaml(pvcName,pvcNameSpace))


	// Generate desired ufile bucket name
	volumeName := sanitizeVolumeID(req.GetName())
	// Check arguments
	if len(volumeName) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Missing Name in CreateVolumeRequest")
	}
	if req.GetVolumeCapabilities() == nil {
		return nil, status.Error(codes.InvalidArgument, "Missing VolumeCapabilities in CreateVolumeRequest")
	}
	if err := cs.Driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME); err != nil {
		klog.Infof("Invalid create volume req: %v", req)
		return nil, err
	}

	klog.Infof("Request to create bucket volume %s", volumeName)

	// todo get bucket name from parameters
	//bucketName := req.GetParameters()["bucket"]
	//if len(bucketName) == 0 {
	//	return nil, status.Error(codes.Internal, "Missing bucket in parameters")
	//}

	bucketName := volumeName


	// check bucket is exist or not
	exist, err := cs.minioClient.BucketExists(bucketName)
	if err != nil || exist {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to check bucket volume %s: %v, or bucket has exist", bucketName, err))
	}

	// create bucket
	err = cs.minioClient.MakeBucket(bucketName, "us-east-1")
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create bucket volume %s: %v", bucketName, err))

	}

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId: bucketName,
			// In fact, a ufile bucket has no storage capacity limitation, we assume it as 256 TiB
			CapacityBytes: 256 * TiB,
			VolumeContext: req.GetParameters(),
		},
	}, nil
}

func (cs *controllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {

	bucketName := req.GetVolumeId()
	// Check arguments
	if len(bucketName) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Missing VolumeId in DeleteVolumeRequest")
	}

	// check bucket is exist or not
	exist, err := cs.minioClient.BucketExists(bucketName)
	if err != nil || exist {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to check bucket volume %s: %v, or bucket has exist", bucketName, err))
	}

	// delete bucket
	err = cs.minioClient.RemoveBucket(bucketName)

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete bucket volume %s: %v", bucketName, err))

	}

	return &csi.DeleteVolumeResponse{}, nil

}

func (cs *controllerServer) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	return &csi.ControllerPublishVolumeResponse{}, nil
}

func (cs *controllerServer) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	return &csi.ControllerUnpublishVolumeResponse{},nil
}

func (cs *controllerServer) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {

	// Check arguments
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Missing VolumeId in ValidateVolumeCapabilitiesRequest")
	}
	if req.GetVolumeCapabilities() == nil {
		return nil, status.Error(codes.InvalidArgument, "Missing VolumeCapabilities in ValidateVolumeCapabilitiesRequest")
	}

	exist, err := cs.minioClient.BucketExists(req.GetVolumeId())
	if err != nil || exist {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to check bucket volume %s: %v, or bucket has exist", req.GetVolumeId(), err))
	}

	supportedAccessMode := &csi.VolumeCapability_AccessMode{
		Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
	}

	for _, cap := range req.VolumeCapabilities {
		if cap.GetAccessMode().GetMode() != supportedAccessMode.GetMode() {
			return &csi.ValidateVolumeCapabilitiesResponse{Message: fmt.Sprintf("Unsupported AccessMode:%v", cap.AccessMode.GetMode())}, nil
		}
	}

	return &csi.ValidateVolumeCapabilitiesResponse{
		Confirmed: &csi.ValidateVolumeCapabilitiesResponse_Confirmed{
			VolumeCapabilities: []*csi.VolumeCapability{
				{
					AccessMode: supportedAccessMode,
				},
			},
		},
	}, nil
}

func (cs *controllerServer) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {

	return &csi.ControllerExpandVolumeResponse{}, status.Error(codes.Unimplemented, "ControllerExpandVolume is not implemented")

}

func (cs *controllerServer) ControllerGetVolume(ctx context.Context, req *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {

	return &csi.ControllerGetVolumeResponse{}, nil
}
