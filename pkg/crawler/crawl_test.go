package crawler

import (
	"strings"
	"testing"
)

func TestExtractRSS(t *testing.T) {
	test := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charSet="utf-8"/>
		<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1"/>
		<title>65 Words: Write daily in your target language</title>
		<meta name="robots" content="index,follow"/>
		<meta name="description" content="Improve your language skills by writing 65+ words daily in your chosen language. Enhance fluency, vocabulary, and proficiency through consistent practice. Join ..."/>
		<meta property="og:title" content="65 Words: Write daily in your target language"/>
		<meta property="og:description" content="Improve your language skills by writing 65+ words daily in your chosen language. Enhance fluency, vocabulary, and proficiency through consistent practice. Join ..."/>
		<meta property="og:url" content="https://65words.com"/>
		<meta property="og:type" content="website"/>
		<meta property="og:image" content="https://65words.com/og-image.png"/>
		<meta property="og:site_name" content="65words.com"/>
		<link rel="canonical" href="https://65words.com"/>
		<meta property="name" content="65 Words: Write daily in your target language"/>
		<meta name="description" content="Improve your language skills by writing 65+ words daily in your chosen language. Enhance fluency, vocabulary, and proficiency through consistent practice. Join ..."/>
		<meta property="image" content="https://65words.com/og-image.png"/>
		<meta name="next-head-count" content="15"/>
		<link rel="preload" href="/_next/static/css/33f175ff2f814971.css" as="style"/>
		<link rel="stylesheet" href="/_next/static/css/33f175ff2f814971.css" data-n-g=""/>
		<noscript data-n-css=""></noscript>
		<script defer="" nomodule="" src="/_next/static/chunks/polyfills-c67a75d1b6f99dc8.js"></script>
		<script src="/_next/static/chunks/webpack-24780b5468e42e63.js" defer=""></script>
		<script src="/_next/static/chunks/framework-4556c45dd113b893.js" defer=""></script>
		<script src="/_next/static/chunks/main-864e2caf7ce338d4.js" defer=""></script>
		<script src="/_next/static/chunks/pages/_app-5011b08c887c2e0e.js" defer=""></script>
		<script src="/_next/static/chunks/521-f51d0dcdeca77210.js" defer=""></script>
		<script src="/_next/static/chunks/61-15af79ac6ffe8eda.js" defer=""></script>
		<script src="/_next/static/chunks/971-1ab3293a6b935540.js" defer=""></script>
		<script src="/_next/static/chunks/pages/index-57090c08c3cc0320.js" defer=""></script>
		<script src="/_next/static/Rc7GVdPFjwsNW_AhdMK2b/_buildManifest.js" defer=""></script>
		<script src="/_next/static/Rc7GVdPFjwsNW_AhdMK2b/_ssgManifest.js" defer=""></script>
	</head>
	<body>
		<div id="__next">
			<div class="is-overlay-center">
				<div class="spinner"></div>
			</div>
		</div>
		<script id="__NEXT_DATA__" type="application/json">{"props":{"pageProps":{}},"page":"/","query":{},"buildId":"Rc7GVdPFjwsNW_AhdMK2b","nextExport":true,"autoExport":true,"isFallback":false,"locale":"en","locales":["en"],"defaultLocale":"en","scriptLoader":[]}</script>
	</body>
	</html>`

	reader := strings.NewReader(test)

	body := extractFeedURL(reader)

	if body != "" {
		t.Fail()
	}

}
