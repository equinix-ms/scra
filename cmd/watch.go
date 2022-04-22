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
	logger, err := getLogger()
	if err != nil {
		panic(fmt.Errorf("error setting up zap logging: %v", err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, gctx := errgroup.WithContext(ctx)
	rootPrefix := viper.GetString("root-prefix")

	for _, address := range viper.GetStringSlice("containerd-address") {
		a, err := containerd.NewAuditor(address, rootPrefix, gctx, logger)
		if err != nil {
			panic(err)
		}

		group.Go(func() error {
			err := a.Watch()
			if err != nil {
				return err
			}

			return nil
		})

	}

	err = group.Wait()
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
