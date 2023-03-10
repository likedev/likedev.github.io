package main

import (
	"fmt"
	"io/ioutil"
	"log"
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
		fmt.Printf("读取文件错误： %s ", path)
		os.Exit(2)
	}
	// Write data to dst
	err = ioutil.WriteFile(dst_dir, data, 0644)
	if err != nil {
		fmt.Printf("写文件错误 %s", dst_dir)
		os.Exit(2)
	}
}

func main() {
	if _, err := os.Stat("./output"); os.IsNotExist(err) {
		err := os.Mkdir("./output", 0755)
		if err != nil {
			fmt.Printf("创建output文件夹失败")
			log.Fatal(err)
		}
	}
	CopyToOut("./"+index_html_path, "./output/"+GetFileName(index_html_path))
	fmt.Print("build success")
}
