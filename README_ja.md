# gicket

**Git リポジトリに埋め込む分散チケット管理ツール**

[English](README.md)

gicket は、Git リポジトリ内に人間が読める YAML テキストファイルとしてチケットを管理するツールです。Web サーバ不要、データベース不要、ベンダーロックインなし — Git だけで完結します。

## 特徴

- **テキストベース**: チケットは誰でも読み書きできるプレーンな YAML ファイル
- **分散型**: 標準の `git push` / `git pull` でチケットを共有
- **インフラ不要**: サーバもデータベースも不要 — ファイルシステムだけで動作
- **ベンダー非依存**: データは自分のリポジトリに存在し、外部プラットフォームに依存しない
- **シングルバイナリ**: 実行ファイル1つ、依存なし、クロスプラットフォーム（Windows / macOS / Linux）
- **短縮ID**: フルIDの代わりにユニークな前方一致でチケットを参照可能

## クイックスタート

### インストール

```bash
go install github.com/gicket/gicket@latest
```

ソースからビルドする場合:

```bash
git clone https://github.com/gicket/gicket.git
cd gicket
go build -o gicket .
```

### 使い方

```bash
# リポジトリに gicket を初期化
gicket init

# 新しいチケットを作成
gicket new -t "ログイン画面のバリデーション追加" -p high -l bug,frontend

# オープン状態のチケット一覧
gicket list

# 全チケット一覧（クローズ済み含む）
gicket list --all

# チケットの詳細表示（フルID または 前方一致）
gicket show 20260416-200633-709268
gicket show 20260416-20          # 前方一致

# チケットの編集
gicket edit <id> -s in-progress -a "dev@example.com"
gicket edit <id> -d "詳細な説明をここに記述"

# コメントの追加
gicket comment <id> -m "対応を開始しました"

# チケットのクローズ
gicket close <id>
```

### チームでの共有

チケットは Git で管理されるプレーンファイルなので:

```bash
git add .gicket/
git commit -m "チケットを追加"
git push
```

他の開発者は `git pull` するだけでチケットの更新を受け取れます。

## データ形式

チケットは `.gicket/issues/` に YAML ファイルとして保存されます:

```yaml
id: 20260416-200633-709268
title: ログイン画面のバリデーション追加
status: open
priority: high
assignee: tanaka@example.com
labels:
    - bug
    - frontend
created: 2026-04-16T20:06:33+09:00
updated: 2026-04-16T20:07:15+09:00
author: 田中太郎 <tanaka@example.com>
description: |
    ログイン画面でメールアドレスの形式チェックが不足している。
comments:
    - author: 鈴木花子 <suzuki@example.com>
      date: 2026-04-16T21:00:00+09:00
      body: 確認しました。対応します。
```

## ディレクトリ構成

```
your-project/
├── .gicket/
│   ├── config.yml        # プロジェクト設定
│   └── issues/
│       ├── 20260416-200633-709268.yml
│       ├── 20260416-200633-f16bab.yml
│       └── ...
├── src/
└── ...
```

## コマンド一覧

| コマンド | 説明 |
|----------|------|
| `gicket init` | 現在のディレクトリに gicket を初期化 |
| `gicket new` | 新しいチケットを作成 |
| `gicket list` | オープン状態のチケット一覧（`--all` で全件） |
| `gicket show <id>` | チケットの詳細表示 |
| `gicket edit <id>` | チケットのフィールドを編集 |
| `gicket comment <id>` | チケットにコメントを追加 |
| `gicket close <id>` | チケットをクローズ |

## ロードマップ

- [x] **Phase 1**: CLI コア（現在）
- [ ] **Phase 2**: REST API + Web UI（`gicket serve`）
- [ ] **Phase 3**: カンバンボード、リアルタイムフィルタ、ダッシュボード
- [ ] **Phase 4**: VS Code 拡張

## 類似プロジェクト

| プロジェクト | 言語 | アプローチ |
|-------------|------|-----------|
| [git-bug](https://github.com/git-bug/git-bug) | Go | Git オブジェクトとして保存（ファイルではない） |
| [git-issue](https://github.com/dspinellis/git-issue) | Shell | `.issues/` にテキストファイルで保存 |
| [SIT](https://github.com/sit-fyi/sit) | Rust | サーバレス情報トラッカー |
| [Bugs Everywhere](http://www.bugseverywhere.org/) | Python | 複数VCS対応 |

**gicket** は、人間が読める YAML ファイルとリッチな Web UI（今後実装予定）を組み合わせ、シングルバイナリで提供することで差別化を図ります。

## ライセンス

MIT License. 詳細は [LICENSE](LICENSE) を参照してください。
