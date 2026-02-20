package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	ciliumClientset "github.com/cilium/cilium/pkg/k8s/client/clientset/versioned"
	"github.com/gogo/protobuf/proto"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	// k8sClient, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	ciliumClientset, err := ciliumClientset.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	nodes, err := ciliumClientset.CiliumV2().CiliumNodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("node: %+v\n", nodes.Items[0])

	fmt.Println("deleted", ciliumClientset.CiliumV2().CiliumNodeSlices().DeleteCollection(context.Background(), metav1.DeleteOptions{IgnoreStoreReadErrorWithClusterBreakingPotential: proto.Bool(true)}, metav1.ListOptions{}))

	curr := v2.CiliumNodeSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cilium-node-slice-1",
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: v2.SchemeGroupVersion.String(),
			Kind:       "CiliumNodeSlice",
		},
	}
	count := 0
	for _, node := range nodes.Items {
		fmt.Println(node.Name)
		curr.Nodes = append(curr.Nodes, ciliumNodeToCiliumNodeSliceItem(node))
		count++
	}

	fmt.Printf("cms: %+v\n", curr)

	cns, err := ciliumClientset.CiliumV2().CiliumNodeSlices().Create(context.Background(), &curr, metav1.CreateOptions{})
	fmt.Println("created", cns, err)
}

func ciliumNodeToCiliumNodeSliceItem(node v2.CiliumNode) v2.CiliumNodeSliceCore {
	return v2.CiliumNodeSliceCore{
		Name:              node.Name,
		InstanceID:        node.Spec.InstanceID,
		BootID:            node.Spec.BootID,
		Addresses:         node.Spec.Addresses,
		HealthAddressing:  node.Spec.HealthAddressing,
		IngressAddressing: node.Spec.IngressAddressing,
		Encryption:        node.Spec.Encryption,
		ENI:               node.Spec.ENI,
		Azure:             node.Spec.Azure,
		AlibabaCloud:      node.Spec.AlibabaCloud,
		IPAM:              node.Spec.IPAM,
		NodeIdentity:      node.Spec.NodeIdentity,
		IPAMStatus:        node.Status.IPAM,
	}
}
