package main

import (
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	// "github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	_, err = cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: "nginx",
			Volumes: map[string]struct{}{
				"/go/src/github.com/hg2c/example": {},
			},
			ExposedPorts: nat.PortSet{
				"80/tcp": struct{}{},
			},
		},
		&container.HostConfig{
			Binds: []string{
				"/Volumes/Kayle/w/hg2c/golang:/go/src/github.com/hg2c/example:rw",
			},
			PortBindings: nat.PortMap{
				"80/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "80",
					},
				},
			},
			// PortBindings: nat.PortMap
			// PublishAllPorts: true,
			// Mounts: []mount.Mount{
			// 	mount.Mount{
			// 		Type:     "bind",
			// 		Source:   "/Volumes/Kayle/w/hg2c/golang",
			// 		Target:   "/go/src/github.com/hg2c/example",
			// 		ReadOnly: false,
			// 		BindOptions: &mount.BindOptions{
			// 			Propagation: mount.PropagationRPrivate,
			// 		},
			// 	},
			// },
		},
		nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	if _, err = cli.ContainerWait(ctx, resp.ID); err != nil {
		panic(err)
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)
}
