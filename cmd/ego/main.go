package main

import (
	"bytes"
	"flag"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/infinitbyte/ego"
	"io/ioutil"
	"strings"
)

// version is set by the makefile during build.
var version string

func main() {
	outfile := flag.String("o", "ego.go", "output file")
	pkgname := flag.String("package", "", "package name")
	versionFlag := flag.Bool("version", false, "print version")
	flag.Parse()
	log.SetFlags(0)

	// If the version flag is set then print the version.
	if *versionFlag {
		//fmt.Printf("ego v%s\n", version)
		return
	}

	// If no paths are provided then use the present working directory.
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// If no package name is set then use the directory name of the output file.
	if *pkgname == "" {
		abspath, _ := filepath.Abs(*outfile)
		*pkgname = filepath.Base(filepath.Dir(abspath))
		*pkgname = regexp.MustCompile(`(\w+).*`).ReplaceAllString(*pkgname, "$1")
	}

	// Recursively retrieve all ego templates
	var v visitor
	//for _, root := range roots {
	//	if err := filepath.Walk(root, v.visit); err != nil {
	//		scanner.PrintError(os.Stderr, err)
	//		os.Exit(1)
	//	}
	//}

	for _, root := range roots {
		v.listAll(root)
	}

}

func processTemplate(templates []*ego.Template,pkgname string,outfile string)  {
	//fmt.Println("process template, output:",outfile,",package:"+pkgname)

	// Write package to output file.
	p := &ego.Package{Templates: templates, Name: pkgname}

	var buf bytes.Buffer
	// Write template to buffer.
	if err := p.Write(&buf); err != nil {
		log.Fatal("template write: ", err)
	}

	result := buf.Bytes()
	var err error
	if result, err = format.Source(result); err != nil {
		log.Fatal("format: ", err)
	}

	f, err := os.Create(outfile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err = f.Write(result); err != nil {
		log.Fatal("write: ", err)
	}
}

// visitor iterates over
type visitor struct {

}

func (v *visitor)listAll(path string) {
	//fmt.Println("process path:",path)

	files, _ := ioutil.ReadDir(path)
	var paths []string
	for _, fi := range files {
		if fi.IsDir() {
			v.listAll(path + "/" + fi.Name())
			//println(path + "/" + fi.Name())
		}
		if (filepath.Ext(fi.Name()) == ".ego" ||filepath.Ext(fi.Name()) == ".html"||filepath.Ext(fi.Name()) == ".htm") {
			//println(path + "/" + fi.Name())
			paths = append(paths, path + "/" + fi.Name())
		}
	}


	//fmt.Println("generate template:",path)
	// Parse every template file.
	var templates []*ego.Template
	for _, path := range paths {
		t, err := ego.ParseFile(path)
		if err != nil {
			log.Fatal("parse file, ",path,", ", err)
		}
		templates = append(templates, t)
	}

	// If we have no templates then exit.
	if len(templates) == 0 {
		return
	}

	ap,_:=filepath.Abs(path)
	as:=strings.Split(ap,"/")

	lastDirName:=as[len(as)-1]

	//fmt.Println("path name: ",lastDirName)

	processTemplate(templates,lastDirName,path+"/ego.go")
}

//func (v *visitor) visit(path string, info os.FileInfo, err error) error {
//	if info == nil {
//		return fmt.Errorf("file not found: %s", path)
//	}
//
//	fmt.Println("walk file: ",path)
//
//	if !info.IsDir() &&(filepath.Ext(path) == ".ego" ||filepath.Ext(path) == ".html"||filepath.Ext(path) == ".htm") {
//		v.paths = append(v.paths, path)
//	}
//	return nil
//}
