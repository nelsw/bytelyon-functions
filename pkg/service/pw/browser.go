package pw

import (
	"bytelyon-functions/pkg/util/ptr"

	"github.com/playwright-community/playwright-go"
)

type Browser struct {
	playwright.Browser
	*playwright.Proxy
}

func (b *Browser) Search(query string) (html string, img []byte, err error) {

	var c *Context
	if c, err = b.NewContext(); err != nil {
		return
	}
	defer c.Close()

	var p *Page
	if p, err = c.NewPage(); err != nil {
		return
	} else if err = p.SearchGoogle(query); err != nil {
		return
	} else if html, err = p.HTML(); err != nil {
		return
	} else if img, err = p.Screenshot(); err != nil {
		return
	}
	return
}

func (b *Browser) NewContext() (*Context, error) {

	storageState, err := GetStorageState()
	if err != nil {
		return nil, err
	}

	var ctx playwright.BrowserContext
	ctx, err = b.Browser.NewContext(playwright.BrowserNewContextOptions{
		AcceptDownloads:   ptr.True(),
		ColorScheme:       playwright.ColorSchemeDark,
		ForcedColors:      playwright.ForcedColorsNone,
		HasTouch:          ptr.False(),
		IsMobile:          ptr.False(),
		JavaScriptEnabled: ptr.True(),
		Locale:            ptr.Of("en-US"),
		Permissions:       []string{"geolocation", "notifications"},
		ReducedMotion:     playwright.ReducedMotionNoPreference,
		StorageState:      storageState,
		TimezoneId:        ptr.Of("America/New_York"),
	})

	if err != nil {
		return nil, err
	}

	ctx.AddInitScript(playwright.Script{Content: ptr.Of(`() => {
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
}`)})

	return &Context{ctx}, nil
}
