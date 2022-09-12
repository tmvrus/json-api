package endpoint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBaseRequest_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		in    baseRequest
		error string
	}{
		{
			name: "ok",
			in: baseRequest{
				Version: "2.0",
				Method:  "nonEmpty",
				Params:  []byte(`non-zero length`),
			},
		},
		{
			name: "bad version",
			in: baseRequest{
				Version: "2.1",
				Method:  "nonEmpty",
				Params:  []byte(`non-zero length`),
			},
			error: "unsuported version 2.1",
		},
		{
			name: "empty method",
			in: baseRequest{
				Version: "2.0",
				Method:  "",
				Params:  []byte(`non-zero length`),
			},
			error: "empty method passed",
		},
		{
			name: "empty params",
			in: baseRequest{
				Version: "2.0",
				Method:  "nonEmpty",
				Params:  []byte(``),
			},
			error: "empty params",
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := test.in.Validate()
			if test.error == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.error)
			}
		})
	}
}
