#!/bin/bash
docker run --rm --entrypoint '' -v `pwd`:/go/src/github.com/nktks/spanner-er --workdir="/go/src/github.com/nktks/spanner-er" ghcr.io/nktks/spanner-er  /bin/sh -c "/bin/spanner-er $*"
