package main

import (
	"fmt"
	docker_utils_factory "github.com/bborbe/docker_utils/factory"
	"github.com/bborbe/docker_utils/model"
	flag "github.com/bborbe/flagenv"
	"github.com/golang/glog"
	"io"
	"os"
	"runtime"
)

const (
	PARAMETER_REGISTRY = "registry"
	PARAMETER_USER     = "username"
	PARAMETER_PASS     = "password"
	PARAMETER_REPO     = "repository"
)

var (
	registryPtr   = flag.String(PARAMETER_REGISTRY, "", "Registry")
	usernamePtr   = flag.String(PARAMETER_USER, "", "Username")
	passwordPtr   = flag.String(PARAMETER_PASS, "", "Password")
	repositoryPtr = flag.String(PARAMETER_REPO, "", "Repository")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	writer := os.Stdout
	if err := do(writer); err != nil {
		glog.Exit(err)
	}
}

func do(writer io.Writer) error {
	registry := model.Registry{
		Name:     model.RegistryName(*registryPtr),
		Username: model.RegistryUsername(*usernamePtr),
		Password: model.RegistryPassword(*passwordPtr),
	}
	repositoryName := model.RepositoryName(*repositoryPtr)
	glog.V(2).Infof("use registry %v and repo %v", registry, repositoryName)
	if err := registry.Validate(); err != nil {
		return err
	}
	factory := docker_utils_factory.New()
	tags, err := factory.Tags().List(registry, repositoryName)
	if err != nil {
		return err
	}
	for _, tag := range tags {
		fmt.Fprintf(writer, "%s\n", tag.String())
	}
	return nil
}
