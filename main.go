package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/rivo/tview"
	cli "github.com/urfave/cli/v2"
)

type container struct {
	id     string
	name   string
	status string
}

const (
	containerIDIndex = 0
	namesIndex       = 1
	statusIndex      = 2
)

func main() {
	app := &cli.App{
		Name:  "def",
		Usage: "manipulate docker containers",
		Action: func(c *cli.Context) error {
			handle()
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func handle() {
	containers, err := getContainers()
	if err != nil {
		panic(err)
	}
	setScreen(containers)
}

// setScreen メインの画面の描画設定と表示を行う
func setScreen(c []container) {
	app := tview.NewApplication()
	page := tview.NewPages()

	// コンテナ一覧を表示するためのリスト
	list := tview.NewList()
	for _, container := range c {
		list.AddItem(container.name, container.status, 'a', func() {
			showModal(container, page, app)
		})
	}
	list.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})

	page.AddPage("list", list, true, true)
	if err := app.SetRoot(page, true).SetFocus(page).Run(); err != nil {
		panic(err)
	}
}

// showModal モーダルをアプリのページに表示する
func showModal(c container, p *tview.Pages, app *tview.Application) {
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
				startContainer(c, p)
			}
			if buttonLabel == "Stop" {
				stopContainer(c, p)
			}
		})
	p.AddPage("modal", modal, true, true)
}

// getContainers 外部コマンドを実行し、Dockerコンテナを取得する
func getContainers() ([]container, error) {

	var containers []container

	// 外部コマンドの実行よりコンテナ一覧を取得
	out, err := exec.Command("docker", "ps", "-a", "--format", "\"{{.ID}} {{.Names}} {{.Status}}\"").Output()
	if err != nil {
		return nil, err
	}

	// 標準出力より情報を取得し、containerモデルに情報を追加する
	slice := strings.Split(string(out), "\n")
	slice = slice[:len(slice)-1] // 改行区切りのため、末尾は空なので削除
	for _, c := range slice {
		container := container{}
		elements := strings.Split(c, " ")
		for i, e := range elements {
			if i == containerIDIndex {
				container.id = e
			}
			if i == namesIndex {
				container.name = e
			}
			if i == statusIndex {
				container.status = e
			}
		}
		containers = append(containers, container)
	}

	return containers, nil
}

func startContainer(c container, p *tview.Pages) error {
	_, err := exec.Command("docker", "start", c.name).Output()
	if err != nil {
		return err
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
	return nil
}

func stopContainer(c container, p *tview.Pages) error {
	_, err := exec.Command("docker", "stop", c.name).Output()
	if err != nil {
		return err
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
	return nil
}
