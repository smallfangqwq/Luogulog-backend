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

type DiscussOverview struct {
	AuthorID int
	AuthorName string
	Content string
	PostID int
	SendTime int64
	Pages int
	Title string
}

func AnalyseDiscussPageForOverview(htmlContent *http.Response, PostID int) (result DiscussOverview, err error) {
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
			if i != 3 {
				return ;
			}
			//Fourth
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

func GetDiscussRepliesOnSinglePage(Page int, PostID int, htmlConfig declare.ConfigRequest) (result []DiscussReply, err error) {
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

func GetDiscussOverview(PostID int, htmlConfig declare.ConfigRequest) (result DiscussOverview, err error) {
	searchURL := "https://www.luogu.com.cn/discuss/" + strconv.Itoa(PostID)
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
		result, err = AnalyseDiscussPageForOverview(htmlContent, PostID);
		if err == nil {
			return 
		}
	}
	return 
}

func GetDiscussReplies(BeginPage int,  EndPage int/*does not count*/, PostID int, htmlConfig declare.ConfigRequest) (result []DiscussReply, overview DiscussOverview, err error) {
	/*
		不行，我必须表扬一下刘同学：在过不了编译的情况下还能fix bug? 顶级。
	*/
	result = make([]DiscussReply, 20)[:0];
	overview, err = GetDiscussOverview(PostID, htmlConfig);
	if err != nil {return ;}
	for i := BeginPage; i < EndPage; i++ {
		var ret []DiscussReply;
		ret, err = GetDiscussRepliesOnSinglePage(i, PostID, htmlConfig);
		if ret == nil || err != nil || len(ret) == 0 {
			return;
		}
		result = append(result, ret...);
	}
	return ;
}

func GetAllDiscussRepliesSince(BeginPage int, PostID int, htmlConfig declare.ConfigRequest) ([]DiscussReply, DiscussOverview, error) {
	//I have no idea why this function exists, but I am just required to write one like this
	return GetDiscussReplies(BeginPage, (-1)>>1/*max int*/, PostID, htmlConfig);
}