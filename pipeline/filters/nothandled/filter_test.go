package nothandled_test

import (
	"testing"

	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/pipeline/filters/nothandled"
	"github.com/MustWin/cmeter/pipeline/messages/base"
)

/*func TestMain() {

}*/

func TestHandleMessage(t *testing.T) {
	f := nothandled.New()
	err := f.HandleMessage(context.Background(), &base.Message{})
	if err == nil {
		t.Error("no error returned for unhandled message")
	}
}
