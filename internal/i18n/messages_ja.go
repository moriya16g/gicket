package i18n

var messagesJA = map[string]string{
	// root
	"root.short": "Git リポジトリ内で動作する分散チケット管理ツール",
	"root.long":  "gicket は Git リポジトリ内にテキスト(YAML)ベースでチケットを管理するツールです。\nWEBサーバ不要で、Git の push/pull だけで開発者間のチケット共有が可能です。",

	// init
	"init.short":   "現在のディレクトリに gicket を初期化する",
	"init.success": "gicket を初期化しました (.gicket/)",

	// new
	"new.short":         "新しいチケットを作成する",
	"new.title.required": "タイトルは必須です (-t フラグで指定)",
	"new.success":       "チケットを作成しました: %s - %s",
	"new.flag.title":    "チケットのタイトル (必須)",
	"new.flag.priority": "優先度 (low/medium/high)",
	"new.flag.label":    "ラベル (複数指定可)",
	"new.flag.assignee": "担当者",

	// list
	"list.short":      "チケットの一覧を表示する",
	"list.no.tickets": "チケットはありません",
	"list.flag.all":   "すべてのステータスのチケットを表示",

	// show
	"show.short": "チケットの詳細を表示する",

	// edit
	"edit.short":            "チケットを編集する",
	"edit.success":          "チケットを更新しました: %s",
	"edit.flag.title":       "タイトル",
	"edit.flag.priority":    "優先度 (low/medium/high)",
	"edit.flag.status":      "ステータス (open/in-progress/closed)",
	"edit.flag.assignee":    "担当者",
	"edit.flag.label":       "ラベル",
	"edit.flag.description": "説明",

	// comment
	"comment.short":         "チケットにコメントを追加する",
	"comment.body.required": "コメント内容は必須です (-m フラグで指定)",
	"comment.success":       "コメントを追加しました: %s",
	"comment.flag.message":  "コメント内容 (必須)",

	// close
	"close.short":          "チケットをクローズする",
	"close.already.closed": "チケット %s は既にクローズされています",
	"close.success":        "チケットをクローズしました: %s - %s",

	// serve
	"serve.short": "Web UI を起動する",
	"serve.long":  "チケット管理用の Web UI サーバーを起動します。ブラウザでチケットの閲覧・作成・編集ができます。",
	"serve.open":  "Opening http://localhost:%d in your browser...",
	"serve.flag.port": "サーバーのポート番号",

	// hook
	"hook.short":              "Git フックを管理する",
	"hook.long":               "commit-msg フックとカスタムマージドライバをインストール/アンインストールします。",
	"hook.install.short":      "Git フックとマージドライバをインストールする",
	"hook.install.long":       "以下をインストールします:\n  - commit-msg フック: コミットメッセージ内のチケットID参照を検証\n  - カスタムマージドライバ: .gicket/issues/*.yml の3-wayマージを自動処理\n  - .gitattributes: マージドライバの適用ルール",
	"hook.install.success":    "Git フックをインストールしました:",
	"hook.install.commitmsg":  "  ✓ commit-msg フック",
	"hook.install.mergedriver":"  ✓ カスタムマージドライバ (merge.gicket)",
	"hook.install.gitattr":    "  ✓ .gitattributes",
	"hook.install.require.id": "チケットID参照を必須にするには:",
	"hook.uninstall.short":    "Git フックとマージドライバをアンインストールする",
	"hook.uninstall.success":  "Git フックをアンインストールしました",

	// log
	"log.short":      "チケットに関連する Git コミット履歴を表示する",
	"log.long":       "コミットメッセージにチケットIDが含まれるコミットを検索して表示します。",
	"log.header":     "チケット: %s - %s",
	"log.no.commits": "関連するコミットが見つかりません",
	"log.count":      "\n%d 件のコミットが見つかりました",
	"log.flag.count": "表示するコミット数の上限",

	// merge-driver
	"merge_driver.short": "カスタムマージドライバ（git が内部的に呼び出す）",

	// store errors
	"store.gicket.not.found":   ".gicket ディレクトリが見つかりません。'gicket init' で初期化してください",
	"store.dir.create.failed":  "ディレクトリ作成に失敗: %w",
	"store.config.marshal":     "config のマーシャルに失敗: %w",
	"store.config.write":       "config の書き込みに失敗: %w",
	"store.ticket.marshal":     "チケットのマーシャルに失敗: %w",
	"store.ticket.save":        "チケットの保存に失敗: %w",
	"store.ticket.read":        "チケットの読み込みに失敗: %w",
	"store.ticket.parse":       "チケットのパースに失敗: %w",
	"store.ticket.list":        "チケット一覧の取得に失敗: %w",
	"store.ticket.search":      "チケットの検索に失敗: %w",
	"store.ticket.not.found":   "チケット '%s' が見つかりません",
	"store.ticket.ambiguous":   "ID '%s' は複数のチケットに一致します: %v",

	// git errors
	"git.repo.not.found":   "Git リポジトリが見つかりません",
	"git.not.installed":    "git コマンドが見つかりません。Git をインストールしてください",
	"git.hook.exists":      "commit-msg フックが既に存在します。手動でマージしてください: %s",
	"git.hook.install.fail":"commit-msg フックのインストールに失敗: %w",
	"git.merge.driver.fail":"マージドライバの設定に失敗: %w",
	"git.gitattr.fail":     ".gitattributes の設定に失敗: %w",
	"git.hook.remove.fail": "commit-msg フックの削除に失敗: %w",
	"git.gitattr.update":   ".gitattributes の更新に失敗: %w",

	// merge errors
	"merge.ancestor.read": "ancestor の読み込みに失敗: %w",
	"merge.ours.read":     "ours の読み込みに失敗: %w",
	"merge.theirs.read":   "theirs の読み込みに失敗: %w",
	"merge.marshal":       "マージ結果のマーシャルに失敗: %w",
	"merge.write":         "マージ結果の書き込みに失敗: %w",
	"merge.conflict":      "CONFLICT: チケット %s のフィールドが両方のブランチで異なる値に変更されました",
}
