package crawler

import (
	"bufio"
	"os"
	"strings"

	"github.com/radurobot/go-markdown-crawler/internal/utils"
	"github.com/spf13/cobra"
)

type Options struct {
	SiteList                 []string
	Proxy                    string
	OutputFolder             string
	UserAgent                string
	Cookie                   string
	Headers                  []string
	Burp                     string
	Blacklist                string
	Whitelist                string
	WhitelistDomain          string
	FilterLength             string
	Threads                  int
	Concurrent               int
	Depth                    int
	Delay                    int
	RandomDelay              int
	Timeout                  int
	Base                     bool
	LinkFinder               bool
	Sitemap                  bool
	Robots                   bool
	OtherSource              bool
	IncludeSubs              bool
	IncludeOtherSourceResult bool
	Subs                     bool
	Debug                    bool
	JSONOutput               bool
	Verbose                  bool
	Quiet                    bool
	NoRedirect               bool
	Length                   bool
	Raw                      bool
	UseInMemory              bool
	Command                  *cobra.Command
}

func SetupFlags(commands *cobra.Command) {
	commands.Flags().StringP("site", "s", "", "Site to crawl")
	commands.Flags().StringP("sites", "S", "", "Site list to crawl")
	commands.Flags().StringP("proxy", "p", "", "Proxy (Ex: http://127.0.0.1:8080)")
	commands.Flags().StringP("output", "o", "", "Output folder")
	commands.Flags().StringP("user-agent", "u", "web", "User Agent to use\n\tweb: random web user-agent\n\tmobi: random mobile user-agent\n\tor you can set your special user-agent")
	commands.Flags().StringP("cookie", "", "", "Cookie to use (testA=a; testB=b)")
	commands.Flags().StringArrayP("header", "H", []string{}, "Header to use (Use multiple flag to set multiple headers)")
	commands.Flags().StringP("burp", "", "", "Load headers and cookie from burp raw http request")
	commands.Flags().StringP("blacklist", "", "", "Blacklist URL Regex")
	commands.Flags().StringP("whitelist", "", "", "Whitelist URL Regex")
	commands.Flags().StringP("whitelist-domain", "", "", "Whitelist Domain")
	commands.Flags().StringP("filter-length", "L", "", "Turn on length filter")

	commands.Flags().IntP("threads", "t", 1, "Number of threads (Run sites in parallel)")
	commands.Flags().IntP("concurrent", "c", 5, "The number of the maximum allowed concurrent requests of the matching domains")
	commands.Flags().IntP("depth", "d", 1, "MaxDepth limits the recursion depth of visited URLs. (Set it to 0 for infinite recursion)")
	commands.Flags().IntP("delay", "k", 0, "Delay is the duration to wait before creating a new request to the matching domains (second)")
	commands.Flags().IntP("random-delay", "K", 0, "RandomDelay is the extra randomized duration to wait added to Delay before creating a new request (second)")
	commands.Flags().IntP("timeout", "m", 10, "Request timeout (second)")

	commands.Flags().BoolP("base", "B", false, "Disable all and only use HTML content")
	commands.Flags().BoolP("js", "", true, "Enable linkfinder in javascript file")
	commands.Flags().BoolP("sitemap", "", false, "Try to crawl sitemap.xml")
	commands.Flags().BoolP("robots", "", true, "Try to crawl robots.txt")
	commands.Flags().BoolP("other-source", "a", false, "Find URLs from 3rd party (Archive.org, CommonCrawl.org, VirusTotal.com, AlienVault.com)")
	commands.Flags().BoolP("include-subs", "w", false, "Include subdomains crawled from 3rd party. Default is main domain")
	commands.Flags().BoolP("include-other-source", "r", false, "Also include other-source's urls (still crawl and request)")
	commands.Flags().BoolP("subs", "", false, "Include subdomains")

	commands.Flags().BoolP("debug", "", false, "Turn on debug mode")
	commands.Flags().BoolP("json", "", false, "Enable JSON output")
	commands.Flags().BoolP("verbose", "v", false, "Turn on verbose")
	commands.Flags().BoolP("quiet", "q", false, "Suppress all the output and only show URL")
	commands.Flags().BoolP("no-redirect", "", false, "Disable redirect")
	commands.Flags().BoolP("length", "l", false, "Turn on length")
	commands.Flags().BoolP("raw", "R", false, "Enable raw output")
	commands.Flags().BoolP("in-memory", "M", false, "Store hashes in memory instead of SQLite")
}

func ParseFlags(cmd *cobra.Command) *Options {
	siteInput, _ := cmd.Flags().GetString("site")
	sitesListInput, _ := cmd.Flags().GetString("sites")
	proxy, _ := cmd.Flags().GetString("proxy")
	outputFolder, _ := cmd.Flags().GetString("output")
	userAgent, _ := cmd.Flags().GetString("user-agent")
	cookie, _ := cmd.Flags().GetString("cookie")
	headers, _ := cmd.Flags().GetStringArray("header")
	burp, _ := cmd.Flags().GetString("burp")
	blacklist, _ := cmd.Flags().GetString("blacklist")
	whitelist, _ := cmd.Flags().GetString("whitelist")
	whitelistDomain, _ := cmd.Flags().GetString("whitelist-domain")
	filterLength, _ := cmd.Flags().GetString("filter-length")
	threads, _ := cmd.Flags().GetInt("threads")
	concurrent, _ := cmd.Flags().GetInt("concurrent")
	depth, _ := cmd.Flags().GetInt("depth")
	delay, _ := cmd.Flags().GetInt("delay")
	randomDelay, _ := cmd.Flags().GetInt("random-delay")
	timeout, _ := cmd.Flags().GetInt("timeout")
	base, _ := cmd.Flags().GetBool("base")
	linkfinder, _ := cmd.Flags().GetBool("js")
	sitemap, _ := cmd.Flags().GetBool("sitemap")
	robots, _ := cmd.Flags().GetBool("robots")
	otherSource, _ := cmd.Flags().GetBool("other-source")
	includeSubs, _ := cmd.Flags().GetBool("include-subs")
	includeOtherSourceResult, _ := cmd.Flags().GetBool("include-other-source")
	subs, _ := cmd.Flags().GetBool("subs")
	debug, _ := cmd.Flags().GetBool("debug")
	jsonOutput, _ := cmd.Flags().GetBool("json")
	verbose, _ := cmd.Flags().GetBool("verbose")
	quiet, _ := cmd.Flags().GetBool("quiet")
	noRedirect, _ := cmd.Flags().GetBool("no-redirect")
	length, _ := cmd.Flags().GetBool("length")
	raw, _ := cmd.Flags().GetBool("raw")
	useInMemory, _ := cmd.Flags().GetBool("in-memory")

	var siteList []string
	if siteInput != "" {
		siteList = append(siteList, siteInput)
	}

	if sitesListInput != "" {
		siteList = append(siteList, utils.ReadLinesFromFile(sitesListInput)...)
	}

	if stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0 {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			target := strings.TrimSpace(sc.Text())
			if err := sc.Err(); err == nil && target != "" {
				siteList = append(siteList, target)
			}
		}
	}

	return &Options{
		SiteList:                 siteList,
		Proxy:                    proxy,
		OutputFolder:             outputFolder,
		UserAgent:                userAgent,
		Cookie:                   cookie,
		Headers:                  headers,
		Burp:                     burp,
		Blacklist:                blacklist,
		Whitelist:                whitelist,
		WhitelistDomain:          whitelistDomain,
		FilterLength:             filterLength,
		Threads:                  threads,
		Concurrent:               concurrent,
		Depth:                    depth,
		Delay:                    delay,
		RandomDelay:              randomDelay,
		Timeout:                  timeout,
		Base:                     base,
		LinkFinder:               linkfinder,
		Sitemap:                  sitemap,
		Robots:                   robots,
		OtherSource:              otherSource,
		IncludeSubs:              includeSubs,
		IncludeOtherSourceResult: includeOtherSourceResult,
		Subs:                     subs,
		Debug:                    debug,
		JSONOutput:               jsonOutput,
		Verbose:                  verbose,
		Quiet:                    quiet,
		NoRedirect:               noRedirect,
		Length:                   length,
		Raw:                      raw,
		UseInMemory:              useInMemory,
		Command:                  cmd,
	}
}
