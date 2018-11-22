package lib

import (
	"fmt"
	"net/http"
)

func Head(url string) (code int,err error){
	client := &http.Client{}

	reqest, err := http.NewRequest("HEAD", url, nil)

	if err != nil {
		return code,fmt.Errorf("%v", err)
	}
	reqest.Header.Add("User-Agent", "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Googlebot/2.1; +http://www.google.com/bot.html) Safari/537.36")
	response, err := client.Do(reqest)

	if err != nil {
		return code,fmt.Errorf("%v", err)
	}

	defer response.Body.Close()
	return response.StatusCode,nil
}