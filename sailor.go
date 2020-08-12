package main

import (
	"github.com/dferng12/sailor/dockerinterface"
	"github.com/dferng12/sailor/yamlgenerator"
)

func main() {
	config := dockerinterface.GetRunningContainersConfig()
	yamlgenerator.CreateComposeFile(config)
}
