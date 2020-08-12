# Sailor
A go based tool to generate docker-compose files from existing Docker containers

## Description

This tool loads all the running containers in the system and generates a docker compose config file including all the configurations active in each of the containers.

Then, you just have to copy the config to a new Docker system and launch it to have everything up and running again.

## Usage

Simply run the program using the Go program

`go run sailor.go`

or compile it and use it without needing Go

`go build`

`./sailor`

## Future work

- [ ] Flags
- [ ] Other types mounts
- [ ] Attachable stdin
- [ ] Complex networking