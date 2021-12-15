# Client

## 開発環境構築
Node.js(16.x推奨)が必要。
また、OpenAPI Generatorを使用している関係で
「無料Javaのダウンロード」する必要がある。
```shell
$ npm i
$ npm run gen-api
```

## コマンド
### ホットリロード環境
```shell
$ npm run dev
```

**実行前に一度`$ npm run gen-api`を実行する必要があります**

### ビルド
```shell
$ npm run build
```
ビルド

```shell
$ npm run gen-api
```
APIクライアントコード生成
