// Copyright 2018 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package etcdfs

import (
	"os"
	"strings"

	"golang.org/x/tools/godoc/vfs"
)

// New returns a new FileSystem from the provided etcd.
// Wtcd keys should be forward slash-separated pathnames
// and not contain a leading slash.
func New(c *EtcdClient) vfs.FileSystem {
	return &etcdFS{c: c}
}

type etcdFS struct {
	c *EtcdClient
}

func (fs *etcdFS) String() string {
	return "etcdfs"
}

func (fs *etcdFS) RootType(path string) vfs.RootType {
	return ""
}

func (fs *etcdFS) Open(name string) (vfs.ReadSeekCloser, error) {
	b, ok := fs.c.Get((name))
	if !ok {
		return nil, os.ErrNotExist
	}

	return nopCloser{strings.NewReader(b)}, nil
}

func (fs *etcdFS) Lstat(name string) (os.FileInfo, error) {
	b, ok := fs.c.Get((name))
	if ok {
		return fileInfo(name, b), nil
	}
	ents, _ := fs.ReadDir(name)
	if len(ents) > 0 {
		return dirInfo(name), nil
	}
	return nil, os.ErrNotExist
}

func (fs *etcdFS) Stat(path string) (os.FileInfo, error) {
	return fs.Lstat(path)
}

func (fs *etcdFS) ReadDir(path string) ([]os.FileInfo, error) {
	m, err := fs.c.GetValuesByPrefix("")
	if err != nil {
		return nil, err
	}

	return mapFS(m).ReadDir(path)
}
