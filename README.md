# RSS-parser

**RSS-parser** is a tool for parsing RSS news from any site and adding it to
 the database, with user-interface access to finding news in the database by entered letters

## Installation

```
go get github.com/daniilsolovey/rss-parser
```


## Usage

##### -c --config \<path>
Read specified config file. [default: config.toml].

##### --debug
Enable debug messages.

##### -v --version
Print version.

#####  -h --help
Show this help.

##### -u --url
rss url in format: 'https://www.news_site.com/world/rss'

## Build
For build program use command:

```
make build
```

## RUN
For running program use command:

```
./news-parser -u https://www.examplesite.com/world/rss
```