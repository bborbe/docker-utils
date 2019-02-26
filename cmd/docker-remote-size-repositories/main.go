package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/bborbe/docker-utils"
	flag "github.com/bborbe/flagenv"
	http_client_builder "github.com/bborbe/http/client_builder"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

var (
	registryPtr            = flag.String("registry", "", "Registry")
	usernamePtr            = flag.String("username", "", "Username")
	passwordPtr            = flag.String("password", "", "Password")
	passwordFilePtr        = flag.String("passwordfile", "", "Password-File")
	credentialsfromfilePtr = flag.Bool("credentialsfromfile", false, "Read Username and Password from ~/.docker/config.json")
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
	password := docker.RegistryPassword(*passwordPtr)
	if len(*passwordFilePtr) > 0 {
		password, err = docker.RegistryPasswordFromFile(*passwordFilePtr)
		if err != nil {
			return err
		}
	}
	registry := docker.Registry{
		Name:     docker.RegistryName(*registryPtr),
		Username: docker.RegistryUsername(*usernamePtr),
		Password: password,
	}
	if *credentialsfromfilePtr {
		if err := registry.ReadCredentialsFromDockerConfig(); err != nil {
			return errors.Wrap(err, "read credentials failed")
		}
	}
	glog.V(2).Infof("use registry %v", registry)
	client := http_client_builder.New().WithoutProxy().Build()
	repositories, err := docker.NewRepositories(client).List(registry)
	if err != nil {
		return err
	}
	for _, repository := range repositories {
		tags, err := docker.NewTags(client).List(registry, repository)
		if err != nil {
			glog.Warningf("list tags %s %s failed\n", registry.Name, repository.String())
			continue
		}
		var size int
		for _, tag := range tags {
			manifest, err := docker.NewTags(client).Manifest(registry, repository, tag)
			if err != nil {
				glog.Warningf("get manifest %s %s %s failed\n", registry.Name, repository.String(), tag.String())
				continue
			}
			size += manifest.Config.Size
			for _, layer := range manifest.Layers {
				size += layer.Size
			}
		}
		fmt.Fprintf(writer, "%s %d MB\n", repository.String(), size/1024/1024)
	}
	return nil
}
