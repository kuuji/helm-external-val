package util

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	clientInstance Client
)

type Client struct {
	Clientset kubernetes.Interface
}

func GetK8sClient() Client {
	if clientInstance.Clientset == nil {
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		}

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			// attempt to use in cluster if failed to get kubeconfig
			config, err = rest.InClusterConfig()
			if err != nil {
				panic(err.Error())
			}
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		clientInstance = Client{
			Clientset: clientset,
		}
	}
	return clientInstance
}
