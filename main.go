package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-bineanshi/crawler-bind-wrapper/modules"
	"io"
	"net/http"
	"os"
)

const (
	url = "https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=%d&mkt=zh-CN"
)

func FetchBindWrapper(n int) ([]modules.ImageItem, error) {
	response, err := http.Get(fmt.Sprintf(url, n))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			recover()
		}
	}(response.Body)
	resp, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	result := modules.BindWrapperResult{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}
	return result.Images, nil
}

func main() {
	images, err := FetchBindWrapper(1)
	if err != nil {
		panic(err)
	}
	for _, image := range images {
		fmt.Println(image.FullUrl())
		filePath := fmt.Sprintf("pictures/%s.json", image.Enddate[0:6])
		err := ExportToJson(image, filePath, true)
		if err != nil {
			panic(err)
		}
	}
	err = ExportToJson(images[0], "pictures/today.json", false)
	if err != nil {
		panic(err)
	}
}

func ExportToJson(image modules.ImageItem, filePath string, isAppend bool) (err error) {
	items := make(map[string]modules.ImageItem)
	if isExist := Exists(filePath); isExist && isAppend {
		imagesFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		decoder := json.NewDecoder(imagesFile)

		err = decoder.Decode(&items)
		if err != nil {
			return err
		}
	}
	items[image.Enddate] = image
	// Convert golang object back to byte
	byteValue, err := json.MarshalIndent(items, "", "    ")
	if err != nil {
		return err
	}

	// Write back to file
	err = os.WriteFile(filePath, byteValue, 0644)
	return
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
