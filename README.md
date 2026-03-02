# whoami

[English](./README.md) | [中文](./README_CN.md)

**Let every AI know who you are.**

Tired of introducing yourself every time you start a new conversation with an AI? **whoami** solves that. Create one identity profile, and every AI agent you use — ChatGPT, Claude, Cursor, Windsurf, or any other — instantly knows your name, skills, preferences, and work style.

No more repeating yourself. Just talk.

## Why whoami?

- **One profile, all AIs** — Write your profile once, share it across every AI tool you use
- **AI does the work** — Your AI agent installs, configures, and manages everything for you
- **Always up to date** — Tell any AI something new about yourself, and all your AIs learn it
- **Privacy first** — You own your data, you control what's shared
- **Zero friction** — No complex setup, no coding required

## Quick Start

### Option A: Tell your AI agent (easiest)

Just paste this to any AI agent (Cursor, Claude, ChatGPT, etc.):

```
Install a cross-AI agent user identity sync skill by running: `npx skills add MorvanZhou/whoami`. After installation, open https://whoamiagent.com/login?redirect=dashboard to register or login to get an API key.
```

Your AI will handle the rest — installation, login, and configuration.

### Option B: Install manually

```bash
npx skills add MorvanZhou/whoami
```

Then visit [whoamiagent.com](https://whoamiagent.com) to sign in with GitHub or Google and get your API key.

## How It Works

```
You ──→ Install whoami skill on your AI agent
         │
         ├─ Sign in at whoamiagent.com → Get API key
         │
         ├─ AI agent fetches your profile at conversation start
         │
         └─ You share new info → AI saves it → All AIs know it
```

1. **Install** — Add the whoami skill to your AI agent
2. **Sign in** — Create an account at [whoamiagent.com](https://whoamiagent.com) and grab an API key
3. **Talk** — Start any conversation. Your AI automatically loads your profile and knows who you are
4. **Update** — Mention something new about yourself. Your AI saves it, and every AI agent you use learns it

## What Gets Saved?

Your profile is a simple Markdown document. It can include anything you want your AIs to know:

- Your name and role
- Technical skills and experience
- Communication preferences
- Project context and workflows
- Anything that helps AI assist you better

You have full control — edit, update, or delete your profile anytime from the [dashboard](https://whoamiagent.com/dashboard).

## Supported Platforms

whoami works with any AI agent that supports skill installation, including:

- Cursor
- Windsurf
- Claude (via MCP or skills)
- ChatGPT (via custom instructions)
- And more...

## License

[MIT](./LICENSE)
