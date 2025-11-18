package pw

import (
	"bytelyon-functions/pkg/util/ptr"
	"bytelyon-functions/pkg/util/random"
	"errors"
	"regexp"

	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
)

var (
	googleSearchInputSelectors = []string{
		"textarea[name='q']",
		"input[name='q']",
		"textarea[title='Search']",
		"input[title='Search']",
		"textarea[aria-label='Search']",
		"input[aria-label='Search']",
		"textarea",
	}
	sorryRegex = regexp.MustCompile("(google.com/sorry|captcha|unusual traffic)")
)

type Page struct {
	playwright.Page
}

func (p *Page) SearchGoogle(query string) error {
	if res, err := p.GoTo("https://google.com"); err != nil {
		return err
	} else if err = p.IsBlocked(res.URL()); err != nil {
		return err
	} else if err = p.Click(googleSearchInputSelectors...); err != nil {
		return err
	} else if err = p.Type(query); err != nil {
		return err
	} else if err = p.Press("Enter"); err != nil {
		return err
	} else if err = p.IsBlocked(); err != nil {
		return err
	}

	return nil
}

func (p *Page) GoTo(url string) (playwright.Response, error) {
	res, err := p.Goto(url, playwright.PageGotoOptions{
		Timeout:   ptr.Float64(60_000),
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *Page) Click(selectors ...string) error {

	var err error
	for _, selector := range selectors {

		locator := p.Locator(selector)
		if locator == nil {
			continue
		} else if n, e := locator.Count(); e != nil || n == 0 {
			continue
		}
		err = locator.Click()
		break
	}
	return err
}

func (p *Page) Type(s string) error {
	err := p.Keyboard().Type(s, playwright.KeyboardTypeOptions{
		Delay: ptr.Float64(random.Between(10, 30)),
	})
	return err
}

func (p *Page) Press(s string) error {

	p.WaitForTimeout(float64(random.Between(100, 300)))
	if err := p.Keyboard().Press(s); err != nil {
		return err
	}

	if err := p.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State:   playwright.LoadStateNetworkidle,
		Timeout: ptr.Float64(60_000),
	}); err != nil {
		return err
	}

	p.WaitForTimeout(float64(1_000))
	return nil
}

func (p *Page) IsBlocked(url ...string) error {
	u := p.URL()
	if len(url) > 0 {
		u = url[0]
	}
	if sorryRegex.MatchString(u) {
		return errors.New("blocked")
	}
	return nil
}

func (p *Page) Screenshot(path ...string) ([]byte, error) {

	opts := playwright.PageScreenshotOptions{FullPage: ptr.True()}
	if len(path) > 0 {
		opts.Path = ptr.Of(path[0] + ".png")
	}

	b, err := p.Page.Screenshot(opts)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (p *Page) HTML() (string, error) {

	html, err := p.Content()
	if err != nil {
		return "", err
	}

	// remove ui elements
	pruned := regexp.MustCompile(`(?is)<style\b[^>]*>(.*?)</style>`).ReplaceAllString(html, "")
	pruned = regexp.MustCompile(`(?is)<svg\b[^>]*>(.*?)</svg>`).ReplaceAllString(pruned, "")
	pruned = regexp.MustCompile(`(?is)<img\b[^>]*>`).ReplaceAllString(pruned, "")
	//remove form elements
	pruned = regexp.MustCompile(`(?is)<textarea\b[^>]*>(.*?)</textarea>`).ReplaceAllString(pruned, "")
	pruned = regexp.MustCompile(`(?is)<input\b[^>]*>(.*?)</input>`).ReplaceAllString(pruned, "")
	pruned = regexp.MustCompile(`(?is)<button\b[^>]*>`).ReplaceAllString(pruned, "")
	// remove script elements
	pruned = regexp.MustCompile(`(?is)<script\b[^>]*>(.*?)</script>`).ReplaceAllString(pruned, "")
	pruned = regexp.MustCompile(`(?is)<noscript\b[^>]*>(.*?)</noscript>`).ReplaceAllString(pruned, "")
	// remove link elements
	pruned = regexp.MustCompile(`(?is)<link\b[^>]*>`).ReplaceAllString(pruned, "")
	// remove empty elements
	pruned = regexp.MustCompile(`(?is)<span\b[^>]*></span>`).ReplaceAllString(pruned, "")
	pruned = regexp.MustCompile(`(?is)<div\b[^>]*></div>`).ReplaceAllString(pruned, "")

	log.Trace().
		Int("full", len(html)).
		Int("pruned", len(pruned)).
		Msg("HTML")

	return pruned, nil
}

func (p *Page) Close() error {
	if err := p.Page.Close(); err != nil {
		return err
	}
	return nil
}
