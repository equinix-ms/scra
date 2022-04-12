package cmd

import (
	"github.com/equinix-ms/scra/internal/runtimes/containerd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	watchCmd = &cobra.Command{
		Use:   "watch",
		Short: "audit all containers once",
		Run:   watch,
	}
)

func init() {
	rootCmd.AddCommand(watchCmd)
}

func watch(cmd *cobra.Command, args []string) {
	a, err := containerd.NewAuditor(viper.GetString("containerd-address"), viper.GetString("root-prefix"))
	if err != nil {
		panic(err)
	}

	err = a.Watch()
	if err != nil {
		panic(err)
	}
}
