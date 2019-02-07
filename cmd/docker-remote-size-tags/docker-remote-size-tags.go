package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/pkg/errors"

	docker_utils_factory "github.com/bborbe/docker-utils/factory"
	"github.com/bborbe/docker-utils/model"
	flag "github.com/bborbe/flagenv"
	"github.com/golang/glog"
)

var (
	registryPtr            = flag.String(model.ParameterRegistry, "", "Registry")
	usernamePtr            = flag.String(model.ParameterUsername, "", "Username")
	passwordPtr            = flag.String(model.ParameterPassword, "", "Password")
	passwordFilePtr        = flag.String(model.ParameterPasswordFile, "", "Password-File")
	repositoryPtr          = flag.String(model.ParameterRepository, "", "Repository")
	credentialsfromfilePtr = flag.Bool(model.ParameterCredentialsFromDockerConfig, false, "Read Username and Password from ~/.docker/config.json")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	writer := os.Stdout
	if err := do(writer); err != nil {
		glog.Exitf("%+v", err)
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
			return errors.Wrap(err, "read credentials failed")
		}
	}
	glog.V(2).Infof("use registry %v", registry)
	if err := registry.Validate(); err != nil {
		return errors.Wrap(err, "validate registry failed")
	}
	factory := docker_utils_factory.New()
	repository := model.RepositoryName(*repositoryPtr)

	tags, err := factory.Tags().List(registry, repository)
	if err != nil {
		return errors.Wrap(err, "list tags failed")
	}
	for _, tag := range tags {
		manifest, err := factory.Tags().Manifest(registry, repository, tag)
		if err != nil {
			glog.Warningf("get manifest %s %s %s failed\n", registry.Name, repository.String(), tag.String())
			continue
		}
		var size int
		size += manifest.Config.Size
		for _, layer := range manifest.Layers {
			size += layer.Size
		}
		fmt.Fprintf(writer, "%s:%s %d MB\n", repository.String(), tag.String(), size/1024/1024)
	}
	return nil
}
