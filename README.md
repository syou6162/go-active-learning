# go-active-learning
[![CircleCI](https://circleci.com/gh/syou6162/go-active-learning.svg?style=shield)](https://circleci.com/gh/syou6162/go-active-learning)
[![Go Report Card](https://goreportcard.com/badge/github.com/syou6162/go-active-learning)](https://goreportcard.com/report/github.com/syou6162/go-active-learning)

go-active-learning is a command line annotation tool for binary classification problem written in Go. It uses simple active learning algorithm to minimize annotation time. 

# Install

```console
% go get github.com/syou6162/go-active-learning
```

## Build from source

```console
% git clone https://github.com/syou6162/go-active-learning.git
% cd go-active-learning
% make build
```

# Usage
go-active-learning has `annotate` (annotate new examples suggested by active learning) mode. To see the detail options, type `./go-active-learning --help`.

## Annotate new examples
To see the detail options, type `./go-active-learning train --help`.

```console
% ./go-active-learning annotate --input-filename tech_input_example.txt --output-filename additionaly_annotated_examples.txt --openurl
Loading cache...
Label this example (Score: 0.600): http://srdk.rakuten.jp/ (それどこ)

p: Label this example as positive.
n: Label this example as negative.
s: Save additionally annotated examples in 'output-filename'.
h: Show this help.
e: Exit.

Label this example (Score: 1.000): http://srdk.rakuten.jp/ (それどこ)
Labeled as negative
Saving cache...
% cat additionaly_annotated_examples.txt
http://srdk.rakuten.jp/ -1
```

# Author
Yasuhisa Yoshida
