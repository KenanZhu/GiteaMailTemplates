# 貢獻指南 — Gitea 郵件模板

## 參與方式

### 新增風格

1. 使用工具腳手架：`cd tools && go run . create <風格名稱>` — 這會建立完整的目錄結構和全部 11 種郵件類型的佔位 `.tmpl` 檔案
2. 在 `themes/<風格名稱>/` 中編寫每個 `.tmpl` 檔案，套用獨特的視覺設計
3. 重新產生預覽：`cd tools && go run . preview all`（建置腳本會自動發現 `themes/` 下的所有主題目錄並動態生成主題選擇器）
4. 提交包含螢幕截圖的 PR（單張圖片 ≤ 50 KiB，建議 10–20 KiB）

### 風格指南

- 每個風格必須包含全部 **11 種模板類型**
- 僅使用 Gitea 內建模板函式
- 翻譯鍵必須來自 Gitea 官方語系檔案（`mail.*` 命名空間）
- **請勿在以下模板中使用 `.DisplayName`**：collaborator、transfer、release、workflow_run、assigned、default
- 針對 600px 寬度的郵件用戶端設計
- 盡可能在 Gmail、Outlook、Apple Mail 中測試

### 回報 Bug

1. 檢查引用的 Go 模板變數是否存在
2. 驗證翻譯鍵是否與 Gitea 語系檔案相符
3. 確認 `.DisplayName` 未在不支援的模板中使用
4. 重新產生預覽：`cd tools && go run . preview all`
5. 提交 issue，註明風格名稱、郵件類型及錯誤描述

### 改善文件

文件改進、預覽螢幕截圖、安裝指南和翻譯始終歡迎。

---

## 開發環境

- **Go 1.21+** 用於模板渲染與 CLI 工具
- **Node.js 18+**（可選）用於即時開發伺服器與 Juice CSS 內聯

### 本機預覽（靜態）

1. 先生成預覽資料：`cd tools && go run . preview all`
2. 在瀏覽器中開啟 `preview/index.html` — 無需伺服器

> 靜態 Gmail/Outlook 模擬僅供參考，使用 dev 模式可獲得準確的 CSS 內聯渲染。

### 開發伺服器（即時重載 + CSS 內聯）

```bash
cd tools && go run . dev
# → http://localhost:3456
```

修改 `.tmpl` 檔案後自動重建並推送至瀏覽器。

### 整合測試

將模板部署至 Gitea 實例，使用管理後台的測試郵件功能：
**Site Administration > Configuration > Mailer > Send Test Email**

---

## 提交規範

- `style(<名稱>):` — 特定風格的模板變更
- `preview:` — 預覽工具變更
- `tools:` — Go 建置腳本變更
- `docs:` — 文件與翻譯
- `fix:` — Bug 修復
- `project:` — README、LICENSE、中繼資料檔案

## 授權條款

參與貢獻即表示您同意將您的貢獻以 MIT 授權條款授權。
