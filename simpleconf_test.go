package sconf

import (
	"flag"
	"fmt"
	"testing"
)

type company struct {
	Name string
}

type testStruct struct {
	Name, Addr string
	Employer   company
}

func (t testStruct) isEqual(peer testStruct) bool {
	if t.Name == peer.Name && t.Addr == peer.Addr {
		if t.Employer.Name == peer.Employer.Name {
			return true
		}
	}
	return false
}

const (
	defpath string = "./testdata/test.yaml"
)

type testSetup struct {
	def        testStruct
	fpath      string
	args       []string
	result     testStruct
	expectFail bool
}

func doTest(t *testing.T, setup testSetup) error {
	fset := flag.NewFlagSet("testflagset", flag.ContinueOnError)
	cnf, err := NewSConf(setup.def, setup.fpath, fset)
	if err != nil {
		return err
	}
	ferr, aerr := cnf.Read(setup.args)
	t.Logf("ferr is %v, aerr is %v", ferr, aerr)
	if !cnf.GetConf().isEqual(setup.result) {
		return fmt.Errorf("actual result is %+v, different from expected result %+v", cnf.GetConf(), setup.result)
	}
	return nil

}
func TestSconf(t *testing.T) {
	defCnf := testStruct{
		Name:     "defName",
		Addr:     "defAddr",
		Employer: company{Name: "defCom"},
	}
	// fs := flag.NewFlagSet("test", flag.ContinueOnError)
	// cnf, err := NewSConf(defCnf, "", fs)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// buf, err := cnf.MarshalYAML()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println(string(buf))
	// return
	caseList := []testSetup{
		{ // case 0, result should be value from file
			def:   defCnf,
			fpath: defpath,
			args:  []string{},
			result: testStruct{
				Name:     "nameFromFile",
				Addr:     "addrFromFile",
				Employer: company{Name: "comFromFile"},
			},
		},
		{ // case 1, specify config file in args, result should be value from file
			def:   defCnf,
			fpath: "",
			args:  []string{"-f", defpath},
			result: testStruct{
				Name:     "nameFromFile",
				Addr:     "addrFromFile",
				Employer: company{Name: "comFromFile"}},
		},
		{
			// case 2,no file, no args, result should be default
			def:   defCnf,
			fpath: "",
			args:  []string{},
			result: testStruct{
				Name:     "defName",
				Addr:     "defAddr",
				Employer: company{Name: "defCom"}},
		},
		{ // case 3, both args and file, arg should win
			def:   defCnf,
			fpath: "",
			args:  []string{"-f", defpath, "-name", "nameFromArg"},
			result: testStruct{
				Name:     "nameFromArg",
				Addr:     "addrFromFile",
				Employer: company{Name: "comFromFile"}},
		},
		{ // case 4, mix arg and default, arg should win
			def:   defCnf,
			fpath: "",
			args:  []string{"-name", "nameFromArg", "-employer-name", "argCom"},
			result: testStruct{
				Name:     "nameFromArg",
				Addr:     "defAddr",
				Employer: company{Name: "argCom"}},
		},
		{ // case 5, specify nonexist config file, result should be default
			def:   defCnf,
			fpath: "",
			args:  []string{"-f", "dosntexist"},
			result: testStruct{
				Name:     "defName",
				Addr:     "defAddr",
				Employer: company{Name: "defCom"}},
		},
		{ // case 6, specify nonexist config file and args, args should win
			def:   defCnf,
			fpath: "",
			args:  []string{"-f", "dosntexist", "-addr", "addrFromArg"},
			result: testStruct{
				Name:     "defName",
				Addr:     "addrFromArg",
				Employer: company{Name: "defCom"}},
		},
	}
	for i, c := range caseList {
		err := doTest(t, c)
		if err != nil {
			t.Logf("case %d fails with err %v", i, err)
			if !c.expectFail {
				t.Fatal()
			}
		} else {
			t.Logf("case %d finished successfully", i)
		}

	}

}
