package pkg

import (
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// get pvc yaml
func GetPvcYaml(pvcName string, pvcNamespace string) (*corev1.PersistentVolumeClaim,error) {

	if pvcName == "" || pvcNamespace == "" {
		return nil,errors.New("params error")
	}
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil,err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil,err
	}

	pvcYaml, err := clientset.CoreV1().PersistentVolumeClaims(pvcNamespace).Get(pvcName, v1.GetOptions{})
	if err != nil {
		return nil,err
	}

	fmt.Println("<--------*********---------pvcYaml>", pvcYaml)
	fmt.Println("<---- pvc annotations------->", pvcYaml.Annotations)
	return pvcYaml,nil
}
