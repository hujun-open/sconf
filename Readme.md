![Build Status](https://github.com/hujun-open/sconf/actions/workflows/main.yml/badge.svg)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hujun-open/sconf)](https://pkg.go.dev/github.com/hujun-open/sconf)

Package sconf provides simple configuration managment for Golang application, with following features:

  - read configuration from flag and/or YAML file, mix&match, into a struct
  - 3 sources: default value, YAML file and flag
  - priority:in case of multiple source returns same config struct field, the preference is flag over YAML over default value

struct field naming:
  
  - YAML: lowercase
  - flag: lowercase, for nested struct field, it is "-<struct_fieldname_level_0>-<struct_fieldname_level_1>...", see example below:



Example:
```
package main

import (
	"fmt"
	"sconf"
)

type Company struct {
	//the usage tag is used for command line usage
	Name string `usage:"company name"`
}

type Employee struct {
	Name     string `usage:"employee name"`
	Addr     string `usage:"employee address"`
	Employer Company
}

func main() {
	//default config
	def := Employee{
		Name: "defName",
		Addr: "defAddr",
		Employer: Company{
			Name: "defCom",
		},
	}
	cnf, err := sconf.NewSConfCMDLine(def, "")
	if err != nil {
		panic(err)
	}
	ferr, aerr := cnf.ReadwithCMDLine()
	fmt.Printf("ferr %v,aerr %v\n", ferr, aerr)
	fmt.Printf("final result is %+v\n", cnf.GetConf())
}
```
Output:

- no command line args, no config file, default is used
```	
.\test.exe
ferr <nil>,aerr <nil>
final result is {Name:defName Addr:defAddr Employer:{Name:defCom}}
```    
- config file via "-f" command args, value from file take procedence
```
.\test.exe -f .\test.yaml
ferr <nil>,aerr <nil>
final result is {Name:nameFromFile Addr:addrFromFile Employer:{Name:comFromFile}}
```
- mix command line args and config file, args to override employee name:
```
.\test.exe -f .\test.yaml -name nameFromArgs
ferr <nil>,aerr <nil>
final result is {Name:nameFromArgs Addr:addrFromFile Employer:{Name:comFromFile}}
```
- mix command line args and config file, args to override company name:
```
.\test.exe -f .\test.yaml -employer-name comFromArgs
ferr <nil>,aerr <nil>
final result is {Name:nameFromFile Addr:addrFromFile Employer:{Name:comFromArgs}}
```
- command line usage
```
.\test.exe -?
flag provided but not defined: -?
Usage:
  -f <filepath> : read from config file <filepath>
  -addr <string> : employee address
        default:defAddr
  -employer-name <string> : company name
        default:defCom
  -name <string> : employee name
        default:defName
```