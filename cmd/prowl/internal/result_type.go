package internal

type ResultType string

const (
	SponsoredResultType      ResultType = "sponsored"
	OrganicResultType        ResultType = "organic"
	VideoResultType          ResultType = "video"
	ForumResultType          ResultType = "forum"
	ArticleResultType        ResultType = "article"
	PopularProductResultType ResultType = "popular_product"
	RelatedQueryResultType   ResultType = "related_query"
	AlsoAskedResultType      ResultType = "also_asked"
)
