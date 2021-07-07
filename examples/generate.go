package examples

//go:generate protoc -I . --twirpjs_out=pathPrefix=/rpc:../pbgen/ ./*.proto
