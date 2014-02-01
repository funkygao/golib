package progress

import "testing"

func TestProgress(t *testing.T) {
    current, total, cols := 2, 5, 20
    bar := progress(current, total, cols)
    expected := "\x1b[1m2 / 5\x1b[0m [====       ] "
    if expected != bar {
        t.Error("expected:", expected, " current:", bar)
    }
}
