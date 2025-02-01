package jigsaw

import "github.com/wangzmgit/jigsaw"

func Jigsaw() {
	j := jigsaw.New()
	j.SetBgDir("./static/jigsaw/bg/")
	j.SetMaskPath("./static/jigsaw/mask.png")
}
