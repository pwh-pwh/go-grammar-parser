package service

import (
	"errors"
	"fmt"
	"grammar_parser/cpp"
	"strings"
)

const ERR_MSG = "解析出错，请检查语法"
const NOT_LL1 = "该文法不是LL(1)文法"

var CP_ERR = errors.New(ERR_MSG)
var NOTLL1_ERR = errors.New(NOT_LL1)

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

func CpGetTable(str string) ([][]string, string, error) {
	result := cpp.GetTable(str)
	if result == ERR_MSG {
		return nil, "", CP_ERR
	}
	orData := removeNilStr(strings.Split(result, "\n"))
	row := len(orData)
	arr := make([][]string, row)
	//log orData[0]
	fmt.Println("ordata[0]", orData[0])
	spHeader := strings.Split(orData[0], "☉")
	//log spHeader
	for _, item := range spHeader {
		fmt.Printf("item:%v  \n", item)
	}

	col := len(spHeader)
	//init arr
	for index := range arr {
		arr[index] = make([]string, col)
	}
	arr[0][0] = "  "
	for i := 1; i < col; i++ {
		arr[0][i] = spHeader[i-1]
	}
	for i := 1; i < row; i++ {
		stt := strings.Split(orData[i], "☉")[0]
		arr[i][0] = stt[:len(stt)-1]
	}
	for i := 1; i < row; i++ {
		split := strings.Split(orData[i], "~~~")
		for j := 1; j < col; j++ {
			if j == 1 {
				lT := strings.Split(split[0], "@")[1]
				if lT == "☉" {
					arr[i][j] = "  "
				} else {
					arr[i][j] = lT
				}
			} else {
				if split[j-1] == "☉" {
					arr[i][j] = "  "
				} else {
					arr[i][j] = split[j-1]
				}
			}
		}
	}
	return arr, result, nil
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

func CpGetTree(str string, token string) (string, error) {
	res := cpp.GetTree(str, token)
	if res == ERR_MSG {
		return "", CP_ERR
	} else if res == NOT_LL1 {
		return "", NOTLL1_ERR
	}
	return res, nil
}
