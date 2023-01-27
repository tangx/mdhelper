package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func main() {
	err := root.Execute()
	if err != nil {
		panic(err)
	}
}

var root = &cobra.Command{
	Use: "mdhelper",
	Run: func(cmd *cobra.Command, args []string) {
		helper := NewMdHelper(config)
		// file := `./content/posts/2023/01/26/devopscamp-cobra-interactive-survey.md`
		// helper.Replace(file)
		helper.Walk(dir)
	},
}

var (
	config    string
	dir       string
	codeblock = false
)

func init() {
	root.PersistentFlags().StringVarP(&config, "config", "c", "mdhelper.yaml", "mdhelper 配置")
	root.Flags().StringVarP(&dir, "dir", "", "./content", "目标目录")
}

type MdHelper struct {
	RemoteHost      string `yaml:"remoteHost"`
	WorkspacePrefix string `yaml:"workspacePrefix"`
	WorkspaceDir    string `yaml:"workspaceDir"`
}

func NewMdHelper(config string) *MdHelper {
	data, err := os.ReadFile(config)
	if err != nil {
		panic(err)
	}

	helper := &MdHelper{}
	err = yaml.Unmarshal(data, helper)
	if err != nil {
		panic(err)
	}

	return helper
}

func (md *MdHelper) Walk(dirname string) {
	entries, err := os.ReadDir(dirname)
	if err != nil {
		return
	}

	for _, entry := range entries {
		name := filepath.Join(dirname, entry.Name())
		if entry.IsDir() {
			md.Walk(name)
		}

		if !strings.HasSuffix(name, ".md") {
			continue
		}

		fmt.Println("replace =>", name)
		md.Replace(name)
	}

}

func (md *MdHelper) Replace(mdfile string) {
	// mdfile = content/posts/2023/01/26/abcdef.md

	newNameSuffix := "mdhelper.md"
	if strings.HasSuffix(mdfile, newNameSuffix) {
		return
	}

	fi, err := os.Open(mdfile)
	if err != nil {
		return
	}
	defer fi.Close()

	buf := bytes.NewBuffer(nil)
	br := bufio.NewReader(fi)
	for {
		line, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		// 判断是否为代码块 ```
		if isCodeBlock(line) {
			codeblock = !codeblock
		}

		// 提取 图片地址
		img := mustMatchImage(line, codeblock)

		// 如果图片地址为 nil， 则不是图片， 直接写入当前数据
		if img == nil {
			md.out(buf, line)
			continue
		}

		dest := img.Dest
		// 如果是图片地址, 且为完全地址
		if strings.HasPrefix(dest, "https://") || strings.HasPrefix(dest, "http://") {
			// continue
			md.out(buf, line)
			continue
		}

		// 如果是相对地址, 且以 /static 开头
		// img = /static/assets/logo/avatar.png
		if strings.HasPrefix(dest, md.WorkspacePrefix) {
			dest = strings.TrimPrefix(dest, md.WorkspacePrefix)
		}

		if strings.HasPrefix(dest, "./") {
			dir := filepath.Dir(mdfile)
			imgPath := filepath.Join(dir, dest)
			imgAbsPath, err := filepath.Abs(imgPath)
			if err != nil {
				panic(err)
			}

			dest = strings.TrimPrefix(imgAbsPath, md.WorkspaceDir)
		}

		dest = strings.TrimLeft(dest, "/")
		newURL := fmt.Sprintf("%s/%s", md.RemoteHost, dest)
		newLine := strings.ReplaceAll(string(line), img.Dest, newURL)
		md.out(buf, []byte(newLine))

		data := buf.Bytes()
		name := fmt.Sprintf("%s.%s", mdfile, newNameSuffix)
		os.WriteFile(name, data, os.ModePerm)
	}
}

func (md *MdHelper) out(buf *bytes.Buffer, line []byte) {
	buf.WriteString("\n")
	buf.Write(line)
}

var (
	codeblockRegex = regexp.MustCompile("```(.*)?")
)

func isCodeBlock(input []byte) bool {
	return codeblockRegex.Match(input)
}

var (
	imagePattern = `^!\[(?P<title>.*)?\]\((?P<dest>.*)\)`
	imageRegex   = regexp.MustCompile(imagePattern)
)

type Image struct {
	Title string
	Dest  string
}

func mustMatchImage(input []byte, codeblock bool) *Image {
	if codeblock {
		return nil
	}

	if !imageRegex.Match(input) {
		return nil
	}

	// match := imageRegex.FindSubmatch(input)
	match := imageRegex.FindStringSubmatch(string(input))
	result := make(map[string]string)

	for i, name := range imageRegex.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	return &Image{
		Title: result["title"],
		Dest:  result["dest"],
	}
}
