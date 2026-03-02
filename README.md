# whoami

Let every AI know who you are.

**whoami** is a cross-AI user identity profile sync service. Maintain one profile, and every AI agent you use can instantly understand your background, preferences, and work style — no more repeating yourself.

## Install

```bash
npx skills add MorvanZhou/whoami
```

Or paste this prompt to any AI agent:

> Install a cross-AI agent user identity sync skill by running: `npx skills add MorvanZhou/whoami`. After installation, open https://whoamiagent.com/login?redirect=dashboard to register or login to get an API key.

## Setup

1. **Get an API Key** — Visit [whoamiagent.com](https://whoamiagent.com) and sign in with GitHub or Google.
2. **Configure** — Run the setup command or let your AI agent handle it:

```bash
python3 <SKILL_DIR>/scripts/whoami_profile.py setup
```

This saves your API key to `~/.whoamiagent`.

## Usage

### Get your profile

```bash
python3 <SKILL_DIR>/scripts/whoami_profile.py get
```

### Update your profile

```bash
python3 <SKILL_DIR>/scripts/whoami_profile.py update "# My Profile\n\n## About\n- Name: Your Name\n..."
```

The `update` command overwrites the entire profile. The server automatically keeps the last 3 versions for rollback.

### View profile info

```bash
python3 <SKILL_DIR>/scripts/whoami_profile.py info
```

## How It Works

1. You install the whoami skill on your AI agent
2. You sign in at [whoamiagent.com](https://whoamiagent.com) and get an API key
3. When any AI agent starts a conversation, it runs `whoami get` to fetch your profile
4. Your profile is injected into the conversation context — the AI already knows you
5. When you share new info, the AI runs `whoami update` to save it for all your AI agents

## Config File

`~/.whoamiagent` (works on macOS, Linux, and Windows):

```
WHOAMI_API_KEY=wai_xxxxxxxxxxxxxxxx
WHOAMI_ENDPOINT=https://whoamiagent.com
```

## License

MIT
