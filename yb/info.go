package yb

import (
	"encoding/json"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/homepage-backend/ybtempl"
)

const (
	square_index_url string = "https://www.yiban.cn/square/index"
	my_school_url           = "http://www.yiban.cn/school/getMyGroupAjax?id=5370538"
)

func getMembers(memberTotal *goquery.Selection) (members string) {
	memberTotal.Children().Each(func(i int, span *goquery.Selection) {
		if i == 1 {
			members = span.Children().First().Next().Text()
		}
	})
	return members
}
func getSchoolIntro(c *http.Client) (schoolIntro ybtempl.SchoolIntroTempl) {
	res, _ := c.Get(square_index_url)
	doc, _ := goquery.NewDocumentFromResponse(res)
	temp := doc.Find(".yiban-my-school")
	temp = temp.ChildrenFiltered(".school-intro").Children()
	temp = temp.ChildrenFiltered(".member-total")
	schoolIntro.Members = getMembers(temp)
	schoolIntro.Group = FetchMyGroup(c)
	return schoolIntro
}
func FetchMyGroup(c *http.Client) (group ybtempl.GroupTempl) {
	res, _ := c.Get(my_school_url)
	decoder := json.NewDecoder(res.Body)
	decoder.Decode(&group)
	for i := range group.Data {
		group.Data[i].Url = "http://www.yiban.cn/Newgroup" + group.Data[i].Url
	}
	return group
}
