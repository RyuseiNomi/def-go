package main

import (
	"os/exec"
	"strings"

	"github.com/rivo/tview"
)

// Container リストに表示するコンテナモデル
type Container struct {
	id     string
	name   string
	status string
}

// Containers 複数のContainerモデルを持つ構造体
type Containers []Container

// GetContainers 外部コマンドを実行し、Dockerコンテナを取得する
func GetContainers() (Containers, error) {

	cs := Containers{}

	// 外部コマンドの実行よりコンテナ一覧を取得
	out, err := exec.Command("docker", "ps", "-a", "--format", "\"{{.ID}} {{.Names}} {{.Status}}\"").Output()
	if err != nil {
		return nil, err
	}

	// 標準出力より情報を取得し、containerモデルに情報を追加する
	slice := strings.Split(string(out), "\n")
	slice = slice[:len(slice)-1] // 改行区切りのため、末尾は空なので削除
	for _, c := range slice {
		container := Container{}
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
		cs = append(cs, container)
	}

	return cs, nil
}

// StartContainer コンテナを起動する外部コマンドを実行する
func StartContainer(c Container, p *tview.Pages) error {
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

// StopContainer コンテナを停止する外部コマンドを実行する
func StopContainer(c Container, p *tview.Pages) error {
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
