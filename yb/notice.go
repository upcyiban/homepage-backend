package yb

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/homepage-backend/configuration"
	"github.com/homepage-backend/ybtempl"
)

const (
	notice1_url string = "https://yiban.cn/school/notice/id/5370538/type/1"
	notice2_url string = "https://yiban.cn/school/notice/id/5370538/type/2"
)

func getNotices(c *http.Client, get_url string) []ybtempl.NoticeContent {
	var notices []ybtempl.NoticeContent
	res, _ := c.Get(get_url)
	doc, _ := goquery.NewDocumentFromResponse(res)
	doc.Find(".fl .title").Each(func(i int, a *goquery.Selection) {
		href, _ := a.Attr("href")
		content := a.Text()
		temp := ybtempl.NoticeContent{Text: content, Href: href}
		notices = append(notices, temp)
	})
	return notices
}
func UpdateNotice(notice_1 []ybtempl.NoticeContent, notice_2 []ybtempl.NoticeContent) {
	var ybData ybtempl.YBData
	fileData, err := ioutil.ReadFile(configuration.DataUrl)
	if err != nil {
		log.Fatalf("open data file error , notice not changed: %v\n", err)
		return
	}
	err = json.Unmarshal(fileData, &ybData)
	if err != nil {
		log.Fatalf("decode data.json failed : %v\n", err)
		return
	}
	for i := 0; i < len(ybData.Notices.Notice1.Content); i++ {
		ybData.Notices.Notice1.Content[i] = notice_1[i]
	}
	for i := 0; i < len(ybData.Notices.Notice2.Content); i++ {
		ybData.Notices.Notice2.Content[i] = notice_2[i]
	}
	file, err := os.OpenFile(configuration.DataUrl, os.O_WRONLY|os.O_TRUNC, 0666)
	defer file.Close()

	if err != nil {
		log.Fatalf("write in data.json failed : %v\n", err)
		return
	}
	dataJson, err := json.Marshal(ybData)
	if err != nil {
		log.Fatalf("write in data.json failed : %v\n", err)
		return
	}
	file.Write(dataJson)
}
