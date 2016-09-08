package context

import (
	"net/http"
	"sync"
)

type ContextManager struct {
	contexts map[*http.Request]Context
	mutex    sync.Mutex
}

var DefaultContextManager = NewContextManager()

func NewContextManager() *ContextManager {
	return &ContextManager{
		contexts: make(map[*http.Request]Context),
	}
}

func (m *ContextManager) Context(parent Context, w http.ResponseWriter, r *http.Request) Context {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if ctx, ok := m.contexts[r]; ok {
		return ctx
	}

	if parent == nil {
		parent = Background()
	}

	ctx := WithRequest(parent, r)
	ctx, w = WithResponseWriter(ctx, w)
	ctx = WithLogger(ctx, GetLogger(ctx))
	m.contexts[r] = ctx
	return ctx
}

func (m *ContextManager) Release(ctx Context) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	r, err := GetRequest(ctx)
	if err != nil {
		GetLogger(ctx).Error("no request found in context at release")
		return
	}

	delete(m.contexts, r)
}
