package pw

import (
	"bytelyon-functions/pkg/util/ptr"

	"github.com/playwright-community/playwright-go"
)

type Context struct {
	playwright.BrowserContext
}

func (ctx *Context) NewPage() (*Page, error) {
	page, err := ctx.BrowserContext.NewPage()
	if err != nil {
		return nil, err
	}

	page.AddInitScript(playwright.Script{Content: ptr.Of(`() => {
  Object.defineProperty(window.screen, "width", { get: () => 1920 });
  Object.defineProperty(window.screen, "height", { get: () => 1080 });
  Object.defineProperty(window.screen, "colorDepth", { get: () => 24 });
  Object.defineProperty(window.screen, "pixelDepth", { get: () => 24 });
}`)})

	return &Page{page}, nil
}

func (ctx *Context) Close() error {
	if state, err := ctx.StorageState(); err != nil {
		return err
	} else if err = SaveStorageState(state); err != nil {
		return err
	} else if err = ctx.BrowserContext.Close(); err != nil {
		return err
	}
	return nil
}
