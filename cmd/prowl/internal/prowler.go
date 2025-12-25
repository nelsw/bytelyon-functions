package internal

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"regexp"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

var (
	pw                         *playwright.Playwright
	blockedRegex               = regexp.MustCompile("(google.com/sorry|captcha|unusual traffic)")
	blockedErr                 = errors.New("blocked")
	googleSearchInputSelectors = []string{
		"input[name='q']",
		"input[title='Search']",
		"input[aria-label='Search']",
		"textarea[title='Search']",
		"textarea[name='q']",
		"textarea[aria-label='Search']",
		"textarea",
	}
)

func InitPW() {
	if err := playwright.Install(&playwright.RunOptions{
		Logger: slog.New(slogzerolog.Option{Level: slog.LevelDebug, Logger: NewLogger()}.NewZerologHandler()),
	}); err != nil {
		log.Panic().Err(err).Msg("pw install")
	} else if pw, err = playwright.Run(); err != nil {
		log.Panic().Err(err).Msg("pw run")
	}
}

type Prowler struct {
	headless bool
	browser  playwright.Browser
	context  playwright.BrowserContext
	page     playwright.Page
}

func NewProwler(search *Search, headless bool) (*Prowler, error) {
	if pw == nil {
		InitPW()
	}
	p := &Prowler{headless: headless}

	if err := p.NewBrowser(); err != nil {
		return nil, err
	} else if err = p.NewContext(search); err != nil {
		return nil, err
	} else if err = p.NewPage(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Prowler) NewBrowser() (err error) {

	p.browser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(p.headless),
		Timeout:  playwright.Float(2 * 60_000),
		Args: []string{
			"--disable-accelerated-2d-canvas",
			"--disable-background-networking",
			"--disable-background-timer-throttling",
			"--disable-backgrounding-occluded-windows",
			"--disable-blink-features=AutomationControlled",
			"--disable-breakpad",
			"--disable-component-extensions-with-background-pages",
			"--disable-dev-shm-usage",
			"--disable-extensions",
			"--disable-features=IsolateOrigins,site-per-process",
			"--disable-features=TranslateUI",
			"--disable-gpu",
			"--disable-ipc-flooding-protection",
			"--disable-renderer-backgrounding",
			"--disable-setuid-sandbox",
			"--disable-site-isolation-trials",
			"--disable-web-security",
			"--enable-features=NetworkService,NetworkServiceInProcess",
			"--force-color-profile=srgb",
			"--hide-scrollbars",
			"--metrics-recording-only",
			"--mute-audio",
			"--no-first-run",
			"--no-sandbox",
			"--no-zygote",
		},
		IgnoreDefaultArgs: []string{
			"--enable-automation",
		},
	})

	log.Err(err).Bool("headless", p.headless).Msg("New Browser")

	return
}

func (p *Prowler) NewContext(search *Search) (err error) {

	p.context, err = p.browser.NewContext(playwright.BrowserNewContextOptions{
		AcceptDownloads:   playwright.Bool(true),
		ColorScheme:       playwright.ColorSchemeDark,
		ForcedColors:      playwright.ForcedColorsNone,
		HasTouch:          playwright.Bool(false),
		IsMobile:          playwright.Bool(false),
		JavaScriptEnabled: playwright.Bool(true),
		Locale:            playwright.String("en-US"),
		Permissions:       []string{"geolocation", "notifications"},
		ReducedMotion:     playwright.ReducedMotionNoPreference,
		StorageState:      search.FindState(),
		TimezoneId:        playwright.String("America/New_York"),
	})

	if err == nil {
		err = p.context.AddInitScript(playwright.Script{Content: playwright.String(`() => {
  // navigator
  Object.defineProperty(navigator, "webdriver", { get: () => false });
  Object.defineProperty(navigator, "plugins", {
	get: () => [1, 2, 3, 4, 5],
  });
  Object.defineProperty(navigator, "languages", {
	get: () => ["en-US", "en", "zh-CN"],
  });

  // window
  window.chrome = {
	runtime: {},
	loadTimes: function () {},
	csi: function () {},
	app: {},
  };

  // WebGL
  if (typeof WebGLRenderingContext !== "undefined") {
	const getParameter = WebGLRenderingContext.prototype.getParameter;
	WebGLRenderingContext.prototype.getParameter = function (
	  parameter: number
	) {
	  // UNMASKED_VENDOR_WEBGL / UNMASKED_RENDERER_WEBGL
	  if (parameter === 37445) {
		return "Intel Inc.";
	  }
	  if (parameter === 37446) {
		return "Intel Iris OpenGL Engine";
	  }
	  return getParameter.call(this, parameter);
	};
  }
}`)})
	}

	log.Err(err).Bool("headless", p.headless).Msg("New Context")
	if err == nil {
		p.context.SetDefaultTimeout(60_000)
	}

	return
}

func (p *Prowler) NewPage() (err error) {
	p.page, err = p.context.NewPage()
	if err == nil {
		err = p.page.AddInitScript(playwright.Script{Content: playwright.String(`() => {
  Object.defineProperty(window.screen, "width", { get: () => 1920 });
  Object.defineProperty(window.screen, "height", { get: () => 1080 });
  Object.defineProperty(window.screen, "colorDepth", { get: () => 24 });
  Object.defineProperty(window.screen, "pixelDepth", { get: () => 24 });
}`)})
	}
	log.Err(err).Bool("headless", p.headless).Msg("New Page")
	return
}

func (p *Prowler) Close() {
	if err := p.context.Close(); err != nil {
		log.Warn().Err(err).Bool("headless", p.headless).Msg("Close Context")
		return
	}
	if err := p.browser.Close(); err != nil {
		log.Warn().Err(err).Bool("headless", p.headless).Msg("Close Browser")
		return
	}
}

func (p *Prowler) GoTo(url string) (playwright.Response, error) {
	log.Trace().Str("url", url).Bool("headless", p.headless).Msg("GoTo")
	res, err := p.page.Goto(url, playwright.PageGotoOptions{
		Timeout:   playwright.Float(5_000),
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})

	l := log.Err(err).Bool("headless", p.headless).Str("url", url)
	if res != nil {
		l = l.Int("status", res.Status())
	}
	l.Msg("Goto")
	return res, err
}

func (p *Prowler) Click(selectors ...string) error {

	var err error
	for _, selector := range selectors {

		locator := p.page.Locator(selector)
		if locator == nil {
			continue
		}

		var n int
		if n, err = locator.Count(); n == 0 {
			continue
		}
		p.Wait(1_000, 1_500)
		err = locator.Click(playwright.LocatorClickOptions{
			Delay: playwright.Float(Between(200, 500.0)),
		})
		break
	}
	log.Err(err).Bool("headless", p.headless).Strs("selectors", selectors).Msg("Click")
	return err
}

func (p *Prowler) Type(s string) error {
	p.Wait(1_000, 1_500)
	err := p.page.Keyboard().Type(s, playwright.KeyboardTypeOptions{
		Delay: playwright.Float(Between(500.0, 1000.0)),
	})
	log.Err(err).Bool("headless", p.headless).Str("s", s).Msg("Type")
	return err
}

func (p *Prowler) Press(s string) error {

	p.Wait(1_000, 1_500)
	if err := p.page.Keyboard().Press(s, playwright.KeyboardPressOptions{
		Delay: playwright.Float(Between(200, 500.0)),
	}); err != nil {
		log.Err(err).Bool("headless", p.headless).Str("s", s).Msg("Press")
		return err
	}

	if err := p.page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State:   playwright.LoadStateNetworkidle,
		Timeout: playwright.Float(60_000),
	}); err != nil {
		log.Err(err).Bool("headless", p.headless).Msg("WaitForLoadState")
		return err
	}

	log.Info().Bool("headless", p.headless).Str("s", s).Msg("Press")
	return nil
}

func (p *Prowler) Title() string {

	title, err := p.context.Pages()[len(p.context.Pages())-1].Title()

	if err != nil {
		log.Warn().Err(err).Bool("headless", p.headless).Msg("Title")
		return ""
	}

	log.Info().Str("title", title).Bool("headless", p.headless).Msg("Title")
	return title
}

func (p *Prowler) Screenshot() []byte {
	b, err := p.context.Pages()[len(p.context.Pages())-1].Screenshot(playwright.PageScreenshotOptions{FullPage: playwright.Bool(true)})
	log.Err(err).
		Bool("headless", p.headless).
		Int("size", len(b)).
		Msg("Screenshot")
	return b
}

func (p *Prowler) Content() string {
	content, err := p.context.Pages()[len(p.context.Pages())-1].Content()
	log.Err(err).
		Bool("headless", p.headless).
		Int("size", len(content)).
		Msg("Content")
	return content
}

func (p *Prowler) Wait(min, max int) {
	time.Sleep(time.Duration(rand.Intn(max-min)+min) * time.Millisecond)
}

func (p *Prowler) Google(search *Search) (err error) {
	defer p.Close()
	isBlocked := func(urls ...string) error {

		for _, u := range urls {
			if blockedRegex.MatchString(u) {
				err = errors.New("blocked: " + u)
				break
			}
		}

		log.Err(err).
			Bool("headless", p.headless).
			Str("q", search.Query).
			Strs("urls", urls).
			Msg("isBlocked")

		return err
	}

	log.Info().
		Bool("headless", p.headless).
		Str("q", search.Query).
		Msg("GoogleSearch")

	var res playwright.Response
	if res, err = p.GoTo("https://www.google.com"); err != nil {
		return
	} else if err = isBlocked(p.page.URL(), res.URL()); err != nil {
		return
	} else if err = p.Click(googleSearchInputSelectors...); err != nil {
		return
	} else if err = p.Type(search.Query); err != nil {
		return
	} else if err = p.Press("Enter"); err != nil {
		return
	} else if err = isBlocked(p.page.URL()); err != nil {
		return
	}

	search.SaveState(p.context.StorageState())
	t := time.Now()
	search.CreatePage(t, p.page.URL(), p.Title(), p.Content(), p.Screenshot())

	if !search.HasTargets() {
		log.Trace().Msg("no targets")
		return
	}

	var locators []playwright.Locator

	if locators, err = p.page.Locator(fmt.Sprintf(`[data-dtld]`), playwright.PageLocatorOptions{}).All(); err != nil {
		return
	}

	var newPage playwright.Page
	var att string
	for _, l := range locators {

		att, err = l.GetAttribute("data-dtld")
		log.Trace().Err(err).Str("att", att).Msg("locator")
		if err != nil {
			continue
		}

		if !search.IsTarget(att) {
			continue
		}

		if newPage, err = p.context.ExpectPage(func() error { return l.Click() }, playwright.BrowserContextExpectPageOptions{
			Timeout: playwright.Float(5_000),
		}); err != nil {
			return
		}

		search.SaveState(p.context.StorageState())
		search.CreatePage(t, newPage.URL(), p.Title(), p.Content(), p.Screenshot())
		_ = newPage.Close()
	}

	return
}
