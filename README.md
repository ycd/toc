<div align="center">
<h1>toc</h1>

[toc](https://github.com/ycd/toc) TOC, table of content generator for Markdown files


![toc gif](assets/toc.gif)

</div>


# Table of Contents

<!--toc-->
- [Table of Contents](#table-of-contents)
    * [Usage](#usage)
    * [Installation](#installation)
        * [Packages](#packages)
            * [Arch Linux](#arch-linux)
            * [Homebrew](#homebrew)
            * [Docker](#docker)
        * [Downloads](#downloads)
        * [Installation from source](#installation-from-source)
            * [Unix/Linux](#unixlinux)
    * [Contributing](#contributing)
    * [Licence](#licence)

<!-- tocstop -->

---

## Usage



```
Usage: toc [options]
Options:
	-p, --path      <path>   Path for the markdown file.                               [REQUIRED]
	-a, --append    <bool>   Append toc after <!--toc-->, or write to stdout.          [Default: true]
	-b, --bulleted  <bool>   Write as bulleted, or write as numbered list.             [Default: true] 
	-s, --skip      <int>    Skip the first given number of headers.                   [Default: 0]
	-d, --depth     <int>    Set the number of maximum heading level to be included.   [Default: 6]
	-h, --help               Show this message and exit.
```

Add `<!--toc-->`  to your markdown to the place where you want to add Table of Contents. That's it.

Give the markdown file as an input with `-p`, `--path` flags.

```
$ toc -p path/to/markdown.md
```

Create numbered list instead of bulleted list.

```
$ toc --bulleted=false
```

Write result to standard output instead of appending.

```
$ toc --append=false
```

Skip the first `n` number of headers via `-s`, `--skip` flags.

```
$ toc --skip 2
```

Set the number of maximum heading level to be included with `-d`, `--depth` flags. 

Set maximum heading level to 3 (h3)

```
$ toc --depth 3
```

---


## Installation


### Packages

#### Arch Linux

* [ ] For Arch Linux, install the [``]() package.

#### Homebrew

* [ ] For Homebrew on macOS, install the [``]() formula.

#### Docker 

It is available via two tags.

You can either use `latest` or `$VERSION`.

```sh
docker run --rm -it yagizcan/toc:latest toc
```


### Downloads

Binary downloads of example are available from [the releases section on GitHub](https://github.com/ycd/toc/releases/) for 64-bit Windows, macOS, and Linux targets. They contain the compiled executable.

| platform     |
| ----------- | 
| [macOS 64 Bit](https://github.com/ycd/toc/releases/download/v0.2.5/toc_0.2.5_darwin_x86_64.tar.gz)   
| [Linux 32-Bit](https://github.com/ycd/toc/releases/download/v0.2.5/toc_0.2.5_linux_i386.tar.gz) 
| [Linux ARM 64 Bit](https://github.com/ycd/toc/releases/download/v0.2.5/toc_0.2.5_linux_arm64.tar.gz)    
| [Linux 64 Bit](https://github.com/ycd/toc/releases/download/v0.2.5/toc_0.2.5_linux_x86_64.tar.gz)    
| [Windows 64 Bit](https://github.com/ycd/toc/releases/download/v0.2.5/toc_0.2.5_windows_x86_64.zip)       
| [Windows 32 Bit](https://github.com/ycd/toc/releases/download/v0.2.5/toc_0.2.5_windows_i386.zip)       



### Installation from source

0. Verify that you have Go 1.13+ installed

   ```
   $ go version
   ```

   If `go` is not installed, follow instructions on [the Go website](https://golang.org/doc/install).

1. Clone this repository

   ```
   $ git clone https://github.com/ycd/toc 
   $ cd toc
   ```

2. Build and install

   #### Unix/Linux
   ```
   $ go build .
   # May require you to use sudo
   $ cp toc /usr/local/bin
   ```
   
3. Verify installation

   ```
   $ toc -h 

   Usage: toc [options]
   Options:
      -p, --path      <path>   Path for the markdown file.                               [REQUIRED]
      -a, --append    <bool>   Append toc after <!--toc-->, or write to stdout.          [Default: true]
      -b, --bulleted  <bool>   Write as bulleted, or write as numbered list.             [Default: true] 
      -s, --skip      <int>    Skip the first given number of headers.                   [Default: 0]
      -d, --depth     <int>    Set the number of maximum heading level to be included.   [Default: 6]
      -h, --help               Show this message and exit.
   ```
---


## Contributing

All kinds of Pull Requests and Feature Requests are welcomed!

## Licence

toc's source code is licenced under the [Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0.txt).
