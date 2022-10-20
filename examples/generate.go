package examples

//go:generate protoc -I . --twirpjs_out=pathPrefix=/rpc:../examples_gen/js ./*.proto
//go:generate protoc -I . --twirpts_out=pathPrefix=/rpc:../examples_gen/ts ./*.proto
//go:generate protoc -I . --twirpcs_out=pathPrefix=/rpc:../examples_gen/cs --csharp_out=../examples_gen/cs ./*.proto
