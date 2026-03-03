# whoami — 市场分析与竞品对比报告

> 报告生成时间：2026-03-03  
> 作者：莫烦（MorvanZhou）  
> 产品官网：[whoamiagent.com](https://whoamiagent.com)  
> GitHub：[github.com/morvanzhou/whoami](https://github.com/morvanzhou/whoami)

---

## 一、产品背景

### 核心痛点

随着 AI Agent 成为主流交互范式，一个人会同时使用多个 AI 工具（Claude、Cursor、ChatGPT、Kimi 等），或在不同设备上与 AI 对话。**每个 AI 都需要重新认识你一遍，每个 AI 对你的记忆都是孤立的。**

具体表现为三个场景：

1. **Switch AI, Lose Identity** — 从 ChatGPT 切换到 Claude，它完全不知道你是谁，重新开始
2. **New Chat, Blank Slate** — 即使同一个 AI，每次新对话都从零开始，上下文消失
3. **Can't Edit What AI Knows** — 想更新 AI 对你的认知，但设置入口深埋或锁死，无法管理

### 解决方案

**One profile. Every AI agent. No more repeating yourself.**

做一个 Agent Skill，提供云端服务：
- 每次不同 AI 被使用时，自动从远端拉取用户的身份档案（Profile）注入上下文
- AI 对用户产生的新认知，可以同步写回远端
- 本质是"**AI 世界的个人名片 / 身份档案**"

---

## 二、产品现状（2026-03-03）

### 产品已上线，核心功能完整

| 维度 | 现状 |
|------|------|
| **官网** | whoamiagent.com — 已上线，含完整 Landing Page |
| **GitHub** | github.com/morvanzhou/whoami — MIT 开源 |
| **安装方式** | `npx skills add MorvanZhou/whoami` |
| **支持平台** | Cursor、Windsurf、Claude（MCP/Skills）、ChatGPT |
| **注册方式** | GitHub / Google 一键登录 |
| **核心功能** | 注册 → 获取 API Key → AI 自动加载 Profile → 对话中更新 Profile |

### 核心工作流

```
用户粘贴一段 Prompt 给 AI
    ↓
AI 自动执行：npx skills add MorvanZhou/whoami
    ↓
用户访问 whoamiagent.com 注册/登录，获取 API Key
    ↓
每次对话开始，AI 自动拉取 Profile 注入上下文
    ↓
对话中提到新信息 → AI 自动同步写回远端
```

### Profile 存储格式

用户 Profile 是一份 Markdown 文档，可包含：
- 姓名与职业
- 技术技能与经验
- 沟通偏好
- 项目背景与工作流
- 任何希望 AI 了解的个人信息

用户可随时在 Dashboard 编辑，或直接告诉 AI 更新，**所有 AI 实时同步**。

---

## 三、市场调研与竞品分析

### 3.1 调研范围

本次调研覆盖以下渠道（2025.03 - 2026.03）：
- skills.sh（AI Agent Skill 生态，81,834 个 skill）
- GitHub（关键词：memory、profile、user context、identity、persona、MCP）
- Indie Hackers（AI memory 方向）
- Product Hunt
- 少数派 / 中文社区

### 3.2 竞品全景图

```
                ← 面向开发者/B端                     面向普通用户/C端 →

存"做了什么"    OpenContext ⭐407     Mem0 ⭐25k+
(对话/项目记忆) memorix ⭐149         recall ⭐160
                engram ⭐626          Supermemory ⭐2k+

存"我是谁"      Memobase (B2D SDK)    ← 【whoami 的目标位置】✅
(用户身份档案)  Kinic (Web3方向)         这里目前仍是空白
                ACME Brains (未发布)
```

### 3.3 竞品详细分析

#### A 类：对话记忆型（存"做了什么"）

| 产品 | Stars | 定位 | 目标用户 | 局限性 |
|------|-------|------|----------|--------|
| **Mem0** | 25k+ | AI 应用记忆层 SDK + MCP | 开发者 | B2D 工具，需集成开发；不是 C 端产品 |
| **engram** | 626 | 本地 AI 记忆系统 | 开发者 | 本地部署，无云端同步 |
| **Supermemory** | 2k+ | 个人知识库 + AI 搜索 | 知识工作者 | 存储的是内容/文章，而非"我是谁" |
| **recall** | 160 | 浏览历史 AI 记忆 | 普通用户 | 仅限浏览器行为，无身份档案 |
| **memorix** | 149 | Redis 兼容记忆系统 | 开发者 | 纯技术工具，无 C 端产品形态 |
| **OpenContext** | 407 | 开发项目上下文持久化 | Cursor/Claude Code 用户 | 存储的是项目背景/代码决策，非用户身份 |

#### B 类：用户身份档案型（存"我是谁"）

| 产品 | 定位 | 目标用户 | 与 whoami 的差异 |
|------|------|----------|-----------------|
| **Memobase** | 面向开发者的用户记忆基础设施 SDK | B2D（开发者集成） | 开发者工具，不是 C 端产品；需要 App 集成 |
| **Kinic** | "Plaid for AI Memory"，AI 记忆数据主权 | 偏 Web3/数据所有权用户 | Web3 方向，区块链背景；月收入 $10k（已验证付费意愿） |
| **ACME Brains** | 跨 AI 身份档案 | 开发者 | 尚未正式发布，GitHub 活跃度低 |
| **Claude 内置记忆** | Claude 平台专属记忆 | Claude 用户 | 不跨平台，锁死在 Claude 生态内 |
| **Papr.ai** | 企业级 AI 记忆基础设施 | 企业 | 偏企业，定价高，个人用户无法直接使用 |

#### C 类：平台内置记忆（生态孤岛）

| 平台 | 记忆功能 | 局限性 |
|------|----------|--------|
| ChatGPT | Memory（自动保存对话信息） | 仅限 ChatGPT，不跨平台 |
| Claude | Projects（项目级上下文） | 仅限 Claude，不跨平台 |
| Cursor | Rules for AI（项目规则） | 仅限代码场景，不跨平台 |
| Gemini | 暂无独立记忆功能 | — |

### 3.4 skills.sh 专项调研

**结论：skills.sh 上没有任何同类产品。**

扫描了所有 memory、profile、user、persona、context、sync、identity、persist、cross 等关键词，**没有发现任何"跨 AI/跨设备用户身份同步"的 skill**。

> ✅ whoami 是 skills.sh 生态中该方向的**首个产品**。

### 3.5 中文市场调研

少数派/中文社区：**完全没有同类中文产品**。

最相关的内容是《我是如何用个人笔记，打造出专属 AI 助理的》（82赞），说明中文用户有这个意识，但**没有现成工具**。

> ✅ 中文市场完全空白，先发优势明确。

---

## 四、竞品对比深度分析

### 4.1 whoami vs 最接近竞品

#### vs OpenContext（最新发现，⭐407）

| 维度 | OpenContext | whoami |
|------|-------------|--------|
| **核心问题** | "AI 不记得我们做了什么项目" | "AI 不知道我是谁" |
| **存储内容** | 项目背景、技术决策、代码上下文 | 用户身份、偏好、工作方式、个人背景 |
| **目标用户** | 开发者（Cursor/Claude Code 用户） | 所有 AI 工具用户（C 端） |
| **使用方式** | CLI + MCP + 本地 Desktop App | 注册即用 + 云端托管 |
| **数据归属** | 本地文件（`contexts/` 目录） | 用户拥有，云端同步 |
| **是否需要技术背景** | ✅ 需要（npm install，CLI 操作） | ❌ 不需要（普通用户可用） |
| **关系** | **互补，非竞争** | — |

#### vs Kinic（商业化最好的竞品）

| 维度 | Kinic | whoami |
|------|-------|--------|
| **核心理念** | "数据主权" / "AI 记忆的 Plaid" | "让每个 AI 认识你" |
| **方向** | Web3 / 区块链 / 数据所有权 | 云端 SaaS / 简单易用 |
| **月收入** | $10,000（已验证付费意愿） | 早期 |
| **目标用户** | Web3 用户 / 数据主权倡导者 | 普通 AI 工具用户 |
| **安装复杂度** | 较高（Web3 背景） | 极低（一行命令） |
| **差异化** | whoami 更轻量、更易用、更 C 端 | — |

#### vs Mem0（最大同类项目）

| 维度 | Mem0 | whoami |
|------|------|--------|
| **Stars** | 25k+ | 早期 |
| **定位** | AI 应用的记忆基础设施 | 用户的 AI 身份档案 |
| **使用者** | 开发者（集成到自己的 App） | 终端用户（直接使用） |
| **存储内容** | 对话记忆、用户行为数据 | 用户主动维护的身份档案 |
| **商业模式** | B2D SaaS | C 端 SaaS |
| **关系** | **不同赛道，不直接竞争** | — |

### 4.2 whoami 的核心差异化

| 维度 | 市场现有产品 | whoami |
|------|-------------|--------|
| **存储对象** | 对话记录 / 项目上下文 | **用户身份档案**（我是谁） |
| **使用者** | 开发者 / 企业 | **普通 C 端用户** |
| **跨平台** | 各自为政（平台孤岛） | **真正跨平台同步** |
| **安装方式** | 复杂配置 | **一行命令 + 粘贴 Prompt** |
| **中文市场** | 无 | **首个中文市场产品** |
| **数据格式** | 向量数据库 / 结构化数据 | **Markdown 文档（人可读）** |

---

## 五、市场机会评估

### 5.1 市场验证信号

| 信号 | 来源 | 结论 |
|------|------|------|
| Kinic 月收入 $10k | Indie Hackers | ✅ 该赛道有付费意愿 |
| Mem0 获 $2.8M 种子轮 | 公开信息 | ✅ 资本认可 AI 记忆赛道 |
| skills.sh 无同类产品 | 直接扫描 | ✅ 生态空白 |
| 中文社区无同类产品 | 少数派调研 | ✅ 中文市场空白 |
| OpenContext ⭐407（新项目） | GitHub | ✅ 开发者有跨 AI 上下文需求 |
| 平台内置记忆普遍存在 | ChatGPT/Claude | ✅ 用户需求被验证，但孤岛问题未解决 |

### 5.2 目标用户画像

**主要用户**：
- 同时使用 2 个以上 AI 工具的知识工作者
- 频繁开新对话的重度 AI 用户
- 不想重复自我介绍的效率控

**次要用户**：
- 开发者（可通过 MCP 集成到自己的工具链）
- 想要管理"AI 眼中的自己"的普通用户

### 5.3 竞争壁垒

1. **先发优势**：skills.sh 生态中首个该方向产品
2. **网络效应**：用户越多 → 更多 AI 工具适配 → 更多用户
3. **数据积累**：用户 Profile 越完善 → 越难迁移
4. **分发渠道**：公众号"rethink比移" + 莫烦 Python 社区（21.6k GitHub followers）

---

## 六、技术实现方案

### 6.1 整体架构

```
用户 → whoamiagent.com（前端）
         ↓ 注册/登录（GitHub/Google OAuth）
         ↓ 获取 API Key
    
    FastAPI 后端（腾讯云 CVM）
    ├── 鉴权模块（SQLite + SQLAlchemy）
    ├── Profile 读写（腾讯云 COS）
    └── RESTful API

AI Agent（Cursor / Claude / ChatGPT 等）
    ├── 安装：npx skills add MorvanZhou/whoami
    ├── GET /profile  → 拉取 Profile 注入 system prompt
    └── PATCH /profile → 将新信息写回远端
```

### 6.2 技术栈

| 层级 | 技术 | 说明 |
|------|------|------|
| 后端框架 | FastAPI | 异步高性能，自带 OpenAPI 文档 |
| 数据库 | SQLite + SQLAlchemy | 轻量，后期可迁移 PostgreSQL |
| Profile 存储 | 腾讯云 COS | 每用户一个 `{user_id}.md` 文件 |
| 鉴权方式 | API Key（Bearer Token）| 适合 MCP/Skill 调用 |
| 前端 | 单页应用 | 注册 + 获取 Key + Profile 编辑 |
| 部署 | 腾讯云 CVM + Nginx | 复用已有服务器 |
| 安装方式 | `npx skills add MorvanZhou/whoami` | 零门槛安装 |

### 6.3 API 接口

```
POST /auth/register          # 注册
POST /auth/login             # 登录
POST /auth/apikey            # 创建 API Key
DELETE /auth/apikey/{id}     # 撤销 Key

GET  /profile                # 获取 profile（Bearer Key 鉴权）
POST /profile                # 覆盖写入 profile
PATCH /profile               # 追加/合并更新 profile
GET  /profile/meta           # 获取元信息（字数、更新时间等）
```

### 6.4 安装与使用方式

**Option A（最简单）— 粘贴 Prompt 给 AI：**

```
Install a cross-AI agent user identity sync skill by running: 
npx skills add MorvanZhou/whoami. 
After installation, open https://whoamiagent.com/login?redirect=dashboard 
to register or login to get an API key.
```

**Option B — 直接命令安装：**

```bash
npx skills add MorvanZhou/whoami
```

---

## 七、产品定位总结

| 维度 | 描述 |
|------|------|
| **一句话定位** | Let every AI know who you are |
| **核心价值** | 一份 Profile，所有 AI 同步，不再重复自我介绍 |
| **目标用户** | 同时使用多个 AI 工具的知识工作者 / 开发者 |
| **差异化** | 存"身份档案"而非"对话记录"；C 端开箱即用；跨平台真同步 |
| **分发渠道** | skills.sh（Skill 生态）+ MCP 生态 + 公众号"rethink比移" |
| **商业模式** | 免费基础版 + Pro 版（更大存储、版本历史、多 Key） |
| **中文市场** | 完全空白，先发优势明确 |

---

## 八、下一步建议

### 近期（1-2 周）
- [ ] 在 skills.sh 正式发布 `whoami` Skill 条目
- [ ] 在公众号"rethink比移"发布产品介绍文章
- [ ] 收集早期用户反馈，完善 Profile 编辑体验

### 中期（1 个月）
- [ ] 支持 MCP 标准协议（扩大 AI 工具覆盖范围）
- [ ] Profile 版本历史功能
- [ ] 多 Key 管理（按设备标签区分）

### 长期
- [ ] 支持本地部署（满足隐私需求，参考 rethink.run 的 GPL 开源模式）
- [ ] Profile 公开/私有分区
- [ ] 探索 Pro 付费功能

---

*报告基于 2025.03 - 2026.03 期间的市场调研，结合产品实际上线状态综合整理。*

---

## 九、SEO 优化策略

> 本章基于竞品（Mem0、Memobase、Supermemory、Letta/MemGPT、OpenContext）的官网 Title、Meta Description、H1/H2 用词进行分析，提炼高频关键词并给出 whoamiagent.com 和 GitHub README 的具体优化建议。

### 9.1 竞品关键词分析

通过抓取主要竞品官网 meta 信息，整理出高频 SEO 关键词：

| 竞品 | Title 关键词 | Description 核心用词 |
|------|-------------|---------------------|
| Mem0 | memory layer, AI Apps | continuously learn, past user interactions, personalization |
| Memobase | AI Memory for LLMs, Personalized Agents | long term memory, structured user profiles, personalized interactions across sessions |
| Supermemory | Universal Memory API | long-term memory, store, recall, personalize, user profiles |
| Letta/MemGPT | stateful agents, advanced memory | learn and self-improve over time |
| OpenContext | persistent memory, AI assistant | project context, cross-session |

**高频词汇提炼**：
- `persistent identity` / `persistent profile`
- `cross-session memory` / `cross-agent`
- `AI memory` / `AI identity`
- `user profile` / `personalized AI`
- `MCP` / `AI agent skill`
- `context injection` / `system prompt`
- `no more repeating yourself`
- `works with Cursor / Claude / ChatGPT`

---

### 9.2 whoamiagent.com 官网 SEO 优化建议

#### 📌 Page Title（`<title>` 标签）

**现状（推测）**：`whoami - AI Agent Identity`

**建议改为**：
```
whoami — Persistent AI Identity Profile | Works with Cursor, Claude & ChatGPT
```

或更简洁版：
```
whoami | Your Persistent Identity for Every AI Agent
```

**原因**：
- 加入 `Cursor`、`Claude`、`ChatGPT` 等品牌词，可截获搜索这些工具 + memory/profile 的长尾流量
- `Persistent` 是竞品高频词，SEO 权重高
- 明确说明场景（AI Agent），符合搜索意图

---

#### 📌 Meta Description

**建议**：
```html
<meta name="description" content="whoami gives every AI agent a persistent identity profile about you — synced across Cursor, Claude, ChatGPT and more. Stop repeating yourself. One profile, every AI." />
```

**关键词覆盖**：
- `persistent identity profile`（高搜索量）
- `synced across Cursor, Claude, ChatGPT`（长尾精准词）
- `stop repeating yourself`（情感共鸣 + 可能的搜索词）
- `one profile, every AI`（品牌 slogan 兼 SEO）

---

#### 📌 H1 主标题

**竞品参考**：
- Mem0：`AI Agents Forget. Mem0 Remembers.`
- Supermemory：`The context engineering infrastructure for your AI agent`
- Letta：`Stateful agents: AI with advanced memory`

**建议 whoami H1**：
```
Every AI Knows You. Finally.
```
或：
```
One Profile. Every AI Agent. No More Repeating Yourself.
```

---

#### 📌 H2 副标题 / Feature 区块建议用词

| 功能点 | 建议用词 |
|--------|---------|
| 跨 AI 同步 | `Cross-agent persistent identity` / `Sync your profile across all AI tools` |
| 自动注入上下文 | `Auto-inject your profile into every conversation` / `Context injection for AI agents` |
| 用户可编辑 | `You control what AI knows about you` / `Edit your AI identity anytime` |
| 支持的工具 | `Works with Cursor, Claude, ChatGPT, Windsurf and any MCP-compatible agent` |
| 隐私安全 | `Your profile, your data` / `API key protected, user-owned` |

---

#### 📌 页面结构化数据（Schema.org）

建议在 `<head>` 中添加 JSON-LD：

```json
{
  "@context": "https://schema.org",
  "@type": "SoftwareApplication",
  "name": "whoami",
  "description": "Persistent AI identity profile that syncs across Cursor, Claude, ChatGPT and all MCP-compatible AI agents.",
  "applicationCategory": "ProductivityApplication",
  "operatingSystem": "Web",
  "url": "https://whoamiagent.com",
  "offers": {
    "@type": "Offer",
    "price": "0",
    "priceCurrency": "USD"
  }
}
```

---

#### 📌 Open Graph / Twitter Card

```html
<meta property="og:title" content="whoami — Your Persistent Identity for Every AI Agent" />
<meta property="og:description" content="One profile that every AI agent can read. Works with Cursor, Claude, ChatGPT. Stop repeating yourself." />
<meta property="og:image" content="https://whoamiagent.com/og-image.png" />
<meta name="twitter:card" content="summary_large_image" />
<meta name="twitter:title" content="whoami — Persistent AI Identity Profile" />
<meta name="twitter:description" content="Let every AI know who you are. One profile, synced across all your AI agents." />
```

---

### 9.3 GitHub README SEO 优化建议

GitHub README 本身会被 Google 索引，且 GitHub 域名权重极高，优化 README 等同于优化一个高权重页面。

#### 📌 Repository Description（仓库简介，一行）

**现状（推测）**：`whoami - AI identity profile for agents`

**建议改为**：
```
Persistent user identity profile for AI agents — works with Cursor, Claude, ChatGPT and any MCP-compatible tool
```

**原因**：这一行会直接出现在 Google 搜索结果中，是 GitHub 项目最重要的 SEO 文本。

---

#### 📌 Repository Topics（标签，最多20个）

建议在 GitHub 仓库 Settings → Topics 中添加：

```
ai-agent  mcp  cursor  claude  chatgpt  memory  user-profile  
persistent-memory  identity  context-injection  llm  
ai-memory  cross-agent  skill  fastapi  python
```

**原因**：GitHub Topics 会被 Google 索引，且 GitHub 内部搜索也依赖 Topics。

---

#### 📌 README 第一段（最关键的 SEO 区域）

**建议 README 开头改为**：

```markdown
# whoami — Persistent AI Identity Profile

> **One profile. Every AI agent. No more repeating yourself.**

**whoami** gives every AI agent (Cursor, Claude, ChatGPT, Windsurf) a persistent identity profile about you — automatically injected into every conversation, synced across all your AI tools.

### Why whoami?

- 🔁 **Cross-agent sync** — Switch from Cursor to Claude? Your profile follows you.
- 🧠 **Context injection** — AI reads your profile at the start of every session automatically.
- ✏️ **User-controlled** — You decide what AI knows about you. Edit anytime at [whoamiagent.com](https://whoamiagent.com).
- 🔒 **API key protected** — Your data, your control.
- ⚡ **1-minute setup** — One command, done.
```

**关键词覆盖**：
- `persistent identity profile`
- `cross-agent sync`
- `context injection`
- `Cursor, Claude, ChatGPT, Windsurf`（品牌词截流）
- `MCP-compatible`

---

#### 📌 README 安装章节建议用词

```markdown
## Installation

Works with any MCP-compatible AI agent. Supports Cursor, Claude, ChatGPT, and Windsurf.

\`\`\`bash
npx skills add MorvanZhou/whoami
\`\`\`

Then visit [whoamiagent.com](https://whoamiagent.com) to register and get your API key.
```

**原因**：明确列出支持的工具名，可截获 `cursor memory plugin`、`claude persistent context`、`chatgpt memory sync` 等长尾搜索词。

---

### 9.4 目标关键词优先级矩阵

| 关键词 | 搜索意图 | 竞争度 | 建议优先级 |
|--------|---------|--------|-----------|
| `ai agent memory` | 功能搜索 | 高 | ⭐⭐⭐ 必须覆盖 |
| `persistent user profile AI` | 精准功能 | 低 | ⭐⭐⭐ 必须覆盖 |
| `cursor memory plugin` | 工具搜索 | 低 | ⭐⭐⭐ 必须覆盖 |
| `claude persistent context` | 工具搜索 | 低 | ⭐⭐⭐ 必须覆盖 |
| `cross-agent AI identity` | 品类词 | 极低（新词） | ⭐⭐ 定义品类 |
| `MCP user profile` | 生态词 | 低 | ⭐⭐ 重要 |
| `AI memory MCP` | 生态词 | 中 | ⭐⭐ 重要 |
| `stop repeating yourself AI` | 痛点词 | 极低 | ⭐ 长尾 |
| `chatgpt memory sync` | 工具搜索 | 中 | ⭐⭐ 重要 |
| `whoami agent` | 品牌词 | 低 | ⭐⭐⭐ 品牌保护 |

---

### 9.5 内容营销 SEO 建议

除了站内 SEO，以下内容可以带来外链和长尾流量：

1. **发布一篇博客文章**（放在 whoamiagent.com/blog）：
   - 标题建议：`"Why Every AI Forgets You — And How to Fix It"`
   - 覆盖关键词：`AI memory problem`、`cross-agent context`、`persistent AI identity`

2. **在 GitHub 发布 Discussion / Wiki**：
   - 标题：`How whoami compares to Mem0, OpenContext, and Claude Memory`
   - 这类对比内容会被搜索 `mem0 alternative`、`claude memory alternative` 的用户发现

3. **在 skills.sh 的 Skill 描述中**：
   - 用词建议：`"Persistent cross-agent user identity profile. Works with Cursor, Claude, ChatGPT. Your AI agents will always know who you are."`

4. **公众号文章 SEO 联动**：
   - 文章中明确提及 `whoamiagent.com`，微信文章虽无直接 SEO 价值，但可以带来高质量外链和品牌搜索量

---

*SEO 建议基于 2026-03-03 竞品分析，建议每季度复查关键词排名并动态调整。*
