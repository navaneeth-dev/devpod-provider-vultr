package vultr

import (
	"context"

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

func (v *Vultr) Create(ctx context.Context, req *govultr.InstanceCreateReq, diskSize int) error {
	// create droplet
	instanceOptions := &govultr.InstanceCreateReq{
		Label:      "awesome-go-app",
		Hostname:   "awesome-go.com",
		Backups:    "enabled",
		EnableIPv6: govultr.BoolToBoolPtr(false),
		ImageID:    "docker",
		Plan:       "vc2-1c-1gb",
		Region:     "blr",
	}

	_, _, err := v.client.Instance.Create(ctx, instanceOptions)

	if err != nil {
		return err
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

// func (d *DigitalOcean) Status(ctx context.Context, name string) (client.Status, error) {
// 	// get droplet
// 	droplet, err := d.GetByName(ctx, name)
// 	if err != nil {
// 		return client.StatusNotFound, err
// 	} else if droplet == nil {
// 		// check if volume exists
// 		volume, err := d.volumeByName(ctx, name)
// 		if err != nil {
// 			return client.StatusNotFound, err
// 		} else if volume != nil {
// 			return client.StatusStopped, nil
// 		}

// 		return client.StatusNotFound, nil
// 	}

// 	// is busy?
// 	if droplet.Status != "active" {
// 		return client.StatusBusy, nil
// 	}

// 	return client.StatusRunning, nil
// }

// func (d *DigitalOcean) GetByName(ctx context.Context, name string) (*godo.Droplet, error) {
// 	droplets, _, err := d.client.Droplets.ListByName(ctx, name, &godo.ListOptions{})
// 	if err != nil {
// 		return nil, err
// 	} else if len(droplets) > 1 {
// 		return nil, fmt.Errorf("multiple droplets with name %s found", name)
// 	} else if len(droplets) == 0 {
// 		return nil, nil
// 	}

// 	return &droplets[0], nil
// }

// func (d *DigitalOcean) Delete(ctx context.Context, name string) error {
// 	// delete volume
// 	volume, err := d.volumeByName(ctx, name)
// 	if err != nil {
// 		return err
// 	} else if volume != nil {
// 		// detach volume
// 		for _, dropletID := range volume.DropletIDs {
// 			_, _, err = d.client.StorageActions.DetachByDropletID(ctx, volume.ID, dropletID)
// 			if err != nil {
// 				return errors.Wrap(err, "detach volume")
// 			}
// 		}

// 		// wait until volume is detached
// 		for len(volume.DropletIDs) > 0 {
// 			time.Sleep(time.Second)

// 			// re-get volume
// 			volume, err = d.volumeByName(ctx, name)
// 			if err != nil {
// 				return err
// 			} else if volume == nil {
// 				break
// 			}
// 		}

// 		// delete volume
// 		if volume != nil {
// 			_, err = d.client.Storage.DeleteVolume(ctx, volume.ID)
// 			if err != nil {
// 				return errors.Wrap(err, "delete volume")
// 			}
// 		}
// 	}

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
