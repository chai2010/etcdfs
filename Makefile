# Copyright 2018 <chaishushan{AT}gmail.com>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

default:

test:
	vgo fmt ./...
	vgo vet ./...
	vgo test ./...

clean:
