package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const index_html_path = "assets/pages/index.html"

func GetFileName(path string) string {
	delg := "\\"
	if strings.Contains(path, "/") {
		delg = "/"
	}
	lastIndex := strings.LastIndex(path, delg)
	if lastIndex < 0 {
		return path
	}
	return path[lastIndex+1:]
}

func CopyToOut(path string, dst_dir string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		os.Exit(2)
	}
	// Write data to dst
	err = ioutil.WriteFile(dst_dir, data, 0644)
	if err != nil {
		os.Exit(2)
	}
}

func main() {
	CopyToOut("./"+index_html_path, "./output/"+GetFileName(index_html_path))
	fmt.Print("build success")
}
