package ui

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"grammar_parser/service"
	"io/ioutil"
)

type SpGram struct {
	gramEntry   *widget.Entry
	openFileBtn *widget.Button
	spBtn       *widget.Button
	rLFBtn      *widget.Button
	rRRBtn      *widget.Button
	firstBtn    *widget.Button
	followBtn   *widget.Button
	result      *fyne.Container
}

func (sp *SpGram) OpenFile(window fyne.Window) func() {
	return func() {
		fileDia := dialog.NewFileOpen(func(closer fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.NewInformation("未选择文件", "取消选择文件", window).Show()
				return
			}
			if closer != nil {
				defer closer.Close()
			} else {
				return
			}
			bytes, err := ioutil.ReadAll(closer)
			if err != nil {
				dialog.NewError(errors.New("文件读取失败"), window).Show()
				return
			}
			sp.gramEntry.SetText(string(bytes))
		}, window)
		fileDia.Show()
	}
}

func (sp *SpGram) SpOnClick() func() {
	return func() {
		sp.result.RemoveAll()
		gramService, err := service.NewSpGramService(sp.gramEntry.Text)
		if err != nil {
			sp.result.Add(widget.NewLabel(err.Error()))
			return
		}
		result, _ := gramService.GetInvalid()
		//todo list result
		var data = []string{}
		for _, item := range result {
			data = append(data, item.Left+"->"+item.Right)
		}
		list := widget.NewList(
			func() int {
				return len(data)
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("template")
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(data[i])
			})
		sp.result.Add(list)
	}
}

func (sp *SpGram) InitUi(window fyne.Window) fyne.CanvasObject {
	sp.gramEntry = widget.NewMultiLineEntry()
	sp.gramEntry.Wrapping = fyne.TextWrapWord
	sp.gramEntry.SetMinRowsVisible(8)
	sp.openFileBtn = widget.NewButton("打开文件", sp.OpenFile(window))
	sp.spBtn = widget.NewButton("化简文法", sp.SpOnClick())
	sp.rLFBtn = widget.NewButton("消除左公因子", func() {

	})
	sp.rRRBtn = widget.NewButton("消除左递归", func() {

	})
	sp.firstBtn = widget.NewButton("first集合", func() {

	})
	sp.followBtn = widget.NewButton("follow集合", func() {

	})
	sp.result = container.NewMax(widget.NewLabel("结果显示"))
	gramC := container.NewVBox(widget.NewLabel("文法输入框"), sp.gramEntry)
	btnsC := container.NewGridWrap(fyne.NewSize(70, 40), sp.openFileBtn, sp.spBtn, sp.rLFBtn, sp.rRRBtn, sp.firstBtn, sp.followBtn)
	leftC := container.NewVSplit(gramC, btnsC)
	return container.NewHSplit(leftC, sp.result)
}
