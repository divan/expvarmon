package main

import (
	docker "github.com/fsouza/go-dockerclient"
	"os"
	"path/filepath"
	"strings"
)

type Container struct {
	Name string
}

func getContainers() ([]string, error) {
	cli, err := connect()
	if err != nil {
		return nil, err
	}

	containers, err := cli.ListContainers(docker.ListContainersOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	var names []string
	for _, cont := range containers {
		names = append(names, cont.Names[0])
	}

	return names, nil
}

// connect establishes connection to the Docker daemon.
func connect() (*docker.Client, error) {
	var (
		client *docker.Client
		err    error
	)
	// TODO: add boot2docker shellinit support
	endpoint := os.Getenv("DOCKER_HOST")
	if endpoint == "" {
		endpoint = "unix:///var/run/docker.sock"
	}
	cert_path := os.Getenv("DOCKER_CERT_PATH")
	if cert_path != "" {
		client, err = docker.NewTLSClient(
			endpoint,
			filepath.Join(cert_path, "cert.pem"),
			filepath.Join(cert_path, "key.pem"),
			filepath.Join(cert_path, "ca.pem"),
		)
	} else {
		client, err = docker.NewClient(endpoint)
	}

	return client, err
}

// cleanContainerName clears leading '/' symbol from container name.
func cleanContainerName(name string) string {
	return strings.TrimLeft(name, "/")
}
