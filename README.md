# go-ddd-sample

GoでDDDの基本構成を試すためのサンプルです。ドメインの書き方を矯正する linter を強化する前提で、最小のモデルを用意しています。

## 目標

- Entity / Value Object / Repository / Usecase の責務分離
- ドメイン不変条件の明確化
- linterで検出しやすい構造の土台作り

## ディレクトリ構成

```
internal/
  domain/
    user/
      user.go
      user_id.go
      user_name.go
      email.go
      repository.go
  usecase/
    register_user.go
  infrastructure/
    memory/
      user_repository.go
```

## 使い方

### テスト

```
make test
```

### フォーマット

```
make fmt
```

### 静的解析

```
make vet
```

## 現在のユースケース

- ユーザー登録
  - Name / Email の Value Object でバリデーション
  - 既存メールの重複チェック

## 今後の拡張アイデア

- linter ルール案の追加（例: Value Object 直参照禁止, ドメイン層外の検証禁止）
- アプリケーションサービスの追加
- 永続化（DB）実装
