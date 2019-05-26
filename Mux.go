package mayday

import (
	"context"
	"net"
	"net/http"
	"sync"
)

// Mux struct implement
type Mux struct {
	//methods trees key is http method
	trees map[string]*node

	//middlewares
	middlewares []http.Handler

	MethodNotAllow http.HandlerFunc
	NotFound       http.HandlerFunc

	cxtPool *sync.Pool
}

// New a router interface
func New() *Mux {
	return &Mux{
		trees:       make(map[string]*node),
		middlewares: make([]http.Handler, 0),
		cxtPool: &sync.Pool{
			New: func() interface{} {
				return context.TODO()
			},
		},
	}
}

func (mux *Mux) Dump() {
	if mux.trees != nil {
		for _, v := range mux.trees {
			v.dumpNode()
		}
	}
}

func defaultMethodNotAllow(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(405)
}

// Use add middleware onto stack
func (mux *Mux) Use(handler http.Handler) {
	mux.middlewares = append(mux.middlewares, handler)
}

func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.URL.Path
	unescape := false

	tree, ok := mux.trees[method]
	if ok {
		handlers, params, tsr := tree.getValue(path, nil, unescape)
		ctx := mux.cxtPool.Get().(context.Context)
		// Copy over default net/http server context keys
		if v, ok := r.Context().Value(http.ServerContextKey).(*http.Server); ok {
			ctx = context.WithValue(ctx, http.ServerContextKey, v)
		}
		if v, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr); ok {
			ctx = context.WithValue(ctx, http.LocalAddrContextKey, v)
		}
		for _, param := range params {
			ctx = context.WithValue(ctx, param.Key, param.Value)
		}

		r = r.WithContext(ctx)

		if tsr {
			for _, handler := range handlers {
				handler.ServeHTTP(w, r)
			}
		} else {
			if mux.NotFound != nil {
				mux.NotFound(w, r)
			}
		}
		mux.cxtPool.Put(ctx)
	} else {
		if mux.MethodNotAllow != nil {
			mux.MethodNotAllow(w, r)
		} else {
			defaultMethodNotAllow(w, r)
		}
	}

}

func (mux *Mux) handle(method, path string, handler http.Handler) {
	tree, ok := mux.trees[method]
	if !ok {
		tree = new(node)
		mux.trees[method] = tree
	}
	// 把middleware添加在path前
	handlers := append(mux.middlewares, handler)
	tree.addRoute(path, handlers)
}

func (mux *Mux) GET(path string, handler http.Handler) {
	mux.handle("GET", path, handler)
}

func (mux *Mux) POST(path string, handler http.Handler) {
	mux.handle("POST", path, handler)
}

func (mux *Mux) PUT(path string, handler http.Handler) {
	mux.handle("PUT", path, handler)
}

func (mux *Mux) DELETE(path string, handler http.Handler) {
	mux.handle("DELETE", path, handler)
}

func (mux *Mux) OPTIONS(path string, handler http.Handler) {
	mux.handle("OPTIONS", path, handler)
}

func (mux *Mux) HEAD(path string, handler http.Handler) {
	mux.handle("HEAD", path, handler)
}

func (mux *Mux) PATCH(path string, handler http.Handler) {
	mux.handle("PATCH", path, handler)
}

func (mux *Mux) TRACE(path string, handler http.Handler) {
	mux.handle("TRACE", path, handler)
}

func (mux *Mux) GetValue(path, method string) []http.Handler {
	if mux.trees != nil {
		if n, ok := mux.trees[method]; ok {
			handlers, _, _ := n.getValue(path, nil, false)
			return handlers
		}
	}
	return nil
}
