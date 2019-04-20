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

## How To Load

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


