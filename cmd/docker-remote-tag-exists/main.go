package main

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/bborbe/docker-utils"
	flag "github.com/bborbe/flagenv"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

var (
	registryPtr     = flag.String("registry", "", "Registry")
	usernamePtr     = flag.String("username", "", "Username")
	passwordPtr     = flag.String("password", "", "Password")
	passwordFilePtr = flag.String("passwordfile", "", "Password-File")
	repositoryPtr   = flag.String("repository", "", "Repository")
	tagPtr          = flag.String("tag", "", "Tag")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := do(context.Background()); err != nil {
		glog.Exitf("%+v", err)
	}
}

func do(ctx context.Context) error {
	registry := &docker.Registry{
		Url:      *registryPtr,
		Username: *usernamePtr,
		Password: *passwordPtr,
	}
	if len(*passwordFilePtr) > 0 {
		if err := registry.RegistryPasswordFromFile(*passwordFilePtr); err != nil {
			return err
		}
	}
	glog.V(2).Infof("use registry %v, repo %v and tag %v", registry, docker.RepositoryName(*repositoryPtr), docker.TagName(*tagPtr))
	client := docker.NewV2Client(docker.NewHttpClient(http.DefaultClient), *registry)
	exists, err := client.ExistsTag(ctx, docker.RepositoryName(*repositoryPtr), docker.TagName(*tagPtr))
	if err != nil {
		return errors.Wrap(err, "check tag exists failed")
	}
	fmt.Printf("%v\n", exists)
	return nil
}
