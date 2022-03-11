package util

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetSecret(namespace string, name string, client Client) (*v1.Secret, error) {

	secret, err := client.Clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func ComposeSecretValues(secret *v1.Secret, dataKey string) (yaml string) {
	yaml = string(secret.Data[dataKey])
	return yaml
}
