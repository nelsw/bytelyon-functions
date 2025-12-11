package model

import (
	"fmt"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestNewRSS(t *testing.T) {
	q := `Hummer+EV+Reviews`
	urls := []string{
		fmt.Sprintf("https://news.google.com/rss/search?q=%s&hl=en-US&gl=US&ceid=US:en", q),
		fmt.Sprintf("https://www.bing.com/news/search?format=rss&q=%s", q),
		fmt.Sprintf("https://www.bing.com/search?format=rss&q=%s", q),
	}
	for _, u := range urls {
		fmt.Println(u)
	}
	for _, u := range urls {
		rss, err := NewRSS(u)
		assert.NoError(t, err)
		assert.NotNil(t, rss)
		for _, i := range rss.Channel.Items {
			log.Debug().EmbedObject(i).Msg("item")
		}
	}

}
