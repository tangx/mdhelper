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
	imageExp     = regexp.MustCompile(imagePattern)

	linkPattern = `\[(?P<title>.*)\]\((?P<dest>http(s)?:\/\/.*)\)`
	linkExp     = regexp.MustCompile(linkPattern)
)

type Link struct {
	Title string
	Dest  string
}

func mustMatchImage(input []byte, codeblock bool) *Link {
	if codeblock {
		return nil
	}

	if !imageExp.Match(input) {
		return nil
	}

	// match := imageRegex.FindSubmatch(input)
	match := imageExp.FindStringSubmatch(string(input))
	result := make(map[string]string)

	for i, name := range imageExp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	return &Link{
		Title: result["title"],
		Dest:  result["dest"],
	}
}

func mustMatchLink(input []byte) *Link {
	if !linkExp.Match(input) {
		return nil
	}

	match := linkExp.FindStringSubmatch(string(input))
	result := make(map[string]string)

	for i, name := range linkExp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	return &Link{
		Title: result["title"],
		Dest:  result["dest"],
	}
}
