package to traverse through json objects 

Install
=======

    go get github.com/pchang211/go-jsontree/exp/jsonpath

Usage
======

Implements a subset of the path extraction tool for JSON discussed in http://goessner.net/articles/JsonPath/

Supports simple path selection and array indexing. For example, given the body

```
{"foo":{"bar":"baz", "array":["hello", "world"]}}
```

selection with the following will yield

```
$.foo.bar
>> "baz"
$.foo.array[0]
>> "hello"
$.foo.array[1]
>> "world"
```

Author
======

Philip Chang [pchang211@gmail.com]
Adapted from: 
Bryan Matsuo [bryan dot matsuo at gmail dot com]

