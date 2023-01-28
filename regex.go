package main

import "regexp"

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
