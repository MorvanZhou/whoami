import os

from pydantic_settings import BaseSettings, SettingsConfigDict

_env_file = ".env.dev" if os.getenv("WHOAMI_ENV") == "dev" else ".env"


class Settings(BaseSettings):
    # Environment: "dev" or "prod"
    env: str = "prod"

    # Database
    database_url: str = "sqlite+aiosqlite:///./whoami.db"

    # JWT
    jwt_secret: str = "change-me-to-a-random-secret"
    jwt_algorithm: str = "HS256"
    jwt_expire_days: int = 7

    # GitHub OAuth
    github_client_id: str = ""
    github_client_secret: str = ""
    github_redirect_uri: str = "https://whoamiagent.com/api/auth/github/callback"

    # Google OAuth
    google_client_id: str = ""
    google_client_secret: str = ""
    google_redirect_uri: str = "https://whoamiagent.com/api/auth/google/callback"

    # Tencent Cloud COS
    cos_secret_id: str = ""
    cos_secret_key: str = ""
    cos_region: str = "ap-hongkong"
    cos_bucket: str = "whoamiagent-xxx"

    # General
    frontend_url: str = "https://whoamiagent.com"
    api_key_prefix: str = "wai_"

    model_config = SettingsConfigDict(env_file=_env_file, env_file_encoding="utf-8")


settings = Settings()
