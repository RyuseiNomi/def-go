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
				if err := StartContainer(c, p); err != nil {
					panic(err)
				}
				p.RemovePage("modal")

				// 起動完了モーダル
				completeModal := tview.NewModal().
					SetText("Completed to start the container!").
					AddButtons([]string{"Ok"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if buttonLabel == "Ok" {
							p.RemovePage("completeModal")
						}
					})
				p.AddPage("completeModal", completeModal, true, true)
			}
			if buttonLabel == "Stop" {
				if err := StopContainer(c, p); err != nil {
					panic(err)
				}
				p.RemovePage("modal")

				// 消去完了モーダル
				completeModal := tview.NewModal().
					SetText("Completed to stop the container!").
					AddButtons([]string{"Ok"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if buttonLabel == "Ok" {
							p.RemovePage("completeModal")
						}
					})
				p.AddPage("completeModal", completeModal, true, true)
			}
		})
	p.AddPage("modal", modal, true, true)
}
