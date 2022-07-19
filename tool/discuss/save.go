package discuss

import (
	"luogulog/declare"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type StructureDiscussReply struct {
	AuthorID string
	AuthorName string
	Content string
	PostID string
	SendTime int
}

type StructureDiscussTitle struct {
	AuthorID string
	AuthorName string
	Content string
	PostID string
	SendTime int
	Pages int
	Title string
}

func AnalyseDiscussPage(htmlContent *http.Response) (result []StructureDiscussReply, err error) {
	HtmlDocument, err := goquery.NewDocumentFromReader(htmlContent.Body)
	HaveData := false
	HtmlDocument.Find(".am-comment-meta").Each(func(i int, Selection *goquery.Selection) {
		if i == 0 {
			return 
		}
		HaveData = true
		AText := Selection.Find("a").First().Text()
		Count := i - 1
		result = append(result, StructureDiscussReply{})
		result[Count].AuthorName = AText
		var regEXP *regexp.Regexp
		regEXP, err = regexp.Compile(`[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}`)
		HereTime, _ := time.Parse("2006-01-02 15:04", regEXP.FindString(Selection.Text()))
		result[Count].SendTime = int(HereTime.Unix())
		result[Count].AuthorID, _ = Selection.Find("a").First().Attr("href")
		result[Count].AuthorID = strings.Trim(result[Count].AuthorID, "/user")
	})
	if !HaveData {
		result = nil
		err = nil
		return 
	}
	HtmlDocument.Find(".am-comment-bd").Each(func(i int, Selection *goquery.Selection) {
		ContentHtmlData, _ := Selection.Html()
		if i == 0 {
			return
		}
		result[i - 1].Content = ContentHtmlData
	})
	err = nil
	return 
}

func GetDiscussReply(Page int, PostID int, htmlConfig declare.ConfigRequest) (result []StructureDiscussReply, err error) {
	searchURL := "https://www.luogu.com.cn/discuss/" + strconv.Itoa(PostID) + "?page=" + strconv.Itoa(Page)
	result = nil
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return 
	}
	req.Header.Set("User-Agent", htmlConfig.UA)
	req.Header.Set("Host", htmlConfig.Host)
	req.Header.Set("Referer", htmlConfig.Referer)
	CookieNumber := 0
	MaxCookieNumber := len(htmlConfig.Cookies)
	for CookieNumber = 0; CookieNumber < MaxCookieNumber; CookieNumber ++ {
		req.Header.Set("Cookie", htmlConfig.Cookies[CookieNumber])
		client := &http.Client{
			Timeout: time.Second * time.Duration(htmlConfig.TimeOut),
		}
		var htmlContent *http.Response;
		htmlContent, err = client.Do(req)
		if err != nil {
			return
		}
		result, err = AnalyseDiscussPage(htmlContent)
		if result != nil {
			return 
		}
	}
	return 
}