package model

type DatumType string

const (
	SponsoredDatumType      DatumType = "sponsored"
	OrganicDatumType                  = "organic"
	VideoDatumType                    = "video"
	ForumDatumType                    = "forum"
	ArticleDatumType                  = "article"
	PopularProductDatumType           = "popular_product"
	RelatedQueryDatumType             = "related_query"
	AlsoAskedDatumType                = "also_asked"
)
