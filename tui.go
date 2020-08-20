package main

import "github.com/rivo/tview"

// GetList コンテナ一覧を表示するためのアイテムを取得
func GetList(cs Containers, page *tview.Pages, app *tview.Application) *tview.List {
	// コンテナ一覧を表示するためのリスト
	list := tview.NewList()
	for _, container := range cs {
		list.AddItem(container.name, container.status, 'a', func() {
			// コンテナが選択された時に出現するモーダル
			showModal(container, page, app)
		})
	}
	list.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})
	return list
}

// ShowModal モーダルをアプリのページに表示する
func showModal(c Container, p *tview.Pages, app *tview.Application) {
	// モーダル
	modal := tview.NewModal().
		SetText("What do you want to next?").
		AddButtons([]string{"Start", "Stop", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Cancel" {
				p.RemovePage("modal")
				return
			}
			if buttonLabel == "Start" {
				StartContainer(c, p)
			}
			if buttonLabel == "Stop" {
				StopContainer(c, p)
			}
		})
	p.AddPage("modal", modal, true, true)
}
