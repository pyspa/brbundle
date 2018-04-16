package brbundle

import (
	"io"
	"strings"

	"github.com/labstack/echo"
)

func EchoMount(path string, bundle *Bundle) func(echo.Context) error {
	return func(c echo.Context) error {
		newPath := strings.TrimPrefix(c.Path(), path)
		if newPath == c.Path() {
			return echo.ErrNotFound
		}

		if !strings.HasPrefix(newPath, "/") {
			newPath = "/" + newPath
		}

		file, err := bundle.Find(newPath)
		if err != nil {
			return echo.ErrNotFound
		}

		if supportBrotli(c.Request().Header.Get("Accept-Encoding")) {
			reader, err := file.BrotliReader()
			c.Response().Header().Set("Content-Encoding", "br")
			defer reader.Close()
			if err != nil {
				return echo.NewHTTPError(500, "Internal Server Error")
			}
			io.Copy(c.Response(), reader)
		} else {
			reader, err := file.Reader()
			defer reader.Close()
			if err != nil {
				return echo.NewHTTPError(500, "Internal Server Error")
			}
			io.Copy(c.Response(), reader)
		}
		return nil
	}
}
