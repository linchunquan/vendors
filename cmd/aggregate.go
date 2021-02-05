package main

import (
	"bufio"
	"bytes"
	"github.com/kardianos/govendor/cliprompt"
	"github.com/kardianos/govendor/run"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sort"
	"log"
	"io/ioutil"
)

func aggregate(projects []string, output string) error {

	gopath := getGoPathDirectory()
	for i, p := range projects {
		if !strings.HasPrefix(p, gopath) {
			projects[i] = filepath.Join(gopath, "src", p) + string(os.PathSeparator)
		}
	}

	excludePrefix := []string{}
	for _, project := range projects {
		//log.Printf("==================>exclude:%s\n",project)
		excludePrefix = append(excludePrefix, project)
	}

	packageMap := map[string]bool{}

	for _, p := range projects {
		listPackage(gopath, p, excludePrefix, packageMap)
	}

	copyPackages(gopath, packageMap, output)

	return nil
}

func getGoPathDirectory() string {
	gopath := os.Getenv("GOPATH")
	if len(gopath) == 0 {
		panic("$GOPATH not set")
	}
	return gopath
}

func listPackage(gopath, prjSource string, excludePrefix []string, packageMap map[string]bool) {
	os.Chdir(prjSource)
	prompt := &cliprompt.Prompt{}
	allArgs := []string{"", "list","-v"}
	buf := &bytes.Buffer{}
	_, err := run.Run(buf, allArgs, prompt)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(buf)
	for {
		line, _, err := r.ReadLine()
		if err != nil && err != io.EOF {
			panic(err)
		}
		if len(line) == 0 {
			break
		}
		line = bytes.TrimSpace(line)
		if bytes.HasPrefix(line, []byte("├")) || bytes.HasPrefix(line, []byte("└")){
			continue
		}

		p := bytes.Split(line, []byte(" "))
		if len(p) >= 2 {
			l := bytes.TrimLeft(bytes.TrimSpace(p[len(p)-1]), "::")
			if len(l) > 0 {
				location := filepath.Join(gopath, "src", string(l))
				if !exclude(location, excludePrefix){
					if _, ok := packageMap[location]; !ok {
						//log.Printf("======> add location:%s\n", location)
						packageMap[location] = true
					}
				}
			}
		}
	}
}

func exclude(location string, excludePrefix []string)bool{
	for _, prefix := range excludePrefix{
		if strings.HasPrefix(location, prefix){
			return true
		}
	}
	return false
}

func copyPackages(gopath string, packageMap map[string]bool, output string){
	var packages []string
	for k, _ := range packageMap{
		packages = append(packages, k)
	}

	sort.Slice(packages, func(i, j int) bool {
		return strings.Compare(packages[i], packages[j]) < 0
	})

	for _, p := range packages{
		//log.Printf("======> ADD location:%s\n", p)
		files, err := ioutil.ReadDir(p)
		if err!=nil{
			log.Printf(" ======>error:%v\n",err)
		//	panic(err)
		}
		for _, file := range files{
			if !file.IsDir(){
				srcPath := filepath.Join(p,file.Name())
				//log.Printf("to copy file:%s\n", srcPath)
				dstPath := filepath.Join(output,srcPath[len(gopath)+4:])
				//log.Printf("to copy file:%s to:%s\n", srcPath, dstPath)
				err := Copy(srcPath, dstPath)
				if err!=nil{
					panic(err)
				}
			}

		}
	}
	log.Printf("packageMap size:%d\n", len(packageMap))
}
