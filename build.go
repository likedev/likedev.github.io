package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
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

func GetFileContent(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func WriteToFile(file string, content string) {
	dst_dir := GetDirByPath(file)
	if _, err := os.Stat(dst_dir); os.IsNotExist(err) {
		err := os.MkdirAll(dst_dir, 0755)
		if err != nil {
			fmt.Printf("创建文件夹失败%s\n", dst_dir)
			log.Fatal(err)
		}
	}
	ioutil.WriteFile(file, []byte(content), 0755)
}

// File copies a single file from src to dst
func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
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

func CopyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

/*
*
get filename => file content mappings in one dir
*/
func LoadDirFileContentMap(dir string) map[string]string {
	res := make(map[string]string)
	files := SearchDirRecursively(dir)
	for _, f := range files {
		fname := GetFileName(f)
		fname = strings.TrimSuffix(fname, ".html")
		data, _ := ioutil.ReadFile(f)
		str := string(data)
		res[strings.TrimSpace(fname)] = str
	}
	return res
}

func GetArticleInfo(src string) (map[string]string, string) {
	res := make(map[string]string)
	str := GetFileContent(src)
	//<!-- [article-meta] edit_time=2023-03-10  -->
	r, _ := regexp.Compile("(<!--\\s?\\[article-meta]\\s?(.*?)\\s?-->)")
	//r.MatchString(str)
	arr := r.FindAllString(str, -1)
	if len(arr) == 0 {
		log.Println("article has no meta info", src)
		return res, ""
	}
	for _, ss := range arr {
		meta_str := strings.TrimSuffix(strings.TrimSpace(strings.Split(ss, "]")[1]), "-->")
		meta_arr := strings.Split(meta_str, "=")
		res[strings.TrimSpace(meta_arr[0])] = strings.TrimSpace(meta_arr[1])
	}
	return res, str
}

func ReplaceSubFrag(srcContent string, fragMap map[string]string) string {
	for frag, content := range fragMap {
		tagStr := fmt.Sprintf("<%s/>", frag)
		if strings.Contains(srcContent, tagStr) {
			srcContent = strings.ReplaceAll(srcContent, tagStr, ReplaceSubFrag(content, fragMap))
		}
	}
	return srcContent
}

func ProcesHtml(src string) {
	//load fragMap
	fragMap := LoadDirFileContentMap("./assets/frag")
	//load layout
	layoutMap := LoadDirFileContentMap("./assets/layout")

	//load meta info
	metaMap, blogContent := GetArticleInfo(src)
	layout := metaMap["layout"]
	if layout == "" {
		layout = "blog_layout"
	}
	layout = strings.TrimSpace(layout)
	layoutContent := layoutMap[layout]

	include_regex, _ := regexp.Compile("<!--\\s?\\[include-content]\\s?-->")
	html_content := include_regex.ReplaceAllString(layoutContent, blogContent)
	html_content = ReplaceSubFrag(html_content, fragMap)

	//replace placeholder
	for metaKey, metaV := range metaMap {
		html_content = strings.ReplaceAll(html_content, "{{"+metaKey+"}}", metaV)
	}

	outFileName := "./output/" + src[strings.LastIndex(src, "pages/posts")+len("pages/"):]

	//outDir
	WriteToFile(outFileName, html_content)
}

func main() {
	CopyFile("./pages/index.html", "./output/index.html")

	ProcesHtml("./pages/posts/2023/how-the-blog-is-built.html")

}
