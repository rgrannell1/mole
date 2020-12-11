# Mole

After [Grant's golden mole](https://en.wikipedia.org/wiki/Grant%27s_golden_mole).

### Stability Index

> 1, Experimental - This project might die, it's undertested and underdocumented, and redesigns and breaking changes are likely

### Usage

mole extracts Chrome's history files as JSON to be consumed by downstream tools like jq and gron. It also performs joins on the various tables to give the most information possible per-entry.

```
downloads
downloads_slices
downloads_url_chains
keyword_search_terms
meta
segment_usage
segments
sqlite_sequence
typed_url_sync_metadata
urls
visit_source
visits
```
```sh
mole ls <tablename> [--db <string>]
```
```sh
mole ls keyword_search_terms
```
Find all mole-related searches
```sh
mole ls keyword_search_terms | jq -r 'select(.term | contains("mole"))'
```
### Files

```
go.mod            go dependency information
main.go           the core code
read-sqlite.go    SQL-reading utility code
```

### License

The MIT License

Copyright (c) 2020 Róisín Grannell

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

