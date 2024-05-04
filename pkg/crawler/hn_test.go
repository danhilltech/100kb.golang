package crawler

import (
	"html"
	"testing"

	"mvdan.cc/xurls/v2"
)

func TestBasic(t *testing.T) {

	text := `"<a href=\"https:&#x2F;&#x2F;en.wikipedia.org&#x2F;wiki&#x2F;Sanpaku\" rel=\"nofollow noreferrer\">https:&#x2F;&#x2F;en.wikipedia.org&#x2F;wiki&#x2F;Sanpaku</a><p>&gt; <i>According to Chinese&#x2F;Japanese medical [...] when the upper sclera is visible it is said to be an indication of mental imbalance in people such as psychotics, murderers, and anyone rageful. In either condition, it is believed that these people attract accidents and violence.</i><p>It might not be scientific but people with this look certainly do freak me out.  (FWIW, I haven&#x27;t seen any images of Sam with these eyes.)<p><a href=\"https:&#x2F;&#x2F;en.wikipedia.org&#x2F;wiki&#x2F;Marshall_Applewhite#&#x2F;media&#x2F;File:Marshall_Applewhite.jpg\" rel=\"nofollow noreferrer\">https:&#x2F;&#x2F;en.wikipedia.org&#x2F;wiki&#x2F;Marshall_Applewhite#&#x2F;media&#x2F;Fil...</a>"`

	txt := html.UnescapeString(text)

	rxStrict := xurls.Strict()

	urls := rxStrict.FindAllString(txt, 1)

	t.Log(urls)
}
