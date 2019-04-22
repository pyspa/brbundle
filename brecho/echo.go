package brecho

import (
	"io"
	"strings"

	"github.com/labstack/echo"
	"github.com/shibukawa/brbundle"
)

func Mount(option ...brbundle.WebOption) func(echo.Context) error {
	var o brbundle.WebOption
	if len(option) > 0 {
		o = option[0]
	}
	if o.Repository == nil {
		o.Repository = brbundle.DefaultRepository
	}

	return func(c echo.Context) error {
		newPath := c.Path()

		if strings.HasPrefix(newPath, "/") {
			newPath = newPath[1:]
		}

		file, err := o.Repository.Find(newPath)
		if err != nil {
			if o.SPAFallback != "" {
				file, err = o.Repository.Find(o.SPAFallback)
				if err != nil {
					return echo.ErrNotFound
				}
			} else {
				return echo.ErrNotFound
			}
		}

		brreader, err := file.BrotliReader()
		if brreader != nil && brbundle.HasSupportBrotli(c.Request().Header.Get("Accept-Encoding")) {
			reader, err := file.BrotliReader()
			c.Response().Header().Set("Content-Encoding", "br")
			if err != nil {
				return echo.NewHTTPError(500, "Internal Server Error")
			}
			defer reader.Close()
			io.Copy(c.Response(), reader)
		} else {
			reader, err := file.Reader()
			if err != nil {
				return echo.NewHTTPError(500, "Internal Server Error")
			}
			defer reader.Close()
			io.Copy(c.Response(), reader)
		}
		return nil
	}
}
