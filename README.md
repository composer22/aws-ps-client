# aws-ps-client

A simple client to pull down AWS Parameter Storage key/values

Written in [Golang.](http://golang.org)

## About

Small client app that will pull down key/values from AWS Parameter Storage
and print them into bash export statements, json, or plain text.

Useful in docker container provisioning environment variables during
docker-entrypoint.sh

## Command Line Usage

```
Client for retrieving values from AWS Parameter store for a given key

Usage:
  aws-ps-client [command]

Available Commands:
  get         Get a value from the AWS Parameter Store for a key
  getpath     Get key/values from a directory in the AWS Parameter Store
  version     Version of the application

Flags:
  -k, --aws-access-key string      AWS IAM access authentication key or full path to a file containing key
  -s, --aws-access-secret string   AWS IAM access authentication secret or full path to a file containing secret
  -r, --aws-region string          AWS region ex: us-west-2
      --config string              config file (default is $HOME/.aws-ps-client.yaml)
  -f, --format string              Return 'bash', 'json', or 'text' (default "bash")

Use "aws-ps-client [command] --help" for more information about a command.

```
## Configuration

A config file is mandatory. You can place it in the same directory as the
application or in your home directory. The name, with period, is:

.aws-ps-client.yaml

This file can be empty, but must exist.

You have two places you can configure AWS

Pass parameters in your application
Put attributes in your config file

Parameters or attributes can be values or can be paths to a file that contains
the value for example:

VALUES by call:

--aws-access-key foo
--aws-access-secret bar
--aws-access-region us-west-2
--aws-access-key /path/to/key.file
--aws-access-secret /path/to/secret.file
--aws-access-region /path/to/region.file

VALUES by config file:

aws-access-key foo
aws-access-secret bar
aws-access-region us-west-2
aws-access-key /path/to/key.file
aws-access-secret /path/to/secret.file
aws-access-region /path/to/region.file

an example config file is included under /examples

## Building

This code currently requires version 1.9.2 or higher of Go.

. build.sh is the tool to create multiple executables. Edit what you need/don't need.

For package management, look to dep for instructions: <https://github.com/golang/dep>

commands:
```
dep ensure -add
dep ensure -update
dep status
```


Information on Golang installation, including pre-built binaries, is available at <http://golang.org/doc/install>.

Run `go version` to see the version of Go which you have installed.

Run `go build` inside the directory to build.

Run `go test ./...` to run the unit regression tests.

A successful build run produces no messages and creates an executable called `aws-ps-client` in this
directory.

Run `go help` for more guidance, and visit <http://golang.org/> for tutorials, presentations, references and more.

## License

(The MIT License)

Copyright (c) 2017 Pyxxel Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to
deal in the Software without restriction, including without limitation the
rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
sell copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
IN THE SOFTWARE.
