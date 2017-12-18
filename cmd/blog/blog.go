package main

import (
	"os"
	"github.com/jordanknott/blog"
	"os/exec"
	"io/ioutil"
	"strings"
	"path/filepath"
	"encoding/json"
	"sort"
	"strconv"
)

func RunEditor(filepath string) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		println("Please set the EDITOR environment variable")
		os.Exit(-1)
	}
	cmd := exec.Command("vim", filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func newDraft(args []string) {
	blog.LoadConfig()
	blogDirectory := blog.GetBlogLocation() + "/drafts"

	commandArgs := args[2:]
	var title string
	flags := make(map[string]string)
	if len(commandArgs) == 1 {
		title = commandArgs[0]
	} else if len(commandArgs) > 1 {
		title = commandArgs[0]
		mapKey := ""
		for _, flag := range commandArgs[1:] {
			if mapKey == "" {
				mapKey = flag
			} else {
				flags[mapKey] = flag
				mapKey = ""
			}
		}
	} else {

	}
	L := blog.InitLuaRuntime("format.lua")
	fileContents := blog.FormatFile(L, title, flags)
	filename := blog.TitleToFileName(title)
	ioutil.WriteFile(blogDirectory + "/" + filename, []byte(fileContents), 0644)
	RunEditor(blogDirectory + "/" + filename)
	os.Exit(0)
}

func listDrafts() {
	blog.LoadConfig()
	blogDirectory := blog.GetBlogLocation() + "/drafts"
	files, err := ioutil.ReadDir(blogDirectory)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		name := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		println(name)
		dat, err := ioutil.ReadFile(blogDirectory + "/" + f.Name())
		if err != nil {
			panic(err)
		}
		isHeader := false
		yamlHeader := make(map[string]string)
		for _, l := range strings.Split(string(dat), "\n") {
			if !isHeader {
				if l == "---" {
					isHeader = true
				}
			} else {
				if l == "---" {
					break
				}
				if strings.TrimSpace(l) != "" {
					parts := strings.Split(l, ":")
					if len(parts) == 2 {
						yamlHeader[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
					}
				}
			}
		}
		for key, value := range yamlHeader {
			println(key + ": " + value)
		}
	}
}
func editDraft() {}
func deleteDraft() {}
func publishDraft() {}

func newPost() {}
func listPosts() {}
func editPost() {}
func deletePost() {}
func unpublishPost() {}

type Post struct {
	Name     string
	Filename string
	Date     string
	Category string
	Tags     []string
}

// Load a manifest of blog posts and IDs from the given filename
func LoadManifest(filename string) (data map[string]Post) {
	draftsJson, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(draftsJson), &data)
	if err != nil {
		panic(err)
	}
	return data
}

// Writes a manifest to a file
func WriteManifest(filename string, data *map[string]Post) {
	dataJson, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename, dataJson, 0644)
	if err != nil {
		panic(err)
	}
}

// Generates a manifest from scratch based on files that exist in the given filepath
func RegenerateManifest(filepath string, filename string) {

}

// Gets the next key to use for a blog post in a manifest
// It is the largest current key + 1
func GetNextKey(manifest *map[string]Post) int {
	keys := make([]string, 0, len(*manifest))
	for k := range *manifest {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	i, err := strconv.Atoi(keys[0])
	if err != nil {
		panic(err)
	}
	return i + 1
}

// Add a post to a manifest
func AddPostToManifest(manifest *map[string]Post) {
	nextKey := GetNextKey(manifest)
	(*manifest)[strconv.Itoa(nextKey)] = Post{Name: "Hello World", Date: "Date", Tags: []string{}, Category: "category"}
}

/*

A command line tool for managing Markdown posts for a Jekyll blog. Easily create, edit, delete, and (un)publish drafts
and posts.

Flags for the new command is parsed by a Lua script to allow for custom user defined YAML headings that can be
set through a flag.

COMMANDS

draft
- new
- edit
- publish
- delete
- list

posts
- new
- edit
- unpublish
- delete
- list

deploy

 */
func main() {
	args := os.Args[1:]

	manifest := LoadManifest("drafts/.drafts")
	AddPostToManifest(&manifest)
	WriteManifest("drafts/.drafts", &manifest)
	if len(args) >= 1 {
		switch(args[0]) {
		case "draft":
			println("Draft command")
			switch(args[1]) {
			case "new":
				newDraft(args)
			case "edit":
			case "publish":
			case "delete":
			case "list":
				if len(args) == 3 {
					if args[3] == "--regenerate" {
						// Regenerate manifest
					}
				}
				listDrafts()
				break
			default:
				println("Unknown subcommand, options are new, edit, publish, delete, and list.")
				os.Exit(-1)
			}
		case "post":
			println("Post command")
			break
		case "deploy":
			println("Deploy")
		}

	} else {
		println("Usage: ")
	}

}

