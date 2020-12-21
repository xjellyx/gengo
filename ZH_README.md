## gengo 
    一个通过定义golang结构体自动生成curd代码，根据mvc模型生成三层结构代码，初始化项目完修改配置文件直接可以运行，无需自己搭建项目，已经生成的代码不会覆盖，
自动生成swagger文档

## 安装
```console
go get -u github.com/olongfen/gengo/cmd/gengo
```

## gengo 命令
```console
gengo --help

NAME:
   gengo - A new cli application

USAGE:
   gengo [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --output value, -o value  The name of schema output to generate output 
   --input value, -i value   The name of the input go file path
   --mod value, -m value     The name of project go module
   --web value, -w value     The name of project web frame (default: "gin")
   --orm value, -r value     The name of project orm frame (default: "gorm")
   --transformError, -t      The name of transform db err (default: true)
   --help, -h                show help (default: false)

```

## Web 框架
- [gin](github.com/gin-gonic/gin)

## orm框架
- [gorm](gorm.io/gorm)

## 使用
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
   
## 示例
- [demo](https://github.com/olongfen/demo)

## 证书
- MIT License
