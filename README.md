# BRBundle

## Flow Chart

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
  
* Exe Bundle

  This tool appends content files to your application. You can add them after compiling your application.
  This is the best for application development.

* Packed Bundle

  This tool generates one single binary that includes content files.

* Folder Bundle

  For debugging. You can access content files without any building tasks.


## Selecting Compression Method

BRBundle chooses compression format Brotli and LZ4 automatically.

That option you can choose is using ``-z`` or not.

* If your application is a web application server, always turn on ``-z``
* Otherwise, decide using ``-z`` from size and booting speed.

``-z`` option makes the content compressed with Brotli.
Brotli is more effective lossless compression algorithm than gzip and deflate.
BRBundle returns compressed content directly for browsers because almost all browsers
[supports ``Content-Encoding: br``](https://caniuse.com/#search=brotli).

But Brotli has some cons.
If the content is already compressed (like PNG, JPEG, OpenOffice formats), compression ratio is not effective.
And loading compressed contents is slower than uncompressed content.
Even if turned off Brotli, BRBundle fall back to LZ4. So the content like JSON becomes smaller than original
and not slower than uncompressed content so much.

Now, current code skip compression if the content size after compression is not enough small:

* var u: int = uncompressed_size
* var c: int = compressed_size
* var enough_small: bool = (u - 1000 > c) || (u > 10000 && (u * 0.90 > c))


