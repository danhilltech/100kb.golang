# 100kb

“There is no such thing as a moral or an immoral book. Books are well written, or badly written.” I want to find all the well written content on the internet.

## What is this?
This is a personal project. It’s a simple feed of articles written by real people with interesting things to say.

I wanted to explore some new technologies I haven't had a chance to try before, mostly rust and LLMs.

## Opinionated
As I browse HN each day I find myself most drawn to the content on personal blogs: people who have interesting things to say, from their own point-of-view, on a whole range of topics. Rather than being interested in particular topics, I’m more interested in reading from interesting people.

To quote Oscar Wilde, “There is no such thing as a moral or an immoral book. Books are well written, or badly written.”

I want to find all the well written content on the internet.

To do so I started digging into features that are predictive of what, in completely my opinion, is well written content. The features I use are things like: is the content written in the first person, is there high content density (usually less than 100kb in page weight - my first feature, hence the name of the project), are there minimal ads/trackers/bloat, how active is the blog, and a dozen or so more.

## How it works
It's a simple pipeline that I can run locally from a `crontab` on my desktop. The pipeline runs in a few steps:

1.  A search phase runs to find new possible blogs to index. Mostly from the latest hacker news and some public datasets of personal blogs.
2.  If I discover a valid RSS feeds on those blogs, I index that content.
3.  I extract features from the new content (page weight, use of ads/trackers, first/third person writing, and about 20 more).
4.  Articles are filtered with a simple logistic regression trained on a few hundred hacker news links I hand labeled.
5.  Static HTML is written and uploaded to the CDN.

The code is mostly written in go with sqlite. I use `rust-bert` and Brave's rust implementation of `adblock`, both with protobuf connecting them to go (cgo/FFI was a nightmare).

I very aggressively cache content - basically, if I get a valid response from a URL it won't be crawled again. My assumption is that for the personal blogs, content doesn't materially change after publication. I took a lot from [Rachel's writing on RSS readers](https://rachelbythebay.com/w/2024/08/02/fs/).

## Feedback/PRs
Would love to hear from you, or feel free to open a PR.

## Inspiration/Prior art
There's a bunch of other great work out there that does similar things. Including a lot of [past hacker news discussions](https://www.google.com/search?q=personal+blogs+site%3Anews.ycombinator.com).

1.  [Search my site](https://searchmysite.net/search/browse/)
2.  [Indie Blog](https://indieblog.page/all)
3.  [Places to find indie web content](https://untested.sonnet.io/Places+to+Find+Indie+Web+Content)
4.  [Marginalia Search](https://github.com/MarginaliaSearch/PublicData/blob/master/sets/random-domains.txt)
5.  [Kagi Search - Small Web](https://github.com/kagisearch/smallweb)
6.  [Bearblog discover](https://bearblog.dev/discover/)