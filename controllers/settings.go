package controllers

type MenuItem struct {
	Href string
	Name string
}

type Params struct {
	MenuItems []MenuItem
}

var defaultParams = Params{
	[]MenuItem{
		{"/", "Home"},
		{"/watch", "Watch"},
	},
}
