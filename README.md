# etcdfs: virtual file system on etcd

[![Build Status](https://travis-ci.org/chai2010/etcdfs.svg)](https://travis-ci.org/chai2010/etcdfs)
[![Go Report Card](https://goreportcard.com/badge/github.com/chai2010/etcdfs)](https://goreportcard.com/report/github.com/chai2010/etcdfs)
[![GoDoc](https://godoc.org/github.com/chai2010/etcdfs?status.svg)](https://godoc.org/github.com/chai2010/etcdfs)

## Install

1. `go get github.com/chai2010/etcdfs`

## Example

Start etcd in local:

	etcd

Put some data to edcd:

	ETCDCTL_API=3 etcdctl put /abc/readme.md abc/aaa-value
	ETCDCTL_API=3 etcdctl put /abc/hello.go  "package main; func main(){}"
	ETCDCTL_API=3 etcdctl get --prefix ""

Create go progrom:

```go
package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/httpfs"

	"github.com/chai2010/etcdfs"
)

var (
	flagEtcdHost = flag.String("ectd-host", "localhost:2379", "set etcd nodes")
)

func main() {
	flag.Parse()

	etcdClient, err := etcdfs.NewEtcdClient(strings.Split(*flagEtcdHost, ","), time.Second/10)
	if err != nil {
		log.Fatal(err)
	}

	ns := vfs.NameSpace{}
	ns.Bind("/", etcdfs.New(etcdClient), "/", vfs.BindReplace)

	http.Handle("/", http.FileServer(httpfs.New(ns)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Run the program ([hello.go](hello.go)):

	go run hello.go

Then open http://127.0.0.1:8080/ in browser.

## BUGS

Report bugs to <chaishushan@gmail.com>.

Thanks!
