package commands

import (
	"github.com/spf13/cobra"
	"io.klector/klector/api"
	"io.klector/klector/storage"
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
		RunE:    runServer,
	}
)

func runServer(cmd *cobra.Command, args []string) error {
	config := updateStorageConfigFromCommandLine(
		storage.NewDefaultStorageConfiguration(),
	)
	storage := storage.Create(config)
	return api.Create(&storage)
}

func updateStorageConfigFromCommandLine(config *storage.StorageConfiguration) *storage.StorageConfiguration {
	return config
}

func Execute() int {
	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		return 1
	}

	return 0
}
