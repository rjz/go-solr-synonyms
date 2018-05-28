package synonyms

import (
	"github.com/kr/pretty"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	docs := map[string]TermGraph{
		"":             TermGraph{},
		`# comment ok`: TermGraph{},
		`sapphire, azure => blue#lagging comment`: TermGraph{
			"sapphire": Term{Replacement: "blue"},
			"azure":    Term{Replacement: "blue"},
		},
		`sapphire, prussian blue => azure, blue,cyan`: TermGraph{
			"sapphire":      Term{Replacement: "azure"},
			"prussian blue": Term{Replacement: "azure"},
			"azure":         Term{Equivalent: "blue"},
			"blue":          Term{Equivalent: "cyan"},
		},
		`
     sapphire => blue
     azure =>blue
     `: TermGraph{
			"sapphire": Term{Replacement: "blue"},
			"azure":    Term{Replacement: "blue"},
		},
		` sapphire   => blue`: TermGraph{
			"sapphire": Term{Replacement: "blue"},
		},
		`azure , blue `: TermGraph{
			"azure": Term{Equivalent: "blue"},
			"blue":  Term{Equivalent: "azure"},
		},
		`azure, blue # comment ok`: TermGraph{
			"azure": Term{Equivalent: "blue"},
			"blue":  Term{Equivalent: "azure"},
		},
	}

	for doc, expected := range docs {
		actual, _ := Parse(doc)
		if !reflect.DeepEqual(expected, *actual) {
			t.Errorf("Expected '%s', got '%s' from '%s'", pretty.Sprint(expected), pretty.Sprint(*actual), doc)
		}
	}
}

func TestParseInvalid(t *testing.T) {
	invalidDocs := []string{
		"azure",
		",azure",
		"azure, blue=",
		"azure, blue>",
		"azure > blue",
		"=azure, blue",
	}

	for _, doc := range invalidDocs {
		if _, err := Parse(doc); err == nil {
			t.Errorf("Expected error for '%s', didn't get it", doc)
		}
	}
}

func TestEquivalents(t *testing.T) {
	sm, _ := Parse("azure, blue, cerulean, cyan")
	expected := []string{"blue", "cerulean", "cyan"}
	actual := sm.Equivalents("azure")
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
	}
}

func TestReplacements(t *testing.T) {
	sm, _ := Parse("azure => blue, cerulean, cyan")
	expected := []string{"cerulean", "cyan", "blue"}
	actual := sm.Replacements("azure")
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
	}
}
