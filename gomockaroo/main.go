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
	"sort"
	"strings"

	"github.com/gregpechiro/mockaroo"
	"github.com/gregpechiro/mockaroo/cli"
)

const (
	packageUsage = "Full package from GOPATH to the struct. You can use main if the\n" +
		"\tstruct is in the main program you are working on.\n" +
		"\tThis requires you to use \"-f\"\n"

	structUsage = "The Name of the struct you want to generate data for.\n" +
		"\tThis is case sensitive\n"

	filesUsage = "A comma separated list (`[]string`) of file names. If the struct is \n" +
		"\tin main this is the file that the struct is located in and all other\n" +
		"\tdependencies that strut may have that are also delared in main.\n" +
		"\tIf -p is anything other than \"main\" this will be ignored\n"

	countUsage = "The number of structs in a slice you want generated. It will always\n" +
		"\tgenerate a slice even if count is 1\n"

	matchUsage = "Whether regular expresion should be used on the struct\n" +
		"\tfield names to determine mockaroo types. If this flag is present \n" +
		"\twithout an argument it will be marked true. This will take more\n" +
		"\ttime and resources\n"

	maxUsage = "The maximum number of indices that should be generated for\n" +
		"\tthe structs nested slices. This will default to 10 if not\n" +
		"\tprovided or set less than 1\n"

	packageDefault = ""
	structDefault  = ""
	filesDefault   = ""
	countDefault   = 10
	matchDefault   = false
	maxDefault     = 0
)

var fullPkg string
var strct string
var strctFile string
var count int
var match bool
var maxSlice int

type FlagSlice []string

func (f *FlagSlice) String() string {
	return strings.Join(*f, " ")
}

func (f *FlagSlice) Set(value string) error {
	*f = strings.Split(value, ",")
	return nil
}

var strctFiles FlagSlice
var flagSet *flag.FlagSet

func init() {

	flagSet = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage of %s:\n\n", os.Args[0])
		m := map[string]string{}
		flagSet.VisitAll(func(f *flag.Flag) {
			name, usage := flag.UnquoteUsage(f)
			s, ok := m[usage]
			s = fmt.Sprintf("  -%s%s", f.Name, s)
			if !ok {
				if len(name) > 0 {
					s += "  " + name
				}
			}
			s += "\n    \t"
			m[usage] = s
		})
		var ss []string
		for u, f := range m {
			ss = append(ss, f+u)
		}
		sort.Strings(ss)
		for _, s := range ss {
			fmt.Fprint(os.Stderr, s, "\n")
		}
	}

	flagSet.StringVar(&fullPkg, "p", packageDefault, packageUsage)
	flagSet.StringVar(&fullPkg, "package", packageDefault, packageUsage)

	flagSet.StringVar(&strct, "s", structDefault, structUsage)
	flagSet.StringVar(&strct, "struct", structDefault, structUsage)

	flagSet.Var(&strctFiles, "f", filesUsage)
	flagSet.Var(&strctFiles, "files", filesUsage)

	flagSet.IntVar(&count, "c", countDefault, countUsage)
	flagSet.IntVar(&count, "count", countDefault, countUsage)

	flagSet.BoolVar(&match, "m", matchDefault, matchUsage)
	flagSet.BoolVar(&match, "match", matchDefault, matchUsage)

	flagSet.IntVar(&maxSlice, "mx", maxDefault, maxUsage)
	flagSet.IntVar(&maxSlice, "max", maxDefault, maxUsage)

}

func main() {

	flagSet.Parse(os.Args[1:])

	if fullPkg == "" {
		fmt.Println("-p (package) cannot be empty")
		return
	}
	if strct == "" {
		fmt.Println("-s (struct name) cannot be empty")
		return
	}
	if fullPkg == "main" && len(strctFiles) == 0 {
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
		"maxSlice":    maxSlice,
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

	if err := exec.Command("go", append([]string{"run", "mockaroo-temp.go"}, strctFiles...)...).Run(); err != nil {
		panic(err)
	}

	if err := exec.Command("gofmt", "-w", finalfile).Run(); err != nil {
		panic(err)
	}

	os.Remove("mockaroo-temp.go")

}
