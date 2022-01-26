# Server

## 開発環境起動

VSCodeのデバッグ機能が使える環境になっているので、
.air.tomlの指定の箇所をコメントアウトして、
以下の方法で立ち上げればブレークポイントなどを
VSCodeから設定できる環境が立ち上がる。
また、airを使ってhot reloadもできるようになっている。

1. serverディレクトリで`go generate ./...`
   - wire,mockgen,oapi-codegenのinstallをしていないとここでエラーが出る可能性がある
2. プロジェクトルートで`docker compose -f dev/compose.yml up`
3. debuggerを使用する場合はVSCodeでRemote Docker Debugを実行

localhost:8000にアクセスするとphpmyadminが使用できます。

## ディレクトリ構成

```bash
.
├── README.md
├── domain アプリケーションの中核を表すドメインを書く
│   └── values
├── handler APIのハンドラー
│   └── v1
│       ├── api.go
│       ├── checker.go 認証のチェック
│       ├── context.go echo.Contextの操作
│       ├── openapi oapi-codegenによる生成コード
│       │   └── openapi.gen.go
│       └── session.go sessionの操作
├── main.go
├── pkg
│   ├── common wireで必要な型の定義
│   │   └── types.go
│   └── context context.Contextのkey
│       └── keys.go
├── repository dbなどでのデータ永続化
│   ├── db.go dbとの接続の抽象化
│   ├── errors.go
│   ├── gorm2 GORM 2.0でのrepositoryの実装
│   │   ├── db.go db接続など諸々の実装
│   │   ├── db_test.go
│   │   └── tables.go databaseのテーブルに関する構造体
│   ├── lock_types.go
│   └── mock
│       └── db.go テスト用のrepository.DBのモック
├── service アプリケーションの本質の処理
│   └── v1 serviceの実装
├── tools.go
├── wire.go wireによるDIの設定
└── wire_gen.go wireの生成コード
```
