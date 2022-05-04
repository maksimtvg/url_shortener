package generator

import (
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestUriGenerator_GenerateUri(t *testing.T) {
	var out []string
	callsNumber := 10000

	tests := struct {
		input string
		error error
	}{
		input: "https://example.com/testurl",
		error: nil,
	}

	generator := NewUriGenerator(0)

	for i := 0; i < callsNumber; i++ {
		uri, err := generator.GenerateUri()
		if !reflect.DeepEqual(err, tests.error) {
			t.Errorf("error test failed")
		}
		out = append(out, uri)
	}

	for j := 0; j < callsNumber; j++ {
		for k := 0; k < callsNumber; k++ {
			if j == k {
				continue
			}
			if strings.Compare(out[k], out[j]) == 0 {
				t.Errorf("slice contains equal strings")
			}
		}
	}
}

// tests correct type error returns
func TestUriGenerator_MaxIntError(t *testing.T) {
	testsTable := []struct {
		input int64
		error error
	}{
		{
			input: math.MaxInt64,
			error: maxNumberError,
		},
		{
			input: math.MaxInt64 - 1,
			error: nil,
		},
	}

	for _, tt := range testsTable {
		g := NewUriGenerator(tt.input)
		_, err := g.GenerateUri()
		if !reflect.DeepEqual(err, tt.error) {
			t.Errorf("expected an error (%v), got error (%v)", tt.error, err)
		}
	}
}
