package main

import (
	"strconv"

	"fyne.io/fyne/theme"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

func NewUi(application fyne.App, hosts *hosts) *ui {
	window := application.NewWindow("hosts manager")
	return &ui{
		window: window,
		hosts:  hosts,
	}
}

type ui struct {
	window fyne.Window
	hosts  *hosts
}

func (ui *ui) assembleTitleRow() fyne.CanvasObject {
	tableLine := widget.NewLabel("line")
	tabItemIp := widget.NewLabel("ip")
	tabItemDomain := widget.NewLabel("domain")
	tabItemAction := widget.NewLabel("action")
	return fyne.NewContainerWithLayout(layout.NewGridLayout(4), tableLine, tabItemIp, tabItemDomain, tabItemAction)
}

func (ui *ui) assembleActionButton(lineNo int) fyne.CanvasObject {
	buttonEdit := widget.NewButton("edit", func() {
		ui.loadEdit(lineNo)
	})
	buttonDelete := widget.NewButton("delete", func() {
		dialog.ShowConfirm("sure?", "really to delete?", func(b bool) {
			if b {
				err := ui.hosts.DeleteHost(lineNo)
				if err != nil {
					dialog.ShowError(err, ui.window)
					return
				}
				err = ui.hosts.Sync()
				if err != nil {
					dialog.ShowError(err, ui.window)
					_ = ui.hosts.Load()
				}
			}
			ui.load()
		}, ui.window)
	})
	return fyne.NewContainerWithLayout(layout.NewGridLayout(2), buttonEdit, buttonDelete)
}

func (ui *ui) assembleRow(line int, record *record, buttons fyne.CanvasObject) fyne.CanvasObject {
	if !record.IsDisplay || record.IsComment {
		return nil
	}
	return fyne.NewContainerWithLayout(layout.NewGridLayout(4), widget.NewLabel(strconv.Itoa(line)), widget.NewLabel(record.Ip), widget.NewLabel(record.Domain), buttons)
}

func (ui *ui) assembleCreateButtonRow() fyne.CanvasObject {
	return fyne.NewContainerWithLayout(layout.NewGridLayout(1), widget.NewButton("+", ui.loadAdd))
}

func (ui *ui) load() {
	var containers []fyne.CanvasObject
	containers = append(containers, ui.assembleTitleRow())
	var line = 1
	for _, record := range ui.hosts.records {
		row := ui.assembleRow(line, record, ui.assembleActionButton(record.LineNo))
		if row == nil {
			continue
		}
		containers = append(containers, row)
		line++
	}
	containers = append(containers, ui.assembleCreateButtonRow())
	ui.window.SetContent(
		fyne.NewContainerWithLayout(layout.NewGridLayout(1), containers...),
	)
}

func (ui *ui) loadAdd() {
	var containers []fyne.CanvasObject
	containers = append(containers, ui.assembleTitleRow())
	var line = 1
	for _, record := range ui.hosts.records {
		row := ui.assembleRow(line, record, ui.assembleActionButton(record.LineNo))
		if row == nil {
			continue
		}
		containers = append(containers, row)
		line++
	}
	containers = append(containers, ui.assembleCreateRow())
	containers = append(containers, ui.assembleCreateButtonRow())
	ui.window.SetContent(
		fyne.NewContainerWithLayout(layout.NewGridLayout(1), containers...),
	)
}

func (ui *ui) assembleCreateRow() fyne.CanvasObject {
	var onSubmit = func(ipEntry *widget.Entry, domainEntry *widget.Entry) {
		ipValue := ipEntry.Text
		if ipValue == "" {
			dialog.ShowInformation("Error", "ip cannot be blank", ui.window)
			return
		}
		domainValue := domainEntry.Text
		if domainValue == "" {
			dialog.ShowInformation("Error", "domain cannot be blank", ui.window)
			return
		}
		ui.hosts.AddHost(ipValue, domainValue, false)
		err := ui.hosts.Sync()
		if err != nil {
			dialog.ShowError(err, ui.window)
			_ = ui.hosts.Load()
		}
		ui.load()
	}

	var onCancel = func() {
		ui.load()
	}
	return ui.assembleActiveRow("", "", onSubmit, onCancel)
}

func (ui *ui) assembleActiveRow(defaultIpValue, defaultDomainValue string, onSubmit func(ipEntry *widget.Entry, domainEntry *widget.Entry), onCancel func()) fyne.CanvasObject {
	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("8.8.8.8")
	ipEntry.SetText(defaultIpValue)

	domainEntry := widget.NewEntry()
	domainEntry.SetPlaceHolder("www.google.com")
	domainEntry.SetText(defaultDomainValue)

	form := &widget.Form{
		OnSubmit: func() {
			onSubmit(ipEntry, domainEntry)
		},
		OnCancel: onCancel,
	}

	form.Append("ip", ipEntry)
	form.Append("domain", domainEntry)
	form.CreateRenderer()

	box := widget.NewHBox()
	box.Append(widget.NewButtonWithIcon("ok", theme.ConfirmIcon(), func() {
		form.OnSubmit()
	}))
	box.Append(widget.NewButtonWithIcon("cancel", theme.CancelIcon(), func() {
		form.OnCancel()
	}))
	return fyne.NewContainerWithLayout(layout.NewGridLayout(4), widget.NewLabel(""), ipEntry, domainEntry, box)
}

func (ui *ui) loadEdit(lineNo int) {
	var containers []fyne.CanvasObject
	containers = append(containers, ui.assembleTitleRow())
	var line = 1
	for _, record := range ui.hosts.records {
		if record.LineNo == lineNo { //需要修改的行
			containers = append(containers, ui.assembleEditRow(record))
		} else {
			row := ui.assembleRow(line, record, ui.assembleActionButton(record.LineNo))
			if row == nil {
				continue
			}
			containers = append(containers, row)
		}
		line++
	}
	containers = append(containers, ui.assembleCreateButtonRow())
	ui.window.SetContent(
		fyne.NewContainerWithLayout(layout.NewGridLayout(1), containers...),
	)
}

func (ui *ui) assembleEditRow(record *record) fyne.CanvasObject {
	var onSubmit = func(ipEntry *widget.Entry, domainEntry *widget.Entry) {
		ipValue := ipEntry.Text
		if ipValue == "" {
			dialog.ShowInformation("Error", "ip cannot be blank", ui.window)
			return
		}
		domainValue := domainEntry.Text
		if domainValue == "" {
			dialog.ShowInformation("Error", "domain cannot be blank", ui.window)
			return
		}
		_ = ui.hosts.EditHost(record.LineNo, ipValue, domainValue, false)
		err := ui.hosts.Sync()
		if err != nil {
			dialog.ShowError(err, ui.window)
			_ = ui.hosts.Load()
		}
		ui.load()
	}

	var onCancel = func() {
		ui.load()
	}
	return ui.assembleActiveRow(record.Ip, record.Domain, onSubmit, onCancel)
}

func (ui *ui) Show() {
	ui.load()
	ui.window.ShowAndRun()
}
