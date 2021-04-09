package declparse

import (
	"reflect"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			input := strings.TrimRight(tt.s, ";")
			p := NewStringParser(input + ";")
			got, err := p.Parse()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.n) {
				t.Errorf("Parse()\n  got: %s\n want: %s", got, tt.n)
			}
		})
	}
	// str := "- (instancetype)initWithContentRect:(NSRect)contentRect styleMask:(NSWindowStyleMask)style backing:(NSBackingStoreType)backingStoreType defer:(BOOL)flag;"
	// p := NewStringParser(str)
	// stmt, err := p.Parse()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println(stmt)
}
