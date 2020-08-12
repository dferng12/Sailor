package dockerinterface

import (
	"context"
	"strconv"
	"strings"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
)

type Volume struct {
}

type Network struct {
}

type ContainerConfig struct {
	Name       string             `yaml:"container_name,omitempty"`
	Image      string             `yaml:"image"`
	Mounts     []string           `yaml:"volumes,omitempty"`
	Ports      []string           `yaml:"ports,omitempty"`
	Network    map[string]Network `yaml:"networks,omitempty"`
	Env        []string           `yaml:"environment,omitempty"`
	Command    string             `yaml:"command,omitempty"`
	Entrypoint string             `yaml:"entrypoint,omitempty"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// GetContainers returns a list of the current running containers
func GetContainers() []types.Container {
	cli, err := docker.NewEnvClient()
	checkErr(err)

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	checkErr(err)

	return containers

}

func GetContainerConfig(container types.Container) ContainerConfig {
	cli, err := docker.NewEnvClient()
	checkErr(err)

	containerConfig := ContainerConfig{}
	containerConfig.Name = container.Names[0][1:len(container.Names[0])]
	containerConfig.Image = container.Image
	containerConfig.Command = container.Command

	for _, mount := range container.Mounts {
		containerConfig.Mounts = append(containerConfig.Mounts, mount.Source+":"+mount.Destination+":"+mount.Mode)
	}

	for _, port := range container.Ports {
		containerConfig.Ports = append(containerConfig.Ports, strconv.Itoa(int(port.PublicPort))+":"+strconv.Itoa(int(port.PrivatePort))+"/"+port.Type)
	}

	containerDockerConfig, _ := cli.ContainerInspect(context.Background(), container.ID)
	containerConfig.Env = containerDockerConfig.Config.Env
	containerConfig.Entrypoint = strings.Join(containerDockerConfig.Config.Entrypoint, " ")

	containerConfig.Network = make(map[string]Network)
	if containerDockerConfig.NetworkSettings.NetworkSettingsBase.Bridge == "" {
		containerConfig.Network["docker0"] = Network{}
	}

	return containerConfig
}

func GetRunningContainersConfig() []ContainerConfig {
	containers := GetContainers()
	config := []ContainerConfig{}
	for _, container := range containers {
		config = append(config, GetContainerConfig(container))
	}
	return config
}
