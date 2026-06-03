# 贡献指南 — Gitea 邮件模板

## 参与方式

### 添加新风格

1. 创建风格目录：`themes/<风格名称>/`
2. 从现有风格复制目录结构
3. 用独特的视觉设计实现全部 11 个 `.tmpl` 文件
4. 重新生成预览：`go run ./tools/build-preview.go`（构建脚本会自动发现 `themes/` 下的所有主题目录）
5. 在 `preview/index.html` 的 `<select id="sel-theme">` 中添加选项
6. 提交包含截图渲染效果的 PR

### 风格指南

- 每个风格必须包含全部 **11 种模板类型**
- 只使用 Gitea 内置模板函数
- 翻译键必须来自 Gitea 官方语言文件（`mail.*` 命名空间）
- **不要在以下模板中使用 `.DisplayName`**：collaborator、transfer、release、workflow_run、assigned、default
- 面向 600px 宽度的邮件客户端设计
- 尽可能在 Gmail、Outlook、Apple Mail 中测试

### 报告 Bug

1. 检查引用的 Go 模板变量是否存在
2. 验证翻译键是否与 Gitea 语言文件匹配
3. 确认 `.DisplayName` 未在不支持的模板中使用
4. 重新生成预览：`go run ./tools/build-preview.go`
5. 提交 issue，注明风格名称、邮件类型及错误描述

### 改进文档

文档改进、预览截图、安装指南和翻译始终欢迎。

---

## 开发环境

无需构建工具或依赖。

### 本地预览

1. 先生成预览数据：`go run ./tools/build-preview.go`
2. 在浏览器中打开 `preview/index.html`
3. 在 Modern、Gmail、Outlook、Raw source 模式间切换验证效果

### 集成测试

将模板部署到 Gitea 实例，使用管理后台的测试邮件功能：
**Site Administration > Configuration > Mailer > Send Test Email**

---

## 提交规范

- `style(<名称>):` — 特定风格的模板变更
- `preview:` — 预览工具变更
- `tools:` — Go 构建脚本变更
- `docs:` — 文档和翻译
- `fix:` — Bug 修复
- `project:` — README、LICENSE、元文件

## 许可协议

参与贡献即表示您同意将您的贡献以 MIT 许可证授权。
