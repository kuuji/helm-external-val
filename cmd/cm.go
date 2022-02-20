package cmd

import (
	"fmt"
	util "helm-external-val/util"
	k8s "helm-external-val/util/kubernetes"
	"os"

	"github.com/spf13/cobra"
)

var kubeNamespace string
var dataKey string
var output string

var cmCmd = &cobra.Command{
	Use:   "cm <name>",
	Short: "Get the content of values from a cm and write it to a file",
	Long:  `Get the content of values from a cm and write it to a file`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmName := args[0]
		client := k8s.GetK8sClient()
		cm, err := k8s.GetConfigMap(kubeNamespace, cmName, client)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
		values := k8s.ComposeValues(cm, dataKey)
		util.WriteValuesToFile(values, output)
		fmt.Printf("%s written to %s\n", cmName, output)
	},
}

func init() {
	rootCmd.AddCommand(cmCmd)
	cmCmd.PersistentFlags().StringVar(&kubeNamespace, "kube_namespace", "default", "The namespace to get the cm from")
	cmCmd.PersistentFlags().StringVar(&dataKey, "dataKey", "values.yaml", "The key to get the data from a cm")
	cmCmd.PersistentFlags().StringVarP(&output, "out", "o", "values-cm.yaml", "The file to output the values to")
}
