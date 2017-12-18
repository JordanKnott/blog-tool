package blog

import (
	"github.com/go-ini/ini"
	"errors"
)

var config *ini.File

func LoadConfig() {
	cfg, err := ini.Load("blog.ini")
	if err != nil {
		panic(err)
	}
	config = cfg
}

func GetBlogLocation() string {
	section, err := config.GetSection("general")
	if err != nil {
		panic(err)
	}

	if section.HasKey("blog_location") {
		return section.Key("blog_location").String()
	} else {
		panic(errors.New("Please set the blog_location option in blog.ini!"))
	}
}