package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"

	"github.com/teekennedy/goldmark-markdown"
)

func indent(level int) string {
	return strings.Repeat("  ", level)
}

func newExtractWalker(src []byte) Walker {
	return NewWalkerWithHandlers([]HandlerSet{
		{ast.KindText, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
			n := node.(*ast.Text)
			fmt.Printf("Text: %+v\n", string(n.Segment.Value(src)))
			return ast.WalkContinue, nil
		}},
	})
}

func Extract(globPatterns []string) error {
	paths, err := FindPaths(globPatterns)
	if err != nil {
		return err
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRenderer(markdown.NewRenderer(markdown.WithHeadingStyle(markdown.HeadingStyleATX))),
	)

	for _, path := range paths {
		fmt.Println(path)

		// Read file
		fileContents, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
			return err
		}

		reader := text.NewReader(fileContents)
		node := md.Parser().Parse(reader)

		err = newExtractWalker(fileContents).Walk(node, true)

		md.Renderer().Render(os.Stdout, fileContents, node)
	}
	return nil
}
