package cmd

import (
	"github.com/equinix-ms/scra/internal/runtimes/containerd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	onceCmd = &cobra.Command{
		Use:   "once",
		Short: "audit all containers once",
		Run:   once,
	}
)

func init() {
	rootCmd.AddCommand(onceCmd)
}

func once(cmd *cobra.Command, args []string) {
	a, err := containerd.NewAuditor(viper.GetString("containerd-address"), viper.GetString("root-prefix"))
	if err != nil {
		panic(err)
	}

	err = a.AuditOnce()
	if err != nil {
		panic(err)
	}
}
