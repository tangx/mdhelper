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

	// All: ![image](https://www.baidu.com/logo.png "ImageName")
	// Title: image
	// Dest: https://www.baidu.com/logo.png "ImageName"
	// Link: https://www.baidu.com/logo.png
	// Name: ImageName
	linkPattern = `(?P<all>!?\[(?P<title>.*)\]\((?P<dest>(?P<link>http(s)?[\S]+)( ?"(?P<name>.*)")?)\))`
	linkExp     = regexp.MustCompile(linkPattern)
)

type Link struct {
	All   string
	Title string
	Dest  string
	Link  string
	Name  string
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

func mustMatchLink(input []byte, codeblock bool) *Link {
	if codeblock {
		return nil
	}

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
		All:   result["all"],
		Title: result["title"],
		Dest:  result["dest"],
		Link:  result["link"],
		Name:  result["name"],
	}
}
