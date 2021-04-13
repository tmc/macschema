package declparse

import (
	"testing"
)

func TestAST_Strings(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got := tt.n.String()
			if got != tt.s {
				t.Errorf("String()\n  got: %s\n want: %s", got, tt.s)
			}
		})
	}
}
