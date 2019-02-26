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
	repositoryPtr          = flag.String("repository", "", "Repository")
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
			return errors.Wrap(err, "read registry password from file")
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
	repositoryName := docker.RepositoryName(*repositoryPtr)
	glog.V(2).Infof("use registry %v and repo %v", registry, repositoryName)
	client := http_client_builder.New().WithoutProxy().Build()
	tags, err := docker.NewTags(client).List(registry, repositoryName)
	if err != nil {
		return errors.Wrap(err, "list tags failed")
	}
	for _, tag := range tags {
		fmt.Fprintf(writer, "%s\n", tag.String())
	}
	return nil
}
