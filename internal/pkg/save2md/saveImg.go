package save2md

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func saveImage(imgUrl string) (string, error) {
	slog.Info("正在保存图片，请稍等")
	imgName := strings.Join(strings.Split(imgUrl[strings.Index(imgUrl, "/uploads/")+len("/uploads/"):], "/"), "-")
	imgPath := viper.GetString("dir.img") + imgName

	os.Mkdir(viper.GetString("dir.img"), os.ModePerm)

	res, err := http.Get(imgUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	file, err := os.Create(imgPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	return imgPath, err
}
