// synonyms provides parsing for Solr synonym files
// Format details at (e.g.): https://lucene.apache.org/core/6_6_0/analyzers-common/org/apache/lucene/analysis/synonym/SolrTermParser.html
package synonyms

import (
	"fmt"
	"strings"
)

// Term represents an edge in a TermGraph
type Term struct {
	Replacement string
	Equivalent  string
}

// TermGraph holds the nodes and edges of a synonym graph
type TermGraph map[string]Term

// Equivalents returns the equivalent values of a given term
func (graph TermGraph) Equivalents(term string) []string {
	eqs := []string{}
	eq := term

	for {
		t, ok := graph[eq]
		if !ok {
			break
		}

		eq = t.Equivalent
		if eq == "" || eq == term {
			break
		}
		eqs = append(eqs, eq)
	}

	return eqs
}

// Replacements return any explicit replacements for a term
func (graph TermGraph) Replacements(term string) []string {
	if syn, ok := graph[term]; ok {
		r := syn.Replacement
		eqs := graph.Equivalents(r)
		if r != "" {
			eqs = append(eqs, r)
		}
		return eqs
	}
	return nil
}

func parseLine(line string) (*TermGraph, error) {
	var replacements []string
	var equivalents []string
	var termList *[]string

	line = strings.TrimSpace(line)
	termList = &equivalents
	term := ""

	for i, c := range line {
		if c == '#' || c == '\n' {
			break
		} else if c == ',' || c == '=' {
			if i == len(line)-1 || term == "" {
				return nil, fmt.Errorf("Invalid character '%c'", c)
			}

			*termList = append(*termList, strings.TrimRight(term, " "))
			term = ""
		} else if c == '>' {
			if i < 1 || line[i-1] != '=' {
				return nil, fmt.Errorf("Invalid character '%c'", c)
			}
			termList = &replacements
		} else if c != ' ' || term != "" {
			// TODO: restrict "legal" character set
			term = term + string(c)
		}
	}

	term = strings.TrimRight(term, " ")
	if term != "" {
		*termList = append(*termList, term)
	}

	graph := make(TermGraph, len(equivalents))
	if len(equivalents) == 0 {
		return &graph, nil
	}

	if len(replacements) > 0 {
		// Map equivalent terms to the first replacement
		for _, p := range equivalents {
			graph[p] = Term{Replacement: replacements[0]}
		}

		// Mark subsequent replacements as "equivalent"
		for i, p := range replacements[:len(replacements)-1] {
			graph[p] = Term{Equivalent: replacements[i+1]}
		}
	} else if len(equivalents) > 1 {
		for i, p := range equivalents[:len(equivalents)-1] {
			graph[p] = Term{Equivalent: equivalents[i+1]}
		}
		graph[equivalents[len(equivalents)-1]] = Term{Equivalent: equivalents[0]}
	} else {
		return nil, fmt.Errorf("Invalid line")
	}

	return &graph, nil
}

// Parse a synonym file
func Parse(doc string) (*TermGraph, error) {
	lines := strings.Split(doc, "\n")
	graph := TermGraph{}
	for i, line := range lines {
		entries, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("invalid string '%s' at line %d", line, i)
		}

		for str, entry := range *entries {
			if _, ok := graph[str]; ok {
				// TODO: handle references to previously-defined terms
				fmt.Errorf("'%s' already set", str)
			}
			graph[str] = entry
		}
	}
	return &graph, nil
}
