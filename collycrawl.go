package main

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

// Jsonyamllink is the output structure for json/yaml links
type Jsonyamllink struct {
	LinkTitle string
	Referrer  string
	URL       string
	Depth     int
	Rank      int
	Baseurl   string
}

// Jsonyamllinks is the output array
type Jsonyamllinks []Jsonyamllink

//byRankDepth implements the sort interface
type byRankDepth []Jsonyamllink

func (links byRankDepth) Len() int { return len(links) }
func (links byRankDepth) Swap(i, j int) {
	links[i], links[j] = links[j], links[i]
}
func (links byRankDepth) Less(i, j int) bool {
	if links[i].Rank < links[j].Rank {
		return true
	}
	if links[i].Rank > links[j].Rank {
		return false
	}
	return links[i].Depth < links[j].Depth
}

// crawlurls crawls the urls in the urllist and output the identified json and yaml links

func crawlurls(urllist []string, cfg *Config) (Jsonyamllinks, error) {

	// Iterate over search results and print the result name and URL.

	maxdepth := cfg.Colly.Maxdepth
	disallowedurlsregex := cfg.Colly.Disallowedurlsregex
	parllels := cfg.Colly.Parallelism
	exclusionsregex := cfg.Colly.Exclusionsregex
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only referral domains
		//colly.AllowedDomains(""),
		colly.MaxDepth(maxdepth),
		colly.Async(),
		colly.DisallowedURLFilters(
			regexp.MustCompile(disallowedurlsregex),
		),
	)
	c.Limit(&colly.LimitRule{Parallelism: parllels})
	/*
		can be used to implement the proxySwitchers
		proxySwitcher, err := proxy.RoundRobinProxySwitcher("socks5://49.12.4.194:58302", "socks5://103.29.156.142:1080")
		if err != nil {
			log.Fatal(err)
		}
		c.SetProxyFunc(proxySwitcher)*/

	extensions.RandomUserAgent(c)

	jsonyamllinks := make([]Jsonyamllink, 0, 200)
	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if (strings.HasSuffix(link, ".json") || strings.HasSuffix(link, ".yaml")) == true {
			linkTitle := e.Text
			linkRef := e.Request.URL.String()
			linkurl := e.Request.AbsoluteURL(link)
			depth := e.Request.Depth
			//rank := 1
			rank, _ := strconv.Atoi(e.Request.Ctx.Get("rank"))

			baseurl := e.Request.Ctx.Get("referer")
			jsonyamllink := Jsonyamllink{
				LinkTitle: linkTitle,
				Referrer:  linkRef,
				URL:       linkurl,
				Depth:     depth,
				Rank:      rank,
				Baseurl:   baseurl,
			}
			jsonyamllinks = append(jsonyamllinks, jsonyamllink)

		}
		// Print link
		//fmt.Printf("Link found: %q -> %s\n%s\n", e.Text, link, e.Request.URL.Path)
		// Visit link found on page
		// Only those links are visited which are relative
		re := regexp.MustCompile(exclusionsregex)
		if !(re.MatchString(link)) {
			if e.Request.Depth < maxdepth {
				if link != "" {
					if strings.HasPrefix(link, e.Request.Ctx.Get("basepath")) {
						err := e.Request.Visit(e.Request.AbsoluteURL(link))
						if err != nil {
							fmt.Println(e.Request.AbsoluteURL(link), err.Error())
						}

					}
				}
			}
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", "\nError:", err)
		if err.Error() == "Too Many Requests" {
			//r.Request.Retry()
		}
	})
	//asynchrosonsly crawls the list of urls in the urllist

	for i, s := range urllist {
		u, _ := url.Parse(s)
		ctx := colly.NewContext()
		ctx.Put("referer", u.String())
		ctx.Put("rank", strconv.Itoa(i))
		ctx.Put("basepath", u.Path)
		fmt.Println("Visiting", u.String())
		// Start scraping
		c.Request("GET", u.String(), nil, ctx, nil)

	}
	// wait for the crawling to complete
	c.Wait()
	// sort results JSON data when the scraping job has finished
	sort.Sort(byRankDepth(jsonyamllinks))
	return jsonyamllinks, nil
}
