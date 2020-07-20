package pkg

import (
	"errors"
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// get pvc yaml
func GetPvcYaml(pvcName string, pvcNamespace string) error{

	if pvcName == "" || pvcNamespace == "" {
		return errors.New("params error")
	}
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil { return err}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil { return err}


	pvcYaml, err := clientset.CoreV1().PersistentVolumeClaims(pvcNamespace).Get(pvcName, v1.GetOptions{})
	if err != nil { return err }

	fmt.Println("<--------*********---------pvcYaml>", pvcYaml)
	fmt.Println("<---- pvc annotations------->",pvcYaml.Annotations)
	return nil
}
