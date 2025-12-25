package model

type Results struct {
	Sponsored []*Result `json:"sponsored"`
	Organic   []*Result `json:"organic"` // WEB_RESULT_INNER
	Videos    []*Result `json:"video"`   // VIDEO_RESULT
	Forums    []*Result `json:"forum"`   // COMMUNITY_MODE_WEB_RESULT (discussions & forums)
	Articles  []*Result `json:"article"` // NEWS_ARTICLE_RESULT (what people are saying)
}
