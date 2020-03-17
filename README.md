# spanner-er [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license] [![CI Status](https://github.com/nakatamixi/spanner-er/workflows/CI/badge.svg)](https://github.com/nakatamixi/spanner-er/actions)


[license]: https://github.com/nakatamixi/spanner-er/blob/master/LICENSE

`spanner-er` is a command-line tool to generate ER diagram from DDL schama file.

# Install

## Install to host

`spanner-er` depends on graphviz.
Install graphviz on your host platform.
```
brew install graphviz
```
or
```
apk add --no-cache graphviz ttf-freefont
```
or
```
apt-get update && apt-get install graphviz
```
Install spanner-er
```
go get -u github.com/nakatamixi/spanner-er
```

## Use docker
you can use docker image
```
docker run --rm -v `pwd`:/go/src/github.com/nakatamixi/spanner-er --workdir="/go/src/github.com/nakatamixi/spanner-er" nakatamixi/spanner-er -h
```
or
```
git clone git@github.com:nakatamixi/spanner-er.git
cd ./spanner-er
./scripts/spanner-er-docker.sh -h
```
In this case, you should use relative path for `-s`, `-o` option.

# Usage
```
spanner-er -h
Usage:
  -T string
    	output file type. default is png(pass to dot option -T) (default "png")
  -h	print help
  -o string
    	output file name.default is spanner_er.<type>(pass to dot option -o)
  -s string
    	spanner schema file
```

# Sample image
![image](https://user-images.githubusercontent.com/7553415/76856135-f949ca80-6895-11ea-88c2-bee9218ee2c3.png)

