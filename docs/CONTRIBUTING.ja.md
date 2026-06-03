# コントリビューション — Gitea メールテンプレート

## 参加方法

### 新しいスタイルの追加

1. スタイルディレクトリを作成: `themes/<style-name>/`
2. 既存スタイルからディレクトリ構造をコピー
3. 11個の `.tmpl` ファイルを独自デザインで実装
4. プレビューを再生成: `go run ./tools/build-preview.go`（ビルドスクリプトが `themes/` 以下の全テーマを自動検出します）
5. `preview/index.html` の `<select id="sel-theme">` に追加
6. スクリーンショット付きで PR を提出

### スタイルガイドライン

- 各スタイルに **全11種類のテンプレート** を含める
- Gitea 組み込み関数のみ使用
- 翻訳キーは Gitea 公式ロケール (`mail.*`) から
- **`.DisplayName` を使用しない** (collaborator, transfer, release, workflow_run, assigned, default)
- 600px 幅のメールクライアント向けにデザイン
- Gmail, Outlook, Apple Mail でテスト

### バグレポート

1. Go 変数が正しいか確認
2. 翻訳キーが Gitea ロケールと一致するか確認
3. `.DisplayName` が誤用されていないか確認
4. プレビューを再生成: `go run ./tools/build-preview.go`
5. スタイル名とメールタイプを明記して issue を作成

---

## 開発セットアップ

### ローカルプレビュー

1. データを生成: `go run ./tools/build-preview.go`
2. `preview/index.html` をブラウザで開く
3. Modern, Gmail, Outlook, Raw source を切り替えて確認

### 結合テスト

Gitea にデプロイ: **Site Administration > Configuration > Mailer > Send Test Email**

---

## コミット規約

- `style(<name>):` — スタイルのテンプレート変更
- `preview:` — プレビューツールの変更
- `tools:` — Go ビルドスクリプトの変更
- `docs:` — ドキュメントと翻訳
- `fix:` — バグ修正
- `project:` — README, LICENSE, メタファイル

## ライセンス

貢献により、MIT ライセンスに同意したものとみなします。
