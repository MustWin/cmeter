package mockapi

import (
	"net/http"

	"github.com/MustWin/cmeter/configuration"
	"github.com/MustWin/cmeter/context"
)

type Server struct {
	context.Context
	addr string
}

func (server *Server) ListenAndServe() error {
	s := http.Server{
		Addr:    server.addr,
		Handler: server,
	}

	context.GetLogger(server).Infof("listening on %v", server.addr)
	return s.ListenAndServe()
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx := context.DefaultContextManager.Context(server, w, r)
	defer context.DefaultContextManager.Release(ctx)
	defer func() {
		status, ok := ctx.Value("http.response.status").(int)
		if ok && status >= 200 && status <= 399 {
			context.GetResponseLogger(ctx).Infoln("response completed")
		}
	}()

	var err error
	w, err = context.GetResponseWriter(ctx)
	if err != nil {
		context.GetLogger(ctx).Warnf("response writer not found in context")
	}

	context.GetRequestLogger(ctx).Infoln("request started")
	eventServe(ctx, w, r)
}

func eventServe(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	context.GetLogger(ctx).Infoln("event published")
	w.WriteHeader(http.StatusOK)
}

func NewServer(ctx context.Context, config *configuration.Config) (*Server, error) {
	return &Server{
		Context: ctx,
		addr:    "localhost:9090",
	}, nil
}
