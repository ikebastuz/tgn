package bot

import (
	"sync"
)

type Store struct {
	count int64
	mu    sync.Mutex
}
