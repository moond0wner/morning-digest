package habr

import (
	"fmt"
	"html"
	"strings"
	"time"
)

type RSS struct {
	Channel struct {
		Items []Item `xml:"item"`
	} `xml:"channel"`
}

type Item struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
}

// example
// <rss>
//   <channel>
//     <item>
//       <title>title</title>
//       <link>link</link>
//       <pubDate>date</pubDate>
//
//     </item>

type Article struct {
	Title string
	Link  string
	Date  time.Time
}

func (a *Article) String() string {
	cleanLink := strings.Split(a.Link, "?")[0]
	safeTitle := html.EscapeString(a.Title)
	return fmt.Sprintf(
		"📰 <a href=\"%s\">%s</a>\n   🕒 %s",
		cleanLink,
		safeTitle,
		a.Date.Format("02 Jan, 15:04"),
	)
}
