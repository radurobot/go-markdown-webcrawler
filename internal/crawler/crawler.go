package crawler

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/jaeles-project/gospider/core"
	"github.com/sirupsen/logrus"

	"github.com/radurobot/go-markdown-crawler/internal/markdown"
	"github.com/radurobot/go-markdown-crawler/internal/storage"
	"github.com/radurobot/go-markdown-crawler/internal/utils"
)

func Crawl(options *Options, hashStore storage.HashStore) {
	var wg sync.WaitGroup
	inputChan := make(chan string, options.Threads)
	if options.OutputFolder == "" {
		logrus.Warn("No output folder specified. Creating a new folder.")
		options.OutputFolder = fmt.Sprintf("results_%d-%d", os.Getpid(), time.Now().Unix())
	}
	err := utils.CreateDirIfNotExists(options.OutputFolder)
	if err != nil {
		logrus.Fatalf("Failed to create output folder: %s", err)
	}

	for i := 0; i < options.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for rawSite := range inputChan {
				site, err := url.Parse(rawSite)
				if err != nil {
					logrus.Errorf("Failed to parse %s: %s", rawSite, err)
					continue
				}

				crawlSite(site, options, hashStore)
			}
		}()
	}

	for _, site := range options.SiteList {
		inputChan <- site
	}
	close(inputChan)
	wg.Wait()

	logrus.Info("Crawling complete.")
}

func crawlSite(site *url.URL, options *Options, hashStore storage.HashStore) {
	var siteWg sync.WaitGroup

	crawler := core.NewCrawler(site, options.Command)
	crawler.C.OnResponse(func(r *colly.Response) {
		if options.OutputFolder != "" {
			siteName := r.Request.URL.Hostname() + r.Request.URL.Path
			err := markdown.ConvertToMarkdown(siteName, r.Body, options.OutputFolder, hashStore)
			if err != nil {
				log.Printf("Failed to convert %s: %s", siteName, err)
			}
		}
	})

	siteWg.Add(1)
	go func() {
		defer siteWg.Done()
		crawler.Start(options.LinkFinder)
	}()

	if options.Sitemap {
		siteWg.Add(1)
		go core.ParseSiteMap(site, crawler, crawler.C, &siteWg)
	}

	if options.Robots {
		siteWg.Add(1)
		go core.ParseRobots(site, crawler, crawler.C, &siteWg)
	}

	if options.OtherSource {
		siteWg.Add(1)
		go func() {
			defer siteWg.Done()
			urls := core.OtherSources(site.Hostname(), options.IncludeSubs)
			for _, url := range urls {
				url = strings.TrimSpace(url)
				if len(url) == 0 {
					continue
				}

				outputFormat := fmt.Sprintf("[other-sources] - %s", url)
				if options.IncludeOtherSourceResult {
					fmt.Println(outputFormat)
					if crawler.Output != nil {
						crawler.Output.WriteToFile(outputFormat)
					}
				}

				_ = crawler.C.Visit(url)
			}
		}()
	}

	siteWg.Wait()
	crawler.C.Wait()
	crawler.LinkFinderCollector.Wait()
}
