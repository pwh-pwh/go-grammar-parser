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
	graStr      string
	getLocalS   func(s string) (*service.SpGramService, error)
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

func (sp *SpGram) listDataLt(f func() ([]service.GramTuple, error)) {
	sp.result.RemoveAll()
	result, err := f()
	if err != nil {
		sp.result.Add(widget.NewLabel(err.Error()))
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
	sp.result.Add(list)
}

func (sp *SpGram) GetSer() func(s string) (*service.SpGramService, error) {
	var ser *service.SpGramService
	var err error
	return func(s string) (*service.SpGramService, error) {
		if s != sp.graStr {
			sp.graStr = s
			ser, err = service.NewSpGramService(s)
		}
		return ser, err
	}
}

func (sp *SpGram) SpOnClick() func() {
	return func() {
		sp.result.RemoveAll()
		gramService, err := sp.getLocalS(sp.gramEntry.Text)
		if err != nil {
			sp.result.Add(widget.NewLabel(err.Error()))
			return
		}
		sp.listDataLt(gramService.GetInvalid)
	}
}

func (sp *SpGram) rLFOnClick() func() {
	return func() {
		sp.result.RemoveAll()
		gramService, err := sp.getLocalS(sp.gramEntry.Text)
		if err != nil {
			sp.result.Add(widget.NewLabel(err.Error()))
			return
		}
		sp.listDataLt(gramService.GetRemoveLeftFactor)
	}
}

func (sp *SpGram) rRROnClick() func() {
	return func() {
		sp.result.RemoveAll()
		gramService, err := sp.getLocalS(sp.gramEntry.Text)
		if err != nil {
			sp.result.Add(widget.NewLabel(err.Error()))
			return
		}
		sp.listDataLt(gramService.GetRemoveLeftRecurse)
	}
}

func (sp *SpGram) tableData(f func() ([]service.GramTuple, error)) {
	first, err := f()
	if err != nil {
		sp.result.Add(widget.NewLabel(err.Error()))
		return
	}
	list := widget.NewTable(
		func() (int, int) {
			return len(first), 2
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Col == 0 {
				o.(*widget.Label).SetText(first[i.Row].Left)
			} else {
				o.(*widget.Label).SetText(first[i.Row].Right)
			}
		})
	sp.result.RemoveAll()
	sp.result.Add(list)
}

func (sp *SpGram) firstOnClick() func() {
	return func() {
		sp.result.RemoveAll()
		gramService, err := sp.getLocalS(sp.gramEntry.Text)
		if err != nil {
			sp.result.Add(widget.NewLabel(err.Error()))
			return
		}
		sp.tableData(gramService.GetFirst)
	}
}

func (sp *SpGram) followOnClick() func() {
	return func() {
		sp.result.RemoveAll()
		gramService, err := sp.getLocalS(sp.gramEntry.Text)
		if err != nil {
			sp.result.Add(widget.NewLabel(err.Error()))
			return
		}
		sp.tableData(gramService.GetFollow)
	}
}

func (sp *SpGram) InitUi(window fyne.Window) fyne.CanvasObject {
	sp.getLocalS = sp.GetSer()
	sp.gramEntry = widget.NewMultiLineEntry()
	sp.gramEntry.Wrapping = fyne.TextWrapWord
	sp.gramEntry.SetMinRowsVisible(8)
	sp.openFileBtn = widget.NewButton("打开文件", sp.OpenFile(window))
	sp.spBtn = widget.NewButton("化简文法", sp.SpOnClick())
	sp.rLFBtn = widget.NewButton("消除左公因子", sp.rLFOnClick())
	sp.rRRBtn = widget.NewButton("消除左递归", sp.rRROnClick())
	sp.firstBtn = widget.NewButton("first集合", sp.firstOnClick())
	sp.followBtn = widget.NewButton("follow集合", sp.followOnClick())
	sp.result = container.NewMax(widget.NewLabel("结果显示"))
	gramC := container.NewVBox(widget.NewLabel("文法输入框"), sp.gramEntry)
	btnsC := container.NewGridWrap(fyne.NewSize(70, 40), sp.openFileBtn, sp.spBtn, sp.rLFBtn, sp.rRRBtn, sp.firstBtn, sp.followBtn)
	leftC := container.NewVSplit(gramC, btnsC)
	return container.NewHSplit(leftC, sp.result)
}
