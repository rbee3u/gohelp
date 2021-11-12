package mate_test

import (
	"testing"

	"github.com/rbee3u/gohelp/mate"
)

func TestPathWalker(t *testing.T) {
	pw := mate.NewPathWalker()
	if got, want := pw.String(), ""; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}

	pw.Enter(mate.StringPath("key"))
	if got, want := pw.String(), "key"; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}

	pw.Enter(mate.IntegerPath(1))
	if got, want := pw.String(), "key_1"; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}

	pw.Enter(mate.StringPath("prop"))
	if got, want := pw.String(), "key_1_prop"; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}

	pw.Exit()
	if got, want := pw.String(), "key_1"; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}

	pw.Exit()
	if got, want := pw.String(), "key"; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}

	pw.Exit()
	if got, want := pw.String(), ""; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}
}
