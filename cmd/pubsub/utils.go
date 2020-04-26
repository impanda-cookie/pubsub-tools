package pubsub

import (
	"os"
	"strings"
)

/**
递归创建文件
*/
func createFile(filePath string) (error, *os.File) {
	pieceArr := strings.Split(filePath, "/")
	length := len(pieceArr)
	folderPath := strings.Join(pieceArr[:length-1], "/")
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return err, nil
	}
	f, err := os.Create(filePath)
	if err != nil {
		return err, nil
	}
	_, err = f.Write([]byte("{}"))
	if err != nil {
		return err, nil
	}
	return nil, f
}

//查看文件是否存在
func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
