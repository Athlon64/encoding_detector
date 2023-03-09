package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

const (
	ReadBytes = 1024
)

func main() {
	start := time.Now()
	// 获取文件名
	fileName := getFileName()

	// 打开文件并读取 1024 字节
	buf := openAndRead(fileName)

	// 尝试用各种字符集解码器解码
	pageName := getPageName(buf)

	// 输出成功匹配的字符集
	fmt.Printf("code.page = %s\n", pageName)
	end := time.Now()
	fmt.Println(end.Sub(start))
}

func exitProgram(err error) {
	fmt.Println(err)
	os.Exit(-1)
}

func getFileName() string {
	args := os.Args
	if args == nil || len(args) < 2 {
		exitProgram(errors.New("please specify a file"))
	}

	return args[1]
}

func openAndRead(fileName string) []byte {
	file, err := os.Open(fileName)
	if err != nil {
		exitProgram(err)
	}
	defer file.Close()

	buf := make([]byte, ReadBytes)
	if _, err := file.Read(buf); err != nil {
		exitProgram(err)
	}

	return buf
}

func getPageName(buf []byte) string {
	// 936 = 简体中文，950 = 繁体中文，932 = 日文，949 = 韩文
	pages := []string{"65001", "936", "950", "932", "949"}
	maxScore := 0
	maxIndex := 0
	for index, pageName := range pages {
		if score := tryDecode(buf, pageName); score > maxScore {
			maxScore = score
			maxIndex = index
		}
	}

	return pages[maxIndex]
}

func tryDecode(buf []byte, pageName string) int {
	var t transform.Transformer

	switch pageName {
	case "65001":
		t = unicode.UTF8.NewDecoder()
	case "936":
		t = simplifiedchinese.GBK.NewDecoder()
	case "950":
		t = traditionalchinese.Big5.NewDecoder()
	case "932":
		t = japanese.ShiftJIS.NewDecoder()
	case "949":
		t = korean.EUCKR.NewDecoder()
	default:
		exitProgram(errors.New("Unsupported encoding!"))
	}

	// 计算置信度得分
	result, _, _ := transform.Bytes(t, buf)
	errors := bytes.Count(result, []byte("\uFFFD"))

	return (ReadBytes - errors)
}
