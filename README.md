# go-clone
[![GoDoc](https://godoc.org/github.com/lmlat/go-clone?status.png)](https://pkg.go.dev/github.com/lmlat/go-clone)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/lmlat/go-clone)
[![Go Report Card](https://goreportcard.com/badge/github.com/lmlat/go-clone)](https://goreportcard.com/report/github.com/lmlat/go-clone)
[![Coverage Status](https://coveralls.io/repos/github/lmlat/go-clone/badge.svg?branch=master)](https://coveralls.io/github/lmlat/go-clone?branch=master)
[![Build Status](https://github.com/lmlat/go-clone/actions/workflows/go-ci.yml/badge.svg)](https://github.com/lmlat/go-clone/actions/workflows/go-ci.yml)
[![Release](https://img.shields.io/github/release/lmlat/go-clone.svg?style=flat-square)](https://github.com/lmlat/go-clone/releases)
![GitHub](https://img.shields.io/github/license/lmlat/go-clone)



This project is based on the Go language development cloning toolkit and supports shallow and deep copy operations of all types of data in the Go language.

When deep or shallow copy struct type data, you can manipulate exported or unexported fields.

## Feature

- Provides shallow and deep copy functions.
- Comprehensive，efficient and reusable，supports shallow and deep copies for struct exported and unexported fields.
- Does not depends on any third-party libraries.
- Provide unit test cases for each exported functions.

## Usage

Install:
```go
go get "github.com/lmlat/go-clone"
```
Import:
```go
import "github.com/lmlat/go-clone"
```
> Note: 
>
> After importing the package, by default all provided functionality is defined in a package named `clone`.

## Example

### Deep Copy

Deep copy the value of a struct type, containing the fields in the struct that are not exported. 

```go
type C struct {
	name     string
	Age      int
	birthday time.Time
	hobby    []string
}

src := C{"aitao", 100, time.Now(), []string{"ping pong", "badminton", "football"}}
dst := Deep(src).(C)
// When you modify the value of non-reference type, the source data is not affected.
dst.name = "哆啦A梦" 
// When you modify the value of reference type, the source data is not affected.
dst.hobby[0] = "乒乓球"
```

Deep copy the value of a struct pointer type, containing the fields in the strcut that are not exported.

```go	
src := &C{"aitao", 100, time.Now(), []string{"ping pong", "badminton", "football"}}
dst := Deep(src).(*C) // deep copy value
```
Deep copy only the exported fields in the structure:

```go	
dst = Deep(src, WithOpFlags(OnlyPublicField)).(C)
// The same effect can be achieved using CopyProperties function.
dst = CopyProperties(src).(C)
```

Some commonly used operation flags:

> 1. OnlyPublicField：only exported field values in the struct are copied, only struct type values.
> 1. OnlyPrivateField：only field values that are not exported in the struct are copied, only struct type values.
> 1. AllFields：copy all field values in the struct, including exported and unexported fields, only for struct type values.
> 1. DeepString：deep copy string type values.
> 1. DeepFunc：deep copy func type values.
> 1. DeepArray：deep copy array type values.

Some special types of shallow copy processing:

1. time.Time
1. reflect.Type：The type system of the Go language is static at compile time, meaning that the type information is determined at compile time. Therefore, when dealing with the 'reflect.Type' interface, it is usually not necessary to make a deep copy (' reflect.rtype 'is immutable).
  Shallow copy the value of a struct type, containing the fields in the struct that are not exported. 

### Shallow Copy

```go
src := C{"aitao", 100, time.Now(), []string{"ping pong", "badminton", "football"}}
dst := Shallow(src).(C)
```
