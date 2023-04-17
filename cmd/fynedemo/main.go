package main

import (
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//go:embed lan.otf
var TTF []byte

//todo 乱码问题
func main() {
	//os.Setenv("FYNE_FONT", "/Users/coderpwh/GolandProjects/grammar_parser/cmd/fynedemo/lan.otf")
	myApp := app.New()
	myWindow := myApp.NewWindow("TabContainer Widget")

	tabs := container.NewAppTabs(
		container.NewTabItem("简单文法解析", widget.NewLabel("Hello")),
		container.NewTabItem("高级文法解析", widget.NewLabel("World!")),
	)

	tabs.SetTabLocation(container.TabLocationTrailing)
	myWindow.Resize(fyne.NewSize(800, 500))
	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}
