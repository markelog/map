# Map [![Build Status](https://travis-ci.org/markelog/map.svg?branch=master)](https://travis-ci.org/markelog/map) [![GoDoc](https://godoc.org/github.com/markelog/map?status.svg)](https://godoc.org/github.com/markelog/map) [![Go Report Card](https://goreportcard.com/badge/github.com/markelog/map)](https://goreportcard.com/report/github.com/markelog/map)

> Generates map site

Simple site map generator, supports couple reporters, depth levels and etc

## Install
```sh
go get github.com/markelog/map
```

## Usage
```sh
# Create map and output it to the terminal1
$ map http://example.com

# Create map and output map in yaml form
$ map http://example.com --reporter=yaml

# Specify maximum amount of links to check
$ map http://example.com -r yaml --max=50

# Pipe it
$ map http://example.com -r yaml -m 50 > example.com.yaml

# Or use "out" flag to pipe (so you can see the spinner comparing with previous command :)
$ map http://example.com -r yaml -m 50 --out=./example.com.yaml
```
