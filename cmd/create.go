package cmd

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/loft-sh/devpod/pkg/log"
	"github.com/loft-sh/devpod/pkg/ssh"
	"github.com/navaneeth-dev/devpod-provider-vultr/pkg/options"
	"github.com/navaneeth-dev/devpod-provider-vultr/pkg/vultr"
	"github.com/spf13/cobra"
	"github.com/vultr/govultr/v3"
)

// CreateCmd holds the cmd flags
type CreateCmd struct{}

// NewCreateCmd defines a command
func NewCreateCmd() *cobra.Command {
	cmd := &CreateCmd{}
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create an instance",
		RunE: func(_ *cobra.Command, args []string) error {
			options, err := options.FromEnv(false)
			if err != nil {
				return err
			}

			return cmd.Run(context.Background(), options, log.Default)
		},
	}

	return createCmd
}

// Run runs the command logic
func (cmd *CreateCmd) Run(ctx context.Context, options *options.Options, log log.Logger) error {
	fmt.Println("ok")
	req, err := buildInstance(options)
	if err != nil {
		return err
	}

	fmt.Println("ok")

	// diskSize, err := strconv.Atoi(options.DiskSize)
	// if err != nil {
	// 	return errors.Wrap(err, "parse disk size")
	// }

	return vultr.NewVultr(options.Token, ctx).Create(ctx, req, 69)
}

func GetInjectKeypairScript(machineFolder, machineID string) (string, error) {
	publicKeyBase, err := ssh.GetPublicKeyBase(machineFolder)
	if err != nil {
		return "", err
	}

	publicKey, err := base64.StdEncoding.DecodeString(publicKeyBase)
	if err != nil {
		return "", err
	}

	resultScript := `#!/bin/sh
# Create DevPod user and configure ssh
useradd devpod -d /home/devpod
if grep -q sudo /etc/groups; then
	usermod -aG sudo devpod
elif grep -q wheel /etc/groups; then
	usermod -aG wheel devpod
fi
echo "devpod ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/91-devpod
mkdir -p /home/devpod/.ssh
echo '` + string(publicKey) + `' > /home/devpod/.ssh/authorized_keys
chmod 0700 /home/devpod/.ssh
chmod 0600 /home/devpod/.ssh/authorized_keys
chown devpod:devpod /home/devpod
chown -R devpod:devpod /home/devpod/.ssh

# Make sure we don't get limited
ufw allow 22/tcp || true
`

	return resultScript, nil
}

func buildInstance(options *options.Options) (*govultr.InstanceCreateReq, error) {
	// generate ssh keys
	userData, err := GetInjectKeypairScript(options.MachineFolder, options.MachineID)
	if err != nil {
		return nil, err
	}

	// generate instance object
	instance := &govultr.InstanceCreateReq{
		Label:      options.MachineID,
		Hostname:   "awesome-go.com",
		Backups:    "disabled",
		EnableIPv6: govultr.BoolToBoolPtr(false),
		OsID:       362,
		Plan:       "vc2-1c-2gb",
		Region:     "blr",
		UserData:   userData,
		Tags:       []string{"devpod"},
	}

	return instance, nil
}
