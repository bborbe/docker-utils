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


var (
	registryPtr     = flag.String(model.ParameterRegistry, "", "Registry")
	usernamePtr     = flag.String(model.ParameterUsername, "", "Username")
	passwordPtr     = flag.String(model.ParameterPassword, "", "Password")
	passwordFilePtr = flag.String(model.ParameterPasswordFile, "", "Password-File")
	repositoryPtr   = flag.String(model.ParameterRepository, "", "Repository")
	tagPtr          = flag.String(model.ParameterTag, "", "Tag")
	credentialsfromfilePtr = flag.Bool(model.ParameterCredentialsFromDockerConfig, false, "Read Username and Password from ~/.docker/config.json")
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
	if *credentialsfromfilePtr {
		if err := registry.ReadCredentialsFromDockerConfig(); err != nil {
			return fmt.Errorf("read credentials failed: %v", err)
		}
	}
	repositoryName := model.RepositoryName(*repositoryPtr)
	tag := model.Tag(*tagPtr)

	glog.V(2).Infof("use registry %v, repo %v and tag %v", registry, repositoryName, tag)
	if err = registry.Validate(); err != nil {
		return fmt.Errorf("validate registry failed: %v", err)
	}
	factory := docker_utils_factory.New()
	if err = factory.Tags().Delete(registry, repositoryName, tag); err != nil {
		glog.V(2).Infof("delete tag exists failed: %v", err)
		return err
	}
	fmt.Fprintf(writer, "tag deleted\n")
	return nil
}