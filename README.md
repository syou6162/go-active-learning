# go-active-learning
[![CircleCI](https://circleci.com/gh/syou6162/go-active-learning.svg?style=shield)](https://circleci.com/gh/syou6162/go-active-learning)
[![Go Report Card](https://goreportcard.com/badge/github.com/syou6162/go-active-learning)](https://goreportcard.com/report/github.com/syou6162/go-active-learning)
[![Coverage Status](https://coveralls.io/repos/github/syou6162/go-active-learning/badge.svg?branch=master)](https://coveralls.io/github/syou6162/go-active-learning?branch=master)

go-active-learning is a command line annotation tool for binary classification problem written in Go. It uses simple active learning algorithm to reduce annotation time.

# Install

```console
% go get github.com/syou6162/go-active-learning
```

## Build from source

```console
% git clone https://github.com/syou6162/go-active-learning.git
% cd go-active-learning
% createdb go-active-learning
% createdb go-active-learning-test
% sql-migrate up -env=local
% sql-migrate up -env=test
% make build
```

# Usage
go-active-learning has `annotate` (annotate new examples suggested by active learning) mode and `diagnose` (check label conflicts in training data) mode. To see the detail options, type `./go-active-learning --help`.

## Annotation model
To see the detail options, type `./go-active-learning annotate --help`.

## Annotate new examples from command line interface
To see the detail options, type `./go-active-learning annotate cli --help`.

```console
% ./go-active-learning annotate cli --open-url
Loading cache...
Label this example (Score: 0.600): http://srdk.rakuten.jp/ (それどこ)

p: Label this example as positive.
n: Label this example as negative.
s: Skip this example.
h: Show this help.
e: Exit.

Label this example (Score: 1.000): http://srdk.rakuten.jp/ (それどこ)
Labeled as negative
```

## Annotate new examples from slack
To see the detail options, type `./go-active-learning annotate cli --help`. To annotate new examples from slack, you need to create slack bot, and obtain token from [here](https://my.slack.com/services/new/bot). You can pass token via environmental variable (`SLACK_TOKEN`).

```console
% export SLACK_TOKEN=xoxb-SLACK-TOKEN
% ./go-active-learning annotate slack --filter-status-code-ok --channel CHANNEL_ID
```

## Diagnosis model
To see the detail options, type `./go-active-learning diagnose --help`.

### Diagnose training data
This subcommand diagnoses label conflicts in training data. 'conflict' means that an annotated label is '-1/1', but a predicted label by model is '1/-1'. In the above example, `http://www3.nhk.or.jp/news/` is a conflict case ('Label' is -1, but 'Score' is positive). You may need to collect such news articles to train a good classifier.

```console
% ./go-active-learning diagnose label-conflict
Loading cache...
Index   Label   Score   URL     Title
0       -1      0.491   http://www3.nhk.or.jp/news/
1       1       0.491   http://blog.yuuk.io/
2       1       0.491   http://www.yasuhisay.info/
3       -1      -3.057  http://r.gnavi.co.jp/g-interview/       ぐるなび みんなのごはん
4       1       4.264   http://hakobe932.hatenablog.com/        hakobe-blog ♨
5       -1      -7.151  http://suumo.jp/town/   SUUMOタウン
6       -1      -26.321 https://www.facebook.com/       ログイン (日本語)
7       1       44.642  http://www.songmu.jp/riji/      おそらくはそれさえも平凡な日々
8       1       121.170 http://motemen.hatenablog.com/  詩と創作・思索のひろば
Saving cache...
```

### Diagnose feature weight
This subcommand list pairs of feature weight and its name.

```console
% ./go-active-learning diagnose feature-weight --filter-status-code-ok | head -n 10
+0.80   BODY:/
+0.80   BODY:ほか
+0.80   BODY:郁
+0.80   BODY:単行本
+0.80   BODY:姿
+0.80   BODY:暗黙
+0.80   BODY:創造
+0.80   BODY:企業
+0.80   BODY:野中
+0.80   BODY:準備
```

# Author
Yasuhisa Yoshida
