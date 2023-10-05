# go-clone

[![GoDoc](https://godoc.org/github.com/lmlat/go-clone?status.png)](https://pkg.go.dev/github.com/lmlat/go-clone)
[![Go Report Card](https://goreportcard.com/badge/github.com/lmlat/go-clone)](https://goreportcard.com/report/github.com/lmlat/go-clone)
[![Coverage Status](https://coveralls.io/repos/github/lmlat/go-clone/badge.svg?branch=master)](https://coveralls.io/github/lmlat/go-clone?branch=master)

This project is based on the Go language development cloning toolkit and supports shallow and deep copy operations of all types of data in the Go language.

When deep or shallow copy struct type data, you can manipulate exported or unexported fields.

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

Example:
deep copy the value of a struct type, containing the fields in the struct that are not exported. 
```go
type C struct {
	name     string
	Age      int
	birthday time.Time
	hobby    []string
}

src := C{"aitao", 100, time.Now(), []string{"ping pong", "badminton", "football"}}
dst := Deep(src).(C) // deep copy value

fmt.Printf("%+v\n", src)
fmt.Printf("\n%+v\n", dst)

dst.name = "小阿梦" // modify an unreferenced type
fmt.Printf("%+v\n", src)
fmt.Printf("\n%+v\n", dst)

dst.hobby[0] = "乒乓球" // modify a referenced type
fmt.Printf("%+v\n", src)
fmt.Printf("\n%+v\n", dst)
```

deep copy the value of a struct pointer type, containing the fields in the strcut that are not exported.

```go	
src := &C{"aitao", 100, time.Now(), []string{"ping pong", "badminton", "football"}}
dst := Deep(src).(C) // deep copy value
```