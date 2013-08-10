# go-autobuild

go-autobuild builds your go project every time a go source file has been added, modified or deleted.

## Motivation

go-autobuild is not the first application that builds your go applications every time a file changes. For example there are

- [go-gb](https://github.com/skelterjohn/go-gb)
- [gowatch](https://bitbucket.org/gotamer/gowatch)
- [buildwatch](https://bitbucket.org/jzs/buildwatch)
- [gobuild](https://code.google.com/p/gobuild/)
- [goautotest](https://github.com/ryanslade/goautotest)
- [GoMon](https://github.com/aaudis/GoMon)
- [go-bldbot](https://github.com/sbinet/go-bldbot)
- and many others ...

Most of them use [fsnotify](https://github.com/howeyc/fsnotify) which has the drawback of not (natively) supporting [recursive watchers](https://github.com/howeyc/fsnotify/issues/56) and operating system limitations for the number of watched files.

Even though I must admit that adding recursive support to fsnotify isn't that hard and the limitition for the maxiumum number
of watches files and folders on some operating system isn't likely to be a problem for small go applications I chose to build my autobuild tool upon [go-fswatch](https://github.com/andreaskoch/go-fswatch) - a filesystem watcher which does not depend on inotify.

## Usage

Just start go-autobuild in your project folder and it will build your project every time you change a file:

```bash
cd ~/dev/gocode/src/<your-project-name>
go-autobuild
```

## Build Status

[![Build Status](https://travis-ci.org/andreaskoch/go-autobuild.png?branch=master)](https://travis-ci.org/andreaskoch/go-autobuild)

## Contribute

If you have an idea how to make this little tool better please send me a message or a pull request.

All contributions are welcome.