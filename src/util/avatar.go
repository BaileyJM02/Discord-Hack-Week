package util

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetAvatar() string {
		resp, err := http.Get("https://github.com/finione/Logo/raw/master/DFA1BD48-23D6-4AB4-A5D1-AC29D3651579.png")
		if err != nil {
			fmt.Println("Error retrieving the file, ", err)
			return ""
		}

		defer func() {
			_ = resp.Body.Close()
		}()

		img, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading the response, ", err)
			return ""
		}

		contentType := http.DetectContentType(img)
		base64img := base64.StdEncoding.EncodeToString(img)

		avatar := fmt.Sprintf("data:%s;base64,%s", contentType, base64img)
		return avatar
}