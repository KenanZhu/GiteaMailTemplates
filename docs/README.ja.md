# Gitea メールテンプレート

[Gitea](https://about.gitea.com) のためのメールテンプレートコレクション。

> **110 ファイル — 10 スタイル、各11種類**

---

## スタイルギャラリー

| プレビュー | スタイル | 対象 | 特徴 |
|---|---|---|---|
| ![Horizon](images/horizon.png) | **Horizon** | エンタープライズ | Blue accent, slate, centred cards |
| ![Terminal](images/terminal.png) | **Terminal** | 開発者 | Dark mode, monospace, green CLI |
| ![Ember](images/ember.png) | **Ember** | コミュニティ | Warm amber, rounded, inclusive |
| ![Bloom](images/bloom.png) | **Bloom** | クリエイティブ | Frosted glass, soft blue light |
| ![Heritage](images/heritage.png) | **Heritage** | 教育/研究 | Navy and gold, serif, classic |
| ![Neon](images/neon.png) | **Neon** | ゲーム/Web3 | Cyberpunk neon, pink and cyan |
| ![Mono](images/mono.png) | **Mono** | デザイン/編集 | Swiss brutalist, black-white-red |
| ![Terra](images/terra.png) | **Terra** | サステナビリティ | Warm earth tones, organic textures |
| ![Ink](images/ink.png) | **Ink** | 出版/ニュース | Editorial print, navy and gold |
| ![Aurora](images/aurora.png) | **Aurora** | プレミアムSaaS | Ethereal gradients, purple and teal |

> 画像は[ライブプレビュー](../preview/index.html)からの600px幅のスクリーンショットです。キャプチャ手順は [images/README.md](images/README.md) を参照してください。

[**ライブプレビュー**](../preview/index.html)

---

## インストール

```bash
cp -r themes/horizon/mail/* /var/lib/gitea/custom/templates/mail/
systemctl restart gitea
```

## プレビュー

```bash
go run ./tools/build-preview.go
```
その後 `preview/index.html` を開く。

## 互換性

- **Gitea 1.21+**, 100%互換, 組み込み関数のみ

## ライセンス

MIT — [LICENSE](../LICENSE).
