/*
	Copyright 2019 whiteblock Inc.
	This file is a part of the genesis.

	Genesis is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	Genesis is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/whiteblock/genesis/pkg/entity"
	"github.com/whiteblock/genesis/pkg/repository"
	"strconv"
)

//DockerService provides a intermediate interface between docker and the order from a command
type DockerService interface {

	//CreateClient creates a new client for connecting to the docker daemon
	CreateClient(conf entity.DockerConfig, host string) (*client.Client, error)

	//CreateContainer attempts to create a docker container
	CreateContainer(ctx context.Context, cli *client.Client, container entity.Container) entity.Result

	//StartContainer attempts to start an already created docker container
	StartContainer(ctx context.Context, cli *client.Client, name string) entity.Result
	RemoveContainer(ctx context.Context, cli *client.Client, name string) entity.Result
	CreateNetwork(ctx context.Context, cli *client.Client, net entity.Network) entity.Result
	AttachNetwork(ctx context.Context, cli *client.Client, network string, container string) entity.Result
	CreateVolume(ctx context.Context, cli *client.Client, volume entity.Volume) entity.Result
	RemoveVolume(ctx context.Context, cli *client.Client, name string) entity.Result
	PlaceFileInContainer(ctx context.Context, cli *client.Client, containerName string, file entity.File) entity.Result
	PlaceFileInVolume(ctx context.Context, cli *client.Client, volumeName string, file entity.File) entity.Result
	Emulation(ctx context.Context, cli *client.Client, netem entity.Netconf) entity.Result
}

type dockerService struct {
	repo repository.DockerRepository
}

//NewDockerService creates a new DockerService
func NewDockerService(repo repository.DockerRepository) (DockerService, error) {
	return dockerService{repo: repo}, nil
}

//CreateClient creates a new client for connecting to the docker daemon
func (ds dockerService) CreateClient(conf entity.DockerConfig, host string) (*client.Client, error) {
	return client.NewClientWithOpts(
		client.WithAPIVersionNegotiation(),
		client.WithHost(host),
		client.WithTLSClientConfig(conf.CACertPath, conf.CertPath, conf.KeyPath),
	)
}

//CreateContainer attempts to create a docker container
func (ds dockerService) CreateContainer(ctx context.Context, cli *client.Client, dContainer entity.Container) entity.Result {
	var envVars []string
	for key, val := range dContainer.Environment {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, val))
	}

	config := &container.Config{
		Cmd:        dContainer.Args,
		Env:        envVars,
		Image:      dContainer.Image,
		Entrypoint: []string{dContainer.EntryPoint},
		Labels:     dContainer.Labels,
	}

	mem, err := dContainer.GetMemory()
	if err != nil {
		return entity.NewFatalResult(err)
	}

	cpus, err := strconv.ParseFloat(dContainer.Cpus, 64)
	if err != nil {
		return entity.NewFatalResult(err)
	}

	hostConfig := &container.HostConfig{}
	hostConfig.NanoCPUs = int64(1000000000 * cpus)
	hostConfig.Memory = mem

	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: dContainer.NetworkConfig.EndpointsConfig,
	}

	_, err = ds.repo.ContainerCreate(ctx, cli, config, hostConfig, networkConfig, dContainer.Name)
	if err != nil {
		return entity.NewFatalResult(err)
	}

	return entity.NewSuccessResult()
}

//StartContainer attempts to start an already created docker container
func (ds dockerService) StartContainer(ctx context.Context, cli *client.Client, name string) entity.Result {
	opts := types.ContainerStartOptions{}
	err := ds.repo.ContainerStart(ctx, cli, name, opts)
	if err != nil {
		return entity.NewErrorResult(err)
	}

	return entity.NewSuccessResult()
}

func (ds dockerService) RemoveContainer(ctx context.Context, cli *client.Client, name string) entity.Result {
	//TODO
	return entity.Result{}
}

func (ds dockerService) CreateNetwork(ctx context.Context, cli *client.Client, net entity.Network) entity.Result {
	//TODO
	return entity.Result{}
}

func (ds dockerService) AttachNetwork(ctx context.Context, cli *client.Client, network string, containerName string) entity.Result {
	//TODO
	return entity.Result{}
}

func (ds dockerService) CreateVolume(ctx context.Context, cli *client.Client, volume entity.Volume) entity.Result {
	//TODO
	return entity.Result{}
}

func (ds dockerService) RemoveVolume(ctx context.Context, cli *client.Client, name string) entity.Result {
	//TODO
	return entity.Result{}
}

func (ds dockerService) PlaceFileInContainer(ctx context.Context, cli *client.Client, containerName string, file entity.File) entity.Result {
	//TODO
	return entity.Result{}
}

func (ds dockerService) PlaceFileInVolume(ctx context.Context, cli *client.Client, volumeName string, file entity.File) entity.Result {
	//TODO
	return entity.Result{}
}

func (ds dockerService) Emulation(ctx context.Context, cli *client.Client, netem entity.Netconf) entity.Result {
	//TODO
	return entity.Result{}
}