/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"helm-external-val/util"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// downloaderCmd represents the downloader command
var downloaderCmd = &cobra.Command{
	Use:   "downloader",
	Short: "Get value from remote source and output it to stdout",
	Long: `Get value from remote source and output it to stdout.
	URL is formatted like below
	<protocol>://<namespace>/<name>

	For example:
	cm://kuuji/helm-values 
	This will get values from a config map named helm-values in the namespace kuuji
	.`,
	Args: cobra.ExactArgs(4),
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

func init() {
	rootCmd.AddCommand(downloaderCmd)
}
