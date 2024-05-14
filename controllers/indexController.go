package controllers

import "github.com/nathanroberts55/beatbattle/common"

type indexProps struct {
	Params
}

func HomeIndex(c *common.Ctx) error {
	return c.Render("index", indexProps{
		DefaultParams,
	})
}
