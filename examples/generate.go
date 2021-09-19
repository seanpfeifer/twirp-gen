package examples

//go:generate protoc -I . --twirpjs_out=pathPrefix=/rpc:../pbgen/ ./*.proto
//go:generate protoc -I . --twirpcs_out=pathPrefix=/rpc:../pbgen/ ./*.proto

// NOTE: If you're doing both JS + C#, you can simplify to one line:
//   protoc -I . --twirpjs_out=pathPrefix=/rpc:../pbgen/ --twirpjs_out=pathPrefix=/rpc:../pbgen/ ./*.proto
