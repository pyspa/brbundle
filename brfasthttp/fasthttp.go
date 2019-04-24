package brfasthttp

import (
	"io"

	"github.com/valyala/fasthttp"
	"github.com/shibukawa/brbundle"
	"github.com/shibukawa/brbundle/websupport"
)

func Mount(option ...brbundle.WebOption) fasthttp.RequestHandler {
	o := websupport.InitOption(option)

	return func (ctx *fasthttp.RequestCtx) {
		p, ok := ctx.UserValue("filepath").(string)
		if !ok {
			p = string(ctx.Path())
		}

		file, found, redirectDir := websupport.FindFile(p, o)
		if redirectDir {
			ctx.Redirect("./", fasthttp.StatusFound)
			return
		} else if !found {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}

		reader, etag, headers, err := websupport.GetContent(file, o, string(ctx.Request.Header.Peek("Accept-Encoding")))
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		}
		defer reader.Close()

		for _, header := range headers {
			ctx.Response.Header.Set(header[0], header[1])
		}
		if string(ctx.Request.Header.Peek("If-None-Match")) == etag {
			ctx.SetStatusCode(fasthttp.StatusNotModified)
			return
		} else {
			defer reader.Close()
			io.Copy(ctx, reader)
		}
	}
}

func MountRouter(option ...brbundle.WebOption) (string, fasthttp.RequestHandler) {
	return "/*filepath", Mount(option...)
}