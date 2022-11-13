# twirp-gen

Simple Protobuf code generators specifically for Twirp services.

This exists because other code generators either required external dependencies or were too complex for my needs.

## Installing

```sh
# For TypeScript generation
go install github.com/seanpfeifer/twirp-gen/cmd/protoc-gen-twirpts@latest

# For C# generation
go install github.com/seanpfeifer/twirp-gen/cmd/protoc-gen-twirpcs@latest

# For JavaScript generation
go install github.com/seanpfeifer/twirp-gen/cmd/protoc-gen-twirpjs@latest
```

Ensure your `~/go/bin` (Linux) or `%USERPROFILE%/go/bin` (Windows) are on your PATH.

## Usage

```sh
# Typical Twirp usage to generate methods from the "examples" dir into "./examples_gen/generated.js"
# For TypeScript
protoc -I ./examples --twirpts_out=./examples_gen/ ./examples/*.proto
# For JavaScript
protoc -I ./examples --twirpjs_out=./examples_gen/ ./examples/*.proto

# If you use twirp.WithServerPathPrefix(), eg `twirp.WithServerPathPrefix("/rpc")`, you can specify the
# prefix with the "pathPrefix" flag
protoc -I ./examples --twirpjs_out=pathPrefix=/rpc:./examples_gen/ ./examples/*.proto

# C# is the same as above two examples, but uses `--twirpcs_out`
protoc -I ./examples --twirpcs_out=./examples_gen/ ./examples/*.proto
protoc -I ./examples --twirpcs_out=pathPrefix=/rpc:./examples_gen/ ./examples/*.proto

# For TypeScript, you can also add `protoCase=true` to generate field names in the same case as the proto file
protoc -I ./examples --twirpts_out=pathPrefix=/rpc,protoCase=true:./examples_gen/ ./examples/*.proto
```

## Generator Options

Option | Language | Description
-------|----------|------------
`pathPrefix` | TypeScript, JavaScript, C# | The path prefix to use for the Twirp service. This is useful if you use `twirp.WithServerPathPrefix()`.
`protoCase` | TypeScript | If true, field names will be generated in the same case as the proto file. If false, field names will be generated in camelCase. This is useful if you DO NOT use `twirp.WithServerJSONCamelCaseNames()`.

## TypeScript

This is an extremely simple generator for TypeScript, which communicates with JSON only. Assumptions:

* You are OK with the `fetch()` browser API
  * That is, you don't need to worry about older browsers that may not support it
* Uses camelCase fields by default (ie, you have the option `twirp.WithServerJSONCamelCaseNames(true)` set)
  * If you don't want fields changed to camelCase, use `protoCase=true` in your `twirpts_out` options, eg if you want to use the same case as the proto file  and to also generate with the path prefix `/rpc`:
    * `--twirpts_out=pathPrefix=/rpc,protoCase=true:./`

## C#

This is an extremely simple generator for C#. This client makes requests using binary serialized Protobuf.

On error, the generated functions will throw a `GeneratedAPI.Exception`.

## JavaScript

Another extremely simple generator, but for JavaScript, which communicates with JSON only. It makes a few assumptions to keep things simple:

* Your Twirp server has the option `twirp.WithServerJSONCamelCaseNames(true)` set
  * This client uses proto3 JSON serialization instead of the snake-case default in Twirp
* You are OK with the `fetch()` browser API
  * That is, you don't need to worry about older browsers that may not support it

## Unity Example

Note that this example does not deal with errors. It is just a simple example of how to call the generated functions.

```cs
using System.Net.Http;
using UnityEngine;
using Google.Protobuf; // For ByteString in this example

public class APICaller : MonoBehaviour {
  public const string SERVER_ADDRESS = "http://localhost:8080/";

  static readonly HttpClient client = new HttpClient();
  static APICaller() {
    client.BaseAddress = new System.Uri(SERVER_ADDRESS);
    client.DefaultRequestHeaders.Accept.Add(GeneratedAPI.CONTENT_TYPE_PROTOBUF);
  }

  // Start is called before the first frame update
  async void Start() {
    var checkoutResult = await GeneratedAPI.CreateCheckoutSession(client, new Shop.CheckoutRequest { ItemId = "item_cool_stuff" });
    Debug.Log(checkoutResult);

    var logoutResult = await GeneratedAPI.Logout(client, new Account.LogoutRequest { AccountId = 42, Token = ByteString.CopyFrom(new byte[]{1,2,3,4,5}) });
    Debug.Log(logoutResult);
  }
}
```
