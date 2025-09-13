package model

type Shopify struct {
	AccessToken string `json:"access_token"`
	GraphQLUrl  string `json:"graphql_url"`
}

type OpenAI struct {
	Token string `json:"api_key"`
}
