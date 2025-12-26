package prowl

import (
	. "bytelyon-functions/pkg/util"

	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
)

const (
	browserContextScript = `() => {
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
}`
)

func (p *Prowler) NewBrowserContext() (err error) {
	p.BrowserContext, err = p.Browser.NewContext(playwright.BrowserNewContextOptions{
		AcceptDownloads:   Ptr(true),
		ColorScheme:       playwright.ColorSchemeDark,
		ForcedColors:      playwright.ForcedColorsNone,
		HasTouch:          Ptr(false),
		IsMobile:          Ptr(false),
		JavaScriptEnabled: Ptr(true),
		Locale:            Ptr("en-US"),
		Permissions:       []string{"geolocation", "notifications"},
		ReducedMotion:     playwright.ReducedMotionNoPreference,
		//StorageState:      search.FindState(),
		TimezoneId: Ptr("America/New_York"),
		UserAgent:  p.BrowserType.RandomUserAgent(),
	})
	if err == nil {
		p.BrowserContext.SetDefaultTimeout(60_000)
		err = p.BrowserContext.AddInitScript(playwright.Script{Content: Ptr(browserContextScript)})
	}

	log.Err(err).Msg("Prowler - NewBrowserContext")

	return
}
