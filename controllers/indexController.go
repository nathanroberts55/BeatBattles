package controllers

import "github.com/nathanroberts55/beatbattle/common"

type pageProps struct {
	Params
}

func HomeIndex(c *common.Ctx) error {
	return c.Render("index", pageProps{
		DefaultParams,
	})
}

func About(c *common.Ctx) error {
	return c.Render("about", pageProps{
		DefaultParams,
	})
}
