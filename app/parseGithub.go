package app

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type GithubItem struct {
	link, title, desc, langColor, lang, stars, forks, starsToday string
}

func parseGithub(id int, channel chan ParseResult) {
	fmt.Println("Parse Github...")
	const url = "https://github.com/trending"
	items := make([]GithubItem, 0)

	// parse wepage and collect information
	fetchHtmlPage(url).Find("article").Each(func(_ int, el *goquery.Selection) {
		item := GithubItem{}
		link, _ := el.Find(".lh-condensed").Find("a").Attr("href")
		item.link = "https://github.com" + link
		item.title = PrettyStr(el.Find("h1").Text())
		item.desc = PrettyStr(el.Find("p").Text())

		// find language element and color
		langColorEl := el.Find(".repo-language-color").Nodes
		langColor := ""
		if len(langColorEl) > 0 {
			colorAttr := langColorEl[0].Attr[1].Val
			langColor = colorAttr[len(colorAttr)-7:]
		}
		item.langColor = langColor

		// find info element
		info := PrettyStr(el.Find(".d-inline-block").Text())
		infoSlice := strings.Fields(info)
		if len(langColor) == 0 {
			infoSlice = append([]string{""}, infoSlice...)
		}

		item.lang = infoSlice[0]
		item.stars = infoSlice[1]
		item.forks = infoSlice[2]
		item.starsToday = infoSlice[5]

		items = append(items, item)
	})

	// create html
	itemsHtml := ""
	for _, item := range items {
		langDiv := ""
		if item.lang != "" {
			langDiv = `
			<div class="lang">
				<div class="icon" style="background-color: %s"></div>
				<div class="text">%s</div>
			</div>
			`
			langDiv = fmt.Sprintf(langDiv, item.langColor, item.lang)
		}
		itemsHtml += fmt.Sprintf(githubItemHtml, item.link, item.title, item.desc, langDiv, item.stars, item.forks, item.starsToday)
	}

	channel <- ParseResult{id, fmt.Sprintf(columnHtml, "Github Trending", url, itemsHtml)}
}
