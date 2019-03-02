package main

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/bborbe/docker-utils"
	flag "github.com/bborbe/flagenv"
	"github.com/golang/glog"
)

var (
	registryPtr     = flag.String("registry", "", "Registry")
	usernamePtr     = flag.String("username", "", "Username")
	passwordPtr     = flag.String("password", "", "Password")
	passwordFilePtr = flag.String("passwordfile", "", "Password-File")
	repositoryPtr   = flag.String("repository", "", "Repository")
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
	glog.V(2).Infof("use registry %v and repo %v", registry, docker.RepositoryName(*repositoryPtr))
	client := docker.NewV2Client(docker.NewHttpClient(http.DefaultClient), *registry)
	tags := make(chan docker.TagName, runtime.NumCPU())
	go func() {
		defer close(tags)
		if err := client.ListTags(ctx, docker.RepositoryName(*repositoryPtr), tags); err != nil {
			glog.Warningf("list tags failed: %v", err)
		}
	}()
	for tag := range tags {
		fmt.Printf("%s\n", tag.String())
	}
	return nil
}
