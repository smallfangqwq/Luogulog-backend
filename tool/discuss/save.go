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

type DiscussReply struct {
	AuthorID int
	AuthorName string
	Content string
	PostID int //PostID of discuss
	SendTime int64
	ReplyID int //unique ID for each Reply given by Luogu, coming from data-report-id attribute of the Report button
}

type Discuss struct {
	AuthorID int
	AuthorName string
	Content string
	PostID int
	SendTime int64
	Pages int
	Title string
}

func AnalyseDiscussPageForOverview(htmlContent *http.Response, PostID int) (result Discuss, err error) {
	HtmlDocument, err := goquery.NewDocumentFromReader(htmlContent.Body)
	if err != nil {
		return 
	}
	Selection := HtmlDocument.Find(".am-comment-meta").First()
	Vertify, err := Selection.Html() 
	if Vertify == "" || err != nil {
		return 
	}
	oldT := Selection.Text()
	regExp, _ := regexp.Compile(`[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}`)
	sendTime, _ := time.Parse("2006-01-02 15:04", regExp.FindString(oldT))
	result.SendTime = sendTime.Unix()
	result.AuthorName = Selection.Find("a").First().Text()
	AuthorId, _ := Selection.Find("a").First().Attr("href")
	result.AuthorID, _ = strconv.Atoi( strings.Trim(AuthorId, "/user") )

	Selection = HtmlDocument.Find(".am-comment-bd").First()
	ContentHtmlData, _ := Selection.Html()
	result.Content = ContentHtmlData
	result.PostID = PostID;
	err = nil
	return 
}

func AnalyseDiscussPageForReplies(htmlContent *http.Response, PostID int) (result []DiscussReply, err error) {
	HtmlDocument, err := goquery.NewDocumentFromReader(htmlContent.Body)
	HaveData := false
	HtmlDocument.Find(".am-comment-meta").Each(func(i int, Selection *goquery.Selection) {
		if i == 0 {
			HaveData = true
			return 
		}
		AText := Selection.Find("a").First().Text()
		Count := i - 1
		result = append(result, DiscussReply{})
		result[Count].AuthorName = AText
		var regEXP *regexp.Regexp
		regEXP, err = regexp.Compile(`[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}`)
		HereTime, _ := time.Parse("2006-01-02 15:04", regEXP.FindString(Selection.Text()))
		result[Count].SendTime = HereTime.Unix()
		UserCentreURL, _ := Selection.Find("a").First().Attr("href")
		result[Count].AuthorID, _ = strconv.Atoi(strings.Trim(UserCentreURL, "/user"));
		var ReportID string;
		Selection.Find("a").Each(func(i int, Selection *goquery.Selection) {
			if i != 2 {
				return ;
			}
			//Third
			ReportID, _ = Selection.Attr("data-report-id");
		});
		result[Count].ReplyID, _ = strconv.Atoi(ReportID);
		result[Count].PostID = PostID;
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

func GetDiscussReply(Page int, PostID int, htmlConfig declare.ConfigRequest) (result []DiscussReply, err error) {
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
		result, err = AnalyseDiscussPageForReplies(htmlContent, PostID)
		if result != nil {
			return 
		}
	}
	return 
}