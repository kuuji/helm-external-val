package cmd

import (
	"fmt"
	util "helm-external-val/util"
	k8s "helm-external-val/util/kubernetes"
	"os"

	"github.com/spf13/cobra"
)

var kubeSecretNamespace string
var dataSecretKey string
var secretOutput string

// secretCmd represents the secret command
var secretCmd = &cobra.Command{
	Use:   "secret <name>",
	Short: "Get the content of values from a secret and write it to a file",
	Long:  `Get the content of values from a secret and write it to a file`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		secretName := args[0]
		client := k8s.GetK8sClient()
		secret, err := k8s.GetSecret(kubeSecretNamespace, secretName, client)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
		values := k8s.ComposeSecretValues(secret, dataSecretKey)
		util.WriteValuesToFile(values, secretOutput)
		fmt.Printf("%s written to %s\n", secretName, secretOutput)
	},
}

func init() {
	rootCmd.AddCommand(secretCmd)
	secretCmd.PersistentFlags().StringVar(&kubeSecretNamespace, "kube_namespace", "default", "The namespace to get the secret from")
	secretCmd.PersistentFlags().StringVar(&dataSecretKey, "dataKey", "values.yaml", "The key to get the data from a secret")
	secretCmd.PersistentFlags().StringVarP(&secretOutput, "out", "o", "values-secret.yaml", "The file to output the values to")
}
