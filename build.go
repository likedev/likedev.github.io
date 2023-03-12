package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

const index_html_path = "assets/pages/index.html"

func processIndex() {
	CopyToOut("./"+index_html_path, "./output/"+GetFileName(index_html_path))
}

var post_paths []string

func visit(path string, di fs.DirEntry, err error) error {
	if !strings.Contains(path, ".html") || strings.Contains(path, "index.html") {
		return nil
	}
	fmt.Printf("process: %s\n", path)
	post_paths = append(post_paths, path)
	return nil
}

func processPostPages() {
	err := filepath.WalkDir("pages", visit)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(post_paths); i++ {
		fmt.Printf(post_paths[i])
	}
}

func main() {
	processPostPages()
}
