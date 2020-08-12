package yamlgenerator

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/dferng12/sailor/dockerinterface"
	"gopkg.in/yaml.v2"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Compose struct {
	Version  string                                     `yaml:"version"`
	Services map[string]dockerinterface.ContainerConfig `yaml:"services"`
	Networks map[string]dockerinterface.Network         `yaml:"networks,omitempty"`
	Volumes  map[string]dockerinterface.Volume          `yaml:"volumes",omitempty`
}

func CreateComposeFile(containerConfigs []dockerinterface.ContainerConfig) {
	networks := make(map[string]dockerinterface.Network)
	volumes := make(map[string]dockerinterface.Volume)

	composeFile := Compose{Version: "3"}
	composeFile.Services = make(map[string]dockerinterface.ContainerConfig)

	for _, containerConfig := range containerConfigs {
		for name, network := range containerConfig.Network {
			if _, ok := networks[name]; !ok {
				networks[name] = network
			}
		}

		for _, volume := range containerConfig.Mounts {
			if volume[0] != '/' {
				volumes[strings.Split(volume, ":")[0]] = dockerinterface.Volume{}
			}
		}

		composeFile.Services[containerConfig.Name] = containerConfig
	}

	composeFile.Networks = networks
	composeFile.Volumes = volumes

	d, err := yaml.Marshal(&composeFile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	f, err := os.Create("docker-compose.yml")
	defer f.Close()
	checkErr(err)

	_, err = f.Write(d)
	checkErr(err)
}

func ProcessComposeFile() Compose {
	dat, err := ioutil.ReadFile("docker-compose.yml")

	t := Compose{}

	err = yaml.Unmarshal(dat, &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return t
}
