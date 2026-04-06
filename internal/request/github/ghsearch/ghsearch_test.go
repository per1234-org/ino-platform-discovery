package ghsearch

import (
	"testing"

	"github.com/google/go-github/v79/github"
	"github.com/stretchr/testify/assert"
)

// Test_intersecting provides coverage for the `intersecting` function.
func Test_intersecting(t *testing.T) {
	resultsA := github.CodeSearchResult{}
	resultsB := github.CodeSearchResult{}

	intersectingRepositoryA := github.Repository{
		HTMLURL: new("foo"),
	}
	intersectingCodeResultA := github.CodeResult{
		Repository: &intersectingRepositoryA,
	}

	intersectingRepositoryB := github.Repository{
		HTMLURL: new("bar"),
	}
	intersectingCodeResultB := github.CodeResult{
		Repository: &intersectingRepositoryB,
	}

	nonIntersectingRepository := github.Repository{
		HTMLURL: new("baz"),
	}
	nonIntersectingCodeResult := github.CodeResult{
		Repository: &nonIntersectingRepository,
	}

	resultsA.CodeResults = append(resultsA.CodeResults, &intersectingCodeResultA, &intersectingCodeResultB)
	resultsB.CodeResults = append(resultsB.CodeResults, &nonIntersectingCodeResult, &intersectingCodeResultB, &intersectingCodeResultA)

	assertion := []*github.CodeResult{
		{
			Repository: &intersectingRepositoryA,
		},
		{
			Repository: &intersectingRepositoryB,
		},
	}

	assert.Equal(t, assertion, intersecting(&resultsA, &resultsB))
}
