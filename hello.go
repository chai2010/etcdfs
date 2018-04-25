// Copyright 2018 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

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
