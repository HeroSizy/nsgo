//go:build js && wasm

package harmonyos

import (
	"encoding/json"
	"testing"

	"github.com/nsgo-dev/nsgo/internal/driver"
)

// TestInputSendsEditability checks that a field's input operation ends with whether it is read-only and disabled, where Index.ets reads them: o[13] and o[14].
func TestInputSendsEditability(t *testing.T) {
	for _, c := range []driver.InputConfig{{Readonly: true}, {Disabled: true}, {}} {
		m := &mounter{views: map[driver.NodeID]*view{}}
		m.input(4, &c)
		var got [][]any
		if err := json.Unmarshal(m.o.bytes(), &got); err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 || len(got[0]) != 15 || got[0][0] != "in" {
			t.Fatalf("input operations %v, want one in with 13 arguments", got)
		}
		if o := got[0]; o[13] != c.Readonly || o[14] != c.Disabled {
			t.Errorf("readonly %v and disabled %v, want %v and %v", o[13], o[14], c.Readonly, c.Disabled)
		}
	}
}
