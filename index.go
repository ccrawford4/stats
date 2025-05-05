package main

type UrlEntry struct {
	TotalWords         int
	Title, Description string
}

type (
	UrlMap  map[string]UrlEntry
	WordMap map[string]int
)

type SearchResult struct {
	UrlMap            UrlMap
	TermFrequency     Frequency
	TotalDocsSearched int
	Found             bool
}

type StatResult struct {
	MostFrequent  WordMap
	LeastFrequent WordMap
}

type Index interface {
	containsUrl(url string) bool
	search(word string) *SearchResult
	getTotalWords(url string) int
	insertCrawlResults(c *CrawlResult)
	getStatResults(amount uint) (*StatResult, error)
}
