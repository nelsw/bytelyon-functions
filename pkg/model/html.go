package model

import (
	"bytelyon-functions/pkg/util/escape"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

type HTML struct {
	SRC string
	Map map[string][]string
}

func (html HTML) Size() int {
	return len(html.SRC)
}

func (html HTML) String() string {
	return html.SRC
}

func MakeHTML(s string) HTML {
	return HTML{
		SRC: escape.Unicode(s),
		Map: make(map[string][]string),
	}
}

func (html HTML) Tags(tag string) []string {

	if val, ok := html.Map[tag]; ok {
		return val
	}

	opened := fmt.Sprintf("<%s", tag)
	closed := fmt.Sprintf("</%s>", tag)
	dead := fmt.Sprintf("<%s></%s>", tag, tag)

	s := html.String()

	var results []string

	for {

		s = strings.ReplaceAll(s, dead, "")

		closedIdx := strings.Index(s, closed)
		if closedIdx == -1 {
			break
		}

		openedIdx := strings.LastIndex(s[:closedIdx], opened)

		if result := s[openedIdx : closedIdx+len(closed)]; result != dead {
			results = append(results, result)
		}

		s = s[:openedIdx+1] + s[closedIdx+len(closed)+1:]
	}

	html.Map[tag] = results
	//
	//log.Debug().
	//	Str("size", pretty.Number(len(html.Map[tag]))).
	//	Str("tag", tag).
	//	Msg("HTML Tags")

	return results
}

func (html HTML) TagsWithAttribute(t string, s string) []string {

	var result []string
	for _, tag := range html.Tags(t) {

		portIdx, stbdIdx := strings.Index(tag, "<"), strings.Index(tag, ">")
		if strings.Contains(tag[portIdx+1:stbdIdx], s+"=") {
			result = append(result, tag)
		}
	}
	//
	//log.Debug().
	//	Str("size", pretty.Number(len(result))).
	//	Str("tag", t).
	//	Str("attr", s).
	//	Msg("HTML Tags w/ Attribute")

	return result
}

func (html HTML) TagsWithClass(t string, s string) []string {

	var result []string
	for _, tag := range html.Tags(t) {

		_, str, ok := strings.Cut(tag, "class=")
		if !ok {
			continue
		}

		idx := strings.Index(str, ">")
		str = strings.TrimSpace(str[:idx])

		if strings.Contains(str, s) {
			result = append(result, tag)
		}
	}
	//
	//log.Debug().
	//	Str("size", pretty.Number(len(result))).
	//	Str("tag", t).
	//	Str("class", s).
	//	Msg("HTML Tags w/ Class")

	return result
}

func (html HTML) AttributeValue(attr string) string {
	attr = attr + "="
	left := strings.Index(html.SRC, attr)
	right := strings.Index(html.SRC[left:], " ")
	return html.SRC[left+len(attr)+1 : left+right-1]
}

func (html HTML) Text() string {
	left := strings.Index(html.SRC, ">")
	right := strings.LastIndex(html.SRC, "<")
	text := html.SRC[left+1 : right]
	if text == "" || text[0] == '<' && text[len(text)-1] == '>' {
		return ""
	}
	return escape.HtmlEntity(text)
}

func (html HTML) TagsWithoutAttributes(t string) []string {
	prefix := fmt.Sprintf("<%s>", t)
	var result []string
	for _, tag := range html.Tags(t) {
		if strings.HasPrefix(tag, prefix) {
			result = append(result, tag)
		}
	}

	//log.Debug().
	//	Str("size", pretty.Number(len(result))).
	//	Str("tag", t).
	//	Msg("HTML Tags w/out Attributes")

	return result
}

func (html HTML) LastTag(t string) string {

	tags := html.Tags(t)
	tag := tags[len(tags)-1]

	//log.Debug().
	//	Str("tag", t).
	//	Str("size", pretty.Number(len(tag))).
	//	Msg("HTML Last Tag")

	return tag
}

func (html HTML) Formatted(skips ...string) {

	var index int
	var breakSkip string
	for _, line := range html.Lines() {

		if breakSkip != "" {

			if strings.HasSuffix(line, breakSkip) {
				breakSkip = ""
			}

			continue
		}

		if len(skips) > 0 {
			for _, s := range skips {
				if strings.HasPrefix(line, fmt.Sprintf("<%s", s)) {
					breakSkip = fmt.Sprintf("</%s>", s)
					break
				}
			}
			if breakSkip != "" {
				continue
			}
		}

		port, stbd := strings.Count(line, "<"), strings.Count(line, ">")
		if strings.HasPrefix(line, "<img") || port+stbd == 4 && strings.Count(line, "/") == 1 {
			fmt.Printf("%s%s\n", strings.Repeat(" ", index), line)
		} else if strings.HasPrefix(line, "</") {
			index -= 2
			fmt.Printf("%s%s\n", strings.Repeat(" ", index), line)
		} else {
			fmt.Printf("%s%s\n", strings.Repeat(" ", index), line)
			index += 2
		}
	}

}

func (html HTML) Lines() []string {

	var result []string
	for _, chunk := range strings.Split(strings.ReplaceAll(html.SRC, "><", ">ø<"), "ø") {
		result = append(result, chunk)
	}
	return result
}

func (html HTML) SerpSponsoredProducts() []any {
	var result []any

	// each div is a product
	divTags := html.TagsWithAttribute("div", "data-aavs")
	//pretty.Print/ln(divTags)
	// we only care about two script tags,
	// one for img, the other text data
	scriptTags := html.Tags("script")
	if len(scriptTags) < 44 {
		panic("not enough script tags fml")
	}

	var imgs []string

	for _, tag := range scriptTags {
		if strings.Contains(tag, "'platop") {

			a := strings.Index(tag, "{")
			z := strings.LastIndex(tag, "}")

			for i, v := range strings.Split(tag[a+1:z], "var ") {
				if i == 1 {
					a = strings.Index(v, "'")
					imgs = append(imgs, v[a+1:])
				}
			}
		}
	}

	scriptTag := scriptTags[44]

	chunks := strings.Split(scriptTag, "})();")
	var lines []string
	var productFragments []HTML
	for _, chunk := range chunks {

		if strings.Contains(chunk, "pla-hovercard-content-ellip") {
			parts := strings.Split(chunk, "','")
			part := parts[1]
			part = strings.TrimSuffix(part, "');")
			part = escape.Hexadecimal(part)
			lines = append(lines, part)
			productFragments = append(productFragments, MakeHTML(part))
		}
	}

	for i := 0; i < len(divTags); i++ {

		divTag := MakeHTML(divTags[i])
		aTag := MakeHTML(productFragments[i].Tags("a")[1])
		spanTags := aTag.Tags("span")
		var spans []string
		for _, span := range spanTags {

			txt := MakeHTML(span).Text()
			if txt != "" {
				spans = append(spans, txt)
			}
		}

		from := aTag.TagsWithAttribute("div", "aria-label")
		var entity string
		if len(from) > 0 {
			entity = MakeHTML(from[0]).Text()
		}
		tit := aTag.TagsWithoutAttributes("div")
		var title string
		if len(tit) > 0 {
			title = MakeHTML(tit[0]).Text()
		}

		var img string
		if i < len(imgs) {
			img = imgs[i]
		}

		log.Trace().
			Int("index", i).
			Str("domain", divTag.AttributeValue("data-dtld")).
			Str("href", aTag.AttributeValue("href")).
			Strs("spans", spans).
			Str("entity", entity).
			Str("title", title).
			Str("img", img).
			Msg("SerpProduct")
	}

	return result
}

func (html HTML) SerpSponsoredResults() []any {
	var results []any
	return results
}

func (html HTML) SerpOrganicResults() map[string]any {

	dataMap := map[string]any{}
	imgMap := map[string]any{}

	scriptTags := html.Tags("script")

	for _, tag := range scriptTags {

		if !strings.Contains(tag, `function(){var m={"`) {
			continue
		}

		chunk := strings.Split(tag, `function(){var m={`)[1]
		chunk = strings.Split(chunk, `};`)[0]

		split := strings.ReplaceAll(chunk, `":[`, `"ø[`)
		splits := strings.Split(split, `ø`)

		var keys []string
		var vals []any

		for i := 0; i < len(splits); i++ {
			s := splits[i]
			if s[0] != '[' {
				keys = append(keys, s)
				continue
			}

			if s[len(s)-1] == ']' {
				vals = append(vals, s)
				continue
			}

			idx := strings.LastIndex(s, `,`)
			vals = append(vals, s[:idx])
			keys = append(keys, s[idx+1:])
		}

		for i, key := range keys {
			dataMap[key] = vals[i]
		}

		break
	}

	for _, tag := range scriptTags {

		if !strings.Contains(tag, "s=") || !strings.Contains(tag, "ii=") {
			continue
		}

		chunks := strings.Split(tag, "=")
		if len(chunks) < 4 {
			continue
		}

		kk := strings.Split(chunks[3], "'")
		if len(kk) < 3 {
			continue
		}

		vv := strings.Split(chunks[2], "'")
		if len(vv) < 3 {
			continue
		}

		imgMap[kk[1]] = vv[1]
	}

	divs := html.TagsWithAttribute("div", "data-rpos")

	metaMap := map[string]any{}

	for i, div := range divs {

		openTag := div[:len(div)-len(`</div>`)]
		openIdx := strings.Index(html.SRC, openTag)

		closeTag := `</div></div></div></div></div></div>`
		if i+1 < len(divs) {
			closeTag = divs[i+1][:len(div)-len(`</div>`)]
		}

		closeIdx := strings.Index(html.SRC[openIdx:], closeTag)
		productContent := html.SRC[openIdx : openIdx+closeIdx]
		if !strings.Contains(productContent, "</h3>") {
			continue
		}

		jsDataAttributeIdxAlpha := strings.Index(productContent, `jsdata=`)
		if jsDataAttributeIdxAlpha == -1 {
			continue
		}
		jsDataAttributeIdxOmega := strings.Index(productContent[jsDataAttributeIdxAlpha+1:], ` `)
		if jsDataAttributeIdxOmega == -1 {
			continue
		}
		jsDataAttributeValue := productContent[jsDataAttributeIdxAlpha+1 : jsDataAttributeIdxAlpha+jsDataAttributeIdxOmega]
		if chunkyBits := strings.Split(jsDataAttributeValue, `;`); len(chunkyBits) > 0 {
			jsDataAttributeValue = chunkyBits[len(chunkyBits)-1]
		}

		imgs := strings.Split(productContent, "<img ")
		var logos []string
		for _, img := range imgs[1:] {
			chunks := strings.Split(img, `id="`)
			if len(chunks) >= 2 {
				id := chunks[1]
				if x := strings.Index(id, ` `); x > -1 {
					id = id[:x-1]
					if _, ok := imgMap[id]; ok {
						ss := strings.Split(chunks[1], `src="`)
						if len(ss) > 1 {
							src := ss[1]
							if x = strings.Index(src, ` `); x > -1 {
								src = src[:x-1]
								imgMap[id] = src
								continue
							}
						}

					}

				}
			}
			chunks = strings.Split(img, `src=`)
			if len(chunks) < 2 {
				continue
			}
			chunk := chunks[1]
			x := strings.Index(chunk, ` `)
			if x < 0 {
				continue
			}
			chunk = chunk[:x]
			if !strings.HasPrefix(chunk, "\"data") {
				continue
			}
			logos = append(logos, chunk[:x])
		}

		m := map[string]any{}
		for _, chunk := range strings.Split(productContent, "</span>") {
			if chunk == "" {
				continue
			}
			chunk = strings.ReplaceAll(chunk, "<em>", "")
			chunk = strings.ReplaceAll(chunk, "</em>", "")
			idx := strings.LastIndex(chunk, ">")
			if s := chunk[idx+1:]; s != "" && s != " · " && !strings.HasPrefix(s, ` › `) {
				s = escape.HtmlEntity(s)
				if _, ok := m[s]; ok {
					continue
				}
				m[s] = s
			}
		}

		metaMap[jsDataAttributeValue] = map[string]any{
			"data":  dataMap,
			"img":   imgMap,
			"logo":  logos,
			"spans": m,
		}
	}

	return metaMap
}

func (html HTML) SerpMoreProducts() []any {
	// todo - see product-viewer-group
	var results []any
	return results
}
