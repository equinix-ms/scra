package cmd

import (
	"context"
	"fmt"

	"github.com/equinix-ms/scra/internal/runtimes/containerd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, gctx := errgroup.WithContext(ctx)
	rootPrefix := viper.GetString("root-prefix")

	for _, address := range viper.GetStringSlice("containerd-address") {
		a, err := containerd.NewAuditor(address, rootPrefix, gctx)
		if err != nil {
			panic(err)
		}

		group.Go(func() error {
			err = a.AuditOnce()
			if err != nil {
				return err
			}

			return nil
		})
	}

	err := group.Wait()
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
