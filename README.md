# Docker-Utils

Package provide some docker utils

## List remote images

`go get github.com/bborbe/docker-utils/cmd/docker-remote-repositories`

```
docker-remote-repositories \
-registry=docker.benjamin-borbe.de \
-username=bborbe \
-password=xxx \
-alsologtostderr \
-v=0
```

## List tags of remote image

`go get github.com/bborbe/docker-utils/cmd/docker-remote-tags`

```
docker-remote-tags \
-registry=docker.benjamin-borbe.de \
-username=bborbe \
-password=xxx \
-repository=bborbe/auth-http-proxy \
-alsologtostderr \
-v=0
```

## Check if remote image with tag exists

`go get github.com/bborbe/docker-utils/cmd/docker-remote-tag-exists`

```
docker-remote-tag-exists \
-registry=docker.benjamin-borbe.de \
-username=bborbe \
-password=xxx \
-repository=bborbe/auth-http-proxy \
-tag=1.0.1 \
-alsologtostderr \
-v=0
```

## Delete image tag

`go get github.com/bborbe/docker-utils/cmd/docker-remote-tag-delete`

```
docker-remote-tag-delete \
-registry=docker.benjamin-borbe.de \
-username=bborbe \
-password=xxx \
-repository=bborbe/auth-http-proxy \
-tag=1.0.1 \
-alsologtostderr \
-v=0
```
