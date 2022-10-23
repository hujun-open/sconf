//Copyright MIT 2022 Hu Jun

/*
Package sconf provides simple configuration managment for Golang application, with following features:

  - read configuration from flag and/or YAML file, mix&match, into a struct
  - 3 sources: default value, YAML file and flag
  - priority:in case of multiple source returns same config struct field, the preference is flag over YAML over default value

struct field naming:
  
  - YAML: lowercase
  - flag: lowercase, for nested struct field, it is "-<struct_fieldname_level_0>-<struct_fieldname_level_1>...", see example below:



Example:

	package main

	import (
		"fmt"
		"sconf"
	)

	type Company struct {
		Name string
	}

	type Employee struct {
		Name, Addr string
		Employer   Company
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
		//create a new SConf instance with command line flags, with no default config file path
		cnf, err := sconf.NewSConfCMDLine(def, "")
		if err != nil {
			panic(err)
		}
		// read configuration from all sources
		ferr, aerr := cnf.ReadwithCMDLine()
		fmt.Printf("ferr %v,aerr %v\n", ferr, aerr)
		//return the configuration in a Employee struct
		fmt.Printf("final result is %+v\n", cnf.GetConf())
	}

Output:

- no command line args, no config file, default is used
	
	.\test.exe
	ferr <nil>,aerr <nil>
	final result is {Name:defName Addr:defAddr Employer:{Name:defCom}}
- config file via "-f" command args, value from file take procedence

	.\test.exe -f .\test.yaml
	ferr <nil>,aerr <nil>
	final result is {Name:nameFromFile Addr:addrFromFile Employer:{Name:comFromFile}}
- mix command line args and config file, args to override employee name:

	.\test.exe -f .\test.yaml -name nameFromArgs
	ferr <nil>,aerr <nil>
	final result is {Name:nameFromArgs Addr:addrFromFile Employer:{Name:comFromFile}}

- mix command line args and config file, args to override company name:

	.\test.exe -f .\test.yaml -employer-name comFromArgs
	ferr <nil>,aerr <nil>
	final result is {Name:nameFromFile Addr:addrFromFile Employer:{Name:comFromArgs}}

*/
package sconf
