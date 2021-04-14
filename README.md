# macschema

[![GoDoc](https://godoc.org/github.com/progrium/macschema?status.svg)](https://godoc.org/github.com/progrium/macschema)
<!--a href="https://github.com/progrium/macschema/actions?workflow=test"><img alt="Test workflow" src="https://img.shields.io/github/workflow/status/progrium/macschema/Test?label=test&logo=github&style=flat-square"></a-->
[![Go Report Card](https://goreportcard.com/badge/github.com/progrium/macschema)](https://goreportcard.com/report/github.com/progrium/macschema)
<a href="https://twitter.com/progriumHQ" title="@progriumHQ on Twitter"><img src="https://img.shields.io/badge/twitter-@progriumHQ-55acee.svg" alt="@progriumHQ on Twitter"></a>
<a href="https://github.com/progrium/macschema/discussions" title="Project Forum"><img src="https://img.shields.io/badge/community-forum-ff69b4.svg" alt="Project Forum"></a>
<a href="https://github.com/sponsors/progrium" title="Sponsor Project"><img src="https://img.shields.io/static/v1?label=sponsor&message=%E2%9D%A4&logo=GitHub" alt="Sponsor Project" /></a>

------
Toolchain for generating JSON definitions for Apple APIs like this:

```json
{
  "Class": {
    "Name": "NSScreen",
    "Description": "An object that describes the attributes of a computer’s monitor or screen.",
    "Declaration": "@interface NSScreen : NSObject",
    "InstanceMethods": [
      {
        "Name": "convertRectFromBacking:",
        "Description": "Converts the rectangle from the device pixel aligned coordinates system of a screen.",
        "Declaration": "- (NSRect)convertRectFromBacking:(NSRect)rect;",
        "Return": {
          "Name": "NSRect"
        },
        "Args": [
          {
            "Name": "rect",
            "Type": {
              "Name": "NSRect"
            }
          }
        ],
        "TopicURL": "https://developer.apple.com/documentation/appkit/nsscreen/1388364-convertrectfrombacking?language=objc"
      },
      ...
    ],
    "InstanceProperties": [ ... ],
    "TypeProperties": [ ... ],
    "Frameworks": [
      "AppKit"
    ],
    "Platforms": [
      "macOS 10.0+"
    ],
    "TopicURL": "https://developer.apple.com/documentation/appkit/nsscreen?language=objc",
  },
  "Kind": "class",
  "PullDate": "2021-04-14T18:22:57.729056-05:00",
  "Version": 2
}
```

## Installing macschema

```
$ go get github.com/progrium/macschema
```

Chrome is required for downloading topic data. You can also use headless Chrome in Docker. We recommend [chromedp/headless-shell](https://github.com/chromedp/docker-headless-shell).

## Using macschema

The `macschema` tool has several subcommands for downloading topics from Apple documentation and parsing topics into schemas. The commands will assume they can use two directories in the working directory: `api` and `doc`, where schemas and topics are downloaded and saved as JSON. 

To pull and show a schema, which will download relevant topics and parse into schema:
```
$ macschema pull appkit/nswindow --show
```

Other commands:
```
$ macschema
Generates JSON definitions for Apple APIs

Usage:
  macschema [command]

Available Commands:
  crawl       Downloads topics linked from a topic to doc dir
  fetch       Download a topic to doc dir
  help        Help about any command
  pull        Generate a schema in api dir fetching topics if needed

Flags:
  -h, --help          help for macschema
      --lang string   use language (default "objc")
      --show          show resulting JSON to stdout
  -v, --version       version for macschema

Use "macschema [command] --help" for more information about a command.
```

## Project Status

Currently able to generate schemas for most classes, but other high level constructs coming soon:

* [x] Classes
* [ ] Functions
* [ ] Typedefs and enums
* [ ] Constants / variables

Currently it focuses on Objective-C APIs, but is designed to support Swift in the future if needed.

## Declaration Parsing / AST

There is a lexer/parser system and AST for Objective-C declarations in `declparse`. This is where
most development will happen to support new language constructs, so if you run into a problem it
may involve a declaration that has not been added to tests.

For debugging parser issues, you can use a lexer tool to see what tokens the parser is working with:

```
$ echo "@interface NSScreen : NSObject" | go run ./tools/lexer/main.go
```

## About

macschema come out of the [macdriver project](https://github.com/progrium/macdriver), primarily for code generation use.

MIT Licensed