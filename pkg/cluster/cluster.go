package cluster

import (
	"context"
	"fmt"
	"github.com/minc-org/minc/pkg/retry"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

func GetPodStatus(kubeConfig []byte) error {
	// Create Kubernetes client
	config, err := clientcmd.RESTConfigFromKubeConfig(kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to build config from kubeconfig bytes: %v", err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	// Define namespaces to check
	namespaces := []string{"kube-flannel", "kube-proxy", "kube-system", "openshift-dns", "openshift-ingress", "openshift-service-ca"}

	for _, ns := range namespaces {
		podStatusFunc := func() error {
			pods, err := clientSet.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				return fmt.Errorf("failed to get pods in namespace %s: %v", ns, err)
			}
			for _, pod := range pods.Items {
				if pod.Status.Phase != "Running" {
					return fmt.Errorf("pod %s in namespace %s is not running. Current status: %s", pod.Name, ns, pod.Status.Phase)
				}
			}
			return nil
		}
		return retry.Retry(podStatusFunc, 5, 2*time.Second)
	}
	return nil
}
