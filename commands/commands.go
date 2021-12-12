package commands

import (
	"github.com/spf13/cobra"
	"log"
)

var (
	rootCmd = &cobra.Command{
		Use:     "klector",
		Example: "klector --help",
		Short:   "klector cli",
	}
	runCmd = &cobra.Command{
		Use:     "run",
		Example: "klector run",
		Short:   "Start klector",
	}
)

func Execute() int {
	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		return 1
	}

	return 0
}
