package model

import "github.com/oklog/ulid/v2"

type SerpSection string

const (
	GeneralResults   SerpSection = "general-results"
	MoreProducts     SerpSection = "more-products"
	PopularProducts  SerpSection = "popular-products"
	SponsoredProduct SerpSection = "sponsored-product"
	SponsoredResult  SerpSection = "sponsored-result"
	SponsoredResults SerpSection = "sponsored-results"
)

type SERP struct {
	ID ulid.ULID `json:"id"`
	// Query is the search keywords
	Query string
	// IMG is a page screenshot
	IMG string `json:"img"`
	// HTML is the page content
	HTML string `json:"html"`
	// Results is a slice of SERP results
	Results   []SERPResult `json:"results"`
	FollowAds []string     `json:"follow_ads"`
	SkipAds   []string     `json:"skip_ads"`
}

func NewSerp(q string) *SERP {
	return &SERP{
		ID:    NewUlid(),
		Query: q,
	}
}

type SERPResult struct {
	// URL is the result link
	URL string `json:"url"`
	// Section is the type of result container
	Section SerpSection `json:"section"`
	// Title is the main heading of the result
	Title string `json:"title"`
	// Index is the array index within a section
	Index int `json:"index"`
	// Entity is the name of the Company or Organization that owns the domain name
	Entity string `json:"entity"`
}

type SERPProduct struct {
	SERPResult
	// IMG is the product preview
	IMG string `json:"img"`
	// Tag is the badge on the top left of the product image
	Tag string `json:"tag"`
	// Price1 is the original price
	Price1 string `json:"price_1"`
	// Price2 is the discounted price
	Price2 string `json:"price_2"`
	// Details are product specific information
	Details string `json:"details"`
}

type SERPPage struct {
	SERPResult
	// Logo is the entity logo
	Logo string `json:"logo"`
	// IMG is the page preview
	IMG string `json:"img"`
	// Description is the summary of the result
	Description string `json:"description"`
	// Details are bulleted information after the description
	Details string `json:"details"`
	// Date is the published date for this page
	Date string `json:"date"`
}

type SERPVideo struct {
	SERPResult
	// IMG is the video preview
	IMG string `json:"img"`
	// Category is the channel, subreddit, path for this video
	Category string `json:"category"`
	// Date is the published date for this video
	Date string `json:"date"`
	// Description is a video summary
	Description string `json:"description"`
}
