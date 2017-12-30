package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kjelly/vimwiki-server/searcher"
	"github.com/kjelly/vimwiki-server/utils"
)

func main() {
	searcherEngine := searcher.New()
	searcherEngine.BuildIndex("/home/kjelly/Dropbox/vimwiki/")

	app := iris.New()
	app.RegisterView(iris.HTML("./views", ".html"))
	app.StaticWeb("/static", "./static")
	app.Get("/search", func(ctx context.Context) {
		fmt.Printf("This is search1\n")
		key := ctx.FormValue("key")
		fmt.Printf("Key=%s\n", key)
		result, _ := searcherEngine.Search(key)
		ctx.ViewData("text", result)
		ctx.View("test.html")
	})
	app.Get("/*", func(ctx context.Context) {
		input, err := utils.ReadFile(ctx.Path())
		if err != nil {
			ctx.StatusCode(404)
			ctx.WriteString(err.Error())
		} else {
			ctx.ViewData("content", utils.MarkdownToHTML(input))
			ctx.View("markdown.html")
		}

	})
	app.Run(iris.Addr("0.0.0.0:8080"))
}
