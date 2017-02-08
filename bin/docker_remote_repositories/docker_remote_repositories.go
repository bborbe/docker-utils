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
)

var (
	registryPtr = flag.String(PARAMETER_REGISTRY, "", "Registry")
	usernamePtr = flag.String(PARAMETER_USER, "", "Username")
	passwordPtr = flag.String(PARAMETER_PASS, "", "Password")
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
	glog.V(2).Infof("use registry %v", registry)
	if err := registry.Validate(); err != nil {
		return err
	}
	factory := docker_utils_factory.New()
	repositories, err := factory.Repositories().List(registry)
	if err != nil {
		return err
	}
	for _, repository := range repositories {
		fmt.Fprintf(writer, "%s\n", repository.String())
	}
	return nil
}
