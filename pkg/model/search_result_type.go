package model

type SearchResultType string

const (
	SponsoredSearchResultType      SearchResultType = "sponsored"
	OrganicSearchResultType                         = "organic"
	VideoSearchResultType                           = "video"
	ForumSearchResultType                           = "forum"
	ArticleSearchResultType                         = "article"
	PopularProductSearchResultType                  = "popular_product"
	RelatedQuerySearchResultType                    = "related_query"
	AlsoAskedSearchResultType                       = "also_asked"
)
