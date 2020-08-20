package main

import "github.com/rivo/tview"

func (wfm *workFlowManager) Handle() {
	// TUIページの初期化
	app := tview.NewApplication()
	page := tview.NewPages()

	cs, err := GetContainers()
	if err != nil {
		panic(err)
	}

	// コンテナ一覧を表示するためのリスト
	list := GetList(cs, page, app)

	page.AddPage("list", list, true, true)
	if err := app.SetRoot(page, true).SetFocus(page).Run(); err != nil {
		panic(err)
	}
}
