/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"helm-external-val/util"

	"github.com/spf13/cobra"
)

var namespace string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "helm-external-val config-map-name",
	Short: "A brief description of your application",
	Args:  cobra.ExactArgs(4),
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		_, ns, cmName, err := ParseUrl(args[3])
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
		client := util.GetK8sClient()
		cm, err := util.GetConfigMap(ns, cmName, client)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
		values := util.ComposeValues(cm)
		fmt.Printf("%s\n", values)
	},
}

func ParseUrl(url string) (protocol string, namespace string, configMapName string, err error) {
	parsedUrl := strings.Split(url, "://")
	protocol = parsedUrl[0]
	err = nil
	if len(parsedUrl) < 2 {
		err = errors.New(":// missing after protocol")
		return
	}
	config := strings.Split(parsedUrl[1], "/")
	if config[0] == "" {
		err = errors.New("no config provided after protocol")
		return
	} else if len(config) == 1 {
		namespace = "default"
		configMapName = config[0]
	} else {
		namespace = config[0]
		configMapName = config[1]
	}
	return
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.helm-external-val.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "name of the config map to fetch")
}
