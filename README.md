# twirp-gen

Simple Protobuf code generators specifically for Twirp services.

## JavaScript

This is an extremely simple generator for JavaScript. It makes a few assumptions to keep things simple:

* Your Twirp server has the option `twirp.WithServerJSONCamelCaseNames(true)` set
  * This client uses proto3 JSON serialization instead of the snake-case default in Twirp
* You are OK with the `fetch()` browser API
  * That is, you don't need to worry about older browsers that may not support it

### Installing

```sh
go install github.com/seanpfeifer/twirp-gen/cmd/protoc-gen-twirpjs
```

Ensure your `~/go/bin` (Linux) or `%USERPROFILE%/go/bin` (Windows) are on your PATH.

### Usage

```sh
# Typical Twirp usage to generate methods from the "examples" dir into "./pbgen/generated.js"
protoc -I ./examples --twirpjs_out=./pbgen/ ./examples/*.proto

# If you use twirp.WithServerPathPrefix(), eg `twirp.WithServerPathPrefix("/rpc")`, you can specify the
# prefix with the "pathPrefix" flag
protoc -I ./examples --twirpjs_out=pathPrefix=/rpc:./pbgen/ ./examples/*.proto
```

