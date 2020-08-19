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

// handle メインの画面の描画設定と表示を行う
func handle() {
	containers, err := getContainers()
	if err != nil {
		panic(err)
	}
	setScreen(containers)
}

func setScreen(c []container) {
	app := tview.NewApplication()
	list := tview.NewList()
	for _, container := range c {
		list.AddItem(container.name, container.status, 'a', func() {
			app.Stop()
		})
	}
	list.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})
	if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}
}

// getContainers 外部コマンドを実行し、Dockerコンテナを取得する
func getContainers() ([]container, error) {

	var containers []container

	// 外部コマンドの実行よりコンテナ一覧を取得
	out, err := exec.Command("docker", "ps", "--format", "\"{{.ID}} {{.Names}} {{.Status}}\"").Output()
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
