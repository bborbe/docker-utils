// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	flag "github.com/bborbe/flagenv"
	"github.com/golang/glog"

	"github.com/bborbe/docker-utils"
)

var (
	registryPtr     = flag.String("registry", "", "Registry")
	usernamePtr     = flag.String("username", "", "Username")
	passwordPtr     = flag.String("password", "", "Password")
	passwordFilePtr = flag.String("passwordfile", "", "Password-File")
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
	repositories := make(chan docker.RepositoryName, runtime.NumCPU())
	go func() {
		defer close(repositories)
		if err := client.ListRepositories(ctx, repositories); err != nil {
			glog.Warningf("read repos failed: %v", err)
		}
	}()
	for repository := range repositories {
		tags := make(chan docker.TagName, runtime.NumCPU())
		go func() {
			defer close(tags)
			if err := client.ListTags(ctx, repository, tags); err != nil {
				glog.Warningf("list tags %s failed\n", repository.String())
			}
		}()
		var size int
		for tag := range tags {
			manifest, err := client.Manifest(ctx, repository, tag)
			if err != nil {
				glog.Warningf("get manifest %s %s failed\n", repository.String(), tag.String())
				continue
			}
			size += manifest.Config.Size
			for _, layer := range manifest.Layers {
				size += layer.Size
			}
		}
		fmt.Printf("%s %d MB\n", repository.String(), size/1024/1024)
	}
	return nil
}
