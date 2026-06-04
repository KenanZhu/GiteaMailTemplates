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

---

## 预览

**静态模式：**
```bash
cd tools && go run . preview all
```
然后打开 `preview/index.html`。

> [!WARNING]
> 静态模式下 Gmail/Outlook 模拟仅供参考，使用 dev 模式可获得相对准确的 CSS 内联渲染。

**开发服务器（实时重载 + Juice CSS 内联 + Gmail/Outlook 模拟）：**
```bash
cd tools && go run . dev     # 需要 Node.js
# → http://localhost:3456
```

> [!NOTE]
> Dev 模拟无法 100% 还原各个邮件客户端的渲染差异，仅供参考，请以实际效果为准。

| 功能 | 静态 | Dev |
|-----------|--------|-----|
| Go 模板渲染 | ✅ | ✅ |
| 主题/模板切换 | ✅ | ✅ |
| Juice CSS 内联 | — | ✅ |
| Gmail/Outlook CSS 剥离 | — | ✅ |
| 实时重载 | — | ✅ |
| 需要 Node.js | — | ✅ |

---

## 兼容性

- **Gitea 1.21+**，100% 变量兼容，仅使用内置函数和官方翻译键

## 许可证

MIT — 详见 [LICENSE](../LICENSE)。
