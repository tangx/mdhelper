package main

import (
	"github.com/spf13/cobra"
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
		helper := NewMdHelper(yamlConfig)
		// file := `./content/posts/2023/01/26/devopscamp-cobra-interactive-survey.md`
		// helper.Replace(file)

		// helper := NewMdHelperFromToml(tomlConfig)
		helper.Walk(dir)
	},
}

var (
	yamlConfig string
	tomlConfig string
	dir        string
	codeblock  = false
)

func init() {
	root.PersistentFlags().StringVarP(&yamlConfig, "yaml-config", "y", "mdhelper.yaml", "mdhelper yaml 配置")
	// root.PersistentFlags().StringVarP(&tomlConfig, "toml-config", "t", "config.toml", "hugo Toml 配置")
	root.Flags().StringVarP(&dir, "dir", "", "./content/posts", "目标目录")
}
