package brecho

import (
	"github.com/shibukawa/brbundle/websupport"
	"io"
	"net/url"
	"strings"

	"github.com/labstack/echo"
	"github.com/shibukawa/brbundle"
)

func Mount(option ...brbundle.WebOption) func(echo.Context) error {
	o := websupport.InitOption(option)
	return func(c echo.Context) error {
		var err error
		p := c.Request().URL.Path
		if strings.HasSuffix(c.Path(), "*") { // When serving from a group, e.g. `/static*`.
			p = c.Param("*")
		}
		p, err = url.PathUnescape(p)
		if err != nil {
			return err
		}

		file, found, redirectDir := websupport.FindFile(p, o)
		if redirectDir {
			return c.Redirect(301, "./")
		} else if !found {
			return echo.ErrNotFound
		}

		reader, etag, headers, err := websupport.GetContent(file, o, c.Request().Header.Get("Accept-Encoding"))
		if err != nil {
			return echo.NewHTTPError(500, "Internal Server Error")
		}
		defer reader.Close()
		wHeaders := c.Response().Header()
		for _, header := range headers {
			wHeaders.Set(header[0], header[1])
		}
		if c.Request().Header.Get("If-None-Match") == etag {
			c.NoContent(304)
		} else {
			defer reader.Close()
			io.Copy(c.Response(), reader)
		}
		return nil
	}
}

