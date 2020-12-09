## gengo 
> gengo It is an automatic golang code  to generate corresponding addition, deletion, modification, and inspection logic and interfaces by defining a structure
> 
> ** Thank you!**

## install
```console
go get -u github.com/olongfen/gengo/cmd/gengo
```

## gengo cli
```console
gengo --help

NAME:
   gengo - A new cli application

USAGE:
   gengo [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --output value, -o value  The name of schema output to generate output for, comma separated
   --input value, -i value   The name of the input go file path
   --mod value, -m value     The name of project go module
   --web value, -w value     The name of project web frame (default: "gin")
   --orm value, -r value     The name of project orm frame (default: "gorm")
   --transformError, -t      The name of transform db err (default: true)
   --help, -h                show help (default: false)

```

## Supported Web Frameworks
- [gin](github.com/gin-gonic/gin)

## Supported Orm Frameworks
- [gorm](gorm.io/gorm)

## Usage
```conlose
- mkdir demo 
- cd demo 
- go mod init demo 
- echo '
  package main
  
  import("gorm.io/gorm")
  
  type User struct {
          gorm.Model
          Name string 
          Age int
  }
  ' >> gen.go
- gengo -i ./gen.go -o . -m demo

```

## Todo list
- add orm frame 
   - xorm 
-  add web frame
   - echo
   - iris
   
## Example
- [demo](https://github.com/olongfen/demo)

## License
- MIT License
