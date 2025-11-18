package pw

import (
	"bytelyon-functions/pkg/service/s3"
	"bytelyon-functions/pkg/util/ptr"
	"encoding/json"

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

func (ctx *Context) SaveStorage() error {

	state, err := ctx.StorageState()
	if err != nil {
		return err
	}

	b, _ := json.MarshalIndent(state, "", "\t")
	_ = s3.New().Put("pw/storage-state/_.json", b)

	return nil
}

func (ctx *Context) Close() error {
	if err := ctx.SaveStorage(); err != nil {
		return err
	} else if err = ctx.BrowserContext.Close(); err != nil {
		return err
	}
	return nil
}
