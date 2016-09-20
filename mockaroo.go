package mockaroo

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var SLICERANGE = "[1-10]"

type Setup struct {
	AbsolutePkgName string
	FullPkgName     string
	PkgName         string
	StrctName       string
	Vars            interface{}
	Import          bool
}

func NewSetup(fullPkg, strct string, ptr interface{}) Setup {
	oStrct := reflect.ValueOf(ptr).Elem().Type()
	s := Setup{
		AbsolutePkgName: oStrct.PkgPath(),
		FullPkgName:     fullPkg,
		StrctName:       strct,
	}
	if fullPkg == "main" {
		s.Import = false
		s.PkgName = GetShortPackage(s.AbsolutePkgName)
	} else {
		s.Import = true
		s.PkgName = GetShortPackage(fullPkg)
	}
	return s
}

func (s Setup) GetPkgPrefix(fullPkg string) string {
	if s.AbsolutePkgName == fullPkg {
		if s.FullPkgName == "main" {
			return ""
		}
		return GetShortPackage(fullPkg) + "."
	}
	return fullPkg + "."
}

type MockTypes struct {
	MTypes   []MockType
	Template string
	Pre      string
	TempName string
	Setup    Setup
}

func NewMockTypes(fullPkg, strct string, ptr interface{}) *MockTypes {
	m := MockTypes{}
	m.Setup = NewSetup(fullPkg, strct, ptr)
	m.Template = "{{ define \"vars\" }}"
	GetFields(ptr, "", &m)
	m.Template += "\n{{ end }}"
	return &m
}

type MockType struct {
	Name         string `json:"name"`
	PercentBlank int    `json:"percentBlank"`
	Formula      string `json:"formula"`
	Type         string `json:"type"`
	Min          int64  `json:"min,omitempty"`
	Max          uint64 `json:"max,omitempty"`
	Decimals     int    `json:"decimals"`
}

func NewMockType(name, typ string) MockType {
	return MockType{
		Name: name,
		Type: typ,
	}
}

func GetFields(ptr interface{}, start string, mockTypes *MockTypes) {
	var strct reflect.Value
	if reflect.TypeOf(ptr) == reflect.TypeOf(reflect.Value{}) {
		strct = ptr.(reflect.Value)
	} else {
		strct = reflect.ValueOf(ptr).Elem()
	}
	strctType := strct.Type()
	for i := 0; i < strct.NumField(); i++ {
		fld := strct.Field(i)
		mockType := MockType{Name: start + strctType.Field(i).Name}
		mockTypes.Template += "\n" + strctType.Field(i).Name + ":"
		if mockTypes.Pre == "" {
			mockTypes.Pre = "."
		}
		if !strings.Contains(mockTypes.TempName, "[") {
			mockTypes.TempName = mockType.Name
		}
		SetMockType(fld.Type(), mockType, mockTypes)
	}
}

func SetMockType(t reflect.Type, mockType MockType, mockTypes *MockTypes) {
	switch t.Kind() {
	case reflect.String:
		mockType.Type = "Full Name"
		mockTypes.Template += "\"{{ " + mockTypes.Pre + mockTypes.TempName + " }}\","
	case reflect.Int:
		mockType.Type = "Number"
		mockType.Max = 9223372036854775807
		mockType.Min = -9223372036854775808
		mockTypes.Template += "{{ " + mockTypes.Pre + mockTypes.TempName + " }},"
	case reflect.Int8:
		mockType.Type = "Number"
		mockType.Max = 127
		mockType.Min = -128
		mockTypes.Template += "{{ " + mockTypes.Pre + mockTypes.TempName + " }},"
	case reflect.Int16:
		mockType.Type = "Number"
		mockType.Max = 32767
		mockType.Min = -32768
		mockTypes.Template += "{{ " + mockTypes.Pre + mockTypes.TempName + " }},"
	case reflect.Int32:
		mockType.Type = "Number"
		mockType.Max = 2147483647
		mockType.Min = -2147483648
		mockTypes.Template += "{{ " + mockTypes.Pre + mockTypes.TempName + " }},"
	case reflect.Int64:
		mockType.Type = "Number"
		mockType.Max = 9223372036854775807
		mockType.Min = -9223372036854775808
		mockTypes.Template += "{{ " + mockTypes.Pre + mockTypes.TempName + " }},"
	case reflect.Uint:
		mockType.Type = "Number"
		mockType.Max = 18446744073709551615
		mockType.Min = 0
		mockTypes.Template += "{{ " + mockTypes.Pre + mockTypes.TempName + " }},"
	case reflect.Uint8:
		mockType.Type = "Number"
		mockType.Max = 255
		mockType.Min = 1
		mockTypes.Template += "{{ " + mockTypes.Pre + mockTypes.TempName + " }},"
	case reflect.Uint16:
		mockType.Type = "Number"
		mockType.Max = 65535
		mockType.Min = 1
		mockTypes.Template += "{{ " + mockTypes.Pre + mockTypes.TempName + " }},"
	case reflect.Uint32:
		mockType.Type = "Number"
		mockType.Max = 4294967295
		mockType.Min = 1
		mockTypes.Template += "{{ " + mockTypes.Pre + mockTypes.TempName + " }},"
	case reflect.Uint64:
		mockType.Type = "Number"
		mockType.Max = 18446744073709551615
		mockType.Min = 0
		mockTypes.Template += "{{ " + mockTypes.Pre + mockTypes.TempName + " }},"
	case reflect.Float32:
		mockType.Type = "Number"
		mockType.Max = 32767
		mockType.Min = -32768
		mockType.Decimals = 4
		mockTypes.Template += "{{ " + mockTypes.Pre + mockTypes.TempName + " }},"
	case reflect.Float64:
		mockType.Type = "Number"
		mockType.Max = 2147483647
		mockType.Min = -2147483648
		mockType.Decimals = 4
		mockTypes.Template += "{{ " + mockTypes.Pre + mockTypes.TempName + " }},"
	case reflect.Bool:
		mockType.Type = "Boolean"
		mockTypes.Template += "{{ " + mockTypes.Pre + mockTypes.TempName + " }},"
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Struct {
			mockType.Type = "JSON Array"
			mockTypes.MTypes = append(mockTypes.MTypes, mockType)
			mockTypes.Template += " []" + mockTypes.Setup.GetPkgPrefix(t.Elem().PkgPath()) + t.Elem().Name() + "{ {{ range $" + mockTypes.TempName + " := ." + mockTypes.TempName + " }}\n"
		} else {
			mockTypes.Template += " " + t.String() + "{ {{ range $" + mockTypes.TempName + " := ." + mockTypes.TempName + " }}\n"
			mockType.Name = mockTypes.TempName + SLICERANGE
		}
		mockTypes.Pre = "$"
		SetMockType(t.Elem(), mockType, mockTypes)
		mockTypes.Template += "{{ end }}\n},"
		mockTypes.Pre = "."
		mockTypes.TempName = ""
		return
	case reflect.Struct:
		mockTypes.Template += mockTypes.Setup.GetPkgPrefix(t.PkgPath()) + t.Name() + "{"
		GetFields(reflect.Indirect(reflect.New(t)), mockType.Name+".", mockTypes)
		mockTypes.Template += "\n},"
		return
	}
	mockTypes.MTypes = append(mockTypes.MTypes, mockType)
}

func GetData(mockTypes []MockType, count int) []byte {
	b, err := json.Marshal(mockTypes)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post("http://www.mockaroo.com/api/generate.json?key=f1876740&count="+strconv.Itoa(count)+"&array=true", "application/json", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func GetDataPretty(mockTypes []MockType, count int) string {
	b, err := json.Marshal(mockTypes)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post("http://www.mockaroo.com/api/generate.json?key=f1876740&count="+strconv.Itoa(count)+"&array=true", "application/json", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	if err := json.Indent(&buf, body, "", "\t"); err != nil {
		panic(err)
	}
	return buf.String()
}

func GenVars(mockTypes *MockTypes) {
	t, err := template.New("temp").Parse(StartTemp)
	if err != nil {
		panic(err)
	}
	t, err = t.Parse(mockTypes.Template)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	if err := t.ExecuteTemplate(buf, "temp", mockTypes); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(strings.Replace(mockTypes.Setup.FullPkgName, "/", ".", -1)+"-"+mockTypes.Setup.StrctName+".go", buf.Bytes(), 0664); err != nil {
		panic(err)
	}
}

var StartTemp = `
package main

{{ if .Setup.Import }}import "{{ .Setup.FullPkgName }}"{{ end }}

var users []{{ if .Setup.Import }}{{ .Setup.PkgName }}.{{ end }}{{ .Setup.StrctName }} = []{{ if .Setup.Import }}{{ .Setup.PkgName }}.{{ end }}{{ .Setup.StrctName }}{ {{ range .Setup.Vars }}
	    { {{ template "vars" . }}
	},{{ end }}
}
`

func GetShortPackage(fullPkg string) string {
	pkgSplit := strings.Split(fullPkg, "/")
	return pkgSplit[len(pkgSplit)-1]
}
