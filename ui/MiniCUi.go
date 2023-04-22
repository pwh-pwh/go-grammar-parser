package ui

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type MiniCUi struct {
	gramEntry *widget.Entry
	subBtn    *widget.Button
	saveFile  *widget.Button
	openFile  *widget.Button
	image     *fyne.Container
}

func (m *MiniCUi) SubOnClick() func() {
	return func() {
		m.image.RemoveAll()
		client := http.DefaultClient
		reader := strings.NewReader(m.gramEntry.Text)
		request, err := http.NewRequest(http.MethodPost, "http://47.106.206.78:8888/minic", reader)
		if err != nil {
			m.image.Add(widget.NewLabel("内部错误"))
			fmt.Println(err)
			return
		}
		response, err := client.Do(request)
		if err != nil {
			m.image.Add(widget.NewLabel("内部错误"))
			fmt.Println(err)
			return
		}
		img := canvas.NewImageFromReader(response.Body, "out.png")
		m.image.Add(img)
	}
}

func (m *MiniCUi) OpenFileOnClick(window fyne.Window) func() {
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
			m.gramEntry.SetText(string(bytes))
		}, window)
		fileDia.Show()
	}
}

func (m *MiniCUi) SaveFileOnClick(window fyne.Window) func() {
	return func() {
		dialog.NewFileSave(func(closer fyne.URIWriteCloser, err error) {
			if closer == nil || err != nil {
				return
			}
			client := http.DefaultClient
			reader := strings.NewReader(m.gramEntry.Text)
			request, err := http.NewRequest(http.MethodPost, "http://47.106.206.78:8888/minic", reader)
			if err != nil {
				log.Println(err)
				return
			}
			resp, err := client.Do(request)
			if err != nil {
				log.Println(err)
				return
			}
			bytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
				return
			}
			has_w, err := closer.Write(bytes)
			fmt.Println("已经写了", has_w)
			if err != nil {
				fmt.Println("写文件出错")
			}
			defer closer.Close()
		}, window).Show()
	}
}

func (m *MiniCUi) InitUi(window fyne.Window) fyne.CanvasObject {
	m.gramEntry = widget.NewMultiLineEntry()
	m.gramEntry.Wrapping = fyne.TextWrapWord
	m.gramEntry.SetMinRowsVisible(30)
	m.subBtn = widget.NewButton("生成语法树", m.SubOnClick())
	m.saveFile = widget.NewButton("保存语法树文件", m.SaveFileOnClick(window))
	m.openFile = widget.NewButton("打开文件", m.OpenFileOnClick(window))
	m.image = container.NewMax(widget.NewLabel("png"))
	gramC := container.NewVBox(widget.NewLabel("请输入Minic语句"), m.gramEntry)
	btns := container.NewHBox(m.subBtn, m.saveFile, m.openFile)
	leftC := container.NewVBox(gramC, btns)
	return container.NewHSplit(leftC, m.image)
}
