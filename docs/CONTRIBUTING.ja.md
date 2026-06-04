# コントリビューション — Gitea メールテンプレート

## 参加方法

### 新しいスタイルの追加

1. ツールでスキャフォールド: `cd tools && go run . create <style-name>` — 全11種類のメールタイプ用のプレースホルダ `.tmpl` ファイルを含む完全なディレクトリ構造を作成します
2. `themes/<style-name>/` 内の各 `.tmpl` ファイルを独自デザインで編集
3. プレビューを再生成: `cd tools && go run . preview all`（ビルドスクリプトが `themes/` 以下の全テーマを自動検出し、テーマセレクターも自動生成されます）
4. スクリーンショット付きで PR を提出（1枚あたり 50 KiB 以下、10–20 KiB 推奨）

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
4. プレビューを再生成: `cd tools && go run . preview all`
5. スタイル名とメールタイプを明記して issue を作成

---

## 開発セットアップ

- **Go 1.21+** テンプレートレンダリングとCLI用
- **Node.js 18+**（任意）開発サーバーとJuice CSSインライン用

### ローカルプレビュー（静的）

1. データを生成: `cd tools && go run . preview all`
2. `preview/index.html` をブラウザで開く

> 静的Gmail/Outlookシミュレーションは参考用です。正確なレンダリングにはdevモードを使用してください。

### 開発サーバー（ライブリロード + CSSインライン）

```bash
cd tools && go run . dev
# → http://localhost:3456
```

`.tmpl` ファイルを編集すると自動的に再構築されブラウザに反映されます。

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
