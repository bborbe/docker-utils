# Docker-Utils

Package provide some docker utils

## List remote images

`go get github.com/bborbe/docker_utils/bin/docker_remote_images`

```
docker_remote_images \
-registry=docker.seibert-media.net \
-username=seibertmedia \
-password=xxx \
-alsologtostderr \
-v=0
```

## List tags of remote image

`go get github.com/bborbe/docker_utils/bin/docker_remote_images_tags`

```
docker_remote_images \
--username=bborbe \
--password=xxx \
--registry=docker.io \
--image=bborbe/auth
```

## Check if remote image with tag exists

`go get github.com/bborbe/docker_utils/bin/docker_remote_tag_exists`

```
docker_remote_images \
--username=bborbe \
--password=xxx \
--registry=docker.io \
--image=bborbe/auth \
--tag=latest
```

## Continuous integration

[Jenkins](https://jenkins.benjamin-borbe.de/job/Go-Docker-Utils/)

## Copyright and license

    Copyright (c) 2016, Benjamin Borbe <bborbe@rocketnews.de>
    All rights reserved.
    
    Redistribution and use in source and binary forms, with or without
    modification, are permitted provided that the following conditions are
    met:
    
       * Redistributions of source code must retain the above copyright
         notice, this list of conditions and the following disclaimer.
       * Redistributions in binary form must reproduce the above
         copyright notice, this list of conditions and the following
         disclaimer in the documentation and/or other materials provided
         with the distribution.

    THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
    "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
    LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
    A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
    OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
    SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
    LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
    DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
    THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
    (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
    OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
