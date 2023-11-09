- [English](README.md)
- [简体中文](README-zh.md)

# protoc-gen-validate
基于proto文件的注解，为每个message生成validate函数。

[文章](https://yflee.in/protoc-gen-validate.html)描述了protoc-gen-validate的前世今生，也欢迎大家加我的微信`lyf987667482`一起讨论。
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

生成的`Numerics`的校验方法如下：

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

# 安装
```bash
go install github.com/giftDad/protoc-gen-validate
```

# 例子
```bash
protoc  examples/foo.proto  --validate_out=Mexamples/foo.proto=examples/\;examples:. \
--go_out=Mexamples/foo.proto=examples/\;examples:.
go test -v ./examples/
```

# 应用
适用于基于protobuf的rpc框架 比如[gi-micro](https://github.com/go-micro/generator),[twirp](https://github.com/twitchtv/twirp)等等。
使用方式有两种,我们拿gi-micro举例：

1. 第一种是修改其 protoc-gen-xxx工具，比如gi-micro的[generator](https://github.com/go-micro/generator)
在`cmd/protoc-gen-micro/plugin/micro/micro.go:480`处增加:
```go
g.P("func (h *", unexport(servName), "Handler) ", methName, "(ctx ", contextPkg, ".Context, in *", inType, ", out *",outType, ") error {")
// 增加以下三行
g.P("if err := in.Validate();err != nil {")
g.P("return err")
g.P("}")

g.P("return h.", serveType, ".", methName, "(ctx, in, out)")
g.P("}")
g.P()
```
就能使其生成的方法变为:
```go
func (h *greeterHandler) Hello(ctx context.Context, in *Request, out *Response) error {
    if err := in.Validate();err != nil {
        return err
    }
    return h.GreeterHandler.Hello(ctx, in, out)
}
```
同时在`protoc`中增加`--validate_out=.`,就能正确应用`protoc-gen-validate`

2. 第二种是在server层手动执行`Validate()`:
```go
func (g *Greeter) Hello(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
    if err := req.Validate(); err != nil {
        return err
    }
    rsp.Greeting = "Hello " + req.Name
    return nil
}
```

第一种方式的优势在于无需任何业务操作，就能将validate嵌入程序中，但是会带有一定的侵入性。
而第二种则就像手动挡，在需要的时候手动调用。

# 支持规则

## 复合类型

### 必要
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

## 数字类型
(float, double, int32, int64, uint32, uint64 , sint32, sint64, fixed32, fixed64, sfixed32, sfixed64)

### 相等
```protobuf
// @eq:1.23
float a = 1;
```

### 大于小于等于
```protobuf
// @lt:20
// @gt:10
int32 b = 2;
// @lte:20
// @gte:10
uint64 c = 3;
```

### 是否在数组中
```protobuf
// @in:[1,2,3]
fixed32 d = 4;
// @not_in:[1,2,3]
float e = 5;
```

### 开闭区间
```protobuf
// @range:(1,5)
float f = 6;
// @range:[1,5]
float g = 7;
```

## 字符串类型

### 相等
```protobuf
// @eq:"bar"
string c = 3;
```

### 是否包含字串
```protobuf
// @contains:"bar"
string a = 1;
// @not_contains:"bar"
string b = 2;
```

### 是否在数组中
```protobuf
// @in:["foo", "bar", "baz"]
string d = 4;
// @not_in:["foo", "bar", "baz"]
string e = 5;
```

### 字符串长度
```protobuf
// @len:5
string f = 6;
```

### 字符串长度区间
```protobuf
// @min_len:5
// @max_len:10
string g = 7;
```

### 正则
```protobuf
// @pattern:"(?i)^[0-9a-f]+$"
string h = 8;
```

### 前缀后缀
```protobuf
// @prefix:"foo"
string i = 9;
// @suffix:"bar"
string j = 10;
```

### 常见类型
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


## 数组类型

### 支持所有单项规则
```protobuf
// @gt:10
// @lt:10
repeated float e = 5;
```

### 支持数组长度控制
```protobuf
// @min_items:1
// @max_items:2
repeated int32 a = 1;
```

### 支持数组判重
```protobuf
// @unique:true
repeated int64 b = 2;
// @unique:true
repeated string c = 3;
```
