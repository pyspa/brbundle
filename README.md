# BRBundle

BRBundle is an asset bundling tool for Go. It is inspired by [go-assets](https://github.com/jessevdk/go-assets),
[go.rice](https://github.com/GeertJohan/go.rice) and so on.

It supports four options to bundle assets to help building libraries, CLI applications,
web applications, mobile applications, JavaScript(Gopher.js), including debugging process. 

## Install

```sh
$ go get github.com/shibukawa/brbundle/...
```

## Bundle Type Flow Chart

```text
+--------------+      Yes
| go gettable? |+------------> Embedded Bundle
+----+---------+                    ^
     |                              |
     | No                           | Yes
     v                              |
+----------------+    Yes      +------------+    No
| Single Binary? |+----------> | Gopher.js? +---------->Exe Bundle
+----+-----------+             +------------+
     |
     | No
     v
+--------+            Yes
| Debug? +--------------------> Folder Bundle
+----+---+
     |
     | No
     v
   Packed Bundle
```

by [asciiflow](http://stable.ascii-flow.appspot.com/#Draw)

## Bundling Options

This tool supports 4 options to bundle assets:

* Embedded Bundle

  This tool generates .go file. You have to generate .go file before compiling your application.
  This option is go-gettable.
  
  ```sh
  brbundle embedded [src-dir]
  ```
  
* Exe Bundle

  This tool appends content files to your application. You can add them after compiling your application.

  ```sh
  brbundle bundle [exe-file] [src-dir]
  ```

* Packed Bundle

  This tool generates one single binary that includes content files.
  It can use for DLC.

  ```sh
  brbundle pack [out-file.pb] [src-dir]
  ```

* Folder Bundle

  For debugging. You can access content files without any building tasks.
  You don't have to prepare with brbundle command except encryption is needed.

## How To Access Content

You can get contents by using ``Find()`` function. If contents are bundled with
embedded bundle or exe bundle, you don't have to call any function to load.

```go
import (
	"github.com/shibukawa/brbundle"
	"image"
	"image/png"
)

func main() {
	file, err := brbundle.Find("file.png")
	reader, err := file.Reader()
	img, err := image.Decode(reader)
}
```

### Getting Contents outside of executable

``RegisterBundle()`` ``RegisterFolder()`` register external contents.

BRBundle searches the contents with the following order:

* folder
* bundle
* exe-bundle
* embedded

```go
import (
	"github.com/shibukawa/brbundle"
)

func main() {
	// load packed content
	brbundle.RegisterBundle("pack.pb")
	
	// load folder content
	brbundle.RegisterFolder("static/public")
}
```

## Compression

BRBundle uses [brotli](https://opensource.googleblog.com/2015/09/introducing-brotli-new-compression.html) by default.
brotli is higher compression ratio with faster decompression speed than gzip.

BRBundle's web application middlewares can send brotli-ed content directly.
Almost all browsers [supports ``Content-Encoding: br``](https://caniuse.com/#search=brotli).

It also supports more faster decompression algorithm [LZ4](https://lz4.github.io/lz4/).

## Encryption

It supports contents encryption by AES. It uses base64 encoded 44 bytes key to encrypto.
You can get your key with key-gen sub command.

Each bundles (embedded, bundle, each packed bundle files, folder bundles) can use separated encryption keys.

```sh
$ brbundle key-gen
yt6TX1eCBuG9GPRl2H6SJMbPNhPLOBxEHpb4kkaWyKUDg/tAZ2aSI3A86fw=

$ brbundle pack -c yt6TX1eCBuG9GPRl2H6SJMbPNhPLOBxEHpb4kkaWyKUDg/tAZ2aSI3A86fw= [src-dir]
```

## Web Application Support

BRBundle support its own middleware for famous web application frameworks.
It doesn't provide ``http.FileSystem`` compatible interface because:

* It returns Brotli compressed content directly when client browser supports brotli.
* It support fallback mechanism for [Single Page Application](https://angular.io/guide/deployment#server-configuration).

### ``net/http``

``"github.com/shibukawa/brbundle/brhttp"`` contains ``http.Handler`` compatible API.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/shibukawa/brbundle"
	"github.com/shibukawa/brbundle/brhttp"
)

// The simplest sample
// The server only returns only brbundle's content
// "/static/index.html" returns "index.html" of bundle.
func main() {
	fmt.Println("Listening at :8080")
	http.ListenAndServe(":8080", brhttp.Mount())
}

// Use ServeMux sample to handle static assets with API handler
func main() {
	m := http.NewServeMux()
	m.Handle("/public/", http.StripPrefix("/public", brhttp.Mount()))
	m.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	fmt.Println("Listening at :8080")
	http.ListenAndServe(":8080", m)
}

// Single Page Application sample
// BRBundle's SPA supports is configured by WebOption of Mount() function
// If no contents found in bundles, it returns the specified content.
func main() {
	m := http.NewServeMux()
	m.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	// Single Page Application is usually served index.html at any location
	// and routing errors are handled at browser.
	//
	// You should mount at the last line, because
	// it consumes all URL requests.
	m.Handle("/",
        brhttp.Mount(brbundle.WebOption{
            SPAFallback: "index.html",
        }),
	)
	fmt.Println("Listening at :8080")
	http.ListenAndServe(":8080", m)
}
```

### Echo

[Echo](https://echo.labstack.com/) is a high performance, extensible, minimalist Go web framework.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/shibukawa/brbundle"
	"github.com/shibukawa/brbundle/brecho"
)

// The simplest sample
func main() {
    e := echo.New()
	// Asterisk is required!
    e.GET("/*", brecho.Mount())
    e.Logger.Fatal(e.Start(":1323"))
}

// Use with echo.Group 
func main() {
	e := echo.New()
    e.GET("/api/status", func (c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })
	g := e.Group("/assets")
	// Asterisk is required!
	g.GET("/*", brecho.Mount())
	e.Logger.Fatal(e.Start(":1323"))
}

// Single Page Application sample
// BRBundle's SPA supports is configured by WebOption of Mount() function
// If no contents found in bundles, it returns the specified content.
//
// Single Page Application is usually served index.html at any location
// and routing errors are handled at browser.
func main() {
	e := echo.New()
	e.GET("/api/status", func (c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	// Use brbundle works as an error handler
	echo.NotFoundHandler = brecho.Mount(brbundle.WebOption{
		SPAFallback: "index.html",
	})
	e.Logger.Fatal(e.Start(":1323"))
}
```

### Chi Router

[Chi router](https://github.com/go-chi/chi) is a lightweight, idiomatic and
composable router for building Go HTTP services.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/shibukawa/brbundle"
	"github.com/shibukawa/brbundle/brchi"
)

// Use with chi.Router
func main() {
	r := chi.NewRouter()
	fmt.Println("You can access index.html at /public/index.html")
	// Asterisk is required!
	r.Get("/public/*", brchi.Mount())
	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	fmt.Println("Listening at :8080")
	http.ListenAndServe(":8080", r)
}

// Single Page Application sample
// BRBundle's SPA supports is configured by WebOption of Mount() function
// If no contents found in bundles, it returns the specified content.
//
// Single Page Application is usually served index.html at any location
// and routing errors are handled at browser.
func main() {
	r := chi.NewRouter()
	r.Get("/api/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	fmt.Println("You can access index.html at any location")
	// Use brbundle as an error handler
	r.NotFound(brchi.Mount(brbundle.WebOption{
		SPAFallback: "index.html",
	}))
	fmt.Println("Listening at :8080")
	http.ListenAndServe(":8080", r)
}
```

### fasthttp / fasthttprouter

[fasthttp](https://github.com/valyala/fasthttp) is a fast http package.
[fasthttprouter](https://github.com/buaazp/fasthttprouter) is a high performance request router that scales well for fasthttp.


```go
package main

import (
	"fmt"

	"github.com/buaazp/fasthttprouter"
	"github.com/shibukawa/brbundle"
	"github.com/shibukawa/brbundle/brfasthttp"
	"github.com/valyala/fasthttp"
)

// The simplest sample
func main() {
	fmt.Println("Listening at :8080")
	fmt.Println("You can access index.html at /index.html")
	fasthttp.ListenAndServe(":8080", brfasthttp.Mount())
}

// Use with fasthttprouter
func main() {
	r := fasthttprouter.New()
	r.GET("/api/status", func (ctx *fasthttp.RequestCtx) {
		ctx.WriteString("Hello, World!")
	})
	// "*filepath" is required at the last fragment of path string
	fmt.Println("You can access index.html at /static/index.html")
	r.GET("/static/*filepath", brfasthttp.Mount())

	fmt.Println("Listening at :8080")
	fasthttp.ListenAndServe(":8080", r.Handler)
}

// Single Page Application sample
// BRBundle's SPA supports is configured by WebOption of Mount() function
// If no contents found in bundles, it returns the specified content.
//
// Single Page Application is usually served index.html at any location
// and routing errors are handled at browser.
func main() {
	r := fasthttprouter.New()
	r.GET("/api/status", func (ctx *fasthttp.RequestCtx) {
		ctx.WriteString("Hello, World!")
	})
	fmt.Println("You can access index.html at any location")
	// Use brbundle works as an error handler
	r.NotFound = brfasthttp.Mount(brbundle.WebOption{
		SPAFallback: "index.html",
	})
	fmt.Println("Listening at :8080")
	fasthttp.ListenAndServe(":8080", r.Handler)
}
```

## Internal Design

### File Format

It uses zip format to make single packed file. Embedded bundles and Exe bundles also use zip format.
It doesn't use Deflate algorithm. It uses Brotli or LZ4 inside it.

### Selecting Compression Method

BRBundle chooses compression format Brotli and LZ4 automatically.

That option you can choose is using ``-f`` (faster) or not.

* If your application is a web application server, always turn off ``-f``
* Otherwise, decide using ``-f`` from size and booting speed.

``-f`` option makes the content compressed with LZ4.

But Brotli has some cons.
If the content is already compressed (like PNG, JPEG, OpenOffice formats), compression ratio is not effective.
And loading compressed contents is slower than uncompressed content.
Even if turned off Brotli, BRBundle fall back to LZ4. So the content like JSON becomes smaller than original
and not slower than uncompressed content so much.

Now, current code skip compression if the content size after compression is not enough small:

* var u: int = uncompressed_size
* var c: int = compressed_size
* var enough_small: bool = (u - 1000 > c) || (u > 10000 && (u * 0.90 > c))


