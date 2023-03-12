package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

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

func GetDirByPath(path string) string {
	delg := "\\"
	if strings.Contains(path, "/") {
		delg = "/"
	}
	lastIndex := strings.LastIndex(path, delg)
	if lastIndex < 0 {
		return path
	}
	return path[:lastIndex]
}

func CopyToOut(src string, dst_file string) {

	dst_dir := GetDirByPath(dst_file)
	if _, err := os.Stat(dst_dir); os.IsNotExist(err) {
		err := os.MkdirAll(dst_dir, 0755)
		if err != nil {
			fmt.Printf("创建output文件夹失败")
			log.Fatal(err)
		}
	}
	data, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Printf("读取文件错误： %s ", src)
		os.Exit(2)
	}
	// Write data to dst
	err = ioutil.WriteFile(dst_file, data, 0644)
	if err != nil {
		fmt.Printf("写文件错误 %s", dst_file)
		os.Exit(2)
	}
}

func SearchDirRecursively(fileFullPath string) []string {
	files, err := ioutil.ReadDir(fileFullPath)
	if err != nil {
		log.Fatal(err)
	}
	var myFile []string
	for _, file := range files {
		path := strings.TrimSuffix(fileFullPath, "/") + "/" + file.Name()
		if file.IsDir() {
			subFile := SearchDirRecursively(path)
			if len(subFile) > 0 {
				myFile = append(myFile, subFile...)
			}
		} else {
			myFile = append(myFile, path)
		}
	}
	return myFile
}

func CopyDir(src string, dst string) {

}

/*
*
get filename => file content mappings in one dir
*/
func getDirFileContentMap(dir string) map[string]string {
	res := make(map[string]string)
	files := SearchDirRecursively(dir)
	for _, f := range files {
		fname := GetFileName(f)
		fname = strings.TrimSuffix(fname, ".html")
		data, _ := ioutil.ReadFile(f)
		str := string(data)
		res[fname] = str
	}
	return res
}
