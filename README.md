![License](https://img.shields.io/github/license/bufbuild/protoc-gen-validate?color=blue)

- [English](README.md)
- [简体中文](README-zh.md)

# protoc-gen-validate

based on proto files' annotation, generating a validate function for each message.

[Article](https://yflee.in/tech/protoc-gen-validate.html) describes the past and present of protoc-gen-validate, and I welcome everyone to add my WeChat `lyf987667482` to discuss together.

```protobuf
syntax = "proto3";

package foo.v1;

service Demo {
    rpc Numerics(NumericsReq) returns (Empty);
    rpc Strings(StringsReq) returns (Empty);
    rpc Repeated(RepeatedReq) returns (Empty);
}

message Empty {}

message RepeatedReq {
    // @min_items:1
    // @max_items:2
    repeated int32 a = 1;
    // @unique:true
    repeated int64 b = 2;
    // @unique:true
    repeated string c = 3;
    repeated NumericsReq d = 4;
    // @eq:1.23
    repeated float e = 5;
}

message NumericsReq {
    // @eq:1.23
    float a = 1;
    // @lt:20
    // @gt:10
    int32 b = 2;
    // @lte:20
    // @gte:10
    uint64 c = 3;
    // @in:[1,2,3]
    fixed32 d = 4;
    // @not_in:[1,2,3]
    float e = 5;
    // @range:(1,5)
    float f = 6;
    // @range:[1,5]
    float g = 7;
}

message StringsReq {
    // @contains:"bar"
    string a = 1;
    // @not_contains:"bar"
    string b = 2;
    // @eq:"bar"
    string c = 3;
    // @in:["foo", "bar", "baz"]
    string d = 4;
    // @not_in:["foo", "bar", "baz"]
    string e = 5;
    // @len:5
    string f = 6;
    // @min_len:5
    // @max_len:10
    string g = 7;
    // @pattern:"(?i)^[0-9a-f]+$"
    string h = 8;
    // @prefix:"foo"
    string i = 9;
    // @suffix:"bar"
    string j = 10;
    // @type:url
    string k = 11;
    // @type:phone
    string l = 12;
    // @type:email
    string m = 13;
    // @type:ip
    string n = 14;
}

message Required {
    // @required:true
    Foo a = 1;
}

message Foo {
    string a = 1;
    string b = 2;
}
```

Created validation method for Numerics:
```go
func (m *NumericsReq) Validate() error {
  if m == nil {
          return nil
  }  
  if m.GetA() != 1.23 {
          return NumericsReqValidationError{
                  field:  "A",
                  reason: "value must equal 1.23",
          }
  }  
  if m.GetB() >= 20 {
          return NumericsReqValidationError{
                  field:  "B",
                  reason: "value must less than 20",
          }
  }  
  if m.GetB() <= 10 {
          return NumericsReqValidationError{
                  field:  "B",
                  reason: "value must greater than 10",
          }
  }  
  if m.GetC() > 20 {
          return NumericsReqValidationError{
                  field:  "C",
                  reason: "value must less than or equal to 20",
          }
  }  
  if m.GetC() < 10 {
          return NumericsReqValidationError{
                  field:  "C",
                  reason: "value must greater than or equal to 10",
          }
  }  
  var NumericsReq_D_In = map[uint32]struct{}{  
          1: {},  
          2: {},  
          3: {},
  }  
  if _, ok := NumericsReq_D_In[m.GetD()]; !ok {
          return NumericsReqValidationError{
                  field:  "D",
                  reason: "value must be in list [1,2,3]",
          }
  }  
  var NumericsReq_E_NotIn = map[float32]struct{}{  
          1: {},  
          2: {},  
          3: {},
  }  
  if _, ok := NumericsReq_E_NotIn[m.GetE()]; ok {
          return NumericsReqValidationError{
                  field:  "E",
                  reason: "value must be not in list [1,2,3]",
          }
  }  
  if m.GetF() <= 1 || m.GetF() >= 5 {
          return NumericsReqValidationError{
                  field:  "F",
                  reason: "value must in range (1,5)",
          }
  }  
  if m.GetG() < 1 || m.GetG() > 5 {
          return NumericsReqValidationError{
                  field:  "G",
                  reason: "value must in range [1,5]",
          }
  }  
  return nil
}
```

# Installation
```bash
go install github.com/giftDad/protoc-gen-validate
```

# Example
```bash
protoc  examples/foo.proto  --validate_out=Mexamples/foo.proto=examples/\;examples:. \
--go_out=Mexamples/foo.proto=examples/\;examples:.
go test -v ./examples/
```

# Application
Suitable for RPC frameworks based on protobuf, such as [gi-micro](https://github.com/go-micro/generator),[twirp](https://github.com/twitchtv/twirp) and more. There are two ways to use it, using gi-micro as an example:

1. The first approach is to modify the protoc-gen-xxx tool, such as gi-micro's [generator](https://github.com/go-micro/generator)
Add the following lines at `cmd/protoc-gen-micro/plugin/micro/micro.go:480`:

```go
g.P("func (h *", unexport(servName), "Handler) ", methName, "(ctx ", contextPkg, ".Context, in *", inType, ", out *",outType, ") error {")
// Add the following three lines
g.P("if err := in.Validate();err != nil {")
g.P("return err")
g.P("}")

g.P("return h.", serveType, ".", methName, "(ctx, in, out)")
g.P("}")
g.P()
```
This will transform the generated method into:
```go
func (h *greeterHandler) Hello(ctx context.Context, in *Request, out *Response) error {
    if err := in.Validate();err != nil {
        return err
    }
    return h.GreeterHandler.Hello(ctx, in, out)
}
```
Additionally, add `--validate_out=.` to your `protoc` command to apply `protoc-gen-validate`.

2. The second approach is to manually call `Validate()` in the server layer:
```go
func (g *Greeter) Hello(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
    if err := req.Validate(); err != nil {
        return err
    }
    rsp.Greeting = "Hello " + req.Name
    return nil
}
```

The advantage of the first approach is that it embeds validation in the program without requiring any business operations, but it introduces a certain level of intrusiveness. The second approach is more like a manual mode, where you call validation when needed.

# Supported Rules

## Struct Types

### Required
```protobuf
message Required {
    // @required:true
    Foo a = 1;
}

message Foo {
    string a = 1;
    string b = 2;
}
```

## Numeric Types
(float, double, int32, int64, uint32, uint64 , sint32, sint64, fixed32, fixed64, sfixed32, sfixed64)

### Equality
```protobuf
// @eq:1.23
float a = 1;
```

### Greater than, less than, equal to
```protobuf
// @lt:20
// @gt:10
int32 b = 2;
// @lte:20
// @gte:10
uint64 c = 3;
```

### Inclusion in an array
```protobuf
// @in:[1,2,3]
fixed32 d = 4;
// @not_in:[1,2,3]
float e = 5;
```

### Open and closed intervals
```protobuf
// @range:(1,5)
float f = 6;
// @range:[1,5]
float g = 7;
```

## String Types

### Equality
```protobuf
// @eq:"bar"
string c = 3;
```

### Substring containment
```protobuf
// @contains:"bar"
string a = 1;
// @not_contains:"bar"
string b = 2;
```

### Inclusion in an array
```protobuf
// @in:["foo", "bar", "baz"]
string d = 4;
// @not_in:["foo", "bar", "baz"]
string e = 5;
```

### String length
```protobuf
// @len:5
string f = 6;
```

### String length range
```protobuf
// @min_len:5
// @max_len:10
string g = 7;
```

### Regular expression pattern
```protobuf
// @pattern:"(?i)^[0-9a-f]+$"
string h = 8;
```

### Prefix and suffix
```protobuf
// @prefix:"foo"
string i = 9;
// @suffix:"bar"
string j = 10;
```

### Common types
```protobuf
// @type:url
string k = 11;
// @type:phone
string l = 12;
// @type:email
string m = 13;
// @type:ip
string n = 14;
```


## Array Types

### Supports all single-item rules
```protobuf
// @gt:10
// @lt:10
repeated float e = 5;
```

### Supports array length control
```protobuf
// @min_items:1
// @max_items:2
repeated int32 a = 1;
```

### Supports array uniqueness check
```protobuf
// @unique:true
repeated int64 b = 2;
// @unique:true
repeated string c = 3;
```
