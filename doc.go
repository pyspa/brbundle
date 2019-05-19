// brbundle package provides asset bundler's runtime
// that uses Brotli (https://github.com/google/brotli).
//
// Source repository is here: https://github.com/pyspa/brbundle
//
// Install
//
// To install runtime and commandline tool, type the following command:
//
//   $ go get go.pyspa.org/brbundle/...
//
// Asset Handling
//
// This package provides four kind of bundles to handle assets:
//
// 1. Embedding. Generate .go file that includes binary representation of asset files.
// This is best for libraries and so on that is needed to be go gettable.
//
// 2. Appended to executable. Generate .zip file internally and appended to executables.
// You can replace assets after building.
//
// 3. Single packed binary file. You can specify and import during runtime. I assume this used
// for DLC.
//
// 4. Folder. This is for debugging. You don't have to do anything to import asset
//
// brbundle searches assets the following orders:
//
//   Folder -> Single binary file -> Assets on executable -> Embedded assets
//
// Generate Bundle
//
// The following command generates .go file:
//
//   $ brbundle embedded <src-dir>
//
// The following command append assets to executable:
//
//   $ brbundle bundle <exec-file-path> <src-dir>
//
// The following command generates single packed file:
//
//   $ brbundle pack <out-file-path> <src-dir>
//
// The following command generates asset folder. You can use regular cp command even if you don't have to encrypto assets:
//
//   $ brbundle folder <dest-dir> <src-dir>
//
// Standard Usage
//
// It is easy to use the assets:
//
//   file, err := brbundle.Find("file.png")
//   reader, err := file.Reader()
//   img, err := image.Decode(reader)
//
// Embedded assets and assets appended to executable are available by default.
// The following functions registers assets in single packed file and local folder:
//
//   brbundle.RegisterBundle("masterdata.pb")
//   brbundle.RegisterFolder("frontend/bist")
//
// Web Framework Middlewares
//
// You can save the earth by using brbundle. brbundle middlewares brotli content directly when browser supports it.
// Currently, more than 90% browsers already support (https://caniuse.com/#feat=brotli). brbundle provides
// the following frameworks' middleware:
//
//   * net/http
//   * echo
//   * gin
//   * fasthttp and fastrouter
//   * chi router
//
// net/http:
//
//   m := http.NewServeMux()
//   m.Handle("/static/",
//     http.StripPrefix("/static",
//     brhttp.Mount())
//   http.ListenAndServe("localhost:8000", m)
//
// These middlewares also support SPA(Single Page Application). More detail information is on Angular web site
// (https://angular.io/guide/deployment#routed-apps-must-fallback-to-indexhtml).
//
// All samples are in examples folder: https://github.com/pyspa/brbundle/tree/master/examples
//
// Compression Option
//
// brbundle uses Brotli by default. If you pass --fast/-f option,
// brbundle uses Snappy (https://github.com/google/snappy) instead of Brotli.
// Snappy has low compression ratio, but very fast.
//
// Encryption
//
// brbundle supports encryption. You can generates encryption/decryption key by the following command:
//
//   $ brbundle key-gen
//   pBJ0IB3x4EogUVNqmlI4I0EV9+aGpozmIQvSfF+PLo0NfzeamIeaeXHoTqs
//
// When creating bundles, you can pass the key via --crypto/-c option:
//
//   $ brbundle pack -c pBJ0IB3x4EogUVNqmlI4I0E... images.pb images
//
// Keys should be passed the following functions:
//
//   // for embedded assets
//   // default name is empty string and you can change by using --name/-n option of brbundle command
//   brbundle.SetDecryptoKeyToEmbeddedBundle("name", "pBJ0IB3x4EogUVNqmlI4I0E...")
//   // for executable
//   brbundle.SetDecryptoKeyToExeBundle("pBJ0IB3x4EogUVNqmlI4I0E...")
//   // for bundle
//   brbundle.RegisterBundle("bundle.pb", brbundle.Option{
//     DecryptoKey: "pBJ0IB3x4EogUVNqmlI4I0E...",
//   })
//   // for folder
//   brbundle.RegisterEncryptedFolder("public", "pBJ0IB3x4EogUVNqmlI4I0E...")
//
package brbundle
