package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
)

func main() {
	application := NewApp()
	application.Run()
}

func NewApp() *hostsApp {
	return &hostsApp{
		application: app.New(),
	}
}

func (a *hostsApp) Run() {
	application := app.New()
	hosts, err := NewHosts(NewHostsFileSource())
	if err != nil {
		panic("load hosts failed " + err.Error())
	}
	a.ui = NewUi(application, hosts)
	a.ui.Show()
}

type hostsApp struct {
	application fyne.App
	ui          *ui
}
