package task

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

// GBKToUTF8 将GBK编码的字节切片转换为UTF-8编码的字节切片
func GBKToUTF8(gbkBytes []byte) []byte {
	if validUTF8(gbkBytes) {
		return gbkBytes
	}
	gbkDecoder := simplifiedchinese.GBK.NewDecoder()
	b, _, err := transform.Bytes(gbkDecoder, gbkBytes)
	if err != nil {
		log.Printf("error convert gbk to utf8. Error %v", err)
	}
	return b
}

func validUTF8(buf []byte) bool {
	nBytes := 0
	for i := 0; i < len(buf); i++ {
		if nBytes == 0 {
			if (buf[i] & 0x80) != 0 { //与操作之后不为0，说明首位为1
				for (buf[i] & 0x80) != 0 {
					buf[i] <<= 1 //左移一位
					nBytes++     //记录字符共占几个字节
				}

				if nBytes < 2 || nBytes > 6 { //因为UTF8编码单字符最多不超过6个字节
					return false
				}

				nBytes-- //减掉首字节的一个计数
			}
		} else { //处理多字节字符
			if buf[i]&0xc0 != 0x80 { //判断多字节后面的字节是否是10开头
				return false
			}
			nBytes--
		}
	}
	return nBytes == 0
}

// GBKFileToUTF8File 将GBK编码的文件转换为UTF-8编码的文件
func GBKFileToUTF8File(inputFilePath, outputFilePath string) error {
	// 读取GBK编码的文件内容
	gbkBytes, err := ioutil.ReadFile(inputFilePath)
	if err != nil {
		return err
	}

	// 将GBK编码转换为UTF-8编码
	utf8Bytes := GBKToUTF8(gbkBytes)

	// 将UTF-8编码的内容写入新文件
	err = ioutil.WriteFile(outputFilePath, utf8Bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
