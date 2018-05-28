# go-solr-synonyms

Parse [Solr synonym files][SolrTermParser].

[![Build Status](https://travis-ci.org/rjz/go-solr-synonyms.svg?branch=master)](https://travis-ci.org/rjz/go-solr-synonyms)

## Usage

```go
import (
  "fmt"
  "gopkg.in/rj/go-solr-synonyms.v0"
)

func main () {
  terms, _ := synonyms.Parse(`
beagle, shepherd, heeler => dog
cabbage, kimchi, sauerkraut
`)
  // Find replacements
  fmt.Println(terms.Replacements("beagle")) // "dog"

  // Find equivalents
  fmt.Println(terms.Equivalents("kimchi")) // "cabbage,sauerkraut"
}
```

## License

MIT

[SolrTermParser]: https://lucene.apache.org/core/6_6_0/analyzers-common/org/apache/lucene/analysis/synonym/SolrTermParser.html
