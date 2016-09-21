package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/gregpechiro/mockaroo"
	"github.com/gregpechiro/mockaroo/cli"
)

var fullPkg string   //= "main"
var strct string     // = "User"
var strctFile string //= "user.go"
var count int
var match bool

func main() {

	flag.StringVar(&fullPkg, "p", "", "Full package from GOPATH to the struct. Your can use main if the struct is in the main program you are working on. This requires you to use -f")
	flag.StringVar(&strct, "s", "", "The Name of the struct you want to generate data for. This is case sensitive")
	flag.StringVar(&strctFile, "f", "", "The file where the struct is located. If -p is anything other than \"main\" this will be ignored")
	flag.IntVar(&count, "c", 0, "The number of structs in a slice you want generated. It will always generate a slice even if count is 1")
	flag.Parse()

	if fullPkg == "" {
		fmt.Println("-p (package) cannot be empty")
		return
	}
	if strct == "" {
		fmt.Println("-s (struct name) cannot be empty")
		return
	}
	if fullPkg == "main" && strctFile == "" {
		fmt.Println("-f (struct file) cannot be empty when -p (package) is set to main")
		return
	}
	if count < 1 {
		fmt.Println("-c (count) must be greater than 0")
		return
	}

	pkg := mockaroo.GetShortPackage(fullPkg)

	t, err := template.New("temp").Parse(cli.FILE)
	if err != nil {
		log.Println("Error parsing")
		panic(err)
	}

	buf := new(bytes.Buffer)

	if err := t.Execute(buf, map[string]interface{}{
		"fullPackage": fullPkg,
		"package":     pkg,
		"struct":      strct,
		"count":       count,
		"match":       match,
	}); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile("mockaroo-temp.go", buf.Bytes(), 0664); err != nil {
		panic(err)
	}

	if fullPkg != "main" {
		strctFile = ""
	}

	finalfile := strings.Replace(fullPkg, "/", ".", -1) + "-" + strct + ".go"

	if err := exec.Command("go", "run", "mockaroo-temp.go", strctFile).Run(); err != nil {
		panic(err)
	}

	if err := exec.Command("gofmt", "-w", finalfile).Run(); err != nil {
		panic(err)
	}

	os.Remove("mockaroo-temp.go")

}
