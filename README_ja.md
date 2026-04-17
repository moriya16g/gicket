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
- **内蔵 Web UI**: カンバンボード、フィルタ、ライト/ダークテーマ対応のブラウザベースインターフェース
- **REST API**: 外部ツール連携のための HTTP API を完備
- **VS Code 拡張**: エディタから直接チケットを管理
- **多言語対応**: 英語（デフォルト）と日本語 UI — 環境変数で切り替え可能

## クイックスタート

### インストール

**バイナリをダウンロード**（Go 環境不要）:

[GitHub Releases](https://github.com/moriya16g/gicket/releases) からお使いのプラットフォーム用の最新版をダウンロードし、PATH の通った場所に配置してください。

**Go でインストール**:

```bash
go install github.com/moriya16g/gicket@latest
```

**ソースからビルドする場合**:

```bash
git clone https://github.com/moriya16g/gicket.git
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

# クローズしたチケットを再オープン
gicket reopen <id>

# キーワードでチケットを検索
gicket search "ログイン"

# チケットの統計情報を表示
gicket stats

# JSON 形式で出力（list, show, search, stats で利用可能）
gicket list --json
gicket show <id> --json
gicket search "bug" --json
gicket stats --json

# Web UI を起動（デフォルト: http://localhost:8080）
gicket serve
gicket serve -p 3000   # ポート指定
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
| `gicket reopen <id>` | クローズしたチケットを再オープン |
| `gicket search <keyword>` | キーワードでチケットを検索 |
| `gicket stats` | チケットの統計情報を表示 |
| `gicket serve` | Web UI サーバを起動（`-p` でポート指定、デフォルト 8080） |
| `gicket hook install` | Git フックとカスタムマージドライバをインストール |
| `gicket hook uninstall` | Git フックとマージドライバをアンインストール |
| `gicket log <id>` | チケットに関連する Git コミット履歴を表示 |

## Web UI

`gicket serve` でブラウザベースのインターフェースを起動できます：

- **ダッシュボード**: チケット数カード（Open / In Progress / Closed / Total）
- **リスト & カンバン表示**: テーブルリストとカンバンボードを切り替え可能
- **フィルタ**: ステータスタブ + 全文検索
- **チケット操作**: 作成・編集・クローズ・コメント — すべてブラウザから実行
- **ライト / ダークテーマ**: ヘッダーで切り替え、設定は localStorage に保存

Web UI は `go:embed` でバイナリに埋め込まれるため、追加ファイルは不要です。

## REST API

`gicket serve` は REST API も提供します：

| メソッド | エンドポイント | 説明 |
|----------|---------------|------|
| `GET` | `/api/tickets` | 全チケットの一覧 |
| `POST` | `/api/tickets` | チケットの新規作成 |
| `GET` | `/api/tickets/{id}` | ID でチケットを取得 |
| `PUT` | `/api/tickets/{id}` | チケットの更新 |
| `DELETE` | `/api/tickets/{id}` | チケットの削除 |
| `POST` | `/api/tickets/{id}/comments` | コメントの追加 |

## Git 連携

### フック & マージドライバ

```bash
# Git フックとマージドライバをインストール
gicket hook install

# アンインストール
gicket hook uninstall
```

`gicket hook install` は以下をセットアップします：

- **commit-msg フック**: コミットメッセージ内のチケットID参照を検証（パターン: `gicket:<ticket-id>`）。`GICKET_HOOK_REQUIRE_ID=1` で必須化可能。
- **カスタムマージドライバ**: `.gicket/issues/*.yml` のマージコンフリクトを3-wayマージで自動解決：
  - 単一フィールドの変更: そのまま適用
  - コメント: 両ブランチのコメントを統合（重複除去＋日付順ソート）
  - ラベル: 削除追跡付き和集合
  - ステータス競合: `closed` > `in-progress` > `open` の優先順位
- **.gitattributes**: チケットファイルへのマージドライバ適用ルール

### コミット履歴

```bash
# チケットに関連するコミットを表示
gicket log <id>
gicket log <id> -n 20   # 20件に制限
```

コミットメッセージにチケットIDが含まれるコミット、またはチケットファイルを変更したコミットを検索します。

## VS Code 拡張

`vscode-extension/` ディレクトリに VS Code 拡張が含まれています：

- **サイドバーツリービュー**: ステータスごとにグループ化（Open / In Progress / Closed）、優先度アイコン付き
- **チケット詳細パネル**: 全フィールド・説明・コメントを表示するリッチな Webview
- **クイックコマンド**: 入力プロンプトでチケットの作成・編集・クローズ・再オープン・コメント追加
- **YAML ファイルアクセス**: 任意のチケットの生 YAML ファイルを開く
- **自動更新**: `.gicket/issues/` の変更をファイルウォッチャーが自動検知
- **ランタイム依存ゼロ**: YAML ファイルを直接読み書き — gicket CLI は不要

### ソースからインストール

```bash
cd vscode-extension
npm install
npm run compile
# VS Code で F5 を押して拡張開発ホストを起動
```

ワークスペースに `.gicket` ディレクトリが存在すると自動的にアクティベートされます。

## 言語設定 / i18n

gicket は英語（デフォルト）と日本語に対応しています。環境変数で言語を切り替えられます：

```bash
# 英語（デフォルト）
gicket list

# 日本語
GICKET_LANG=ja gicket list

# またはシステム全体で設定
export GICKET_LANG=ja
```

検出優先順位: `GICKET_LANG` > `LANG` > 英語。

## ロードマップ

- [x] **Phase 1**: CLI コア
- [x] **Phase 2**: REST API + Web UI（`gicket serve`）
- [x] **Phase 3**: Git 連携（フック、マージコンフリクト解決）
- [x] **Phase 4**: VS Code 拡張
- [x] **v1.0.0**: 入力バリデーション、検索、再オープン、統計、JSON出力、設定ファイル

## 類似プロジェクト

| プロジェクト | 言語 | アプローチ |
|-------------|------|-----------|
| [git-bug](https://github.com/git-bug/git-bug) | Go | Git オブジェクトとして保存（ファイルではない） |
| [git-issue](https://github.com/dspinellis/git-issue) | Shell | `.issues/` にテキストファイルで保存 |
| [SIT](https://github.com/sit-fyi/sit) | Rust | サーバレス情報トラッカー |
| [Bugs Everywhere](http://www.bugseverywhere.org/) | Python | 複数VCS対応 |

**gicket** は、人間が読める YAML ファイルと内蔵 Web UI を組み合わせ、シングルバイナリで提供することで差別化を図ります。

## ライセンス

MIT License. 詳細は [LICENSE](LICENSE) を参照してください。
