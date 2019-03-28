package controller

import (
	"bytes"
	"cronmail/config"
	"cronmail/model"
	"cronmail/template"
	"cronmail/tools"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/goquery"
)

var tip = ""
var list []model.Weather
var sList []model.News
var bList []model.News
var word = model.Word{}
var group = new(sync.WaitGroup)

func GetData(city string) (html string) {
	group.Add(4)
	go getOneData()
	go getWeatherData(city)
	go BHot()
	go SHot()
	group.Wait()
	buffer := new(bytes.Buffer)
	template.WeaList(word, tip, list, sList, bList, buffer)
	list = list[0:0]
	sList = sList[0:0]
	bList = bList[0:0]
	html = buffer.String()
	return
}
func getOneData() {
	url := "http://wufazhuce.com"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	s := doc.Find(".item").First()
	imgURL, _ := s.Find(".fp-one-imagen").Attr("src")
	title := s.Find("a").Text()
	href, _ := s.Find("a").Attr("href")
	word.Title = title
	word.ImgURL = imgURL
	word.Href = href
	group.Done()
}
func getWeatherData(city string) {
	url := "https://tianqi.moji.com/weather/china/" + city + "/" + city
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	tip = doc.Find(".wea_tips").Find("em").Text()
	doc.Find(".days").Each(func(i int, s *goquery.Selection) {
		m := model.Weather{}
		s0 := s.Find("li").Eq(0).Text()
		m.Day = strings.Trim(s0, " \n\t\r")
		s1 := s.Find("li").Eq(1)
		m.ImgUrl, _ = s1.Find("span").Find("img").Attr("src")
		m.Text = strings.Trim(s1.Text(), " \n\t\r")
		s2 := s.Find("li").Eq(2).Text()
		m.Temperature = strings.Trim(s2, " \n\t\r")
		s3 := s.Find("li").Eq(3)
		m.WindDirection = s3.Find("em").Text()
		m.WindLevel = s3.Find("b").Text()
		s4 := s.Find("li").Eq(4).Find("strong")
		m.Pollution = strings.Trim(s4.Text(), " \n\t\r")
		m.PollutionLevel, _ = s4.Attr("class")
		list = append(list, m)
	})
	group.Done()
}
func SHot() {
	url := config.Conf.Get("url.sina").(string)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".td-02").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		title := []byte(s.Find("a").Text())
		href, _ := s.Find("a").Attr("href")
		m := model.News{}
		m.Title = string(title)
		m.Link = "https://s.weibo.com" + href
		m.Type = "sina"
		sList = append(sList, m)
	})
	group.Done()
}
func BHot() {
	// Request the HTML page.
	url := config.Conf.Get("url.baidu").(string)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".keyword").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		title := []byte(s.Find(".list-title").Text())
		href, _ := s.Find(".list-title").Attr("href")
		if t, err := tools.GbkToUtf8(title); err != nil {
			panic(err)
		} else {
			m := model.News{}
			m.Title = string(t)
			m.Link = href
			m.Type = "baidu"
			bList = append(bList, m)
		}
	})
	group.Done()
}
