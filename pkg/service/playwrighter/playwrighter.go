package playwrighter

import (
	"bytelyon-functions/pkg/db"
	"bytelyon-functions/pkg/service/s3"
	"bytelyon-functions/pkg/util/ptr"
	"bytelyon-functions/pkg/util/random"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
)

type Device string

const (
	Chrome  Device = "Desktop Chrome"
	Edge           = "Desktop Edge"
	Firefox        = "Desktop Firefox"
	Safari         = "Desktop Safari"
)

type State struct{}

func (s *State) Key() string {
	return "playwrighter/storage-state/_1.json"
}

func findStorageState() *playwright.OptionalStorageState {
	var s struct {
		State                            `json:"-"`
		*playwright.OptionalStorageState `json:",inline"`
	}
	if err := db.Find(&s); err != nil {
		log.Panic().Err(err).Send()
	}
	return s.OptionalStorageState
}

func saveStorageState(t *playwright.StorageState) {
	if _, err := db.Save(&struct {
		State                    `json:"-"`
		*playwright.StorageState `json:",inline"`
	}{StorageState: t}); err != nil {
		log.Err(err).Send()
	}
}

var pw *playwright.Playwright

func init() {
	var err error
	if pw, err = playwright.Run(); err != nil {
		panic(err)
	}
}

type Service interface {
	Search(string) error
}

type Client struct {
	playwright.Browser
}

func (c *Client) Search(q string) error {

	ctx, err := c.NewContext(playwright.BrowserNewContextOptions{
		AcceptDownloads:   ptr.True(),
		ColorScheme:       playwright.ColorSchemeDark,
		ForcedColors:      playwright.ForcedColorsNone,
		HasTouch:          ptr.False(),
		IsMobile:          ptr.False(),
		JavaScriptEnabled: ptr.True(),
		Locale:            ptr.Of("en-US"),
		Permissions:       []string{"geolocation", "notifications"},
		ReducedMotion:     playwright.ReducedMotionNoPreference,
		StorageState:      findStorageState(),
		TimezoneId:        ptr.Of("America/New_York"),
	})
	if err != nil {
		return err
	}
	browserScriptContent := `() => {
  // navigator
  Object.defineProperty(navigator, "webdriver", { get: () => false });
  Object.defineProperty(navigator, "plugins", {
	get: () => [1, 2, 3, 4, 5],
  });
  Object.defineProperty(navigator, "languages", {
	get: () => ["en-US", "en", "zh-CN"],
  });

  // window
  // @ts-ignore - chrome
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
}
`
	ctx.AddInitScript(playwright.Script{
		Content: &browserScriptContent,
	})

	var page playwright.Page
	if page, err = ctx.NewPage(); err != nil {
		return err
	}

	pageScriptContent := `() => {
  Object.defineProperty(window.screen, "width", { get: () => 1920 });
  Object.defineProperty(window.screen, "height", { get: () => 1080 });
  Object.defineProperty(window.screen, "colorDepth", { get: () => 24 });
  Object.defineProperty(window.screen, "pixelDepth", { get: () => 24 });
}
`
	page.AddInitScript(playwright.Script{Content: &pageScriptContent})

	var res playwright.Response
	if res, err = page.Goto("https://google.com", playwright.PageGotoOptions{
		Timeout:   ptr.Float64(60_000),
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		return err
	}

	sorryPatterns := []string{
		"google.com/sorry/index",
		"google.com/sorry",
		"recaptcha",
		"captcha",
		"unusual traffic",
	}

	for _, pattern := range sorryPatterns {
		if strings.Contains(res.URL(), pattern) {
			return errors.New("blocked")
		}
	}

	searchInputSelectors := []string{
		"textarea[name='q']",
		"input[name='q']",
		"textarea[title='Search']",
		"input[title='Search']",
		"textarea[aria-label='Search']",
		"input[aria-label='Search']",
		"textarea",
	}

	var searchInput playwright.Locator
	for _, selector := range searchInputSelectors {
		searchInput = page.Locator(selector)
		if searchInput == nil {
			continue
		}
		if n, e := searchInput.Count(); e != nil || n == 0 {
			continue
		}
		break
	}

	if err = searchInput.Click(); err != nil {
		return err
	}

	if err = page.Keyboard().Type(q, playwright.KeyboardTypeOptions{
		Delay: ptr.Float64(random.Between(10, 30)),
	}); err != nil {
		return err
	}

	page.WaitForTimeout(float64(random.Between(100, 300)))

	if err = page.Keyboard().Press("Enter"); err != nil {
		return err
	}

	page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		Timeout: ptr.Float64(60_000),
	})

	page.WaitForTimeout(float64(1_000))

	for _, pattern := range sorryPatterns {
		if strings.Contains(res.URL(), pattern) {
			return errors.New("blocked")
		}
	}

	var fullHtml string
	if fullHtml, err = page.Content(); err != nil {
		return err
	}

	filename := time.Now().Format("20060102150405") + q

	_, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String("tmp/" + filename + ".png"),
		FullPage: ptr.True(),
	})

	fullHtml = regexp.MustCompile(`(?is)<style\b[^>]*>(.*?)</style>`).ReplaceAllString(fullHtml, "")
	fullHtml = regexp.MustCompile(`(?is)<script\b[^>]*>(.*?)</script>`).ReplaceAllString(fullHtml, "")

	if err = s3.New().Put(filename+".html", []byte(fullHtml)); err != nil {
		return err
	}

	var state *playwright.StorageState
	if state, err = ctx.StorageState(); err != nil {
		return err
	}

	saveStorageState(state)

	return c.Close()
}

func New(device ...Device) (Service, error) {

	if len(device) == 0 {
		device = []Device{random.Element(Chrome, Edge, Firefox, Safari)}
	}

	var browserType playwright.BrowserType

	switch device[0] {
	case Chrome:
		browserType = pw.Chromium
	case Firefox:
		browserType = pw.Firefox
	default:
		browserType = pw.WebKit
	}

	browser, err := browserType.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: ptr.True(),
		Timeout:  ptr.Float64(2 * 60_000),
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
	if err != nil {
		return nil, err
	}

	return &Client{browser}, nil
}
