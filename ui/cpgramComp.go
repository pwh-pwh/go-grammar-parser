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

type CpGram struct {
	gramEntry   *widget.Entry
	tokenEntry  *widget.Entry
	openFileBtn *widget.Button
	spBtn       *widget.Button
	rrBtn       *widget.Button
	rrAlBtn     *widget.Button
	firstBtn    *widget.Button
	followBtn   *widget.Button
	tableBtn    *widget.Button
	treeBtn     *widget.Button
	result      *fyne.Container
}

func (cp *CpGram) FirstOnClick() func() {
	return func() {
		cp.listDataGT(service.CpGetFirst, cp.gramEntry.Text)
	}
}

func (cp *CpGram) FollowOnClick() func() {
	return func() {
		cp.listDataGT(service.CpGetFollow, cp.gramEntry.Text)
	}
}

func (cp *CpGram) listDataGT(f func(string) ([]service.GramTuple, error), st string) {
	cp.result.RemoveAll()
	result, err := f(st)
	if err != nil {
		cp.result.Add(widget.NewLabel(err.Error()))
		return
	}
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
	cp.result.Add(list)
}

func (cp *CpGram) SpAndRROnClick() func() {
	return func() {
		result, err := service.CpGetRR(cp.gramEntry.Text)
		cp.result.RemoveAll()
		if err != nil {
			cp.result.Add(widget.NewLabel(err.Error()))
			return
		}
		cp.listStrData(result)
	}
}

func (cp *CpGram) RRAlOnClick() func() {
	return func() {
		result, err := service.CpGetRRAl(cp.gramEntry.Text)
		cp.result.RemoveAll()
		if err != nil {
			cp.result.Add(widget.NewLabel(err.Error()))
			return
		}
		cp.listStrData(result)
	}
}

func (cp *CpGram) listStrData(data []string) {
	list := widget.NewList(func() int {
		return len(data)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("temp")
	}, func(id widget.ListItemID, object fyne.CanvasObject) {
		object.(*widget.Label).SetText(data[id])
	})
	cp.result.Add(list)
}

func (cp *CpGram) OpenFile(window fyne.Window) func() {
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
			cp.gramEntry.SetText(string(bytes))
		}, window)
		fileDia.Show()
	}
}

func (cp *CpGram) InitUi(window fyne.Window) fyne.CanvasObject {
	cp.gramEntry = widget.NewMultiLineEntry()
	cp.gramEntry.Wrapping = fyne.TextWrapWord
	cp.gramEntry.SetMinRowsVisible(8)
	cp.tokenEntry = widget.NewEntry()
	cp.tokenEntry.Resize(fyne.Size{
		Width:  200,
		Height: 10,
	})
	cp.openFileBtn = widget.NewButton("打开文件", cp.OpenFile(window))
	cp.spBtn = widget.NewButton("化简文法", cp.SpAndRROnClick())
	cp.rrBtn = widget.NewButton("消除左公因子", cp.SpAndRROnClick())
	cp.rrAlBtn = widget.NewButton("消除左递归", cp.RRAlOnClick())
	cp.firstBtn = widget.NewButton("first集合", cp.FirstOnClick())
	cp.followBtn = widget.NewButton("follow集合", cp.FollowOnClick())
	cp.tableBtn = widget.NewButton("LL(1)文法表", func() {})
	cp.treeBtn = widget.NewButton("语法树", func() {

	})
	cp.result = container.NewMax(widget.NewLabel("结果显示"))
	gramC := container.NewVBox(widget.NewLabel("文法输入框"), cp.gramEntry)
	//加入token框
	tC := container.NewVBox(widget.NewLabel("语句输入框"), cp.tokenEntry)
	gCC := container.NewVBox(gramC, tC)

	btnsC := container.NewGridWrap(fyne.NewSize(70, 40),
		cp.openFileBtn, cp.spBtn, cp.rrBtn, cp.rrAlBtn, cp.firstBtn, cp.followBtn,
		cp.tableBtn, cp.treeBtn)
	leftC := container.NewVSplit(gCC, btnsC)
	return container.NewHSplit(leftC, cp.result)
}
