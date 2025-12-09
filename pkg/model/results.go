package model

type Results struct {
	Sponsored []*Result `json:"sponsored"`
	Organic   []*Result `json:"organic"`  // WEB_RESULT_INNER
	Videos    []*Result `json:"videos"`   // VIDEO_RESULT
	Forums    []*Result `json:"forums"`   // COMMUNITY_MODE_WEB_RESULT (discussions & forums)
	Articles  []*Result `json:"articles"` // NEWS_ARTICLE_RESULT (what people are saying)
}
