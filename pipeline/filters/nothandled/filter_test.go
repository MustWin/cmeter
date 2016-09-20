package nothandled_test

import (
	"testing"

	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/pipeline"
	"github.com/MustWin/cmeter/pipeline/filters/nothandled"
)

func TestHandleMessage(t *testing.T) {
	f := nothandled.New()
	err := f.HandleMessage(context.Background(), pipeline.NewNullMessage())
	if err == nil {
		t.Error("no error returned for unhandled message")
	}
}
