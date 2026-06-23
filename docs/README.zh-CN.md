# Gitea 邮件模板

为自托管 [Gitea](https://about.gitea.com) 实例精心设计、面向不同受众的邮件模板集合。

> **110 个模板文件 — 10 种视觉风格，每种 11 种邮件类型**

---

## 设计理念

大多数自托管 Gitea 实例使用默认的纯文本邮件模板。本项目提供了**开箱即用、视觉精美的替代方案**——每种方案都针对特定社区或受众设计，您可以选择最适合您用户的风格。

每个模板都是直接替换。所有 Go 模板变量、翻译键和 Gitea 数据上下文完全兼容。**无需补丁、无需插件、无需 fork。**

---

## 风格画廊

| 预览 | 风格 | 受众 | 特点 |
|---|---|---|---|
| ![Horizon](images/horizon.png) | **Horizon** | 企业/公司 | 蓝色强调色、石板灰排版、居中卡片 |
| ![Terminal](images/terminal.png) | **Terminal** | 开发者/技术 | 暗色模式、等宽字体、绿色命令行风格 |
| ![Ember](images/ember.png) | **Ember** | 社区/开源 | 暖琥珀色、圆角、人文主义、包容 |
| ![Bloom](images/bloom.png) | **Bloom** | 创意/初创 | 磨砂玻璃、柔和蓝色光线、虹彩点缀 |
| ![Heritage](images/heritage.png) | **Heritage** | 教育/研究 | 海军蓝与金色、衬线字体、经典、权威 |
| ![Neon](images/neon.png) | **Neon** | 游戏/Web3/创意科技 | 赛博朋克霓虹、粉红与青色、合成波能量 |
| ![Mono](images/mono.png) | **Mono** | 设计工作室/编辑 | 瑞士粗野主义、黑白红强调、零圆角 |
| ![Terra](images/terra.png) | **Terra** | 可持续/健康 | 温暖大地色调、有机质感、人文衬线 |
| ![Ink](images/ink.png) | **Ink** | 出版/新闻/文学 | 编辑印刷、深蓝与金色、报纸排版 |
| ![Aurora](images/aurora.png) | **Aurora** | 高端SaaS/正念 | 空灵光效渐变、深紫与青绿、大气光晕 |

> 图片为 600px 宽截图，来自[在线预览](../preview/index.html)。截图方法参见 [images/README.md](images/README.md)。

[**在线预览画廊**](../preview/index.html)

---

## 安装

选择一种风格，将 `mail/` 目录复制到 Gitea 自定义模板路径：

```bash
cp -r themes/horizon/mail/* /var/lib/gitea/custom/templates/mail/
systemctl restart gitea
```

切换风格只需覆盖文件，无需更改配置。

### 确认生效

管理后台的测试邮件不会使用自定义模板。要验证模板是否生效，请触发一次真实的
邮件通知。最快的方式是密码重置：退出登录，点击登录页的**"忘记密码"**，查看
重置邮件即可——它将使用你的自定义样式渲染。

---

## 预览

**静态模式：**
```bash
cd tools && go run . preview all
```
然后打开 `preview/index.html`。

**开发服务器（实时重载）：**
```bash
cd tools && go run . dev
# → http://localhost:3456
```

| 功能 | 静态 | Dev |
|-----------|--------|-----|
| Go 模板渲染 | ✅ | ✅ |
| 主题/模板切换 | ✅ | ✅ |
| 实时重载 | — | ✅ |

---

## 兼容性

- **Gitea 1.25+** — v1.25 引入的邮件模板目录结构
- **最新测试：** Gitea 1.26.4<!-- TRACKER:LATEST-TESTED -->
- 与 Gitea 官方模板 100% 兼容 — 详见 [COMPATIBILITY.md](COMPATIBILITY.md)

## 许可证

MIT — 详见 [LICENSE](../LICENSE)。
