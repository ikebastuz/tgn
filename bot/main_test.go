package bot

import (
	"testing"

	"github.com/ikebastuz/tgn/types"
)

func TestEvaluator(t *testing.T) {
	t.Run("should update store by reference", func(t *testing.T) {
		store := Store{
			count: 0,
		}
		update := types.TelegramUpdate{}
		want := int64(1)
		HandleMessage(update, &store)
		if store.count != want {
			t.Errorf("expected %v, got %v", want, store.count)
		}
	})
}
