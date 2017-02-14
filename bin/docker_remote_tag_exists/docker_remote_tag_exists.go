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
	PARAMETER_TAG      = "tag"
)

var (
	registryPtr   = flag.String(PARAMETER_REGISTRY, "", "Registry")
	usernamePtr   = flag.String(PARAMETER_USER, "", "Username")
	passwordPtr   = flag.String(PARAMETER_PASS, "", "Password")
	repositoryPtr = flag.String(PARAMETER_REPO, "", "Repository")
	tagPtr        = flag.String(PARAMETER_TAG, "", "Tag")
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
	tag := model.Tag(*tagPtr)

	glog.V(2).Infof("use registry %v, repo %v and tag %v", registry, repositoryName, tag)
	if err := registry.Validate(); err != nil {
		return err
	}
	factory := docker_utils_factory.New()
	exists, err := factory.Tags().Exists(registry, repositoryName, tag)
	if err != nil {
		glog.V(2).Infof("check tag exists failed: %v", err)
		return err
	}
	fmt.Fprintf(writer, "%v\n", exists)
	return nil
}
