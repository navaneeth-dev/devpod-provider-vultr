package cmd

import (
	"context"

	"github.com/loft-sh/devpod/pkg/log"
	"github.com/navaneeth-dev/devpod-provider-vultr/pkg/options"
	"github.com/navaneeth-dev/devpod-provider-vultr/pkg/vultr"
	"github.com/spf13/cobra"
)

// InitCmd holds the cmd flags
type InitCmd struct{}

// NewInitCmd defines a command
func NewInitCmd() *cobra.Command {
	cmd := &InitCmd{}
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Init an instance",
		RunE: func(_ *cobra.Command, args []string) error {
			// options, err := options.FromEnv(true)
			// if err != nil {
			// 	return err
			// }

			return cmd.Run(context.Background(), &options.Options{}, log.Default)
		},
	}

	return initCmd
}

// Run runs the command logic
func (cmd *InitCmd) Run(ctx context.Context, options *options.Options, log log.Logger) error {
	return vultr.NewVultr(options.Token, ctx).Init(ctx)
}
