git submodule add https://github.com/skeskinen/bert.cpp.git

cd bert.cpp
git submodule update --init --recursive

## TODO
[x] generalize diskv cache
[] seperate bad urls log
[] domain failure table tracking

NEXT: MOVE HEAD ETC INTO HTTP


## Flow
1. HN gets any missing hn posts
2. URLs are added to candidate_urls
3. Check candidate urls against known bad urls, also domains we already have a feed for (feed needs domain col)
4. get new candidate urls
4a. make head request (including redir) - check content type and length
4b. extract feed urls
5. add feed urls to domain table (with domain)
6. refresh feeds, adding articles

Notes
1. retitle feeds to domains, add every domain and last get attempt result

2. requests table
- url
- domain
- content type
- cache file

