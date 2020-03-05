#!/bin/bash
docker run --rm --entrypoint '' -v `pwd`:/go/src/github.com/nakatamixi/spanner-er --workdir="/go/src/github.com/nakatamixi/spanner-er" nakatamixi/spanner-er  /bin/sh -c "/bin/spanner-er $*"
