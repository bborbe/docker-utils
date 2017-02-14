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
	parameterRegistry     = "registry"
	parameterUsername     = "username"
	parameterPassword     = "password"
	parameterPasswordFile = "passwordfile"
)

var (
	registryPtr     = flag.String(parameterRegistry, "", "Registry")
	usernamePtr     = flag.String(parameterUsername, "", "Username")
	passwordPtr     = flag.String(parameterPassword, "", "Password")
	passwordFilePtr = flag.String(parameterPasswordFile, "", "Password-File")
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
	var err error
	password := model.RegistryPassword(*passwordPtr)
	if len(*passwordFilePtr) > 0 {
		password, err = model.RegistryPasswordFromFile(*passwordFilePtr)
		if err != nil {
			return err
		}
	}
	registry := model.Registry{
		Name:     model.RegistryName(*registryPtr),
		Username: model.RegistryUsername(*usernamePtr),
		Password: password,
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
