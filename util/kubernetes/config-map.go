package util

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetConfigMap(namespace string, name string, client Client) (*v1.ConfigMap, error) {

	cm, err := client.Clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func ComposeValues(configmap *v1.ConfigMap) (yaml string) {
	yaml = configmap.Data["values.yaml"]
	return yaml
}
