[![GoDoc](https://godoc.org/github.com/NBR41/gocrontranslate/translator?status.svg)](https://godoc.org/github.com/NBR41/gocrontranslate/translator)
[![Build Status](https://travis-ci.org/NBR41/gocrontranslate.svg?branch=master)](https://travis-ci.org/NBR41/gocrontranslate)
[![Coverage Status](http://codecov.io/gh/NBR41/gocrontranslate/branch/master/graph/badge.svg)](http://codecov.io/gh/NBR41/gocrontranslate)


# gocrontranslate

A simple tool to translate a crontab format to text

## Installation

```
go get github.com/NBR41/gocrontranslate
```

## Usage

```
./gocrontranslate "1 2 3 4 *"
```
displays
```
at 02h01 of every 3 of April
```

The package github.com/NBR41/gocrontranslate/translator can also be used in stand alone to get the translation.
