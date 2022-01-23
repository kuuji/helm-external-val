package util

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func GetK8sClient() kubernetes.Interface {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func GetConfigMap(namespace string, name string, clientset kubernetes.Interface) (*v1.ConfigMap, error) {

	cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func ComposeValues(configmap *v1.ConfigMap) (yaml string) {
	yaml = configmap.Data["values.yaml"]
	return yaml
}

func WriteValuesToFile(yaml string, output string) {
	err := ioutil.WriteFile(output, []byte(yaml), 0600)
	if err != nil {

		log.Fatal(err)
	}
}
