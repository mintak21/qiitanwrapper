|ci| |cov| |licence|

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
│  └── dockerfile/    # Dockerfile
├── gen/               # Swagger自動生成コード
├── go.mod
├── go.sum
├── LICENSE
├── Makefile
├── README.md
└── swagger/           # ドキュメント
```

.. |ci| image:: https://circleci.com/gh/mintak21/texas-holdem-poker.svg?style=shield&circle-token=dc9af5b436e25a00bb0c3dd4e12cdc8c7aeb2904
   :target: https://circleci.com/gh/mintak21/texas-holdem-poker

.. |cov| image:: https://codecov.io/gh/mintak21/qiitanWrapper/branch/master/graph/badge.svg
  :target: https://codecov.io/gh/mintak21/qiitanWrapper

.. |licence| image:: https://img.shields.io/badge/License-Apache%202.0-blue.svg
  :target: https://github.com/mintak21/qiitanWrapper/blob/master/LICENSE
