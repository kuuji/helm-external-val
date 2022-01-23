package util

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"sigs.k8s.io/yaml"
)

func TestWriteValuesToFile(t *testing.T) {
	type args struct {
		yaml string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should write a file",
			args: args{
				yaml: "replicas: \"3\"\ndeployments:\n  server:\n    replicas: \"3\"\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteValuesToFile(tt.args.yaml, "values.yaml")
			yfile, err := ioutil.ReadFile("values.yaml")

			if err != nil {

				log.Fatal(err)
			}

			values := make(map[string]interface{})
			yaml.Unmarshal(yfile, &values)
			if values["deployments"].(map[string]interface{})["server"].(map[string]interface{})["replicas"] != "3" {
				t.Errorf("TestWriteValuesToFile()")
			}
			err = os.Remove("values.yaml")
			if err != nil {
				log.Fatal(err)
			}
		})
	}
}
