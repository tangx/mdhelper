package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	inCodeBlock = false
)

func main() {

	name := "readme.md"

	fi, err := os.Open(name)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	buf := bytes.NewBuffer(nil)

	br := bufio.NewReader(fi)
	for {
		buf.WriteString("\n")

		line, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		if isCodeBlock(line) {
			inCodeBlock = !inCodeBlock
		}

		img := mustMatchImage(line, inCodeBlock)
		if img == nil {
			buf.Write(line)
			continue
		}

		// fmt.Println(img.Title, img.Dest)
		new := "https://cdn.exampla.com/asdf/afsd/image.png"
		dest := strings.ReplaceAll(string(line), img.Dest, new)

		// fmt.Println(dest)
		buf.WriteString(dest)

	}

	err = os.WriteFile("new.readme.md", buf.Bytes(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func codeblock() bool {
	return false
}

var (
	imagePattern = `!\[(?P<title>.*)?\]\((?P<dest>.*)\)`
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

var (
	codeblockRegex = regexp.MustCompile("```(.*)?")
)

func isCodeBlock(input []byte) bool {
	return codeblockRegex.Match(input)
}
