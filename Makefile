# Copyright 2018 <chaishushan{AT}gmail.com>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# cat $HOME/.netrc
# machine api.github.com
# login username
# password SECRET
# protocol https

default:

dev:
	https_proxy="socks5://127.0.0.1:2080" vgo fmt ./...
	https_proxy="socks5://127.0.0.1:2080" vgo vet ./...
	https_proxy="socks5://127.0.0.1:2080" vgo test ./...

test:
	go fmt ./...
	go vet ./...
	go test ./...

clean:
