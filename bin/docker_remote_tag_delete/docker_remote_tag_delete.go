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
	parameterRepository   = "repository"
	parameterTag          = "tag"
)

var (
	registryPtr     = flag.String(parameterRegistry, "", "Registry")
	usernamePtr     = flag.String(parameterUsername, "", "Username")
	passwordPtr     = flag.String(parameterPassword, "", "Password")
	passwordFilePtr = flag.String(parameterPasswordFile, "", "Password-File")
	repositoryPtr   = flag.String(parameterRepository, "", "Repository")
	tagPtr          = flag.String(parameterTag, "", "Tag")
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
	repositoryName := model.RepositoryName(*repositoryPtr)
	tag := model.Tag(*tagPtr)

	glog.V(2).Infof("use registry %v, repo %v and tag %v", registry, repositoryName, tag)
	if err = registry.Validate(); err != nil {
		return err
	}
	factory := docker_utils_factory.New()
	if err = factory.Tags().Delete(registry, repositoryName, tag); err != nil {
		glog.V(2).Infof("delete tag exists failed: %v", err)
		return err
	}
	fmt.Fprintf(writer, "tag deleted\n")
	return nil
}
