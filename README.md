[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/mintak21/qiitanwrapper/Unit%20Test?label=unit%20test&style=plastic)](https://github.com/mintak21/qiitanwrapper/actions/workflows/ci.yml)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/mintak21/qiitanWrapper/Build%20And%20Push%20Docker%20Image?label=build&style=plastic)](https://github.com/mintak21/qiitanwrapper/actions/workflows/build.yml)
![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/mintak21/qiitanwrapper/master?label=go&style=plastic)
![Codecov](https://img.shields.io/codecov/c/github/mintak21/qiitanWrapper?style=plastic)
![GitHub](https://img.shields.io/github/license/mintak21/qiitanWrapper?color=aquamarine&label=license&style=plastic)
![GitHub last commit (branch)](https://img.shields.io/github/last-commit/mintak21/qiitanwrapper/master?style=plastic)

# QiitaAPIのWrapper

k8sクラスタ上で動くAPIのサンプルという意味合いで作成
以下を提供

- 指定タグがついた記事を特定数取得するAPI
- 指定タグがついた記事を投稿日付ごとに特定数取得するAPI
- 指定月に登録された記事のうち、ストック数の多い順に50記事取得するAPI

## ディレクトリ構成

```text
├── api/               # apiコード
├── cmd/               # エントリポイント
├── deployment/
│  └── dockerfile/     # Dockerfile
├── gen/               # Swagger自動生成コード
├── go.mod
├── go.sum
├── LICENSE
├── Makefile
├── README.md
└── swagger/           # ドキュメント
```
