# gengo 
> gengo It is an automatic code generation tool that uses go:generate to generate corresponding addition, deletion, modification, and inspection logic and interfaces by defining a structure
> 
> ** Thank you!**

## install
```console
go get -u github.com/olongfen/gengo/cmd/gengo
```

## gengo clu
```console
gengo --help

Usage of gengo:
  -imports string
        [Required] The name of the import  to import package
  -input string
        [Required] The name of the input go file path
  -mod string
        [Required] The name of project go module
  -orm string
        [Option] The name of project orm frame
  -output string
        [Required] The name of schema output to generate output for, comma separated
  -tfErr
        [Option] The name of transform db err
  -web string
        [Option] The name of project web frame
```