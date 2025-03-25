<!-- Code generated by gomarkdoc. DO NOT EDIT -->

[![Go](https://github.com/engmtcdrm/go-zippy/actions/workflows/test.yml/badge.svg)](https://github.com/engmtcdrm/go-zippy/actions/workflows/test.yml)
[![Release](https://img.shields.io/github/v/release/engmtcdrm/go-zippy.svg?label=Latest%20Release)](https://github.com/engmtcdrm/go-zippy/releases/latest)

# zippy

```go
import "github.com/engmtcdrm/go-zippy"
```

## Index

- [func Contents\(zipFile string\) \(\[\]\*zip.File, error\)](<#Contents>)
- [type Unzippy](<#Unzippy>)
  - [func NewUnzippy\(path string\) \*Unzippy](<#NewUnzippy>)
  - [func \(u \*Unzippy\) Extract\(\) \(\[\]\*zip.File, error\)](<#Unzippy.Extract>)
  - [func \(u \*Unzippy\) ExtractFiles\(files ...string\) \(\[\]\*zip.File, error\)](<#Unzippy.ExtractFiles>)
  - [func \(u \*Unzippy\) ExtractFilesTo\(dest string, files ...string\) \(\[\]\*zip.File, error\)](<#Unzippy.ExtractFilesTo>)
  - [func \(u \*Unzippy\) ExtractTo\(dest string\) \(\[\]\*zip.File, error\)](<#Unzippy.ExtractTo>)
- [type UnzippyInterface](<#UnzippyInterface>)
- [type Zippy](<#Zippy>)
  - [func NewZippy\(path string\) \*Zippy](<#NewZippy>)
  - [func \(z \*Zippy\) Add\(files ...string\) error](<#Zippy.Add>)
  - [func \(z \*Zippy\) Copy\(dest string, files ...string\) error](<#Zippy.Copy>)
  - [func \(z \*Zippy\) Delete\(files ...string\) error](<#Zippy.Delete>)
  - [func \(z \*Zippy\) Freshen\(\) error](<#Zippy.Freshen>)
  - [func \(z \*Zippy\) Update\(files ...string\) error](<#Zippy.Update>)
- [type ZippyInterface](<#ZippyInterface>)


<a name="Contents"></a>
## func [Contents](<https://github.com/engmtcdrm/go-zippy/blob/master/contents.go#L10>)

```go
func Contents(zipFile string) ([]*zip.File, error)
```

Contents returns a list of files in the zip archive.

zipFile is the path to the zip archive.

<a name="Unzippy"></a>
## type [Unzippy](<https://github.com/engmtcdrm/go-zippy/blob/master/unzip.go#L19-L23>)



```go
type Unzippy struct {
    Path      string // Path to the zip archive.
    Junk      bool   // Junk specifies whether to junk the path of files when extracting.
    Overwrite bool   // Overwrite specifies whether to overwrite files when extracting.
}
```

<a name="NewUnzippy"></a>
### func [NewUnzippy](<https://github.com/engmtcdrm/go-zippy/blob/master/unzip.go#L28>)

```go
func NewUnzippy(path string) *Unzippy
```

NewUnzippy creates a new Unzippy instance.

path is the path to the zip archive.

<a name="Unzippy.Extract"></a>
### func \(\*Unzippy\) [Extract](<https://github.com/engmtcdrm/go-zippy/blob/master/unzip.go#L91>)

```go
func (u *Unzippy) Extract() ([]*zip.File, error)
```

Extract all files from zip archive to the same directory as the archive.

<a name="Unzippy.ExtractFiles"></a>
### func \(\*Unzippy\) [ExtractFiles](<https://github.com/engmtcdrm/go-zippy/blob/master/unzip.go#L108>)

```go
func (u *Unzippy) ExtractFiles(files ...string) ([]*zip.File, error)
```

ExtractFiles extracts the specified files from the zip archive.

files to be extracted. If no files are specified, all files will be extracted. Glob patterns are supported.

<a name="Unzippy.ExtractFilesTo"></a>
### func \(\*Unzippy\) [ExtractFilesTo](<https://github.com/engmtcdrm/go-zippy/blob/master/unzip.go#L120>)

```go
func (u *Unzippy) ExtractFilesTo(dest string, files ...string) ([]*zip.File, error)
```

ExtractFilesTo extracts the specified files from the zip archive to a destination directory. The destination directory will be created if it does not exist. The file modification times will be preserved.

dest is the destination directory.

files to be extracted. If no files are specified, all files will be extracted. Glob patterns are supported.

<a name="Unzippy.ExtractTo"></a>
### func \(\*Unzippy\) [ExtractTo](<https://github.com/engmtcdrm/go-zippy/blob/master/unzip.go#L100>)

```go
func (u *Unzippy) ExtractTo(dest string) ([]*zip.File, error)
```

ExtractTo extracts all files from the zip archive to a destination directory. The destination directory will be created if it does not exist. The file modification times will be preserved.

dest is the destination directory.

<a name="UnzippyInterface"></a>
## type [UnzippyInterface](<https://github.com/engmtcdrm/go-zippy/blob/master/unzip.go#L12-L17>)



```go
type UnzippyInterface interface {
    Extract() ([]*zip.File, error)
    ExtractTo(dest string) ([]*zip.File, error)
    ExtractFiles(files ...string) ([]*zip.File, error)
    ExtractFilesTo(dest string, files ...string) ([]*zip.File, error)
}
```

<a name="Zippy"></a>
## type [Zippy](<https://github.com/engmtcdrm/go-zippy/blob/master/zip.go#L35-L41>)



```go
type Zippy struct {
    Path string // Path is the path to the zip archive.
    Junk bool   // Junk specifies whether to junk the path when archiving.
    // contains filtered or unexported fields
}
```

<a name="NewZippy"></a>
### func [NewZippy](<https://github.com/engmtcdrm/go-zippy/blob/master/zip.go#L43>)

```go
func NewZippy(path string) *Zippy
```



<a name="Zippy.Add"></a>
### func \(\*Zippy\) [Add](<https://github.com/engmtcdrm/go-zippy/blob/master/zip.go#L54>)

```go
func (z *Zippy) Add(files ...string) error
```

Add adds files or directories to a zip archive.

files are the files or directories to archive. Glob patterns are supported.

<a name="Zippy.Copy"></a>
### func \(\*Zippy\) [Copy](<https://github.com/engmtcdrm/go-zippy/blob/master/zip.go#L216>)

```go
func (z *Zippy) Copy(dest string, files ...string) error
```

Copy copies files from existing zip archive to a new zip archive.

dest is the new zip archive path. files are the files to copy. If no files are provided, all files will be copied.

<a name="Zippy.Delete"></a>
### func \(\*Zippy\) [Delete](<https://github.com/engmtcdrm/go-zippy/blob/master/zip.go#L126>)

```go
func (z *Zippy) Delete(files ...string) error
```

Delete deletes files or directories from an existing zip archive.

files are the files or directories to delete. Glob patterns are supported.

<a name="Zippy.Freshen"></a>
### func \(\*Zippy\) [Freshen](<https://github.com/engmtcdrm/go-zippy/blob/master/zip.go#L207>)

```go
func (z *Zippy) Freshen() error
```

Freshen updates files in a zip archive if the source file is newer.

<a name="Zippy.Update"></a>
### func \(\*Zippy\) [Update](<https://github.com/engmtcdrm/go-zippy/blob/master/zip.go#L201>)

```go
func (z *Zippy) Update(files ...string) error
```

Update updates files in a zip archive.

<a name="ZippyInterface"></a>
## type [ZippyInterface](<https://github.com/engmtcdrm/go-zippy/blob/master/zip.go#L12-L33>)

ZippyInterface defines the methods for working with zip archives.

```go
type ZippyInterface interface {
    // Add adds files or directories to a zip archive.
    //
    // files are the files or directories to archive. Glob patterns are supported.
    Add(files ...string) error

    // Delete deletes files or directories from an existing zip archive.
    //
    // files are the files or directories to delete. Glob patterns are supported.
    Delete(files ...string) error

    // Update updates files in a zip archive.
    //
    // files are the files or directories to update.
    Update(files ...string) error

    // Freshen updates files in a zip archive if the source file is newer.
    Freshen() error

    // Copy copies files from a zip archive to a destination directory.
    Copy(dest string, files ...string) error
}
```

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
