package main

import (
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"grammar_parser/ui"

	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var spComp ui.SpGram

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(&ui.MyTheme{})
	myWindow := myApp.NewWindow("Grammar parser")
	tabs := container.NewAppTabs(
		container.NewTabItem("简单文法解析", spComp.InitUi(myWindow)),
		container.NewTabItem("高级文法解析", widget.NewLabel("World!")),
	)
	tabs.SetTabLocation(container.TabLocationTrailing)
	myWindow.Resize(fyne.NewSize(800, 500))
	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}
