package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "debug-cli",
	Short:         "debug-cli is a command line tool for debug in k8s clusters",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Exec() {
	rootCmd.AddCommand(agentCmd)
	rootCmd.AddCommand(proxyCmd)
	rootCmd.AddCommand(mockCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
