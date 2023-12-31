package vultr

import (
	"context"
	"fmt"
	"time"

	"github.com/loft-sh/devpod/pkg/client"
	"github.com/pkg/errors"
	"github.com/vultr/govultr/v3"
	"golang.org/x/oauth2"
)

func NewVultr(token string, ctx context.Context) *Vultr {
	config := &oauth2.Config{}
	ts := config.TokenSource(ctx, &oauth2.Token{AccessToken: token})
	return &Vultr{
		client: govultr.NewClient(oauth2.NewClient(ctx, ts)),
	}
}

type Vultr struct {
	client *govultr.Client
}

func (v *Vultr) Init(ctx context.Context) error {
	_, _, _, err := v.client.Instance.List(ctx, &govultr.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "list droplets")
	}

	return nil
}

func (v *Vultr) Create(ctx context.Context, reqOpts *govultr.InstanceCreateReq, diskSize int) error {
	// create droplet
	instance, _, err := v.client.Instance.Create(ctx, reqOpts)
	if err != nil {
		return err
	}

	// wait till active, to fix NotFound error in status
	for instance.Status != "active" {
		instance, _, err = v.client.Instance.Get(ctx, instance.ID)
		if err != nil {
			return err
		}

		fmt.Println("Waiting for pending instance:", instance.Status)
		time.Sleep(time.Second)
	}

	return nil
}

// func (d *DigitalOcean) volumeByName(ctx context.Context, name string) (*godo.Volume, error) {
// 	volumes, _, err := d.client.Storage.ListVolumes(ctx, &godo.ListVolumeParams{Name: name})
// 	if err != nil {
// 		return nil, err
// 	} else if len(volumes) > 1 {
// 		return nil, fmt.Errorf("multiple volumes with name %s found", name)
// 	} else if len(volumes) == 0 {
// 		return nil, nil
// 	}

// 	return &volumes[0], nil
// }

// func (d *DigitalOcean) Stop(ctx context.Context, name string) error {
// 	droplet, err := d.GetByName(ctx, name)
// 	if err != nil {
// 		return err
// 	} else if droplet == nil {
// 		return nil
// 	}

// 	_, err = d.client.Droplets.Delete(ctx, droplet.ID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (v *Vultr) Status(ctx context.Context, name string) (client.Status, error) {
	// get instance
	instance, err := v.GetByName(ctx, name)
	// instance nil, so error
	if err != nil {
		return client.StatusNotFound, err
	}

	// is busy?
	if instance.Status == "pending" {
		return client.StatusBusy, nil
	}

	if instance.Status == "active" {
		return client.StatusRunning, nil
	}

	return client.StatusNotFound, fmt.Errorf("instance status error: %v", client.StatusRunning)
}

func (v *Vultr) GetByName(ctx context.Context, name string) (*govultr.Instance, error) {
	listOptions := &govultr.ListOptions{}
	for {
		instances, meta, _, err := v.client.Instance.List(ctx, listOptions)
		if err != nil {
			return nil, err
		}
		for _, instance := range instances {
			if instance.Label == name {
				return &instance, nil
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			listOptions.Cursor = meta.Links.Next
			continue
		}
	}

	return nil, fmt.Errorf("instance with name %s not found", name)
}

func (v *Vultr) Delete(ctx context.Context, name string) error {
	instance, err := v.GetByName(ctx, name)
	if err != nil {
		return err
	} else if instance == nil {
		return nil
	}

	err = v.client.Instance.Delete(ctx, instance.ID)
	if err != nil {
		return err
	}

	return nil
}
