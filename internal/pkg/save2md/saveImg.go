package save2md

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func saveImage(imgUrl string) (string, error) {
	imgName := strings.Join(strings.Split(imgUrl[strings.Index(imgUrl, "/uploads/")+len("/uploads/"):], "/"), "-")
	imgPath := viper.GetString("dir.img") + imgName

	go func() {
		os.Mkdir(viper.GetString("dir.img"), os.ModePerm)
	}()

	res, err := http.Get(imgUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)

	file, err := os.Create(imgPath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)

	_, err = io.Copy(writer, reader)
	return imgPath, err
}
