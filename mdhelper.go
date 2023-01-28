package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

type MdHelper struct {
	RemoteHost      string `yaml:"remoteHost" toml:"image_cdn_host"`
	WorkspacePrefix string `yaml:"workspacePrefix" toml:"workspace_prefix"`
	WorkspaceDir    string `yaml:"workspaceDir,omitempty"`

	filename string
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

	helper.SetDefaults()

	return helper
}

func NewMdHelperFromToml(config string) *MdHelper {
	data, err := os.ReadFile(config)
	if err != nil {
		panic(err)
	}

	strc := &struct {
		Params struct {
			ImageHandler *MdHelper `toml:"image_handler"`
		} `toml:"params"`
	}{
		Params: struct {
			ImageHandler *MdHelper `toml:"image_handler"`
		}{
			ImageHandler: &MdHelper{},
		},
	}

	err = toml.Unmarshal(data, strc)
	if err != nil {
		panic(err)
	}

	helper := strc.Params.ImageHandler
	helper.SetDefaults()

	return helper
}

func (md *MdHelper) SetDefaults() {
	if md.WorkspaceDir == "" {
		dir2, err2 := os.Getwd()
		if err2 != nil {
			panic(err2)
		}
		md.WorkspaceDir = filepath.Join(dir2, "content")

		fmt.Println("探测当前工作目录: ", md.WorkspaceDir)
	}
}

func (md *MdHelper) Copy() *MdHelper {
	return &MdHelper{
		RemoteHost:      md.RemoteHost,
		WorkspacePrefix: md.WorkspacePrefix,
		WorkspaceDir:    md.WorkspaceDir,
	}
}

func (md *MdHelper) withFile(name string) *MdHelper {
	helper := md.Copy()
	helper.filename = name

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

	helper := md.Copy().withFile(mdfile)

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
		if img != nil {
			imgURL := helper.replaceImage(img)
			newLine := strings.ReplaceAll(string(line), img.Dest, imgURL)
			md.output(buf, []byte(newLine))

			continue
		}

		md.output(buf, line)
	}

	data := buf.Bytes()
	name := fmt.Sprintf("%s.%s", mdfile, newNameSuffix)
	os.WriteFile(name, data, os.ModePerm)

}

func (md *MdHelper) replaceImage(img *Image) string {

	dest := img.Dest
	// 如果是图片地址, 且为完全地址
	if strings.HasPrefix(dest, "https://") || strings.HasPrefix(dest, "http://") {
		return dest
	}

	// 如果是相对地址, 且以 /static 开头
	// img = /static/assets/logo/avatar.png
	if strings.HasPrefix(dest, md.WorkspacePrefix) {
		dest = strings.TrimPrefix(dest, md.WorkspacePrefix)
	}

	if strings.HasPrefix(dest, "./") {
		dir := filepath.Dir(md.filename)
		imgPath := filepath.Join(dir, dest)
		imgAbsPath, err := filepath.Abs(imgPath)
		if err != nil {
			panic(err)
		}

		dest = strings.TrimPrefix(imgAbsPath, md.WorkspaceDir)
	}

	dest = strings.TrimLeft(dest, "/")
	newURL := fmt.Sprintf("%s/%s", md.RemoteHost, dest)
	return newURL
	// newLine := strings.ReplaceAll(string(line), img.Dest, newURL)
	// md.output(buf, []byte(newLine))

}

func (md *MdHelper) output(buf *bytes.Buffer, line []byte) {
	buf.WriteString("\n")
	buf.Write(line)
}
