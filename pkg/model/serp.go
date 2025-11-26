package model

type SerpSection string

const (
	GeneralResults   SerpSection = "general-results"
	MoreProducts     SerpSection = "more-products"
	PopularProducts  SerpSection = "popular-products"
	SponsoredProduct SerpSection = "sponsored-product"
	SponsoredResult  SerpSection = "sponsored-result"
	SponsoredResults SerpSection = "sponsored-results"
)

type SERPResult struct {
	Page
	// Section is the type of result container
	Section SerpSection `json:"section"`
	// Index is the array index within a section
	Index int `json:"index"`
	// Entity is the name of the Company or Organization that owns the domain name
	Entity string `json:"entity"`
}

type SERPProduct struct {
	SERPResult
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
	// Description is the summary of the result
	Description string `json:"description"`
	// Details are bulleted information after the description
	Details string `json:"details"`
	// Date is the published date for this page
	Date string `json:"date"`
}

type SERPVideo struct {
	SERPResult
	// Category is the channel, subreddit, path for this video
	Category string `json:"category"`
	// Date is the published date for this video
	Date string `json:"date"`
	// Description is a video summary
	Description string `json:"description"`
}
