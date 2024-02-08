package internal

import "github.com/yuin/goldmark/ast"

type Walker interface {
	Walk(node ast.Node, entering bool) error
}

type WalkerWithHandlers struct {
	handlers map[ast.NodeKind]func(ast.Node, bool) (ast.WalkStatus, error)
}

type HandlerSet struct {
	Kind ast.NodeKind
	Func func(ast.Node, bool) (ast.WalkStatus, error)
}

func NewWalkerWithHandlers(handlers []HandlerSet) Walker {
	w := &WalkerWithHandlers{
		handlers: map[ast.NodeKind]func(ast.Node, bool) (ast.WalkStatus, error){},
	}
	for _, h := range handlers {
		w.handlers[h.Kind] = h.Func
	}
	return w
}

func (w *WalkerWithHandlers) Walk(node ast.Node, entering bool) error {
	return ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		kind := n.Kind()
		if handler, ok := w.handlers[kind]; ok {
			return handler(n, entering)
		}
		return ast.WalkContinue, nil
	})
}
