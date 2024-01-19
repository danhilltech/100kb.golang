package parsing

import (
	"fmt"
	"strings"
	"testing"
)

var testA = `<!doctype html>
<html lang="en">

  <head>
    <meta charset="utf-8">
    <meta name="description" content="Personal blog of sussman@">
    <meta name="author" content="Ben Collins-Sussman">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Writings of sussman@</title>

    <link rel="stylesheet" href="/blog/bundle.css">

    <link rel="alternate" type="application/atom+xml"
          href="https://social.clawhammer.net/blog/feed/feed.xml"
          title="Writings from sussman@">
  </head>

  <body>
    
    <header>
      <h1>FAQ on leaving Google</h1>

      <div class="site-id">
      <p>
        This is
      the <a href="https://social.clawhammer.net/blog">personal
      blog</a> of <a href="https://www.red-bean.com/sussman">Ben
      Collins-Sussman</a>.
      <ul><li> Also see
      <a href="https://www.debuggingteams.com/">Debugging Teams</a>,
      his book on collaboration & culture in software
      engineering.</li>
      <li>His <a href="https://docs.google.com/presentation/d/10WFQUE7jBgyMNRnP_FPq3F7FkEtpsvXikkIp2nsI10k/edit?usp=sharing">How
      to Leader</a> talk (available
      in <a href="https://abseil.io/resources/swe-book/html/ch06.html">prose
      form</a>.)</li></ul>
      </p>
      </div>
                                
    </header>

    <main>
      <p><em>Context: When I was laid off from Google, I knew I'd be deluged with
questions.  I wrote this FAQ to share with friends and family, to
prevent repeated explanation.  But my other goal was to help so many
of my co-workers process and understand the repeated waves of mass
layoffs.</em></p>
<p><strong>What happened?</strong></p>
<p>Google just did another big round of layoffs.  I was part of them,
along with hundreds of others.  Many of us had long tenure or
seniority; my run was 18 years!</p>
<p><strong>Oh no!  But why were <em>you</em> targeted?</strong></p>
<p>I wasn’t personally targeted, I didn’t mess up.  In fact these layoffs
were extremely impersonal.  Google seems to be carrying out generic
initiatives to save operational cost.  I was an Engineering Director
with “only” 35 reports (rather than a typical 80+ people), and so it’s
likely that some heuristic decided that the business could do fine
without me.</p>
<p><strong>This is unfair!  After all you’ve done, how could Google do this to
you?</strong></p>
<p>Please understand: <em>Google is not a person.</em> It’s many groups of
people following locally-varying processes, rules, and culture.  To
that end, it makes no sense to either love or be angry at “Google”;
it’s not a consciousness, and it has no sense of duty nor debt.</p>
<p><strong>Are you OK?  I’m so sorry!  How are you coping?</strong></p>
<p>I’m fine.  :-) Google culture changed dramatically last year with its
first major round of layoffs, and I saw the writing on the wall.  I’ve
been preparing myself for this (increasingly inevitable) event for
months now – which included plenty of time for all the stages of
grief.  If anything, I have a mixed set of emotions:</p>
<ul>
<li>
<p>enormous pride in building a Chicago Engineering office over
decades, and achieving really cool things in the Developer, Ads, and
Search divisions;</p>
</li>
<li>
<p>deep gratitude in getting to work with some of the most intelligent,
creative people in the world;</p>
</li>
<li>
<p>a sense of relief.  The conflict between “uncomfortable culture” and
“golden handcuffs” was becoming intolerable.</p>
</li>
</ul>
<p><strong>What happens next?</strong></p>
<p>I’ve seen long-tenured leaders exit Google and go into an identity
crisis; that’s not me.  :-)</p>
<p>I have a <a href="https://www.red-bean.com/sussman/">zillion hobbies and shadow
careers</a> – plenty of things to do
and paths to follow.  The <em>first</em> order of business, however, is
probably a long-overdue sabbatical.  After 25+ years in tech, I need a
few months to rest and recover!</p>
<p>I’ll soon publish a couple of ‘post-mortem’ stories.  The first will
be about my own career at Google, and the second will be about how
I’ve seen Google culture change over time.</p>
<p><em>image: the first three software engineers at Google Chicago, 2006</em></p>
<p><img src="/blog/images/eng-chi-2006.jpg" alt="the first three software engineers at Google Chicago,
2006"></p>
<p><em>published January 10, 2024</em></p>

    </main>

  </body>
</html>
`

var testB = `



<!DOCTYPE html>
<html lang="en" class="no-js">

<head>
    


<script src="https://www.cdn.privado.ai/b230e0658a2d4f23bd9374dc1930f2c9.js" type="text/javascript" ></script>



                
        <title>From Vexing Uncertainty to Intellectual Humility | Schizophrenia Bulletin | Oxford Academic</title>

    <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.7.1/jquery.min.js" type="text/javascript"></script>
<script>window.jQuery || document.write('<script src="//oup.silverchair-cdn.com/Themes/Silver/app/js/jquery.3.7.1.min.js" type="text/javascript">\x3C/script>')</script>
<script src="//oup.silverchair-cdn.com/Themes/Silver/app/vendor/v-638397273706282681/jquery-migrate-1.4.1.min.js" type="text/javascript"></script>

    <script type='text/javascript' src='https://platform-api.sharethis.com/js/sharethis.js#property=643701de45aa460012e1032e&amp;product=sop' async='async' class='optanon-category-C0004'></script>


    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=10" />

    
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=Edge" />
    <!-- Turn off telephone number detection. -->
    <meta name="format-detection" content="telephone=no" />

<!-- Bookmark Icons -->
  <link rel="apple-touch-icon" sizes="180x180" href="//oup.silverchair-cdn.com/UI/app/img/v-638397272891870606/apple-touch-icon.png">
  <link rel="icon" type="image/png" href="//oup.silverchair-cdn.com/UI/app/img/v-638397272891870606/favicon-32x32.png" sizes="32x32">
  <link rel="icon" type="image/png" href="//oup.silverchair-cdn.com/UI/app/img/v-638397272891870606/favicon-16x16.png" sizes="16x16">
  <link rel="mask-icon" href="//oup.silverchair-cdn.com/UI/app/img/v-638397272892020660/safari-pinned-tab.svg" color="#001C54">
  <link rel="icon" href="//oup.silverchair-cdn.com/UI/app/img/v-638397272891920635/favicon.ico">
  <link rel="manifest" href="//oup.silverchair-cdn.com/UI/app/img/v-638397272891920635/manifest.json">
  <meta name="msapplication-config" content="//oup.silverchair-cdn.com/UI/app/img/v-638397272891870606/browserconfig.xml">
  <meta name="theme-color" content="#002f65">


    




<link rel="stylesheet" type="text/css" href="//oup.silverchair-cdn.com/UI/app/fonts/icons.css" />

<link rel="stylesheet" type="text/css" href="//oup.silverchair-cdn.com/Themes/Client/app/css/v-638411137764840886/site.min.css" />


    <link rel="preload" 
                  href="https://fonts.googleapis.com/css?family=Merriweather:300,400,400italic,700,700italic|Source+Sans+Pro:400,400italic,700,700italic" 
                  as="style" 
                  onload="this.onload=null;this.rel='stylesheet'">


        <link href="//oup.silverchair-cdn.com/data/SiteBuilderAssetsOriginals/Live/CSS/journals/v-638373778184040619/global.css" rel="stylesheet" type="text/css" />
            <link href="//oup.silverchair-cdn.com/data/SiteBuilderAssets/Live/CSS/schizophreniabulletin/v-637435441280303859/Site.css" rel="stylesheet" type="text/css" />



<script> var dataLayer = [{"full_title":"From Vexing Uncertainty to Intellectual Humility","short_title":"From Vexing Uncertainty to Intellectual Humility","authors":"Michael Dickson","type":"oration","online_publication_date":"2024-01-11","access_type":"Free","license_type":"publisher-standard","event_type":"full-text","discipline_ot_level_1":"Medicine and Health","discipline_ot_level_2":"Psychiatry","supplier_tag":"SC_Journals","object_type":"Article","taxonomy":"taxId%3a39%7ctaxLabel%3aAcademicSubjects%7cnodeId%3aMED00810%7cnodeLabel%3aChild+and+Adolescent+Psychiatry%7cnodeLevel%3a3","siteid":"schbul","authzrequired":"false","doi":"10.1093/schbul/sbad173"}]; </script>
            <script>
                (function (w, d, s, l, i) {
                    w[l] = w[l] || []; w[l].push({
                        'gtm.start':
                            new Date().getTime(), event: 'gtm.js'
                    }); var f = d.getElementsByTagName(s)[0],
                        j = d.createElement(s), dl = l != 'dataLayer' ? '&l=' + l : '';
                        j.setAttribute('type', 'text/javascript');
                        j.setAttribute('class', 'optanon-category-C0002');
                        j.async = true; j.src =
                        'https://www.googletagmanager.com/gtm.js?id=' + i + dl; f.parentNode.insertBefore(j, f);
                })(window, document, 'script', 'dataLayer', 'GTM-W6DD7HV');
            </script>

    

        <script type="text/javascript">
            var App = App || {};
            App.LoginUserInfo = {
                isInstLoggedIn: 0,
                isIndividualLoggedIn: 0
            };

            App.CurrentSubdomain = 'schizophreniabulletin';
            App.SiteURL = 'academic.oup.com/schizophreniabulletin';
        </script>
    
    
    <link href="https://cdn.jsdelivr.net/chartist.js/latest/chartist.min.css" media="print" onload="this.onload=null;this.removeAttribute('media');" rel="stylesheet" type="text/css" />
    <script type="application/ld+json">
        {"@context":"https://schema.org","@type":"ScholarlyArticle","@id":"https://academic.oup.com/schizophreniabulletin/advance-article/doi/10.1093/schbul/sbad173/7517011","name":"From Vexing Uncertainty to Intellectual Humility","datePublished":"2024-01-11","isPartOf":{"@id":"https://academic.oup.com/schizophreniabulletin/schizophreniabulletin","@type":"Periodical","name":"Schizophrenia Bulletin","issn":["1745-1701"]},"url":"https://dx.doi.org/10.1093/schbul/sbad173","inLanguage":"en","copyrightHolder":"Maryland Psychiatric Research Center and Oxford University Press","copyrightYear":"2024","publisher":"Oxford University Press","sameAs":"","author":[{"name":"Dickson, Michael","@type":"Person"}],"description":"I am a 55-year-old husband, father, friend, and professional philosopher. In 1992, as a graduate student at Cambridge University, a porter found me amongst the ","pageStart":"sbad173","pageEnd":"","siteName":"OUP Academic","thumbnailURL":"https://academic.oup.com/data/sitebuilderassetsoriginals/live/images/schizophreniabulletin/schizophreniabulletin_ogimage.png","headline":"From Vexing Uncertainty to Intellectual Humility","image":"https://academic.oup.com/data/sitebuilderassetsoriginals/live/images/schizophreniabulletin/schizophreniabulletin_ogimage.png","image:alt":"site image"}
    </script>
<meta property="og:site_name" content="OUP Academic" />
<meta property="og:title" content="From Vexing Uncertainty to Intellectual Humility" />
<meta property="og:description" content="I am a 55-year-old husband, father, friend, and professional philosopher. In 1992, as a graduate student at Cambridge University, a porter found me amongst the " />
<meta property="og:type" content="article" />
<meta property="og:url" content="https://dx.doi.org/10.1093/schbul/sbad173" />
<meta property="og:updated_time" content="" />
<meta property="og:image" content="https://academic.oup.com/data/sitebuilderassetsoriginals/live/images/schizophreniabulletin/schizophreniabulletin_ogimage.png" />
<meta property="og:image:url" content="https://academic.oup.com/data/sitebuilderassetsoriginals/live/images/schizophreniabulletin/schizophreniabulletin_ogimage.png" />
<meta property="og:image:secure_url" content="https://academic.oup.com/data/sitebuilderassetsoriginals/live/images/schizophreniabulletin/schizophreniabulletin_ogimage.png" />
<meta property="og:image:alt" content="site image" />
<meta name="twitter:card" content="summary_large_image" />
<meta name="citation_author" content="Dickson, Michael" /><meta name="citation_title" content="From Vexing Uncertainty to Intellectual Humility" /><meta name="citation_doi" content="10.1093/schbul/sbad173" /><meta name="citation_journal_title" content="Schizophrenia Bulletin" /><meta name="citation_journal_abbrev" content="Schizophr Bull" /><meta name="citation_pdf_url" content="https://academic.oup.com/schizophreniabulletin/advance-article-pdf/doi/10.1093/schbul/sbad173/55466528/sbad173.pdf" /><meta name="description" content="I am a 55-year-old husband, father, friend, and professional philosopher. In 1992, as a graduate student at Cambridge University, a porter found me amongst the " /><meta name="citation_xml_url" content="https://academic.oup.com/schizophreniabulletin/advance-article-xml/doi/10.1093/schbul/sbad173/7517011" />


    <link rel="canonical" href="https://academic.oup.com/schizophreniabulletin/advance-article/doi/10.1093/schbul/sbad173/7517011" />












    



<script async="async" src="https://securepubads.g.doubleclick.net/tag/js/gpt.js"></script>
<script>
    var googletag = googletag || {};
    googletag.cmd = googletag.cmd || [];
</script>
    
        <script type='text/javascript'>
            var gptAdSlots = [];
            googletag.cmd.push(function() {
                    
                    var mapping_ad1 = googletag.sizeMapping()
                        .addSize([1024, 0], [[970, 90], [728, 90]])
                        .addSize([768, 0], [728, 90])
                        .addSize([0, 0], [320, 50])
                        .build();
                    gptAdSlots["ad1"] = googletag.defineSlot('/116097782/schizophreniabulletin_Behind_Ad1',
                            [[970, 90], [728, 90], [320, 50]], 'adBlockHeader')
                        .defineSizeMapping(mapping_ad1)
                        .addService(googletag.pubads());
                    

                    
                    var mapping_ad2 = googletag.sizeMapping()
                        .addSize([768, 0], [[300, 250], [300, 600], [160, 600]])
                        .build();
                    gptAdSlots["ad2"] = googletag.defineSlot('/116097782/schizophreniabulletin_Behind_Ad2',
                            [[300, 250], [160, 600], [300, 600]], 'adBlockMainBodyTop')
                        .defineSizeMapping(mapping_ad2)
                        .addService(googletag.pubads());
                    

                    
                    var mapping_ad3 = googletag.sizeMapping()
                        .addSize([768, 0], [[300, 250], [300, 600], [160, 600]])
                        .build();
                    gptAdSlots["ad3"] = googletag.defineSlot('/116097782/schizophreniabulletin_Behind_Ad3',
                        [[300, 250], [160, 600], [300, 600]], 'adBlockMainBodyBottom')
                        .defineSizeMapping(mapping_ad3)
                        .addService(googletag.pubads());
                    

                    
                var mapping_ad4 = googletag.sizeMapping()
                        .addSize([0,0], [320, 50])
                        .addSize([768, 0], [728, 90])
                        .build();
                    gptAdSlots["ad4"] = googletag.defineSlot('/116097782/schizophreniabulletin_Behind_Ad4',
                        [728, 90], 'adBlockFooter')
                        .defineSizeMapping(mapping_ad4)
                        .addService(googletag.pubads());
                    

                    
                var mapping_ad6 = googletag.sizeMapping()
                        .addSize([1024, 0], [[970, 90], [728, 90]])
                        .addSize([768, 0], [728, 90])
                        .addSize([0, 0], [320, 50])
                        .build();
                    gptAdSlots["ad6"] = googletag.defineSlot('/116097782/schizophreniabulletin_Behind_Ad6',
                        [[728, 90], [970, 90]], 'adBlockStickyFooter')
                        .defineSizeMapping(mapping_ad6)
                        .addService(googletag.pubads());
                    

                    
                    gptAdSlots["adInterstital"] = googletag.defineOutOfPageSlot('/116097782/schizophreniabulletin_Interstitial_Ad',
                        googletag.enums.OutOfPageFormat.INTERSTITIAL)
                        .addService(googletag.pubads());
                                        

                googletag.pubads().addEventListener('slotRenderEnded', function (event) {
                    if (!event.isEmpty) {
                        $('.js-' + event.slot.getSlotElementId()).each(function () {
                            if ($(this).find('iframe').length) {
                                $(this).removeClass('hide');
                            }
                        });
                    }
                });

                googletag.pubads().addEventListener('impressionViewable', function (event) {
                    if (!event.isEmpty) {
                        $('.js-' + event.slot.getSlotElementId()).each(function () {
                            var $adblockDiv = $(this).find('.js-adblock');
                            var $adText = $(this).find('.js-adblock-advertisement-text');
                            if ($adblockDiv && $adblockDiv.is(':visible') && $adblockDiv.find('*').length > 1) {
                                $adText.removeClass('hide');
                                App.CenterAdBlock.Init($adblockDiv, $adText);
                            }
                            else {
                                $adText.addClass('hide');
                            }

                            //Initialize logic for Sticky Footer Ad
                            var $stickyFooterDiv = $(this).parents('.js-sticky-footer-ad');
                            if ($stickyFooterDiv && $stickyFooterDiv.is(':visible') && $stickyFooterDiv.find('*').length > 1) {
                                App.StickyFooterAd.Init();
                            }
                        });
                    }
                });

                googletag.pubads().setTargeting("jnlspage", "advance-article");
                googletag.pubads().setTargeting("jnlsurl", "schizophreniabulletin/advance-article/doi/10.1093/schbul/sbad173/7517011");
                googletag.pubads().enableSingleRequest();
                googletag.pubads().collapseEmptyDivs();
            });
        </script>
    
<input type="hidden"
       class="hfInterstitial"
       data-interstitiallinks="schizophreniabulletin/issue,schizophreniabulletin/advance-articles,schizophreniabulletin/advance-article,schizophreniabulletin/supplements,schizophreniabulletin/article,schizophreniabulletin/article-abstract,schizophreniabulletin/pages"
       data-subdomain="schizophreniabulletin"/>
    
    
        <script type="text/javascript">
                googletag.cmd.push(function () {
                    
                    googletag.pubads().setTargeting("jnlsdoi", "10.1093/schbul/sbad173");
                    googletag.enableServices();
                });
        </script>




    


    <script type="text/javascript">
        var NTPT_PGEXTRA= 'event_type=full-text&discipline_ot_level_1=Medicine and Health&discipline_ot_level_2=Psychiatry&supplier_tag=SC_Journals&object_type=Article&taxonomy=taxId%3a39%7ctaxLabel%3aAcademicSubjects%7cnodeId%3aMED00810%7cnodeLabel%3aChild+and+Adolescent+Psychiatry%7cnodeLevel%3a3&siteid=schbul&authzrequired=false&doi=10.1093/schbul/sbad173';
    </script>




    <script src="https://scholar.google.com/scholar_js/casa.js" async></script>
</head>

<body data-sitename="schizophreniabulletin" class="off-canvas pg_Article pg_article   " theme-schizophreniabulletin data-sitestyletemplate="Journal" >
            <noscript>
                <iframe title="AIP Publishing Google Tag Manager iframe" src="https://www.googletagmanager.com/ns.html?id=GTM-W6DD7HV"
                        height="0" width="0" style="display:none;visibility:hidden"></iframe>
            </noscript>
            <a href="#skipNav" class="skipnav">Skip to Main Content</a>
<input id="hdnSiteID" name="hdnSiteID" type="hidden" value="5240" /><input id="hdnAdDelaySeconds" name="hdnAdDelaySeconds" type="hidden" value="8000" /><input id="hdnAdConfigurationTop" name="hdnAdConfigurationTop" type="hidden" value="scrolldelay" /><input id="hdnAdConfigurationRightRail" name="hdnAdConfigurationRightRail" type="hidden" value="sticky" />
    






<div class="master-container js-master-container">
<section class="master-header row js-master-header vt-site-page-header">
    <div class="widget widget-SitePageHeader widget-instance-SitePageHeader">
        

    <div class="ad-banner js-ad-banner-header">
    <div class="widget widget-AdBlock widget-instance-HeaderAd">
        <div class="js-adBlock-parent-wrap adblock-parent-wrap">
    <div class="adBlockHeader-wrap js-adBlockHeader hide">

        <div id="adBlockHeader" class="js-adblock at-adblock">
            <script>
            googletag.cmd.push(function () {
                googletag.display('adBlockHeader');
            });
            </script>
        </div>
            <div class="advertisement-text at-adblock js-adblock-advertisement-text hide">Advertisement</div>
            </div>
</div> 
    </div>

    </div>

    <div class="oup-header sigma ">
        <div class="center-inner-row">

<div class="oup-header-logo">

        <a href="/">
            <img src="//oup.silverchair-cdn.com/UI/app/svg/umbrella/oxford-academic-logo.svg" alt="Oxford Academic"
                 class="oup-header-image at-oup-header-image " />
        </a>

</div>                <div class="widget widget-CustomNavLinks widget-instance-CustomNavLinksDeskTop">
            <div class="custom-nav-links-box">
            <div class="custom-nav-link">
                <a href="/journals">Journals</a>
            </div>
            <div class="custom-nav-link">
                <a href="/books">Books</a>
            </div>

    </div>
 
    </div>
            <ul class="oup-header-menu account-menu sigma-account-menu ">

                
                <li class="oup-header-menu-item mobile">
                    <a href="javascript:;" class="mobile-dropdown-toggle mobile-search-toggle">
                        <i class="icon-menu_search"><span class="screenreader-text">Search Menu</span></i>
                    </a>
                </li>

                
                        <li class="oup-header-menu-item mobile info-icon-menu-item">
            <a href="/pages/information" target="_blank" class="at-info-button sigma-info-wrapper" role="button">
                <img class="sigma-info-icon" src="//oup.silverchair-cdn.com/UI/app/svg/i.svg" alt="Information" />
            </a>
        </li>
    <li class="oup-header-menu-item mobile account-icon-menu-item">
        <a href="javascript:;" class="account-button js-account-button at-account-button "
           role="button" data-turnawayparams="journal%3dschizophreniabulletin">
            <img class="sigma-account-icon" src="//oup.silverchair-cdn.com/UI/app/svg/account.svg" alt="Account" />
        </a>
    </li>


                
                <li class="oup-header-menu-item mobile">
                    <a href="javascript:;" class="mobile-dropdown-toggle mobile-nav-toggle">
                        <i class="icon-menu_hamburger"><span class="screenreader-text">Menu</span></i>
                    </a>
                </li>

                

                        <li class="oup-header-menu-item desktop info-icon-menu-item">
            <a href="/pages/information" target="_blank" class="at-info-button sigma-info-wrapper" role="button">
                <img class="sigma-info-icon" src="//oup.silverchair-cdn.com/UI/app/svg/i.svg" alt="Information" />
            </a>
        </li>
    <li class="oup-header-menu-item desktop account-icon-menu-item">
        <a href="javascript:;" class="account-button js-account-button at-account-button sigma-logo-wrapper"
           role="button" data-turnawayparams="journal%3dschizophreniabulletin">
            <img class="sigma-account-icon" src="//oup.silverchair-cdn.com/UI/app/svg/account.svg" alt="Account" />
        </a>
    </li>


                

            </ul>

            <div class="login-box-placeholder js-login-box-placeholder hide">
                <div class="spinner"></div>
            </div>
        </div>
    </div>
<div class="dropdown-panel-wrap">

    
    <div class="dropdown-panel mobile-search-dropdown">
        <div class="mobile-search-inner-wrap">


<div class="navbar-search">
    <div class="mobile-microsite-search">
            <label for="SitePageHeader-mobile-navbar-search-filter" class="screenreader-text js-mobile-navbar-search-filter-label">
                Navbar Search Filter
            </label>
            <select class="mobile-navbar-search-filter js-mobile-navbar-search-filter at-navbar-search-filter" id="SitePageHeader-mobile-navbar-search-filter">

<option class="navbar-search-filter-option at-navbar-search-filter-option" value="">Schizophrenia Bulletin</option><option class="navbar-search-filter-option at-navbar-search-filter-option" value="Parent">Schizophrenia Bulletin Journals</option>
                    <optgroup class="navbar-search-optgroup" label="Search across Oxford Academic">
<option class="navbar-search-filter-option at-navbar-search-filter-option" value="AcademicSubjects/MED00810">Child and Adolescent Psychiatry</option><option class="navbar-search-filter-option at-navbar-search-filter-option" value="Books">Books</option><option class="navbar-search-filter-option at-navbar-search-filter-option" value="Journals">Journals</option><option class="navbar-search-filter-option at-navbar-search-filter-option" value="Umbrella">Oxford Academic</option>                    </optgroup>
            </select>

        <label for="SitePageHeader-mobile-microsite-search-term" class="screenreader-text js-mobile-microsite-search-term-label">
            Mobile Enter search term
        </label>
        <input class="mobile-search-input mobile-microsite-search-term js-mobile-microsite-search-term at-microsite-search-term" type="text"
               maxlength="255" placeholder="Search" id="SitePageHeader-mobile-microsite-search-term">

        <a href="javascript:;" class="mobile-microsite-search-icon mobile-search-submit icon-menu_search">
            <span class="screenreader-text">Search</span>
        </a>

    </div>
</div>


        </div>
    </div>
   
    <div class="dropdown-panel mobile-nav-dropdown">


    <ul class="site-menu site-menu-lvl-0 at-site-menu">
        <li class="site-menu-item site-menu-lvl-0 at-site-menu-item" id="site-menu-item-1575628">

                <a href="/schizophreniabulletin/issue" class="nav-link">
                    Issues
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-0 at-site-menu-item" id="site-menu-item-1575625">

                <a href="javascript:;" class="nav-link js-nav-dropdown at-nav-dropdown" role="button" aria-expanded="false">
                    More Content
                    <i class="desktop-nav-arrow icon-general-arrow-filled-down arrow-icon"></i>
                </a>
                <i class="mobile-nav-arrow icon-general_arrow-down"></i>
                <ul class="site-menu site-menu-lvl-1 at-site-menu">
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575629">

                <a href="/schizophreniabulletin/advance-articles" class="nav-link">
                    Advance articles
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575630">

                <a href="/schizophreniabulletin/search-results?page=1&amp;f_OUPSeries=Editor%27s+Choice" class="nav-link">
                    Editor's Choice
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575631">

                <a href="https://academic.oup.com/schizophreniabulletin/supplements" class="nav-link">
                    Supplements
                </a>
            
        </li>
            </ul>

        </li>
        <li class="site-menu-item site-menu-lvl-0 at-site-menu-item" id="site-menu-item-1575626">

                <a href="javascript:;" class="nav-link js-nav-dropdown at-nav-dropdown" role="button" aria-expanded="false">
                    Submit
                    <i class="desktop-nav-arrow icon-general-arrow-filled-down arrow-icon"></i>
                </a>
                <i class="mobile-nav-arrow icon-general_arrow-down"></i>
                <ul class="site-menu site-menu-lvl-1 at-site-menu">
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575632">

                <a href="https://mc.manuscriptcentral.com/szbltn" class="nav-link">
                    Submission Site
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575633">

                <a href="/schizophreniabulletin/pages/General_Instructions" class="nav-link">
                    Author Guidelines
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575634">

                <a href="/schizophreniabulletin/pages/open-access" class="nav-link">
                    Open Access
                </a>
            
        </li>
            </ul>

        </li>
        <li class="site-menu-item site-menu-lvl-0 at-site-menu-item" id="site-menu-item-1575635">

                <a href="/schizophreniabulletin/subscribe" class="nav-link">
                    Purchase
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-0 at-site-menu-item" id="site-menu-item-1575636">

                <a href="/my-account/email-alerts" class="nav-link">
                    Alerts
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-0 at-site-menu-item" id="site-menu-item-1575627">

                <a href="javascript:;" class="nav-link js-nav-dropdown at-nav-dropdown" role="button" aria-expanded="false">
                    About
                    <i class="desktop-nav-arrow icon-general-arrow-filled-down arrow-icon"></i>
                </a>
                <i class="mobile-nav-arrow icon-general_arrow-down"></i>
                <ul class="site-menu site-menu-lvl-1 at-site-menu">
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575637">

                <a href="/schizophreniabulletin/pages/About" class="nav-link">
                    About Schizophrenia Bulletin
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575638">

                <a href="http://medschool.umaryland.edu" class="nav-link">
                    About the University of Maryland School of Medicine
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575639">

                <a href="http://www.mprc.umaryland.edu/default.asp" class="nav-link">
                    About the Maryland Psychiatric Research Center
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575640">

                <a href="/schizophreniabulletin/article/34/5/799/1886770/Schizophrenia-Bulletin-and-the-Revised-NIH-Public" class="nav-link">
                    About the NIH Public Access Policy
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575641">

                <a href="/schizophreniabulletin/pages/Editorial_Board" class="nav-link">
                    Editorial Board
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575642">

                <a href="https://academic.oup.com/advertising-and-corporate-services/pages/schizophreniabulletin-media-kit" class="nav-link">
                    Advertising and Corporate Services
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575643">

                <a href="http://medicine-and-health-careernetwork.oxfordjournals.org" class="nav-link">
                    Journals Career Network
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575644">

                <a href="https://academic.oup.com/journals/pages/access_purchase/rights_and_permissions/self_archiving_policy_b" class="nav-link">
                    Self-Archiving Policy
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575645">

                <a href="http://www.oxfordjournals.org/en/access-purchase/dispatch-dates.html" class="nav-link">
                    Dispatch Dates
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575646">

                <a href="http://medschool.umaryland.edu/" class="nav-link">
                    Contact Us
                </a>
            
        </li>
            </ul>

        </li>
                <li class="site-menu-item site-menu-lvl-0 at-site-menu-item" id="site-menu-item-custom">
            <a href="/journals" class="nav-link">Journals on Oxford Academic</a>
        </li>
        <li class="site-menu-item site-menu-lvl-0 at-site-menu-item" id="site-menu-item-custom">
            <a href="/books" class="nav-link">Books on Oxford Academic</a>
        </li>
    </ul>
    </div>

</div>

    <div class="journal-header journal-bg">
        <div class="center-inner-row">
                <div class="site-parent-link-wrap">
                    <a href="//academic.oup.com/schizbulljournals" class="site-parent-link">Schizophrenia Bulletin Journals</a>
                </div>
                            <a href="/schizophreniabulletin" class="journal-logo-container">
                    <img id="logo-SchizophreniaBulletin" class="journal-logo" src="//oup.silverchair-cdn.com/data/SiteBuilderAssets/Live/Images/schizophreniabulletin/title307073958.svg" alt="Schizophrenia Bulletin:The Journal of Psychoses and Related Disorders" />
                </a>

            <div class="society-logo-block">
                    <div class="society-block-inner-wrap">
                            <a href="http://medschool.umaryland.edu/" target="" class="society-logo-container">
                                <img id="logo-UniversityofMarylandSchoolofMedicine" class="society-logo" src="//oup.silverchair-cdn.com/data/SiteBuilderAssets/Live/Images/schizophreniabulletin/h1-1614261906.png" alt="University of Maryland School of Medicine" />
                            </a>
                                                    <a href="http://www.mprc.umaryland.edu/" target="" class="society-logo-container">
                                <img id="logo-MarylandPsychiatricResearchCenter" class="society-logo" src="//oup.silverchair-cdn.com/data/SiteBuilderAssets/Live/Images/schizophreniabulletin/h2-1901022237.png" alt="Maryland Psychiatric Research Center" />
                            </a>
                    </div>            </div> 
        </div>
    </div>
<div class="navbar">
    <div class="center-inner-row">
        <nav class="navbar-menu">


    <ul class="site-menu site-menu-lvl-0 at-site-menu">
        <li class="site-menu-item site-menu-lvl-0 at-site-menu-item" id="site-menu-item-1575628">

                <a href="/schizophreniabulletin/issue" class="nav-link">
                    Issues
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-0 at-site-menu-item" id="site-menu-item-1575625">

                <a href="javascript:;" class="nav-link js-nav-dropdown at-nav-dropdown" role="button" aria-expanded="false">
                    More Content
                    <i class="desktop-nav-arrow icon-general-arrow-filled-down arrow-icon"></i>
                </a>
                <i class="mobile-nav-arrow icon-general_arrow-down"></i>
                <ul class="site-menu site-menu-lvl-1 at-site-menu">
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575629">

                <a href="/schizophreniabulletin/advance-articles" class="nav-link">
                    Advance articles
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575630">

                <a href="/schizophreniabulletin/search-results?page=1&amp;f_OUPSeries=Editor%27s+Choice" class="nav-link">
                    Editor's Choice
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575631">

                <a href="https://academic.oup.com/schizophreniabulletin/supplements" class="nav-link">
                    Supplements
                </a>
            
        </li>
            </ul>

        </li>
        <li class="site-menu-item site-menu-lvl-0 at-site-menu-item" id="site-menu-item-1575626">

                <a href="javascript:;" class="nav-link js-nav-dropdown at-nav-dropdown" role="button" aria-expanded="false">
                    Submit
                    <i class="desktop-nav-arrow icon-general-arrow-filled-down arrow-icon"></i>
                </a>
                <i class="mobile-nav-arrow icon-general_arrow-down"></i>
                <ul class="site-menu site-menu-lvl-1 at-site-menu">
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575632">

                <a href="https://mc.manuscriptcentral.com/szbltn" class="nav-link">
                    Submission Site
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575633">

                <a href="/schizophreniabulletin/pages/General_Instructions" class="nav-link">
                    Author Guidelines
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575634">

                <a href="/schizophreniabulletin/pages/open-access" class="nav-link">
                    Open Access
                </a>
            
        </li>
            </ul>

        </li>
        <li class="site-menu-item site-menu-lvl-0 at-site-menu-item" id="site-menu-item-1575635">

                <a href="/schizophreniabulletin/subscribe" class="nav-link">
                    Purchase
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-0 at-site-menu-item" id="site-menu-item-1575636">

                <a href="/my-account/email-alerts" class="nav-link">
                    Alerts
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-0 at-site-menu-item" id="site-menu-item-1575627">

                <a href="javascript:;" class="nav-link js-nav-dropdown at-nav-dropdown" role="button" aria-expanded="false">
                    About
                    <i class="desktop-nav-arrow icon-general-arrow-filled-down arrow-icon"></i>
                </a>
                <i class="mobile-nav-arrow icon-general_arrow-down"></i>
                <ul class="site-menu site-menu-lvl-1 at-site-menu">
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575637">

                <a href="/schizophreniabulletin/pages/About" class="nav-link">
                    About Schizophrenia Bulletin
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575638">

                <a href="http://medschool.umaryland.edu" class="nav-link">
                    About the University of Maryland School of Medicine
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575639">

                <a href="http://www.mprc.umaryland.edu/default.asp" class="nav-link">
                    About the Maryland Psychiatric Research Center
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575640">

                <a href="/schizophreniabulletin/article/34/5/799/1886770/Schizophrenia-Bulletin-and-the-Revised-NIH-Public" class="nav-link">
                    About the NIH Public Access Policy
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575641">

                <a href="/schizophreniabulletin/pages/Editorial_Board" class="nav-link">
                    Editorial Board
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575642">

                <a href="https://academic.oup.com/advertising-and-corporate-services/pages/schizophreniabulletin-media-kit" class="nav-link">
                    Advertising and Corporate Services
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575643">

                <a href="http://medicine-and-health-careernetwork.oxfordjournals.org" class="nav-link">
                    Journals Career Network
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575644">

                <a href="https://academic.oup.com/journals/pages/access_purchase/rights_and_permissions/self_archiving_policy_b" class="nav-link">
                    Self-Archiving Policy
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575645">

                <a href="http://www.oxfordjournals.org/en/access-purchase/dispatch-dates.html" class="nav-link">
                    Dispatch Dates
                </a>
            
        </li>
        <li class="site-menu-item site-menu-lvl-1 at-site-menu-item" id="site-menu-item-1575646">

                <a href="http://medschool.umaryland.edu/" class="nav-link">
                    Contact Us
                </a>
            
        </li>
            </ul>

        </li>
            </ul>
        </nav>
            <div class="navbar-search-container js-navbar-search-container">
                <a href="javascript:;" class="navbar-search-close js_close-navsearch">Close</a>


<div class="navbar-search">
    <div class="microsite-search">
            <label for="SitePageHeader-navbar-search-filter" class="screenreader-text js-navbar-search-filter-label">
                Navbar Search Filter
            </label>
            <select class="navbar-search-filter js-navbar-search-filter at-navbar-search-filter" id="SitePageHeader-navbar-search-filter">

<option class="navbar-search-filter-option at-navbar-search-filter-option" value="">Schizophrenia Bulletin</option><option class="navbar-search-filter-option at-navbar-search-filter-option" value="Parent">Schizophrenia Bulletin Journals</option>
                    <optgroup class="navbar-search-optgroup" label="Search across Oxford Academic">
<option class="navbar-search-filter-option at-navbar-search-filter-option" value="AcademicSubjects/MED00810">Child and Adolescent Psychiatry</option><option class="navbar-search-filter-option at-navbar-search-filter-option" value="Books">Books</option><option class="navbar-search-filter-option at-navbar-search-filter-option" value="Journals">Journals</option><option class="navbar-search-filter-option at-navbar-search-filter-option" value="Umbrella">Oxford Academic</option>                    </optgroup>
            </select>

        <label for="SitePageHeader-microsite-search-term" class="screenreader-text js-microsite-search-term-label">
            Enter search term
        </label>
        <input class="navbar-search-input microsite-search-term js-microsite-search-term at-microsite-search-term" type="text"
               maxlength="255" placeholder="Search" id="SitePageHeader-microsite-search-term">

        <a href="javascript:;" class="microsite-search-icon navbar-search-submit icon-menu_search">
            <span class="screenreader-text">Search</span>
        </a>

    </div>
</div>

<input id="hfCurrentBookSearch" name="hfCurrentBookSearch" type="hidden" value="" /><input id="hfCurrentBookScope" name="hfCurrentBookScope" type="hidden" value="CurrentBook" /><input id="hfBookSiteScope" name="hfBookSiteScope" type="hidden" value="Books" /><input id="hfSeriesScope" name="hfSeriesScope" type="hidden" value="taxWithOr" /><input id="hfParentSiteName" name="hfParentSiteName" type="hidden" value="Schizophrenia Bulletin Journals" /><input id="hfParentSiteUrl" name="hfParentSiteUrl" type="hidden" value="academic.oup.com/schizbulljournals" /><input id="hfSiteID" name="hfSiteID" type="hidden" value="5240" /><input id="hfParentSiteID" name="hfParentSiteID" type="hidden" value="6260" /><input id="hfJournalSiteScope" name="hfJournalSiteScope" type="hidden" value="Journals" /><input id="hfParentSiteScope" name="hfParentSiteScope" type="hidden" value="Parent" /><input id="hfDefaultSearchURL" name="hfDefaultSearchURL" type="hidden" value="search-results?page=1&amp;q=" /><input id="hfIssueSearch" name="hfIssueSearch" type="hidden" value="" /><input id="hfIssueSiteScope" name="hfIssueSiteScope" type="hidden" value="Issue" /><input id="hfUmbrellaScope" name="hfUmbrellaScope" type="hidden" value="Umbrella" /><input id="hfUmbrellaSiteUrl" name="hfUmbrellaSiteUrl" type="hidden" value="academic.oup.com" /><input id="hfUmbrellaSiteId" name="hfUmbrellaSiteId" type="hidden" value="191" /><input id="hfDefaultAdvancedSearchUrl" name="hfDefaultAdvancedSearchUrl" type="hidden" value="advanced-search?page=1&amp;q=" /><input id="hfTaggedCollectionScope" name="hfTaggedCollectionScope" type="hidden" value="" />                <div class="navbar-search-advanced"><a href="/schizophreniabulletin/advanced-search" class="advanced-search js-advanced-search">Advanced Search</a></div>
            </div>            <div class="navbar-search-collapsed"><a href="javascript:;" class="icon-menu_search js_expand-navsearch"><span class="screenreader-text">Search Menu</span></a></div>
    </div>
</div>
<input id="hfEnableOupOnlineProductsFeatures" name="hfEnableOupOnlineProductsFeatures" type="hidden" value="True" />

    <input type="hidden" name="searchScope" id="hfSolrJournalID" value="" />
    <input type="hidden" id="hfSolrJournalName" value="" />
    <input type="hidden" id="hfSolrMaxAllowSearchChar" value="100" />
    <input type="hidden" id="hfJournalShortName" value="" />
    <input type="hidden" id="hfSearchPlaceholder" value="" />
    <input type="hidden" name="hfGlobalSearchSiteURL" id="hfGlobalSearchSiteURL" value="" />
    <input type="hidden" name="hfSearchSiteURL" id="hfSiteURL" value="academic.oup.com/schizophreniabulletin" />
    <input type="hidden" name="RedirectSiteUrl" id="RedirectSiteUrl" value="httpszazjzjacademiczwoupzwcom" />
    <script type="text/javascript">
        (function () {
            var hfSiteUrl = document.getElementById('hfSiteURL');
            var siteUrl = hfSiteUrl.value;
            var subdomainIndex = siteUrl.indexOf('/');

            hfSiteUrl.value = location.host + (subdomainIndex >= 0 ? siteUrl.substring(subdomainIndex) : '');
        })();
    </script>
<input id="routename" name="RouteName" type="hidden" value="schizophreniabulletin" /> 
    </div>

</section>
        <div class="widget widget-SitewideBanner widget-instance-">
        

 
    </div>


    <div id="main" class="content-main js-main ui-base">
            <section class="master-main row">
                <div class="center-inner-row no-overflow">
                    <div id="skipNav" tabindex="-1"></div>
                    


<div class="page-column-wrap">
<div id="InfoColumn" class="page-column page-column--left js-left-nav-col">
    <div class="mobile-content-topbar hide">
        <button class="toggle-left-col toggle-left-col__article">Article Navigation</button>
    </div>
    <div class="info-inner-wrap js-left-nav">
        <button class="toggle-left-col__close btn-as-icon icon-general-close">
            <span class="screenreader-text">Close mobile search navigation</span>
        </button>
        <div class="responsive-nav-title">Article Navigation</div>
        <div class="info-widget-wrap">

            <div class="content-nav">
    <div class="widget widget-ArticleJumpLinks widget-instance-OUP_ArticleJumpLinks_Widget">
        

<h3 class="contents-title" style='display: none;'>Article Contents</h3>
<ul class="jumplink-list js-jumplink-list">
</ul>
 
    </div>

            </div>

        </div>
    </div>
</div>
<div class="sticky-toolbar js-sticky-toolbar"></div>
<div id="ContentColumn" class="page-column page-column--center">

    <div class="article-browse-top article-browse-mobile-nav js-mobile-nav">
    <div class="article-browse-mobile-nav-inner js-mobile-nav-inner">
        <button class="toggle-left-col toggle-left-col__article btn-as-link">
            Article Navigation
        </button>
    </div>
</div>
<div class="article-browse-top article-browse-mobile-nav mobile-sticky-toolbar js-mobile-nav-sticky">
    <div class="article-browse-mobile-nav-inner">
        <button class="toggle-left-col toggle-left-col__article btn-as-link">
            Article Navigation
        </button>
    </div>
</div>

    <div class="content-inner-wrap">
    <div class="widget widget-ArticleTopInfo widget-instance-OUP_ArticleTop_Info_Widget">
        <div class="module-widget article-top-widget">

    

    
        <div class="access-state-logos all-viewports">
            <span class="journal-info__format-label">Journal Article</span>
                    <span class="article-pubstate article-flag"></span>
                </div>


    <div class="widget-items">
                <div class="title-wrap">
                    <h1 class="wi-article-title article-title-main accessible-content-title at-articleTitle">
                        From Vexing Uncertainty to Intellectual Humility
<i class='icon-availability_free' title='Free' ></i>                    </h1>

                </div>
                <div class="wi-authors at-ArticleAuthors">
                    <div class="al-authors-list">
                                <span class="al-author-name-more js-flyout-wrap">




                                    <button type="button" class="linked-name js-linked-name-trigger btn-as-link">Michael Dickson</button><span class='delimiter'></span>

                                    <span class="al-author-info-wrap arrow-up">
<div class="info-card-author authorInfo_OUP_ArticleTop_Info_Widget">
    <div class="name-role-wrap">
        <div class="info-card-name">
    Michael Dickson
            <span class="info-card-footnote"><span class="xrefLink" id="jumplink-c1"></span><a href="javascript:;" reveal-id="c1" data-open="c1" class="link link-ref link-reveal xref-default"><!----></a></span>
</div>

    </div>


    <div class="info-author-correspondence">
        <div content-id="c1"><a href="mailto:dickson@sc.edu" target="_blank">dickson@sc.edu</a></div>
    </div>
    <div class="info-card-search-label">
        Search for other works by this author on:
    </div>

<div class="info-card-search info-card-search-internal">
    <a href="/schizophreniabulletin/search-results?f_Authors=Michael+Dickson" rel="nofollow">Oxford Academic</a>
</div>

    <div class="info-card-search info-card-search-pubmed">
        <a href="http://www.ncbi.nlm.nih.gov/pubmed?cmd=search&amp;term=Dickson M">PubMed</a>
    </div>
    <div class="info-card-search info-card-search-google">
        <a href="http://scholar.google.com/scholar?q=author:%22Dickson Michael%22">Google Scholar</a>
    </div>
</div>                                    </span>
                                </span>

                    </div>
                </div>
<div class="pub-history-wrap clearfix js-history-dropdown-wrap">

        <div class="pub-history-row clearfix">
            <div class="ww-citation-primary"><em>Schizophrenia Bulletin</em>, sbad173, <a href='https://doi.org/10.1093/schbul/sbad173'>https://doi.org/10.1093/schbul/sbad173</a></div>
        </div>
        <div class="pub-history-row clearfix">
            <div class="ww-citation-date-wrap">
                <div class="citation-label">Published:</div>
                <div class="citation-date">11 January 2024</div>
            </div>
                    </div>
            </div>
    </div>
</div>

<script>
    $(document).ready(function () {
        $('.article-top-widget').on('click', '.ati-toggle-trigger', function () {
            $(this).find('.icon-general-add, .icon-minus').toggleClass('icon-minus icon-general-add');
            $(this).siblings('.ati-toggle-content').toggleClass('hide');
        });

        // In Chrome, an anchor tag with target="_blank" and a "mailto:" href opens a new tab/window as well as the email client
        // I suspect this behavior will be corrected in the future
        // Remove the target="_blank"
        $('ul.wi-affiliationList').find('a[href^="mailto:"]').each(function () {
            $(this).removeAttr('target');
        });
    });
</script>

 
    </div>
    <div class="widget widget-ArticleLinks widget-instance-OUP_Article_Links_Widget">
         
    </div>


        <div class="article-body js-content-body">
            <div class="toolbar-wrap js-toolbar-wrap">
                <div class="toolbar-inner-wrap">
                    <ul id="Toolbar" role="navigation">
    <li class="toolbar-item item-pdf js-item-pdf">
        <a class="al-link pdf article-pdfLink" data-article-id="7517011" href="/schizophreniabulletin/advance-article-pdf/doi/10.1093/schbul/sbad173/55466528/sbad173.pdf">
            <img src=//oup.silverchair-cdn.com/UI/app/svg/pdf.svg alt="pdf" /><span class="pdf-link-text">PDF</span>
        </a>
    </li>
                                                    <li class="toolbar-item item-link item-split-view">
                                <a href="javascript:;"
                                   class="split-view js-split-view st-split-view at-split-view"
                                   target="">
                                    <i class="icon-menu-split"></i>
                                    Split View
                                </a>
                            </li>

                        <li class="toolbar-item item-with-dropdown item-views">
                            <a class="at-views-dropdown drop-trigger" href="javascript:;" data-dropdown="FilterDrop" aria-haspopup="true">
                                <i class="icon-menu_views"></i>
                                <div class="toolbar-label">
                                    <div class="toolbar-text">Views</div>
                                    <i class="icon-general-arrow-filled-down arrow-icon"></i>
                                </div>
                            </a>
                            <ul id="ViewsDrop" class="f-dropdown" data-dropdown-content aria-label="submenu">
                                <div class="arrow-up"></div>
                                <li class="article-content-filter js-article-content-filter" data-content-filter="article-content">
                                    <a href="javascript:;"><span>Article contents</span></a>
                                </li>
                                <li class="at-figures-tables article-content-filter js-article-content-filter" data-content-filter="figures-tables">
                                    <a href="javascript:;"><span>Figures &amp; tables</span></a>
                                </li>
                                <li class="article-content-filter js-article-content-filter" data-content-filter="video">
                                    <a href="javascript:;"><span>Video</span></a>
                                </li>
                                <li class="article-content-filter js-article-content-filter" data-content-filter="audio">
                                    <a href="javascript:;"><span>Audio</span></a>
                                </li>
                                <li class="article-content-filter js-article-content-filter" data-content-filter="supplementary-data">
                                    <a href="javascript:;"><span>Supplementary Data</span></a>
                                </li>
                            </ul>
                        </li>


                        <li class="toolbar-item item-cite js-item-cite">
    <div class="widget widget-ToolboxGetCitation widget-instance-OUP_Get_Citation">
        <a href="#" class="js-cite-button at-CiteButton" data-reveal-id="getCitation" data-reveal>
    <i class="icon-read-more"></i>
    <span>Cite</span>
</a>

<div id="getCitation" class="reveal-modal js-citation-modal" data-reveal>
    <h3 class="modal-title">Cite</h3>
        <div class="oxford-citation-text">
            <p>Michael Dickson,  From Vexing Uncertainty to Intellectual Humility, <em>Schizophrenia Bulletin</em>, 2024;, sbad173, <a href="https://doi.org/10.1093/schbul/sbad173">https://doi.org/10.1093/schbul/sbad173</a></p>
        </div>

    <div class="citation-download-wrap">
        <form action="/Citation/Download" method="get" id="citationModal">
            <input type="hidden" name="resourceId" value="7517011" />
            <input type="hidden" name="resourceType" value="3" />
            <label for="selectFormat" class="hide js-citation-format-label">Select Format</label>
            <select required name="citationFormat" class="citation-download-format js-citation-format" id="selectFormat">
                <option selected disabled >Select format</option>
                        <option value="0" >.ris (Mendeley, Papers, Zotero)</option>
                        <option value="1" >.enw (EndNote)</option>
                        <option value="2" >.bibtex (BibTex)</option>
                        <option value="3" >.txt (Medlars, RefWorks)</option>

            </select>
            <button class="btn citation-download-link disabled" type="submit">Download citation</button>
        </form>
    </div>

    <a class="close-reveal-modal" href="javascript:;"><i class="icon-general-close"></i><span class="screenreader-text">Close</span></a>
</div>
 
    </div>

                        </li>
                        <li class="toolbar-item item-tools">
    <div class="widget widget-ToolboxPermissions widget-instance-OUP_Get_Permissions">
            <div class="module-widget">
        <a href="https://s100.copyright.com/AppDispatchServlet?publisherName=OUP&amp;publication=1745-1701&amp;title=From%20Vexing%20Uncertainty%20to%20Intellectual%20Humility&amp;publicationDate=2024-01-11&amp;volumeNum=&amp;issueNum=&amp;author=Dickson%2C%20Michael&amp;startPage=1&amp;endPage=2&amp;contentId=10.1093%2Fschbul%2Fsbad173&amp;oa=&amp;copyright=%C2%A9%20The%20Author%28s%29%202024.%20Published%20by%20Oxford%20University%20Press%20on%20behalf%20of%20the%20Maryland%20Psychiatric%20Research%20Center.%20All%20rights%20reserved.%20For%20permissions%2C%20please%20email%3A%20journals.permissions%40oup.com&amp;orderBeanReset=True" id="PermissionsLink" class="" target="_blank">
            <i class="icon-menu_permissions">
                <span class="screenreader-text">Permissions Icon</span>
            </i>
            Permissions
        </a>
    </div>
 
    </div>

                        </li>
                        
    <li class="toolbar-item item-with-dropdown item-share">
        <a href="javascript:;" class="drop-trigger js-toolbar-dropdown at-ShareButton" data-dropdown="ShareDrop">
            <i class="icon-menu_share"><span class="screenreader-text">Share Icon</span></i>
            <span class="toolbar-label">
                <span class="toolbar-text">Share</span>
                <i class="arrow-icon icon-general-arrow-filled-down js-toolbar-arrow-icon"></i>
            </span>
        </a>
        <ul id="ShareDrop" class="addthis_toolbox addthis_default_style addthis_20x20_style f-dropdown" data-dropdown-content>
    <li>
        <a class="st-custom-button addthis_button_facebook"
           data-network="facebook"
           data-title="From Vexing Uncertainty to Intellectual Humility"
           data-url="https://academic.oup.com/schizophreniabulletin/article/doi/10.1093/schbul/sbad173/7517011"
           data-email-subject="From Vexing Uncertainty to Intellectual Humility"
           href="javascript:;"><span>Facebook</span></a>
    </li>
    <li>
        <a class="st-custom-button addthis_button_twitter"
           data-network="twitter"
           data-title="From Vexing Uncertainty to Intellectual Humility"
           data-url="https://academic.oup.com/schizophreniabulletin/article/doi/10.1093/schbul/sbad173/7517011"
           data-email-subject="From Vexing Uncertainty to Intellectual Humility"
           href="javascript:;"><span>Twitter</span></a>
    </li>
    <li>
        <a class="st-custom-button addthis_button_linkedin"
           data-network="linkedin"
           data-title="From Vexing Uncertainty to Intellectual Humility"
           data-url="https://academic.oup.com/schizophreniabulletin/article/doi/10.1093/schbul/sbad173/7517011"
           data-email-subject="From Vexing Uncertainty to Intellectual Humility"
           href="javascript:;"><span>LinkedIn</span></a>
    </li>
    <li>
        <a class="st-custom-button addthis_button_email"
           data-network="email"
           data-title="From Vexing Uncertainty to Intellectual Humility"
           data-url="https://academic.oup.com/schizophreniabulletin/article/doi/10.1093/schbul/sbad173/7517011"
           data-email-subject="From Vexing Uncertainty to Intellectual Humility"
           href="javascript:;"><span>Email</span></a>
    </li>



        </ul>
    </li>

                        
                    </ul>
                    <div class="toolbar-search">
    <div class="widget widget-SitePageHeader widget-instance-OUP_ArticleToolbarSearchBox">
        



<div class="dropdown-panel-wrap">

    
    <div class="dropdown-panel mobile-search-dropdown">
        <div class="mobile-search-inner-wrap">


<div class="navbar-search">
    <div class="mobile-microsite-search">
            <label for="OUP_ArticleToolbarSearchBox-mobile-navbar-search-filter" class="screenreader-text js-mobile-navbar-search-filter-label">
                Navbar Search Filter
            </label>
            <select class="mobile-navbar-search-filter js-mobile-navbar-search-filter at-navbar-search-filter" id="OUP_ArticleToolbarSearchBox-mobile-navbar-search-filter">

<option class="navbar-search-filter-option at-navbar-search-filter-option" value="">Schizophrenia Bulletin</option><option class="navbar-search-filter-option at-navbar-search-filter-option" value="Parent">Schizophrenia Bulletin Journals</option>
                    <optgroup class="navbar-search-optgroup" label="Search across Oxford Academic">
<option class="navbar-search-filter-option at-navbar-search-filter-option" value="AcademicSubjects/MED00810">Child and Adolescent Psychiatry</option><option class="navbar-search-filter-option at-navbar-search-filter-option" value="Books">Books</option><option class="navbar-search-filter-option at-navbar-search-filter-option" value="Journals">Journals</option><option class="navbar-search-filter-option at-navbar-search-filter-option" value="Umbrella">Oxford Academic</option>                    </optgroup>
            </select>

        <label for="OUP_ArticleToolbarSearchBox-mobile-microsite-search-term" class="screenreader-text js-mobile-microsite-search-term-label">
            Mobile Enter search term
        </label>
        <input class="mobile-search-input mobile-microsite-search-term js-mobile-microsite-search-term at-microsite-search-term" type="text"
               maxlength="255" placeholder="Search" id="OUP_ArticleToolbarSearchBox-mobile-microsite-search-term">

        <a href="javascript:;" class="mobile-microsite-search-icon mobile-search-submit icon-menu_search">
            <span class="screenreader-text">Search</span>
        </a>

    </div>
</div>


        </div>
    </div>
   
    <div class="dropdown-panel mobile-nav-dropdown">

    </div>

</div>


<div class="navbar">
    <div class="center-inner-row">
        <nav class="navbar-menu">

        </nav>
            <div class="navbar-search-container js-navbar-search-container">
                <a href="javascript:;" class="navbar-search-close js_close-navsearch">Close</a>


<div class="navbar-search">
    <div class="microsite-search">
            <label for="OUP_ArticleToolbarSearchBox-navbar-search-filter" class="screenreader-text js-navbar-search-filter-label">
                Navbar Search Filter
            </label>
            <select class="navbar-search-filter js-navbar-search-filter at-navbar-search-filter" id="OUP_ArticleToolbarSearchBox-navbar-search-filter">

<option class="navbar-search-filter-option at-navbar-search-filter-option" value="">Schizophrenia Bulletin</option><option class="navbar-search-filter-option at-navbar-search-filter-option" value="Parent">Schizophrenia Bulletin Journals</option>
                    <optgroup class="navbar-search-optgroup" label="Search across Oxford Academic">
<option class="navbar-search-filter-option at-navbar-search-filter-option" value="AcademicSubjects/MED00810">Child and Adolescent Psychiatry</option><option class="navbar-search-filter-option at-navbar-search-filter-option" value="Books">Books</option><option class="navbar-search-filter-option at-navbar-search-filter-option" value="Journals">Journals</option><option class="navbar-search-filter-option at-navbar-search-filter-option" value="Umbrella">Oxford Academic</option>                    </optgroup>
            </select>

        <label for="OUP_ArticleToolbarSearchBox-microsite-search-term" class="screenreader-text js-microsite-search-term-label">
            Enter search term
        </label>
        <input class="navbar-search-input microsite-search-term js-microsite-search-term at-microsite-search-term" type="text"
               maxlength="255" placeholder="Search" id="OUP_ArticleToolbarSearchBox-microsite-search-term">

        <a href="javascript:;" class="microsite-search-icon navbar-search-submit icon-menu_search">
            <span class="screenreader-text">Search</span>
        </a>

    </div>
</div>

                <div class="navbar-search-advanced"><a href="/schizophreniabulletin/advanced-search" class="advanced-search js-advanced-search">Advanced Search</a></div>
            </div>            <div class="navbar-search-collapsed"><a href="javascript:;" class="icon-menu_search js_expand-navsearch"><span class="screenreader-text">Search Menu</span></a></div>
    </div>
</div>
<input id="hfEnableOupOnlineProductsFeatures" name="hfEnableOupOnlineProductsFeatures" type="hidden" value="True" />

<input id="routename" name="RouteName" type="hidden" value="schizophreniabulletin" /> 
    </div>

                    </div>
                </div>
            </div>
            <div id="ContentTab" class="content active">
    <div class="widget widget-ArticleFulltext widget-instance-OUP_Article_FullText_Widget">
        <div class="module-widget">
    <div class="widget-items" data-widgetname="ArticleFulltext">





<p class="chapter-para">I am a 55-year-old husband, father, friend, and professional philosopher. In 1992, as a graduate student at Cambridge University, a porter found me amongst the cows in the meadows of King’s College, after being there for 2 or 3 days. I was in bad physical shape, having eaten nothing, and apparently getting water from the river. He asked what I was doing. I replied: “I’m solving a problem about stochastic calculus.” This statement was true, but did not answer his question. He took me to the hospital, where I remained for some weeks.</p><p class="chapter-para">It wasn’t the first time that I was psychotic, but it was, maybe, the first time that anybody noticed, the first time that I was unable to hide it from others, and therefore from myself. What follows is an abbreviated account of how I learned—haltingly, with setbacks, over the years—to cope with chronic schizophrenia. There have been near-collapses, but I have managed to keep a job for 30 years. I have not, until recently, been open about my diagnosis (excepting my wife and a close friend).</p><p class="chapter-para">The most important part of my story is people. The reason that I am not in prison, homeless, or dead, is a few people who genuinely respect and care for me, and I them, not least through what some philosophers call “hermeneutical justice.” Without these people, there would be no “coping,” and the rest of what follows could never have happened. I will focus on just two of the symptoms that I experience, symptoms that have not been dislodged by medication (although medication can be helpful in other ways): auditory hallucination, and two recurring delusion-like experiences.</p><p class="chapter-para">I have hallucinated music since childhood. Voices came later, and visual hallucinations later still. Often the voices are distant, a conversation that does not involve me, and I can ignore it. Sometimes the voices are closer, and sometimes they speak to me. These voices are commonly of people I know, but sometimes they are strangers. Sometimes they are critical. Sometimes they comment on what is happening. Sometimes they blather. Occasionally they are encouraging.</p><p class="chapter-para">Musical hallucinations do not distress me. Voices are a different story—they are rarely intrinsically disturbing, but uncertainty about their origin is. For some time, it felt important to figure out whether the voices were coming from people who are physically present. My doctor called it “reality-check.” Sometimes reality check is easy, eg, if there is a voice whispering in my ear but nobody near my ear. But often reality check is very difficult. In a crowded place, hearing a conversation, does one ask people whether they just said anything? Does one snoop around to find the source of the talking? Does one stare at people’s mouths to see whether they are talking?</p><p class="chapter-para">The doctor was half-right: I felt better when I was sure about the origin of voices, and anxious when not (especially when they were directed at me). But sometimes it is awkward, difficult, or practically impossible, to gain that assurance. My frequent inability or unwillingness to do a reality check caused anxiety, which makes symptoms worse, and things can spiral out of control. Once I figured these things out, I typically avoided situations where it would be a problem. There are a lot of those situations, so this solution is not great.</p><p class="chapter-para">Some years ago, a funny situation changed my approach to hallucinations. The scene is a cold, dark, morning, in a coffee shop. There are no other customers. I order my coffee and pastry, sit down, and start working. Soon I hear a conversation. Normally I would have done my reality check, and doing so would have been easy (it’s a small shop), but I felt confident that nobody apart from the sole employee and myself were present, and the voices were not inherently disturbing, so I kept working. Then one of the voices said my name, directly to me. Hearing my name almost always gets my attention, and I turned around, although still expecting to see nothing, but there were two people behind me—real people!—and I knew one of them; he had recognized me and was saying hello.</p><p class="chapter-para">I suppose that sort of thing had happened before, but in that moment I realized something that I had not realized before: It is not important to <em>know</em> where the voices are coming from. It had just been demonstrated to me that prior to turning around I did <em>not</em> know their origin, and yet I was comfortable having taken on the “mere belief,” and as it turned out the <em>false</em> belief, that nobody was there. I realized in that moment that the comfort that came from successful reality-checks came not from knowledge or certainty, but from a <em>clear belief</em> about the voices. In this situation, that belief, even though it turned out to be false, was enough. And after I was forced to change my belief, it was still fine. I was able to turn back around and continue working, now believing that the voices were coming from people behind me. “And what if,” I thought, “those people quietly left, but I kept hearing the conversation, believing it to come from them?” Well, I’d probably eventually discover that they weren’t there, that the conversation was no longer real, and that would be fine too.</p><p class="chapter-para">As trivial as these events might seem, they were life-changing. A similar pattern has played out with other symptoms. Here are two examples.</p><p class="chapter-para">The first is close to “thought-broadcasting,” and for some time I did worry that others might hear my thoughts. I tried hard to think nice thoughts, or to think nothing. After extensive self-reflection, I realized that something slightly different is going on. I realized that it is difficult to tell the difference between speaking out loud and thinking. When I’m focused, I can tell the difference by paying careful attention to my body—especially my lips and throat—but one cannot always focus in that manner, and the resulting uncertainty about what has, or has not, been said out loud can kindle anxiety. Many of my conversations are laced with uncertainty about what I have said out loud, versus merely thought to myself.</p><p class="chapter-para">After I realized what is going on, I tried to avoid this uncertainty, either by trying not to think or say anything (which is difficult), or by frequently repeating myself (which is obnoxious). More recently, I’ve accepted that it rarely matters whether others have heard me. If I happen to mention (or merely to think?) that I’m allergic to eggplant, it matters very little whether you heard. So these days, most of the time, I just make my own determination about whether the other person heard, just as one might do after making an off-hand remark on the periphery of a conversation, and that determination is good enough. I don’t double-check, repeat myself, or ask whether you heard, unless it really matters. This habit produces some false positives and some false negatives. It turns out that most of the time, it just doesn’t matter.</p><p class="chapter-para">The second example concerns mirrors. It often seems to me that there are cameras, or persons, behind mirrors. I used to check mirrors (and still do sometimes), but I have come to realize and to accept that most of the time, it doesn’t matter. If there are voyeurs on the other side, that’s their wretched problem, not mine. For me, the path of least resistance is to allow that there probably is something on the other side. As long as I’m clear with myself, all is well.</p><p class="chapter-para">There is a common theme to these strategies. It’s definite belief, not certainty, that allows me to get along. It’s not that certainty, or something like it, <em>never</em> matters. If you are fixing dinner for me I’ll try to be clear about the eggplant allergy, and I might repeat myself. And as I do when I teach students, I’ll monitor you for a sign that you have heard and understood, and I might even ask you to confirm it. I might, in other words, be a little obnoxious about it, and I hope that you’ll be patient with me. But most of the time, just having a definite, if unconfirmed and possibly false, belief about the situation is fine. It allows one to get along.</p><p class="chapter-para">I think of this attitude as a kind of “intellectual humility” because although I <em>do care</em> about truth—and as a consequence of caring about truth, I do form beliefs about what is true—I no longer agonize about whether my judgments are wrong. For me, living relatively free from debilitating anxiety is incompatible with relentless pursuit of truth. Instead, I need clear beliefs and a willingness to change them when circumstances and evidence demand, without worrying about, or getting upset about, being wrong. This attitude has made life better and has made the “near-collapses” much rarer.</p>    <!-- /foreach in Model.Sections -->
    <div class="widget widget-ArticlePubStateInfo widget-instance-OUP_ArticlePubStateInfo">
         
    </div>



        <div class="article-metadata-standalone-panel clearfix"></div>

        
<div class="copyright copyright-statement">© The Author(s) 2024. Published by Oxford University Press on behalf of the Maryland Psychiatric Research Center. All rights reserved. For permissions, please email: journals.permissions@oup.com</div><div class="license"><div class="license-p">This article is published and distributed under the terms of the Oxford University Press, Standard Journals Publication Model (<a class="link link-uri openInAnotherWindow" href="https://academic.oup.com/pages/standard-publication-reuse-rights" target="_blank">https://academic.oup.com/pages/standard-publication-reuse-rights</a>)</div></div><!-- /foreach -->

        <span id="UserHasAccess" data-userHasAccess="True"></span>

    </div><!-- /.widget-items -->
</div><!-- /.module-widget -->





 
    </div>
    <div class="widget widget-SolrResourceMetadata widget-instance-OUP_Article_ResourceMetadata_Widget">
        
    <div class="article-metadata-panel solr-resource-metadata js-metadata-panel at-ContentMetadata">
                <div class="article-metadata-tocSections">
                    <div class="article-metadata-tocSections-title">Issue Section:</div>
                        <a href="/schizophreniabulletin/search-results?f_TocHeadingTitle=First+Person+Account">First Person Account</a>
                </div>

    </div>


 
    </div>
    <div class="widget widget-EditorInformation widget-instance-OUP_Article_EditorInformation_Widget">
        
 
    </div>

                <div id="ContentTabFilteredView"></div>
                <div class="downloadImagesppt js-download-images-ppt st-download-images-ppt">
                    <a id="lnkDownloadAllImages"
                       class="js-download-all-images-link btn"
                       href="/DownloadFile/DownloadImage.aspx?image=&amp;PPTtype=SlideSet&amp;ar=7517011&amp;xsltPath=~/UI/app/XSLT&amp;siteId=5240">Download all slides</a>
                </div>
                    <div class="widget widget-ArticleDataRepositories widget-instance-Article_DryadLink">
         
    </div>

                <div class="comments">
    <div class="widget widget-UserCommentBody widget-instance-UserCommentBody_Article">
        

 
    </div>
    <div class="widget widget-UserComment widget-instance-OUP_UserComment_Article">
         
    </div>

                </div>
            </div>
        </div>

    </div>
</div>
<div id="Sidebar" class="page-column page-column--right">
    <div class="widget widget-AdBlock widget-instance-ArticlePageTopSidebar">
        <div class="js-adBlock-parent-wrap adblock-parent-wrap">
    <div class="adBlockMainBodyTop-wrap js-adBlockMainBodyTop hide">

        <div id="adBlockMainBodyTop" class="js-adblock at-adblock">
            <script>
            googletag.cmd.push(function () {
                googletag.display('adBlockMainBodyTop');
            });
            </script>
        </div>
            <div class="advertisement-text at-adblock js-adblock-advertisement-text hide">Advertisement</div>
            </div>
</div> 
    </div>
<div class="widget widget-dynamic " data-count="1"> 

    <div class="widget-dynamic-inner-wrap">

<div class="widget widget-dynamic " data-count="8"> 

    <div class="widget-dynamic-inner-wrap">

    <div class="widget widget-ArticleLevelMetrics widget-instance-Article_RightRailB0Article_RightRail_ArticleLevelMetrics">
        



        <div class="artmet-wrapper horizontal-artmet">

<div class="contentmet-border">
    <div class="contentmet-wrapper horizontal-contentmet">
            <div class="contentmet-citations contentmet-badges-wrap js-contentmet-citations hide">
                <h3 class="contentmet-text">Citations</h3>
                    <div class="contentmet-item contentmet-dimensions">
    <div class="widget widget-DimensionsBadge widget-instance-ArticleLevelMetrics_DimensionsBadge">
        <span class="__dimensions_badge_embed__" 
      data-doi="10.1093/schbul/sbad173" 
      data-legend="never" 
      data-style="small_circle" 
      data-hide-zero-citations="false"></span>
<script async src="https://badge.dimensions.ai/badge.js" charset="utf-8"></script>
 
    </div>

                    </div>
                            </div>
                    <div class="contentmet-views contentmet-badges-wrap js-contentmet-views">
                <h3 class="contentmet-text">Views</h3>
                <div class="contentmet-item circle-text circle-text-views">
                    <div class="contentmet-number">105</div>
                </div>
            </div>
                    <div class="contentmet-item contentmet-badges-wrap">
                    <h3 class="contentmet-text">Altmetric</h3>
                    <div class="contentmet-item contentmet-altmetric">
    <div class="widget widget-AltmetricLink widget-instance-ArticleLevelMetrics_AltmetricLinkSummary">
            <!-- Altmetrics -->
    <div id="altmetricEmbedId"
         runat="server"
         class='altmetric-embed'
         data-badge-type="donut"
         data-hide-no-mentions="false"
         data-doi="10.1093/schbul/sbad173"

        

        ></div>
         <script type='text/javascript' src='https://d1bxh8uas1mnw7.cloudfront.net/assets/embed.js'></script>
 
    </div>

                    </div>


            </div>
                    <div class="contentmet-modal-trigger-wrap clearfix">
                <a href="javascript:;" class="artmet-modal-trigger js-artmet-modal-trigger at-alm-metrics-modal-trigger" data-resource-id="7517011" data-resource-type="3">
                    <img class="contentmet-info-icon" src="//oup.silverchair-cdn.com/UI/app/svg/i.svg" alt="Information">
                    <span class="contentmet-more-info">More metrics information</span>
                </a>
            </div>
    </div>
</div>

                <div class="artmet-modal js-artmet-modal" id="MetricsModal">
                    <div class="artmet-modal-contents js-metric-modal-contents at-alm-modal-contents">




    <div class="artmet-full-wrap clearfix js-metric-full-wrap hide">
            <div class="widget-title-1 artmet-widget-title-1">Metrics</div>

            <div class="artmet-item artmet-views-wrap">
                <div class="artmet-views clearfix">
                    <div class="artmet-total-views">
                        <span class="artmet-text">Total Views</span>
                        <span class="artmet-number">105</span>
                    </div>
                    <div class="artmet-breakdown-views-wrap">
                        <div class="artmet-breakdown-view breakdown-border">
                            <span class="artmet-number">93</span>
                            <span class="artmet-text">Pageviews</span>
                        </div>
                        <div class="artmet-breakdown-view">
                            <span class="artmet-number">12</span>
                            <span class="artmet-text">PDF Downloads</span>
                        </div>
                                            </div>
                </div>
                <div class="artmet-views-since">Since 1/1/2024</div>
            </div>

            <script>
                var ChartistData = ChartistData || {};

                ChartistData.data = {
                    labels: ['Jan 2024'],
                    series: [[105]]
                };
            </script>
            <div class="artmet-item artmet-chart">
                <div class="ct-chart ct-octave js-ct-chart"></div>

                <div class="artmet-table">
                    <table>
                        <thead>
                            <tr>
                                <th>Month:</th>
                                <th>Total Views:</th>
                            </tr>
                        </thead>
                        <tbody>
                                <tr>
                                    <td>January 2024</td>
                                    <td>105</td>
                                </tr>
                        </tbody>
                    </table>
                </div>
            </div>

<div class="artmet-stats-wrap clearfix">
        <div class="artmet-item artmet-citations hide">
            <div class="widget-title-2 artmet-widget-title-2">Citations</div>
                <div class="artmet-badges-wrap artmet-dimensions">
    <div class="widget widget-DimensionsBadge widget-instance-ArticleLevelMetrics_DimensionsBadgeDetails">
        <span class="__dimensions_badge_embed__" 
      data-doi="10.1093/schbul/sbad173" 
      data-legend="always" 
      data-style="" 
      data-hide-zero-citations="false"></span>
<script async src="https://badge.dimensions.ai/badge.js" charset="utf-8"></script>
 
    </div>

                    <span class="artmet-dimensions-text">Powered by Dimensions</span>
                </div>
                    </div>


        <div class="artmet-item artmet-altmetric js-show-if-unknown">
            <div class="widget-title-2 artmet-widget-title-2">Altmetrics</div>
            <div class="artmet-badges-wrap js-artmet-badges">

    <div class="widget widget-AltmetricLink widget-instance-ArticleLevelMetrics_AltmetricLinkDetails">
            <!-- Altmetrics -->
    <div id="altmetricEmbedId"
         runat="server"
         class='altmetric-embed'
         data-badge-type="donut"
         data-hide-no-mentions="false"
         data-doi="10.1093/schbul/sbad173"

        
            
                data-badge-details = "right"
            

        ></div>
         <script type='text/javascript' src='https://d1bxh8uas1mnw7.cloudfront.net/assets/embed.js'></script>
 
    </div>
            </div>
        </div>

</div>


    </div>
                        <a class="artmet-close-modal js-artmet-close-modal">&#215;</a>
                    </div>
                </div>
        </div>

 
    </div>
    <div class="widget widget-Alerts widget-instance-Article_RightRailB0Article_RightRail_Alerts">
        
    <div class="alertsWidget">
        <h3>Email alerts</h3>
                <div class="userAlert alertType-1">
                    <a href="javascript:;" class="js-user-alert" role="button"
                       data-userLoggedIn="False"
                       data-alertType="1" href="javascript:;">Article activity alert</a>
                </div>
                <div class="userAlert alertType-3">
                    <a href="javascript:;" class="js-user-alert" role="button"
                       data-userLoggedIn="False"
                       data-alertType="3" href="javascript:;">Advance article alerts</a>
                </div>
                <div class="userAlert alertType-5">
                    <a href="javascript:;" class="js-user-alert" role="button"
                       data-userLoggedIn="False"
                       data-alertType="5" href="javascript:;">New issue alert</a>
                </div>
                    <div class="userAlert alertType-MarketingLink">
                <a href="javascript:;" class="js-user-alert" role="button"
                   data-userLoggedIn="False"
                   data-additionalUrl="/my-account/communication-preferences" href="javascript:;">Receive exclusive offers and updates from Oxford Academic</a>
            </div>
        <div class="userAlertSignUpModal reveal-modal small" data-reveal>
            <div class="userAlertSignUp at-userAlertSignUp"></div>
            <a href="javascript:;" role="button" aria-label="Close" class="close-reveal-modal" href="javascript:;">
                <i class="icon-general-close"></i>
            </a>
        </div>
    </div>
 
    </div>
    <div class="widget widget-TrendMD widget-instance-Article_RightRailB0trendmd">
        <script type='text/javascript' defer src='//js.trendmd.com/trendmd.min.js' data-trendmdconfig='{"journal_id":"62162","element":"#trendmd-suggestions"}' class='optanon-category-C0002'></script>

<div id="trendmd-suggestions"></div> 
    </div>
    <div class="widget widget-ArticleCitedBy widget-instance-Article_RightRailB0Article_RightRail_ArticleCitedBy">
        <div class="rail-widget_wrap vt-articles-cited-by">
    <h3 class="article-cited-title">Citing articles via</h3>
    <div class="widget-links_wrap">
                    <div class="article-cited-link-wrap google-scholar-url">
                <a href="http://scholar.google.com/scholar?q=link:https%3A%2F%2Facademic.oup.com%2Fschizophreniabulletin%2Fadvance-article%2Fdoi%2F10.1093%2Fschbul%2Fsbad173%2F7517011" target="_blank">Google Scholar</a>
            </div>
            </div>
</div> 
    </div>
    <div class="widget widget-ArticleListNewAndPopular widget-instance-Article_RightRailB0Article_RightRail_ArticleNewPopularCombined">
            <ul class="articleListNewAndPopularCombinedView">
            <li>
                <h3>
                    <a href="javascript:;" class="articleListNewAndPopular-mode active" data-mode="MostRecent">Latest</a>
                </h3>
            </li>
            <li>
                <h3>
                    <a href="javascript:;" class="articleListNewAndPopular-mode " data-mode="MostRead">Most Read</a>
                </h3>
            </li>
            <li>
                <h3>
                    <a href="javascript:;" class="articleListNewAndPopular-mode " data-mode="MostCited">Most Cited</a>
                </h3>
            </li>
    </ul>
        <section class="articleListNewPopContent articleListNewAndPopular-ContentView-MostRecent hasContent">




<div id="newPopularList-Article_RightRailB0Article_RightRail_ArticleNewPopularCombined" class="fb">

    



<div class="widget-layout widget-layout--vert ">
            <div class="widget-columns widget-col-1">
                    <div class="col">

<div class="widget-dynamic-entry journalArticleItem at-widget-entry-SCL">
    <span class="hfDoi" data-attribute="10.1093/schbul/sbad185" aria-hidden="true"></span>


    


<a class="journal-link" href="/schizophreniabulletin/advance-article/doi/10.1093/schbul/sbad185/7560288?searchresult=1">
        <div class="widget-dynamic-journal-title">
Pseudoneurotic Symptoms in the Schizophrenia Spectrum: A Longitudinal Study of Their Relation to Psychopathology and Clinical Outcomes    </div>
</a>

 





<div class="widget-dynamic-journal-image-synopsis">
        <div class="widget-dynamic-journal-synopsis">
            
        </div>
</div>

</div>
<div class="widget-dynamic-entry journalArticleItem at-widget-entry-SCL">
    <span class="hfDoi" data-attribute="10.1093/schbul/sbad181" aria-hidden="true"></span>


    


<a class="journal-link" href="/schizophreniabulletin/advance-article/doi/10.1093/schbul/sbad181/7517012?searchresult=1">
        <div class="widget-dynamic-journal-title">
Coping With the Inner Turbulence    </div>
</a>

 





<div class="widget-dynamic-journal-image-synopsis">
        <div class="widget-dynamic-journal-synopsis">
            
        </div>
</div>

</div>
<div class="widget-dynamic-entry journalArticleItem at-widget-entry-SCL">
    <span class="hfDoi" data-attribute="10.1093/schbul/sbad173" aria-hidden="true"></span>


    


<a class="journal-link" href="/schizophreniabulletin/advance-article/doi/10.1093/schbul/sbad173/7517011?searchresult=1">
        <div class="widget-dynamic-journal-title">
From Vexing Uncertainty to Intellectual Humility    </div>
</a>

 





<div class="widget-dynamic-journal-image-synopsis">
        <div class="widget-dynamic-journal-synopsis">
            
        </div>
</div>

</div>
<div class="widget-dynamic-entry journalArticleItem at-widget-entry-SCL">
    <span class="hfDoi" data-attribute="10.1093/schbul/sbad176" aria-hidden="true"></span>


    


<a class="journal-link" href="/schizophreniabulletin/advance-article/doi/10.1093/schbul/sbad176/7516965?searchresult=1">
        <div class="widget-dynamic-journal-title">
Hypothalamic Subunit Volumes in Schizophrenia and Bipolar Spectrum Disorders    </div>
</a>

 





<div class="widget-dynamic-journal-image-synopsis">
        <div class="widget-dynamic-journal-synopsis">
            
        </div>
</div>

</div>
<div class="widget-dynamic-entry journalArticleItem at-widget-entry-SCL">
    <span class="hfDoi" data-attribute="10.1093/schbul/sbad178" aria-hidden="true"></span>


    


<a class="journal-link" href="/schizophreniabulletin/advance-article/doi/10.1093/schbul/sbad178/7504704?searchresult=1">
        <div class="widget-dynamic-journal-title">
Development and Evaluation of a Cognitive Battery for People With Schizophrenia in Ethiopia    </div>
</a>

 





<div class="widget-dynamic-journal-image-synopsis">
        <div class="widget-dynamic-journal-synopsis">
            
        </div>
</div>

</div>
                    </div>
            </div>

</div></div>
        </section>
        <section class="articleListNewPopContent articleListNewAndPopular-ContentView-MostRead hide">
        </section>
        <section class="articleListNewPopContent articleListNewAndPopular-ContentView-MostCited hide">
        </section>
 
    </div>
    <div class="widget widget-RelatedTaxonomies widget-instance-Article_RightRailB0Article_RightRail_RelatedTaxonomies">
            <div class="widget-related-taxonomies-wrap vt-related-taxonomies">
        <div class="widget-related-taxonomies_title">More from Oxford Academic</div>
            <div class="widget-related-taxonomies">
                <a id="more-from-oa-AcademicSubjects_MED00810" class="related-taxonomies-link" href="/search-results?tax=AcademicSubjects/MED00810">Child and Adolescent Psychiatry</a>
            </div>
            <div class="widget-related-taxonomies">
                <a id="more-from-oa-AcademicSubjects_MED00010" class="related-taxonomies-link" href="/search-results?tax=AcademicSubjects/MED00010">Medicine and Health</a>
            </div>
            <div class="widget-related-taxonomies">
                <a id="more-from-oa-AcademicSubjects_MED00800" class="related-taxonomies-link" href="/search-results?tax=AcademicSubjects/MED00800">Psychiatry</a>
            </div>

            <div class="widget-related-taxonomies">
                <a id="more-from-oa-format-Books" class="related-taxonomies-link" href="/books">Books</a>
            </div>
            <div class="widget-related-taxonomies">
                <a id="more-from-oa-format-Journals" class="related-taxonomies-link" href="/journals">Journals</a>
            </div>
    </div>
 
    </div>

    </div>

</div>
    </div>

</div>    <div class="widget widget-AdBlock widget-instance-ArticlePageTopMainBodyBottom">
        <div class="js-adBlock-parent-wrap adblock-parent-wrap">
    <div class="adBlockMainBodyBottom-wrap js-adBlockMainBodyBottom hide">

        <div id="adBlockMainBodyBottom" class="js-adblock at-adblock">
            <script>
            googletag.cmd.push(function () {
                googletag.display('adBlockMainBodyBottom');
            });
            </script>
        </div>
            <div class="advertisement-text at-adblock js-adblock-advertisement-text hide">Advertisement</div>
            </div>
</div> 
    </div>

</div>

</div>
<input id="hfArticleTitle" name="hfArticleTitle" type="hidden" value="From Vexing Uncertainty to Intellectual Humility | Schizophrenia Bulletin | Oxford Academic" />
<input id="hfLeftNavStickyOffset" name="hfLeftNavStickyOffset" type="hidden" value="29" />
<input id="hfAreOpFeaturesEnabled" name="hfAreOpFeaturesEnabled" type="hidden" value="True" />



                </div><!-- /.center-inner-row no-overflow -->
            </section>
    </div>

        <div class="mobile-mask">
        </div>

        <section class="footer_wrap vt-site-footer">
            


    <div class="ad-banner-footer sticky-footer-ad js-sticky-footer-ad">
    <div class="widget widget-AdBlock widget-instance-StickyFooterAd">
        <div class="js-adBlock-parent-wrap adblock-parent-wrap">
    <div class="adBlockStickyFooter-wrap js-adBlockStickyFooter hide">

        <div id="adBlockStickyFooter" class="js-adblock at-adblock">
            <script>
            googletag.cmd.push(function () {
                googletag.display('adBlockStickyFooter');
            });
            </script>
        </div>
            <div class="advertisement-text at-adblock js-adblock-advertisement-text hide">Advertisement</div>
                    <button type="button" class="close-footer-ad js-close-footer-ad">
                <span class="screenreader-text">close advertisement</span>
            </button>
    </div>
</div> 
    </div>

    </div>

    <div class="widget widget-SitePageFooter widget-instance-SitePageFooter">
            <div class="ad-banner ad-banner-footer">
    <div class="widget widget-AdBlock widget-instance-FooterAd">
        <div class="js-adBlock-parent-wrap adblock-parent-wrap">
    <div class="adBlockFooter-wrap js-adBlockFooter hide">

        <div id="adBlockFooter" class="js-adblock at-adblock">
            <script>
            googletag.cmd.push(function () {
                googletag.display('adBlockFooter');
            });
            </script>
        </div>
            <div class="advertisement-text at-adblock js-adblock-advertisement-text hide">Advertisement</div>
            </div>
</div> 
    </div>

    </div>

<div class="journal-footer journal-bg">
    <div class="center-inner-row">

<div class="journal-footer-menu">

    <ul>
<li class="link-1">
    <a href="/schizophreniabulletin/pages/About">About Schizophrenia Bulletin</a>
</li> <li class="link-2">
    <a href="/schizophreniabulletin/pages/Editorial_Board">Editorial Board</a>
</li> <li class="link-3">
    <a href="/schizophreniabulletin/pages/General_Instructions">Author Guidelines</a>
</li> <li class="link-4">
    <a href="http://medschool.umaryland.edu/">Contact Us</a>
</li> <li class="link-5">
    <a href="https://www.facebook.com/OUPAcademic">Facebook</a>
</li> </ul><ul><li class="link-1">
    <a href="https://twitter.com/OxfordJournals">Twitter</a>
</li> <li class="link-2">
    <a href="/schizophreniabulletin/subscribe">Purchase</a>
</li> <li class="link-3">
    <a href="http://www.oxfordjournals.org/en/library-recommendation-form.html">Recommend to your Library</a>
</li> <li class="link-4">
    <a href="https://academic.oup.com/advertising-and-corporate-services/pages/schizophreniabulletin-media-kit">Advertising and Corporate Services</a>
</li> <li class="link-5">
    <a href="http://medicine-and-health-careernetwork.oxfordjournals.org">Journals Career Network</a>
</li> 
    </ul>


</div><!-- /.journal-footer-menu -->
        <div class="journal-footer-affiliations">
            <!-- <h3>Affiliations</h3> -->
<a href="http://schizophreniaresearchsociety.org/" target="">
    <img id="footer-logo-SchizophreniaInternationalResearchSociety" class="journal-footer-affiliations-logo" src="//oup.silverchair-cdn.com/data/SiteBuilderAssets/Live/Images/schizophreniabulletin/f1-591729497.png" alt="Schizophrenia International Research Society" />
</a>                        </div><!-- /.journal-footer-affiliations -->

        <div class="journal-footer-colophon">
            <ul>
<li>Online ISSN 1745-1701</li>                <li>Print ISSN 0586-7614</li>                <li>Copyright &#169; 2024 Maryland Psychiatric Research Center and Oxford University Press</li>
            </ul>
        </div><!-- /.journal-footer-colophon -->


    </div><!-- /.center-inner-row -->
</div><!-- /.journal-footer --> 
    </div>

    <div class="oup-footer">
        <div class="center-inner-row">
    <div class="widget widget-SelfServeContent widget-instance-OupUmbrellaFooterSelfServe">
        



    <input type="hidden" class="SelfServeContentId" value="GlobalFooter_Links" />
    <input type="hidden" class="SelfServeVersionId" value="0" />

<div class="oup-footer-row journal-links">
<div class="global-footer selfservelinks">
<ul>
    <li><a href="/pages/about-oxford-academic">About Oxford Academic</a></li>
    <li><a href="/pages/about-oxford-academic/publish-journals-with-us">Publish journals with us</a></li>
    <li><a href="/pages/about-oxford-academic/our-university-press-partners">University press partners</a></li>
    <li><a href="/pages/what-we-publish">What we publish</a></li>
    <li><a href="/pages/new-features">New features</a>&nbsp;</li>
</ul>
</div>
<div class="global-footer selfservelinks">
<ul>
    <li><a href="/pages/authoring">Authoring</a></li>
    <li><a href="/pages/open-research/open-access">Open access</a></li>
    <li><a href="/pages/purchasing">Purchasing</a></li>
    <li><a href="/pages/institutional-account-management">Institutional account management</a></li>
    <li><a href="https://academic.oup.com/pages/purchasing/rights-and-permissions">Rights and permissions</a></li>
</ul>
</div>
<div class="global-footer selfservelinks">
<ul>
    <li><a href="/pages/get-help-with-access">Get help with access</a></li>
    <li><a href="/pages/about-oxford-academic/accessibility">Accessibility</a></li>
    <li><a href="/pages/contact-us">Contact us</a></li>
    <li><a href="/pages/advertising">Advertising</a></li>
    <li><a href="/pages/media-enquiries">Media enquiries</a></li>
</ul>
</div>
<div class="global-footer selfservelinks global-footer-external">
<ul>
    <li><a href="https://global.oup.com/">Oxford University Press</a></li>
    <li><a href="https://www.mynewsdesk.com/uk/oxford-university-press">News</a></li>
    <li><a href="https://languages.oup.com/">Oxford Languages</a></li>
    <li><a href="https://www.ox.ac.uk/">University of Oxford</a></li>
</ul>
</div>
<div class="OUP-mission">
<p>Oxford University Press is a department of the University of Oxford. It furthers the University's objective of excellence in research, scholarship, and education by publishing worldwide</p>
<img class="journal-footer-logo" src="//oup.silverchair-cdn.com/UI/app/svg/umbrella/oup-logo.svg" alt="Oxford University Press" width="150" height="42" />
</div>
</div>
<div class="oup-footer-row">
<div class="oup-footer-row-links">
<ul>
    <li>Copyright © 2024 Oxford University Press</li>
    <li><button id="Change-Preferences" type="button" onclick="window.top.document.dispatchEvent(new Event('changeConsent'))">Cookie settings</button></li>
    <li><a href="https://global.oup.com/cookiepolicy">Cookie policy</a></li>
    <li><a href="https://global.oup.com/privacy">Privacy policy</a></li>
    <li><a href="/pages/legal-and-policy/legal-notice">Legal notice</a></li>
</ul>
</div>
</div>
<style type="text/css">
    /* Issue right column fix for tablet/mobile */
    @media (max-width: 1200px) {
    .pg_issue .widget-instance-OUP_Issue {
    width: 100%;
    }
    }
    .sf-facet-list .sf-facet label,
    .sf-facet-list .taxonomy-label-wrap label {
    font-size: 0.9375rem;
    }
    .issue-pagination-wrap .pagination-container {
    float: right;
    }
</style>



 
    </div>

        </div>
    </div>
    <div class="ss-ui-only">

    </div>

        </section>
</div>






    



    <div class="widget widget-SiteWideModals widget-instance-SiteWideModals">
        <div id="revealModal" class="reveal-modal" data-reveal>
    <div id="revealContent"></div>
    <a class="close-reveal-modal" href="javascript:;"><i class="icon-general-close"></i><span class="screenreader-text">Close</span></a>
</div>

<div id="globalModalContainer" class="modal-global-container">
    <div id="globalModalContent">
        <div class="js-globalModalPlaceholder"></div>
    </div>
    <a class="close-modal js-close-modal" href="javascript:;"><i class="icon-general-close"><span class="screenreader-text">Close</span></i></a>
</div>
<div id="globalModalOverlay" class="modal-overlay js-modal-overlay"></div>

<div id="NeedSubscription" class="reveal-modal small" data-reveal>
    <div class="subscription-needed">
        <h5>This Feature Is Available To Subscribers Only</h5>
        <p><a href="/sign-in">Sign In</a> or <a href="/my-account/register?siteId=5240&amp;returnUrl=%2fschizophreniabulletin%2fadvance-article%2fdoi%2f10.1093%2fschbul%2fsbad173%2f7517011">Create an Account</a></p>
    </div>
    <a class="close-reveal-modal" href="javascript:;"><i class="icon-general-close"></i><span class="screenreader-text">Close</span></a>
</div>

<div id="noAccessReveal" class="reveal-modal tiny" data-reveal>
    <p>This PDF is available to Subscribers Only</p>
    <a class="hide-for-article-page" id="articleLinkToPurchase" data-reveal><span>View Article Abstract & Purchase Options</span></a>
    <div class="issue-purchase-modal">
        <p>For full access to this pdf, sign in to an existing account, or purchase an annual subscription.</p>
    </div>
    <a class="close-reveal-modal" href="javascript:;"><i class="icon-general-close"></i><span class="screenreader-text">Close</span></a>
</div>
 
    </div>

    



<script type="text/javascript">
    MathJax = {
        tex: {
            inlineMath: [['|$', '$|'], ['\\(', '\\)']]
        }
    };
</script>
<script id="MathJax-script" async src="//cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js"></script>


    <!-- CookiePro Default Categories -->
    <!-- When the Cookie Compliance code loads, if cookies for the associated group have consent...
         it will dynamically change the tag to: script type=text/JavaScript...
         the code inside the tags will then be recognized and run as normal. -->














    <script>
        var NTPT_PGEXTRA = 'event_type=full-text\u0026discipline_ot_level_1=Medicine and Health\u0026discipline_ot_level_2=Psychiatry\u0026supplier_tag=SC_Journals\u0026object_type=Article\u0026taxonomy=taxId%3a39%7ctaxLabel%3aAcademicSubjects%7cnodeId%3aMED00810%7cnodeLabel%3aChild+and+Adolescent+Psychiatry%7cnodeLevel%3a3\u0026siteid=schbul\u0026authzrequired=false\u0026doi=10.1093/schbul/sbad173';
    </script>
    <!-- Copyright 2001-2010, IBM Corporation All rights reserved. -->
    <script type="text/javascript" src="//ouptag.scholarlyiq.com/ntpagetag.js" class="optanon-category-C0002"></script>
        <noscript>
            <img src="//ouptag.scholarlyiq.com/ntpagetag.gif?js=0" height="1" width="1" border="0" hspace="0" vspace="0" alt="Scholarly IQ" />
        </noscript>










    <script type="text/javascript" src="//oup.silverchair-cdn.com/Themes/Client/app/jsdist/v-638411137930894175/site.min.js"></script>


    
    
    
    <script type="text/javascript" src="https://cdn.jsdelivr.net/chartist.js/latest/chartist.min.js"></script>

    <script type="text/javascript" src="//oup.silverchair-cdn.com/Themes/Client/app/jsdist/v-638411137794541408/article-page.min.js"></script>



    

    
    




    <div class="ad-banner js-ad-riser ad-banner-riser">

    <div class="widget widget-AdBlock widget-instance-RiserAd">
            

 
    </div>

    </div>



            
    <div class="end-of-page-js"></div>
</body>
</html>`

var testC = `<!DOCTYPE html>
<html lang="en">
  <head>
  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1">

  <title>the rust project has a burnout problem</title>
  <meta name="description" content="">
  <meta name="author" content="jyn">
  <meta name="og:image" content="/assets/burned%20out%20rust%20club.png">

  <link rel="canonical" href="https://jyn.dev/2024/01/16/the-rust-project-has-a-burnout-problem.html">
  <link rel="alternate" type="application/atom+xml" title="the website of jyn" href="/feed.xml">
  <link rel="icon" type="image/x-icon" href="/assets/favicon.jpg">

</head>

<body>
    <link rel="stylesheet" href="/assets/main.css">

    <header class="site-header" role="banner">

  <div class="wrapper">
    
    
    <a class="site-title" href="/">the website of jyn</a>

    
      <nav class="site-nav">
        <input type="checkbox" id="nav-trigger" class="nav-trigger" />
        <label for="nav-trigger">
          <span class="menu-icon">
            <svg viewBox="0 0 18 15" width="18px" height="15px">
              <path fill="#424242" d="M18,1.484c0,0.82-0.665,1.484-1.484,1.484H1.484C0.665,2.969,0,2.304,0,1.484l0,0C0,0.665,0.665,0,1.484,0 h15.031C17.335,0,18,0.665,18,1.484L18,1.484z"/>
              <path fill="#424242" d="M18,7.516C18,8.335,17.335,9,16.516,9H1.484C0.665,9,0,8.335,0,7.516l0,0c0-0.82,0.665-1.484,1.484-1.484 h15.031C17.335,6.031,18,6.696,18,7.516L18,7.516z"/>
              <path fill="#424242" d="M18,13.516C18,14.335,17.335,15,16.516,15H1.484C0.665,15,0,14.335,0,13.516l0,0 c0-0.82,0.665-1.484,1.484-1.484h15.031C17.335,12.031,18,12.696,18,13.516L18,13.516z"/>
            </svg>
          </span>
        </label>

        <div class="trigger">
          
            
            
          
            
            
            <a class="page-link" href="/about/">about</a>
            
          
            
            
          
            
            
            <a class="page-link" href="/talks/">talks</a>
            
          
            
            
          
            
            
          
        </div>
      </nav>
    
  </div>
</header>


    <main class="page-content" aria-label="Content">
      <div class="wrapper">
        <article class="post" itemscope itemtype="http://schema.org/BlogPosting">

  <header class="post-header">
    <h1 class="post-title" itemprop="name headline">the rust project has a burnout problem</h1>
    <p class="post-meta">
      <time datetime="2024-01-16T00:00:00-05:00" itemprop="datePublished">
        
        Jan 16, 2024
      </time>
      
        • <span itemprop="author" itemscope itemtype="http://schema.org/Person"><span itemprop="name">jyn</span></span>
      
      • audience - developers
    </p>
  </header>

  <div class="post-content" itemprop="articleBody">
    <p><img src="/assets/burned%20out%20rust%20club.png" alt="a melting, smiling, ferris. it's surrounded by the cursive text &quot;burned out rust kid club&quot;." /></p>

<p>the number of people who have left the rust project due to burnout is shockingly high. the number of people in the project who are close to burnout is also shockingly high.</p>

<p>this post is about myself, but it’s not just about myself. i’m not going to name names because either you know what i’m talking about, in which case you know <em>at least</em> five people matching this description, or you don’t, in which case sorry but you’re not the target audience. consider, though, that the project has been around for 15 years, and compare that to the average time a maintainer has been active …</p>

<h2 id="what-does-this-look-like">what does this look like</h2>

<p>(i apologize in advance if this story does not match your experience; hopefully the suggestions on what to do about burnout will still be helpful to you.)</p>

<p>the pattern usually goes something like this:</p>
<ul>
  <li>you want to work on rust. you go to look at the issue tracker. you find something <em>you</em> care about, since the easy/mentored issues are taken. it’s hard to find a mentor because all the experienced people are overworked and burned out, so you end up doing a lot of the work independently.</li>
</ul>

<p>guess what you’ve already learned at this point: work in this project doesn’t happen unless <em>you personally</em> drive it forward. that issue you fixed was opened for years; the majority of issues you will work on as you start will have been open for months.</p>

<ul>
  <li>you become a more active contributor. the existing maintainer is too burned out to do regular triage, so you end up going through the issue backlog (usually, you’re the first person to have done so in years). this reinforces the belief work doesn’t happen unless <em>you</em> do it <em>personally</em>.</li>
  <li>the existing maintainer recognizes your work and turns over a lot of the responsibilities to you, especially reviews. new contributors make PRs. they make silly simple mistakes due to lack of experience; you point them out and they get fixed. this can be fun, for a time. what it’s teaching you is that <em>you personally</em> are responsible for catching mistakes.</li>
  <li>you get tired. you’ve been doing this for a while. people keep making the same mistakes, and you’re afraid to trust other reviewers; perhaps you’re the <em>only</em> reviewer, or other reviewers have let things slip before and you don’t trust their judgement as much as you used to. perhaps you’re assigned too many PRs and you can’t keep up. you haven’t worked on the things you <em>want</em> to work on in weeks, and no one else is working on them because you said you were going to (“they won’t happen unless <em>you do them personally</em>”, a voice says). you want a break, but you have a voice in the back of your head: “the project would be worse without you”.</li>
</ul>

<p>i’m going to stop here; i think everyone gets the idea.</p>

<h2 id="what-can-i-do-about-it">what can i do about it</h2>

<p>“it won’t get done if i don’t do it” and “i need to review everything or stuff will slip through” is exactly the mindset of my own burnout from rust. it doesn’t matter if it’s true, it will cause you pain. if the project cannot survive without <em>you personally</em> putting in unpaid overtime, perhaps it does not deserve to survive.</p>

<p>if you are paid to work on rust, you likely started as an unpaid contributor and got the job later. <em>treat it like a job now</em>. do not work overtime; do not volunteer at every turn; do not work on things far outside your job description.</p>

<p>the best way to help the project is to keep contributing for it for years. to do that, you have to avoid burning out, which means you have to <em>treat yourself well</em>.</p>

<h2 id="what-can-team-leads-do-about-it">what can team leads do about it</h2>

<p>have documentation for “what to do about burnout”; give it just as much priority as technical issues or moderation conflicts.</p>

<p>rotate responsibilities. don’t have the same person assigned to the majority of PRs. if they review other people’s PRs unsolicited, talk to them 1-1 about why they feel the need to do so. if someone is assigned to the review queue and never reviews PRs, talk to them; take them off the queue; give them a vacation or different responsibilities as appropriate.</p>

<p>ask people why they leave. i know at least one person whose burnout story does not match the one in this post. i am sure there are others. you cannot solve a problem if you don’t understand what causes it.</p>

<p><em>take these problems seriously</em>. <a href="https://jyn.dev/2023/12/04/How-to-maintain-an-open-source-project.html">prioritize growing the team and creating a healthy environment over solving technical issues</a>. <strong>the issues will still be there in a few months; your people may not be</strong>.</p>

<h2 id="what-can-the-rust-project-do-about-it">what can the rust project do about it</h2>

<p>one thing bothering me as i wrote this post is how much of this still falls on individuals within the project. i don’t think this is an individual problem; i think it is a cultural, organizational, and resource problem. i may write more about this once i have concrete ideas about what the project could do.</p>

<h2 id="be-well-be-kind-to-each-other-i-love-you">be well. be kind to each other. i love you.</h2>

<p>remember:</p>

<blockquote class="twitter-tweet"><p lang="en" dir="ltr">EMPATHY WITHOUT BOUNDARIES IS SELF DESTRUCTION <a href="https://t.co/HbBwEj4hc3">pic.twitter.com/HbBwEj4hc3</a></p>&mdash; 𖤐ARCH BUDZAR𖤐 (@ArchBudzar) <a href="https://twitter.com/ArchBudzar/status/1313572660048269315?ref_src=twsrc%5Etfw">October 6, 2020</a></blockquote>
<script async="" src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>

<h3 id="acknowledgements">acknowledgements</h3>

<p>thank you <strong>@QuietMisdreavus</strong> for the <em>burned out rust kid club</em> art.</p>

<p>thank you <strong>@Gankra</strong>, <strong>@QuietMisdreavus</strong>, <strong>@alercah</strong>, <strong>@ManishEarth</strong>, <strong>@estebank</strong>, <strong>@workingjubilee</strong> and <strong>@yaahc</strong> for discussion and feedback on early drafts of this post. any errors are my own.</p>

  </div>

  
</article>

      </div>
    </main>

    <footer class="site-footer">

  <div class="wrapper">

    <h2 class="footer-heading">the website of jyn</h2>

    <div class="footer-col-wrapper">
      <div class="footer-col footer-col-1">
        <ul class="contact-list">
          <li>
            
              jyn
            
            </li>
            
            <li><a href="/cdn-cgi/l/email-protection#e2889b8cd7d3d6a2858f838b8ecc818d8f"><span class="__cf_email__" data-cfemail="462c3f2873777206212b272f2a6825292b">[email&#160;protected]</span></a></li>
            
            
            <li><a href="https://pool.sks-keyservers.net/pks/lookup?search=0xE655823D8D8B1088&op=vindex"><span class="username">Public key</span></a>
</li>
            
            <li><a href="/assets/Resume.pdf">Resume</a>
            

        </ul>
      </div>

      <div class="footer-col footer-col-2">
        <ul class="social-media-list">
          
          <li><a href="https://github.com/jyn514"><span class="icon icon--github"><svg viewBox="0 0 16 16" width="16px" height="16px"><path fill="#828282" d="M7.999,0.431c-4.285,0-7.76,3.474-7.76,7.761 c0,3.428,2.223,6.337,5.307,7.363c0.388,0.071,0.53-0.168,0.53-0.374c0-0.184-0.007-0.672-0.01-1.32 c-2.159,0.469-2.614-1.04-2.614-1.04c-0.353-0.896-0.862-1.135-0.862-1.135c-0.705-0.481,0.053-0.472,0.053-0.472 c0.779,0.055,1.189,0.8,1.189,0.8c0.692,1.186,1.816,0.843,2.258,0.645c0.071-0.502,0.271-0.843,0.493-1.037 C4.86,11.425,3.049,10.76,3.049,7.786c0-0.847,0.302-1.54,0.799-2.082C3.768,5.507,3.501,4.718,3.924,3.65 c0,0,0.652-0.209,2.134,0.796C6.677,4.273,7.34,4.187,8,4.184c0.659,0.003,1.323,0.089,1.943,0.261 c1.482-1.004,2.132-0.796,2.132-0.796c0.423,1.068,0.157,1.857,0.077,2.054c0.497,0.542,0.798,1.235,0.798,2.082 c0,2.981-1.814,3.637-3.543,3.829c0.279,0.24,0.527,0.713,0.527,1.437c0,1.037-0.01,1.874-0.01,2.129 c0,0.208,0.14,0.449,0.534,0.373c3.081-1.028,5.302-3.935,5.302-7.362C15.76,3.906,12.285,0.431,7.999,0.431z"/></svg>
</span><span class="username">jyn514</span></a>
</li>

          

          

          
          <li><a href="https://www.linkedin.com/in/jynelson514"><span class="icon icon--linkedin"><?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!-- Created with Inkscape (http://www.inkscape.org/) -->

<svg
   xmlns:dc="http://purl.org/dc/elements/1.1/"
   xmlns:cc="http://creativecommons.org/ns#"
   xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
   xmlns:svg="http://www.w3.org/2000/svg"
   xmlns="http://www.w3.org/2000/svg"
   xmlns:xlink="http://www.w3.org/1999/xlink"
   xmlns:sodipodi="http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd"
   xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape"
   version="1.1"
   width="16"
   height="16"
   viewBox="0 0 16 16"
   sodipodi:docname="icon-linkedin.svg"
   inkscape:version="0.92.1 r15371">
  <metadata>
    <rdf:RDF>
      <cc:Work
         rdf:about="">
        <dc:format>image/svg+xml</dc:format>
        <dc:type
           rdf:resource="http://purl.org/dc/dcmitype/StillImage" />
        <dc:title></dc:title>
        <cc:license
           rdf:resource="http://creativecommons.org/licenses/by-sa/4.0/" />
      </cc:Work>
      <cc:License
         rdf:about="http://creativecommons.org/licenses/by-sa/4.0/">
        <cc:permits
           rdf:resource="http://creativecommons.org/ns#Reproduction" />
        <cc:permits
           rdf:resource="http://creativecommons.org/ns#Distribution" />
        <cc:requires
           rdf:resource="http://creativecommons.org/ns#Notice" />
        <cc:requires
           rdf:resource="http://creativecommons.org/ns#Attribution" />
        <cc:permits
           rdf:resource="http://creativecommons.org/ns#DerivativeWorks" />
        <cc:requires
           rdf:resource="http://creativecommons.org/ns#ShareAlike" />
      </cc:License>
    </rdf:RDF>
  </metadata>
  <sodipodi:namedview
     pagecolor="#ffffff"
     bordercolor="#666666"
     borderopacity="1"
     objecttolerance="10"
     gridtolerance="10"
     guidetolerance="10"
     inkscape:pageopacity="0"
     inkscape:pageshadow="2"
     inkscape:window-width="667"
     inkscape:window-height="429"
     showgrid="false"
     inkscape:zoom="4"
     inkscape:cx="1.6884357e-08"
     inkscape:cy="-8"
     inkscape:window-x="182"
     inkscape:window-y="326"
     inkscape:window-maximized="0"
     inkscape:current-layer="svg2" />
  <image
     width="16"
     height="16"
     preserveAspectRatio="none"
     style="image-rendering:optimizeQuality"
     xlink:href="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEAAAABACAYAAACqaXHeAAAABHNCSVQICAgIfAhkiAAAAydJREFU eJztmz9MGlEcx7/vghY6gJbBQEmKJl3aREgYbKJGhk6dTIhjUybHykQ76tYwMbsUJ4emqUPbyQET GRxIpUnbgQTP5CJxwEMHz3/xdbiYatJ67+jJ78HdZ4T33n3f53jv3nvhGAAg/zEGDBYANgPGouhn ON8D+AZwlkchozGz8/dqYOwBdbauwnkLOE0qwGDBdZ0HAMbCwGBBAViaOgsdbEYBYxHqGGQwFlWo M1DjCaAOQI3rBfhECw4FBpCbjCM3NYqQ36x2eHKB4uYOihUVbeP8zkLeJQxvvnCrQkOBAZTnJ5CI BP/6fa15hPTyVk9KsBwCVp0HgEQkiPL8BIYCA46G6waWAnKT8Vs7f0UiEkRuMu5Epq5iLWBqVLgx O2VlwVLA1YQnQsjv67lh4PhjsNcmQksBhycXwo3ZKSsLlgKKmzvCjdkpKwvWAioqas0jy4ZqzSMU K6oTmbqKpYC2cY708tatEvp6IQT8kbC0Xr8xzg9PLrC0Xu/ZzgOCS+F+xvW7QdcLEF/mEZAeCyOb eoj48H3MjJkH17u6AVU3oOrHKDcOsPZz/7/mH6E5gL97IdbY268d179eNxkNojQ3LrQJA4CVqobF 9TpU3RAqfx3phkA2FcO311PCnQeAV6kYthemMft0xPb1pBKQTcXwfm68o7ohvw+fXqaQjIqLAyQS kB4Ld9z565Tnn9nakUojoORA5wHzl2DnYEYaAY+GA461ZedgRhoBThLy+4QnxL4UAJhzighSLoQ2 GgcoVnaw9mMfgHkyPftkBIvPHwsPlaTgY1Q6AStVDdkP32981jbOUapqKDda2F6YFjqnFH0cSjUE dnUDuc+//vm9qhsoVTWhtkQPc6USILKuvxoWTiGVgHKj1fVrSiWgbVifKqv6saPXlEqACJ3s+G5D KgFO310RJBPg7N0VQSoBFHgCqANQ4wmgDkCNJ4A6ADWeAOoA1HgCqANQ4wmgDkCN6wV4f5GhDkCN J4A6ADWK+S6tS+F8TzFfJHYrfEMBzvLg/IA6StfhvAWc5RUUMhpwmgC/XAXnTepcdw7nTfDLVeA0 iUJG+w1B+v3Clyog6wAAAABJRU5ErkJggg== "
     id="image10"
     x="0"
     y="0" />
</svg>
</span><span class="username">jynelson514</span></a>
</li>

        </ul>
      </div>

      <div class="footer-col footer-col-3">
        <p><p><img src="/assets/burned%20out%20rust%20club.png" alt="a melting, smiling, ferris. it's surrounded by the cursive text &quot;burned out rust kid club&quot;." /></p>

</p>
      </div>
    </div>

  </div>

</footer>


<script data-cfasync="false" src="/cdn-cgi/scripts/5c5dd728/cloudflare-static/email-decode.min.js"></script><script>
  console.info("hi there! thanks for visiting.")
  console.info("have a snake.")
  console.info("%c\n                 _           \n                | |          \n  ___ _ __   ___| | __       \n / __| '_ \\ / _ \\ |/ /       \n \\__ \\ | | |  __/   <        \n |___/_| |_|\\___|_|\\_\\       \n                             \n Web Development is no joke\u2122 \n                             ", 'background:#000;color:#0f0')
  console.info("(stolen with love from https://gus.host/ )")
</script>
</body>
</html>
`

var testD = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
    <title>Passing nothing is surprisingly difficult</title>
    <link rel="stylesheet" href="/style.css">

    <script async src="https://www.googletagmanager.com/gtag/js?id=G-ZR09DJQ0FR"></script>
    <script>
      window.dataLayer = window.dataLayer || [];
      function gtag(){dataLayer.push(arguments);}
      gtag('js', new Date());

      gtag('config', 'G-ZR09DJQ0FR');
    </script>
  </head>
  <body>

<h1>Passing nothing is surprisingly difficult</h1>

<p class="byline">
<a href="/">David Benjamin</a> (2024-01-15)
</p>

<p>
My day job is in <a href="https://www.chromium.org/">browsers</a> and <a href="https://boringssl.googlesource.com/boringssl/">cryptography</a>, not compilers, yet I often find that I need to spend more of my time working through the semantics of programming languages than using them. This post discusses a thorny cross-language issue between C, C++, and Rust. In short:
</p>
<ul>

<li>C’s rules around pointers and <code>memcpy</code> leave no good ways to represent an empty slice of memory.

<li>C++’s pointer rules are fine, but <code>memcpy</code> in C++ inherits the C behavior.

<li>Rust FFI is not <a href="https://blog.rust-lang.org/2015/04/24/Rust-Once-Run-Everywhere.html">zero-cost</a>. Rust picked a C/C++-incompatible slice representation, requiring a conversion in each direction. Forgetting the conversion is an easy mistake and unsound.

<li>Rust slices appear to <em>also</em> be incompatible with Rust pointer arithmetic, to the point that the standard library’s slice iterator is unsound. <span class="update">(Update 2024-01-16: It sounds like this is <a href="https://github.com/rust-lang/rust/issues/117945">in the process of being fixed</a>!)</span>
</li>
</ul>

<p>As FFI challenges inherently span multiple languages, I wrote this mostly to have one common reference that describes the mismatch.</p>

<h2 id="slices">Slices</h2>

<p>
All three languages allow working with <em>slices</em>, or contiguous sequences of objects in memory. (Also called <a href="https://en.cppreference.com/w/cpp/container/span">spans</a>, but we’ll use “slices” here.) A slice is typically a pointer and a length, <code>(start, count)</code>, giving <code>count</code> objects from <code>start</code>, of some type <code>T</code>.
</p>
<p>
A slice can also be specified by two pointers, <code>(start, end)</code>, giving the objects from <code>start</code> (inclusive) to <code>end</code> (exclusive). This is better for iteration because only one value needs to be adjusted to advance, but the length is less available. C++ iterator pairs are a generalization of this form, and Rust slice iterators use this internally. The two forms can be converted with <code>end = start + count</code> and <code>count = end - start</code>, using C-style pointer arithmetic where everything is scaled by the object size. We’ll primarily discuss <code>(start, count)</code>, but this duality means slices are closely related to pointer arithmetic.
</p>
<p>
In C and C++, slices are, at best, library types built out of pointers and lengths. Sometimes functions just take or return pointer and length separately, but still use it to represent a slice of memory. In Rust, slices are language primitives, but the underlying components are exposed for unsafe code and FFIs. To work with each of these, we must understand what combinations of pointers and lengths are valid.
</p>
<p>
This is straightforward for a positive-length slice: <code>start</code> must point within an allocation where there are at least <code>count</code> objects of type <code>T</code>. But suppose we want an empty (length zero) slice. <strong>What are the valid representations of an empty slice?</strong>
</p>
<p>
Certainly we can point <code>start</code> within (or just past) some array of <code>T</code>s and set <code>count</code> to zero. But we may want to make an empty slice without an existing array. For example, a default-constructed <code>std::span&lt;T>()</code> in C++ or <code>&[]</code> in Rust. What are our options then? In particular:
</p>
<ol>

<li>Can an empty slice be <code>(nullptr, 0)</code>?

<li>Can an empty slice be <code>(alignof(T), 0)</code>, or some other aligned address that doesn’t correspond to an allocation?
</li>
</ol>
<p>
The second question may seem odd to C and C++ folks, but Rust folks may recognize it as <code><a href="https://doc.rust-lang.org/std/ptr/struct.NonNull.html#method.dangling">std::ptr::NonNull::dangling</a></code>.

<h2 id="c++-slices">C++</h2>

<p>
C++ is the easiest to discuss, as it has a formal specification and is (almost) self-consistent.
</p>
<p>
First, <code>(nullptr, 0)</code> is a valid empty slice in C++. The STL’s types <a href="https://eel.is/c++draft/span.cons#2">routinely return it</a>, and the language is compatible with it:
</p>
<ul>

<li><code>(T*)nullptr + 0</code> is <a href="https://eel.is/c++draft/expr.add#4.1">defined</a> to be <code>(T*)nullptr</code>

<li><code>(T*)nullptr - (T*)nullptr</code> is <a href="https://eel.is/c++draft/expr.add#5.1">defined</a> to be <code>0</code>
</li>
</ul>
<p>
C++ defines APIs like <code>std::span</code> in terms of <a href="https://eel.is/c++draft/span.cons#4.1">pointer addition</a> for the <code>(start, count)</code> form, and iterator pairs in terms of <a href="https://eel.is/c++draft/iterator.operations#5">pointer subtraction</a> for the <code>(start, end)</code> form, so this is both necessary and sufficient.
</p>
<p>
Moreover, it would be impractical for C++ to forbid <code>(nullptr, 0)</code>. C++ code routinely needs to interact with APIs that specify slices as individual components. Given such an API, <em>no one</em> writes code like this:
</p>

<pre class="code">void takes_a_slice(const uint8_t *in, size_t len);

uint8_t placeholder;
takes_a_slice(&placeholder, 0);
</pre>

<p>
It is much more natural to use <code>nullptr</code>:
</p>

<pre class="code">void takes_a_slice(const uint8_t *in, size_t len);

takes_a_slice(nullptr, 0);
</pre>

<p>
This means, to be practical, functions like <code>takes_a_slice</code> must accept <code>(nullptr, 0)</code>. For implementing such functions to be practical, the underlying language primitives must then also accept <code>(nullptr, 0)</code>.
</p>

<p>
As for the <code>(alignof(T), 0)</code> question, pointer <a href="https://eel.is/c++draft/expr.add#4.2">addition</a> and <a href="https://eel.is/c++draft/expr.add#5.2">subtraction</a> require the pointers point to some object and that the operation stays within the bounds of that object (or one past the end). C++ does not define there to be an object at <code>alignof(T)</code>, so this is not allowed, instead producing Undefined Behavior. This has no immediate usability concern (no one is going to write <code>reinterpret_cast&lt;uint8_t*>(1)</code> to call <code>takes_a_slice</code>), but we’ll see later that it has some consequences for Rust FFIs.
</p>

<p class="update">(Update 2024-01-16: Added the following paragraph, as this shorthand seems to have been unclear.)</p>

<p>
Of course, <em>in principle</em>, nothing stops <code>takes_a_slice</code> from defining its own unique rules for these corner cases. Beyond what the type system naturally provides, user code will rarely formally define semantics, and we must instead look to fuzzier conventions and expectations. Sadly, in C and C++, these fuzzier aspects often include well-definedness. When all the slice-adjacent language primitives consistently use the same preconditions for slice components, a naively written function will match. This is then a reasonable default interpretation for <code>takes_a_slice</code>’s preconditions. When this post discusses “the” rules for slices for C++ or C, it is in part a shorthand for this emergent convention.
</p>

<p>
However, C++ is merely <em>almost</em> self-consistent. C++ picks up <code>memcpy</code> and other functions from C’s standard library, complete with C’s semantics…
</p>

<h2 id="c-slices">C</h2>

<p>
C is messier. For the same reasons discussed above, it is impractical to reject <code>(NULL, 0)</code>. However, C’s language rules make it difficult for a function to accept it. C does not have C++’s special cases for <code>(T*)NULL + 0</code> and <code>(T*)NULL - (T*)NULL</code>. See clauses 8 and 9 of section 6.5.6 of <a href="https://www.open-std.org/jtc1/sc22/wg14/www/docs/n2310.pdf">N2310</a>. <code>memcpy</code> and the rest of the C standard library similarly <a href="https://www.imperialviolet.org/2016/06/26/nonnull.html">forbid</a> <code>(NULL, 0)</code>.
</p>
<p>
I think <strong>this should be considered a bug in the C specification</strong>, and compilers should not <a href="https://gcc.gnu.org/gcc-4.9/porting_to.html">optimize based on it</a>. In <a href="https://boringssl.googlesource.com/boringssl/">BoringSSL</a>, we found the C rules so unusable that we resorted to <a href="https://boringssl.googlesource.com/boringssl/+/17cf2cb1d226b0ba2401304242df7ddd3b6f1ff2%5E%21/">wrapping the standard library</a> with <code>n != 0</code> checks. The pointer arithmetic rules are similarly a <a href="https://boringssl.googlesource.com/boringssl/+/6be491b7bb57c3950d4fbb97fdd4a141e3fa4d63%5E%21/">tax</a> <a href="https://boringssl.googlesource.com/boringssl/+/4984e4a6325e9c6302f846c7bf2b75e8ea3fd9dd%5E%21/">on</a> <a href="https://boringssl.googlesource.com/boringssl/+/3c6085b6ae982a80633bf5369c274036702c6848%5E%21/">development</a>. Moreover, C++ inherits C’s standard library (but not its pointer rules), including this behavior. In Chromium’s C++ code, <code>memcpy</code> has been the single biggest blocker to <a href="https://bugs.chromium.org/p/chromium/issues/detail?id=1394755">enabling UBSan</a>.
</p>
<p>
Fortunately, there is hope. <a href="https://www.npopov.com/2024/01/01/This-year-in-LLVM-2023.html#zero-length-operations-on-null">Nikita Popov</a> and Aaron Ballman have written a <a href="https://docs.google.com/document/d/1guH_HgibKrX7t9JfKGfWX2UCPyZOTLsnRfR6UleD1F8/edit">proposal</a> to fix this in C. (Thank you!) While it won’t make C and C++ safe by any stretch of imagination, this is an easy step to fix an unforced error.
</p>
<p>
Note that, apart from contrived examples with deleted null checks, the current rules do not actually help the compiler meaningfully optimize code. A <code>memcpy</code> implementation cannot rely on pointer validity to speculatively read because, even though <code>memcpy(NULL, NULL, 0)</code> is undefined, slices at the end of a buffer are fine:
</p>

<pre class="code">char buf[16];
memcpy(dst, buf + 16, 0);
</pre>

<p>
If <code>buf</code> were at the end of a page with nothing allocated afterwards, a speculative read from <code>memcpy</code> would break.
</p>

<h2 id="rust-slices">Rust</h2>

<p>
Rust does <em>not</em> allow <code>(nullptr, 0)</code>. Functions like <code>std::slice_from_raw_parts</code> <a href="https://doc.rust-lang.org/std/slice/fn.from_raw_parts.html">require the pointer to be non-null</a>. This comes from Rust treating types like <code>&[T]</code> and <code>*[T]</code> as analogous to <code>&T</code> and <code>*T</code>. They are “references” and “pointers” that are represented as <code>(start, count)</code> pairs. Rust requires every pointer type to have a “null” value outside its reference type. This is used in <code>enum</code> layout optimizations. For example, <code>Option::&lt;&[T]></code> has the same size as <code>&[T]</code> because <code>None</code> uses this null value.
</p>
<p>
Unfortunately, Rust chose <code>(nullptr, 0)</code> for the null slice pointer, which means the empty slice, <code>&[]</code>, cannot use it. That left Rust having to invent an unusual convention: some non-null, aligned, but otherwise dangling pointer, usually <code>(alignof(T), 0)</code>.
</p>
<p>
Is pointer arithmetic defined for this slice? From what I can tell, the answer appears to be no! <span class="update">(Update 2024-01-16: It sounds like this is <a href="https://github.com/rust-lang/rust/issues/117945">in the process of being defined</a>!)</span>
</p>
<p>
Pointer arithmetic in Rust is spelled with the methods <code><a href="https://doc.rust-lang.org/std/primitive.pointer.html#method.add">add</a></code>, <code><a href="https://doc.rust-lang.org/std/primitive.pointer.html#method.sub_ptr">sub_ptr</a></code>, and <code><a href="https://doc.rust-lang.org/std/primitive.pointer.html#method.offset_from">offset_from</a></code>, which the standard library defines in terms of <a href="https://doc.rust-lang.org/std/ptr/index.html#allocated-object">allocated objects</a>. That means, for pointer arithmetic to work with <code>alignof(T)</code>, there must be zero-size slices allocated at every non-zero address. Moreover, <code>offset_from</code> requires the two dangling pointers derived from the same slice to point to the “same” of these objects. While the third bullet <a href="https://doc.rust-lang.org/std/ptr/index.html#safety">here</a>, second sentence, says casting literals gives a pointer that is “valid for zero-sized accesses”, it says nothing about allocated objects or pointer arithmetic.

<p>
Ultimately, these semantics come from LLVM. The Rustonomicon has <a href="https://doc.rust-lang.org/nomicon/vec/vec-alloc.html#:~:text=The%20other%20corner%2Dcase%20we%20need%20to%20worry%20about%20is%20empty%20allocations">more to say on this</a> (beginning “The other corner-case…”). It concludes that, while there are infinitely many <em>zero-size</em> types at <code>0x01</code>, Rust conservatively assumes alias analysis does <em>not</em> allow offsetting <code>alignof(T)</code> with zero for <em>positive-sized</em> types. This means <strong>Rust pointer arithmetic rules are incompatible with Rust empty slices.</strong> But recall that slice iteration and pointer arithmetic are deeply related. The Rustonomicon’s <a href="https://doc.rust-lang.org/nomicon/vec/vec-into-iter.html">sample iterator</a> uses pointer arithmetic, but needs to guard addition with <code>cap == 0</code> in <code>into_iter</code> and cast to <code>usize</code> in <code>size_hint</code>.
</p>
<p>
This is too easy for programmers to forget. Indeed the real Rust slice iterator does pointer arithmetic unconditionally (<a href="https://github.com/rust-lang/rust/blob/76101eecbe9aa80753664bbe637ad06d1925f315/library/core/src/slice/iter.rs#L94">pointer addition</a>, <a href="https://github.com/rust-lang/rust/blob/76101eecbe9aa80753664bbe637ad06d1925f315/library/core/src/slice/iter/macros.rs#L57">pointer subtraction</a>, behind <a href="https://github.com/rust-lang/rust/blob/76101eecbe9aa80753664bbe637ad06d1925f315/library/core/src/slice/iter/macros.rs#L141">some</a> <a href="https://github.com/rust-lang/rust/blob/76101eecbe9aa80753664bbe637ad06d1925f315/library/core/src/slice/iter.rs#L132">macros</a>). This suggests <strong>Rust slice iterators are unsound.</strong>
</p>

<h2 id="ffis">FFIs</h2>

<p>
Beyond self-consistency concerns, all this means Rust and C++ slices are incompatible. Not all Rust <code>(start, count)</code> pairs can be passed into C++ and vice versa. C’s issues make its situation less clear, but the natural fix is to bring it in line with C++.
</p>
<p>
This means Rust FFI is not <a href="https://blog.rust-lang.org/2015/04/24/Rust-Once-Run-Everywhere.html">“zero-cost”</a>. <strong>Passing slices between C/C++ and Rust requires conversions in both directions to avoid Undefined Behavior.</strong>
</p>
<p>
More important (to me) than performance, this is a safety and ergonomics problem. Programmers cannot be expected to memorize language specifications. If given a <code>&[T]</code> and trying to call a C API, the natural option is to use <code><a href="https://doc.rust-lang.org/std/primitive.slice.html#method.as_ptr">as_ptr</a></code>, but that will return a C/C++-incompatible output. Most Rust crates I’ve seen which wrap C/C++ APIs do not convert and are unsound as a result.

<p>
This is particularly an issue because C and C++’s (more serious) safety problems cause <a href="https://security.googleblog.com/2021/09/an-update-on-memory-safety-in-chrome.html">real user harm</a> and need addressing. But there is half a century of existing C and C++ code. We cannot realistically address this with a new language without good FFI. What makes for good FFI? At a bare minimum, I think <strong>calling a C or C++ function from Rust should not be dramatically less safe than calling that function from C or C++.</strong>
</p>

<h2 id="wishlist">Wishlist</h2>

<p>
Empty lists should not be so complicated. We could change C, C++, and Rust in a few ways to improve things: 
</p>

<h3 id="make-c-accept-nullptr">Make C accept <code>nullptr</code></h3>

<p>
See Nikita Popov and Aaron Ballman’s <a href="https://docs.google.com/document/d/1guH_HgibKrX7t9JfKGfWX2UCPyZOTLsnRfR6UleD1F8/edit">proposal</a>.
</p>

<h3 id="fix-rust-slices">Fix Rust’s slice representation</h3>

<p>
All the <code>alignof(T)</code> problems ultimately come from Rust’s unusual empty slice representation. This falls out of Rust’s need for a “null” <code>*[T]</code> value that is not a <code>&[T]</code> value. Rust could have chosen any of a number of unambiguously unused representations, such as <code>(nullptr, 1)</code>, <code>(nullptr, -1)</code>, or <code>(-1, -1)</code>.
</p>
<p>
While this would be a significant change now, with compatibility implications to work through, I think it is worth seriously considering. It would address the root cause of this mess, fixing a soundness hazard in not just Rust FFI, but Rust on its own. This hazard is real enough that Rust’s standard library hits it.
</p>
<p>
This is also the only option I see that fully meets Rust’s “zero-cost” FFI <a href="https://blog.rust-lang.org/2015/04/24/Rust-Once-Run-Everywhere.html">goals</a>. Even if we make C and C++ accept <code>(alignof(T), 0)</code> from Rust (see below), any slices passed from C/C++ to Rust may still be <code>(nullptr, 0)</code>.
</p>

<h3 id="define-invalid-pointers">Define pointer arithmetic for invalid pointers</h3>

<p>
If Rust leaves its slice representation as is, we instead should define pointer arithmetic for <code>NonNull::dangling()</code>. Expecting low-level Rust code to guard all pointer offsets is impractical.
</p>

<p class="update">Update 2024-01-16: Happily, it sounds like this is already <a href="https://github.com/rust-lang/rust/issues/117945">in the process of being defined</a>!</p>

<p>
Where <code>nullptr</code> is a single value which could just be special-cased, there are many <code>alignof(T)</code> values. It seems one would need to define it in terms of the actual allocated objects. This is well beyond my expertise, so I’ve likely gotten all the details and terminology wrong, but one possibility is, in the vein of <a href="https://www.open-std.org/jtc1/sc22/wg14/www/docs/n2364.pdf">PNVI-ae-udi</a>:
</p>
<ul>

<li><code>cast_ival_to_ptrval</code> returns a special <code>@empty</code> provenance when casting garbage values (unchanged from PNVI-ae-udi)

<li>Adding zero to a pointer with the <code>@empty</code> provenance is valid and gives back the original pointer

<li>Two pointers with <code>@empty</code> provenance can be subtracted to give zero if they have the same address
</li>
</ul>
<p>
One subtlety, however, is that <code>cast_ival_to_ptrval</code> might not give back <code>@empty</code> if there is an object at that address. Giving back a concrete provenance means picking up <a href="https://www.open-std.org/jtc1/sc22/wg14/www/docs/n2443.pdf">pointer-zapping</a> semantics, which would be undesirable here. For <code>alignof(T)</code>, that shouldn’t happen if the maximum alignment is under a page and the bottom page is never allocated. But Rust allows not just <code>alignof(T)</code> but <a href="https://doc.rust-lang.org/std/ptr/index.html#safety">any non-zero integer literal</a>, <em>even if some allocation happens to exist at that address</em>. (Perhaps we could use the “user-disambiguation” aspect and say all integer-to-pointer casts may additionally disambiguate to <code>@empty</code>? Would that impact the compiler’s aliasing analysis?)
</p>
<p>
I think this complexity demonstrates why <code>nullptr</code> is a much better choice for an empty slice than a dangling pointer. Pointer arithmetic with <code>nullptr</code> is easy to define, and <code>nullptr</code> cannot alias a real allocation.
</p>
<p>
If Rust (and LLVM) accepted invalid pointers, it would fix the soundness issues within Rust, but not with FFIs. If the C and C++ standards <em>also</em> picked this up, it would <em>partially</em> fix FFIs. We could then directly pass Rust slices into C and C++, but not in the other direction. Directly passing C and C++ slices into Rust can only be fixed by changing Rust to accept <code>(nullptr, 0)</code> form.</p>

<p><s>(Outside of Rust FFI, there’s no reason to use <code>alignof(T)</code> as a pointer in C/C++, so I do not know how plausible it would be for C/C++ to accept it.)</s> <span class="update">Update 2024-01-16: Nelson Elhage reminded me that non-null sentinel pointers are sometimes used to <a href="https://github.com/torvalds/linux/blob/ffc253263a1375a65fa6c9f62a893e9767fbebfa/include/linux/slab.h#L167-L178">allocate zero bytes</a>. While C forbids <code>malloc</code> from doing this (<code>malloc(0)</code> must return either <code>nullptr</code> or a <em>unique</em> non-null pointer), other allocators might reasonably pick this option. It makes error checks more uniform without actually reserving address space. So there is a non-Rust reason to allow these pointers in C and C++.
</span></p>

<h3 id="ffi-helpers">FFI helpers in Rust standard library</h3>

<p>
If the languages’ slice representations cannot be made compatible, we’re still left with safety hazards in Rust FFI. In that case, Rust’s standard library should do more to help programmers pick the right operations: Add analogs of <code>slice::from_raw_parts</code>, <code>slice::as_ptr</code>, etc., that use the C and C++ representation, converting internally as needed. Document existing functions very clear warnings that they cannot be used for FFI. Finally, audit all existing calls in crates.io, as the majority of existing calls are likely for FFI.
</p>
<p>
For <code>slice::from_raw_parts</code>, we could go further and fix the existing function itself. This would be backwards-compatible, but adds unnecessary conversions to non-FFI uses. That said, if the crates.io audit reveals mostly FFI uses, that conversion may be warranted. For non-FFI uses, a type signature incorporating <code><a href="https://doc.rust-lang.org/std/ptr/struct.NonNull.html">std::ptr::NonNull</a></code> would have been more appropriate anyway.

<p>
This would improve things, but it’s an imperfect solution. We’d still sacrifice zero-cost FFI, and we’d still rely on programmers to read the warnings and realize the natural options are incorrect.
</p>

<h2 id="acknowledgements">Acknowledgements</h2>

<p>
Thanks to Alex Chernyakhovsky, Alex Gaynor, Dana Jansens, Adam Langley, and Miguel Young de la Sota for reviewing early iterations of this post. Any mistakes in here are my own.
</p>

</body>
</html>`

func TestLong(t *testing.T) {

	htmls := []string{testA, testB, testC, testD}

	for _, html := range htmls {

		reader := strings.NewReader(html)

		body, _, _, err := HtmlToText(reader)

		if err != nil {
			t.FailNow()
		}

		t.Logf("%+q", body)
	}

}

func TestBasic(t *testing.T) {

	html := `
		<html>
			<head>
				<title>Dan</title>
				<meta name="description" content="desc">
			</head>
		<body>
			<div>
				<p>this is <a>some</a> text</p>
			</div>
		</body>
		</html>	
	`

	reader := strings.NewReader(html)

	body, title, description, err := HtmlToText(reader)

	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	if len(body) != 1 {
		t.Error()
	}
	if title != "Dan" {
		t.Error()
	}
	if description != "desc" {
		fmt.Println(description)
		t.Error()
	}

}
