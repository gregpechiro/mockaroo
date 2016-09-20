package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/gregpechiro/mockaroo"
	"github.com/gregpechiro/mockaroo/cli"
)

var fullPkg = "main"
var strct = "User"
var strctFile = "user.go"
var count = 10

func main() {

	pkg := mockaroo.GetShortPackage(fullPkg)

	t, err := template.New("temp").Parse(cli.FILE)
	if err != nil {
		log.Println("Error parsing")
		panic(err)
	}
	buf := new(bytes.Buffer)

	err = t.Execute(buf, map[string]interface{}{
		"fullPackage": fullPkg,
		"package":     pkg,
		"struct":      strct,
		"count":       count,
	})

	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile("mockaroo-temp.go", buf.Bytes(), 0664); err != nil {
		panic(err)
	}

	if fullPkg != "main" && fullPkg != "this" {
		strctFile = ""
	}

	finalfile := strings.Replace(fullPkg, "/", ".", -1) + "-" + strct + ".go"

	err = exec.Command("go", "run", "mockaroo-temp.go", strctFile).Run()
	if err != nil {
		panic(err)
	}

	err = exec.Command("gofmt", "-w", finalfile).Run()
	if err != nil {
		panic(err)
	}

	os.Remove("mockaroo-temp.go")

}
