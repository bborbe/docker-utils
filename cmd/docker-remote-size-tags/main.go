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
	glog.V(2).Infof("use registry %v", registry)
	client := docker.NewV2Client(docker.NewHttpClient(http.DefaultClient), *registry)
	tags := make(chan docker.TagName, runtime.NumCPU())
	go func() {
		defer close(tags)
		if err := client.ListTags(ctx, docker.RepositoryName(*repositoryPtr), tags); err != nil {
			glog.Warningf("list tags failed: %v", err)
		}
	}()
	for tag := range tags {
		manifest, err := client.Manifest(ctx, docker.RepositoryName(*repositoryPtr), tag)
		if err != nil {
			glog.Warningf("get manifest %s %s failed\n", docker.RepositoryName(*repositoryPtr).String(), tag.String())
			continue
		}
		var size int
		size += manifest.Config.Size
		for _, layer := range manifest.Layers {
			size += layer.Size
		}
		fmt.Printf("%s:%s %d MB\n", docker.RepositoryName(*repositoryPtr).String(), tag.String(), size/1024/1024)
	}
	return nil
}
