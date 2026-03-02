# whoami

[English](./README.md) | [中文](./README_CN.md)

**让每个 AI 都认识你。**

是不是每次和新的 AI 对话，都要重新介绍一遍自己？**whoami** 帮你解决这个问题。只需创建一份身份档案，你使用的所有 AI 工具 —— ChatGPT、Claude、Cursor、Windsurf 等 —— 都能立刻了解你的名字、技能、偏好和工作风格。

不用再重复自我介绍，直接开聊。

## 为什么选 whoami？

- **一份档案，所有 AI 通用** —— 写一次个人资料，所有 AI 工具共享
- **AI 帮你搞定一切** —— 安装、配置、管理全部由你的 AI 助手完成
- **实时同步** —— 告诉任何一个 AI 你的新信息，所有 AI 都会学到
- **隐私优先** —— 数据归你所有，你决定分享什么
- **零门槛** —— 不需要复杂配置，不需要写代码

## 快速开始

### 方式 A：让 AI 帮你安装（最简单）

把下面这段话发给你的 AI 助手（Cursor、Claude、ChatGPT 等）：

```
Install a cross-AI agent user identity sync skill by running: `npx skills add MorvanZhou/whoami`. After installation, open https://whoamiagent.com/login?redirect=dashboard to register or login to get an API key.
```

AI 会自动帮你完成安装、登录和配置，全程不需要你操心。

### 方式 B：手动安装

```bash
npx skills add MorvanZhou/whoami
```

然后访问 [whoamiagent.com](https://whoamiagent.com)，使用 GitHub 或 Google 登录，获取 API Key。

## 工作原理

```
你 ──→ 在 AI 助手上安装 whoami skill
        │
        ├─ 在 whoamiagent.com 登录 → 获取 API Key
        │
        ├─ AI 在对话开始时自动加载你的档案
        │
        └─ 你分享新信息 → AI 保存 → 所有 AI 都能看到
```

1. **安装** —— 把 whoami skill 添加到你的 AI 助手
2. **登录** —— 在 [whoamiagent.com](https://whoamiagent.com) 注册账号，拿到 API Key
3. **对话** —— 开始任意对话，AI 会自动加载你的档案，立刻了解你是谁
4. **更新** —— 聊天中提到新的个人信息，AI 会自动保存，你的所有 AI 助手都能学到

## 档案里写什么？

你的档案是一个简单的 Markdown 文档，可以包含任何你希望 AI 了解的信息：

- 你的名字和职业
- 技术栈和工作经验
- 沟通偏好
- 项目背景和工作流程
- 任何能帮助 AI 更好地协助你的信息

你拥有完全的控制权 —— 随时在 [控制台](https://whoamiagent.com/dashboard) 查看、编辑或删除你的档案。

## 支持的平台

whoami 可以和任何支持 skill 安装的 AI 助手配合使用，包括：

- Cursor
- Windsurf
- Claude（通过 MCP 或 skills）
- ChatGPT（通过自定义指令）
- 更多...

## 开源协议

[MIT](./LICENSE)
