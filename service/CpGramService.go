package service

import (
	"errors"
	"grammar_parser/cpp"
	"strings"
)

const ERR_MSG = "解析出错，请检查语法"

var CP_ERR = errors.New(ERR_MSG)

func CpGetRR(str string) ([]string, error) {
	result := cpp.GetRR(str)
	if result == ERR_MSG {
		return nil, CP_ERR
	}
	return removeNilStr(strings.Split(result, "\n")), nil
}

func removeNilStr(data []string) []string {
	var res []string
	for _, item := range data {
		if len(item) != 0 {
			res = append(res, item)
		}
	}
	return res
}

func CpGetRRAl(str string) ([]string, error) {
	result := cpp.GetRRAl(str)
	if result == ERR_MSG {
		return nil, CP_ERR
	}
	return removeNilStr(strings.Split(result, "\n")), nil
}

func CpGetFirst(str string) ([]GramTuple, error) {
	result := cpp.GetFirst(str)
	if result == ERR_MSG {
		return nil, CP_ERR
	}
	return spFirAndFol(removeNilStr(strings.Split(result, "\n"))), nil
}

func CpGetFollow(str string) ([]GramTuple, error) {
	result := cpp.GetFollow(str)
	if result == ERR_MSG {
		return nil, CP_ERR
	}
	return spFirAndFol(removeNilStr(strings.Split(result, "\n"))), nil
}

func CpGetTable(str string) ([][]string, error) {
	result := cpp.GetTable(str)
	if result == ERR_MSG {
		return nil, CP_ERR
	}
	orData := removeNilStr(strings.Split(result, "\n"))
	row := len(orData)
	arr := make([][]string, row)
	spHeader := strings.Split(orData[0], "===")
	col := len(spHeader)
	//init arr
	for index := range arr {
		arr[index] = make([]string, col)
	}
	arr[0][0] = "  "
	for i := 1; i < col; i++ {
		arr[0][i] = spHeader[i-1]
	}
	//todo
	for i := 1; i < row; i++ {

	}

	return arr, nil
}

func spFirAndFol(data []string) []GramTuple {
	var res []GramTuple
	for _, item := range data {
		split := strings.Split(item, "===")
		var gt GramTuple
		gt.Left = split[0]
		split[1] = split[1][:len(split[1])-1]
		gt.Right = strings.ReplaceAll(split[1], " ", ",")
		res = append(res, gt)
	}
	return res
}
