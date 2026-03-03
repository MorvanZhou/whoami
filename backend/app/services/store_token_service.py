import logging
import secrets
from datetime import datetime, timedelta, timezone

from sqlalchemy import delete, select
from sqlalchemy.ext.asyncio import AsyncSession

from app.models.store_token import StoreApiToken

logger = logging.getLogger("whoami")

TOKEN_EXPIRE_MINUTES = 15


async def cleanup_expired_tokens(db: AsyncSession) -> int:
    """Delete tokens that are used or expired. Returns count of deleted rows."""
    cutoff = datetime.now(timezone.utc) - timedelta(minutes=TOKEN_EXPIRE_MINUTES)
    result = await db.execute(
        delete(StoreApiToken).where(
            (StoreApiToken.is_used.is_(True)) | (StoreApiToken.created_at < cutoff)
        )
    )
    await db.commit()
    count: int = result.rowcount  # type: ignore[attr-defined]
    return count


def generate_store_token() -> str:
    return secrets.token_urlsafe(32)


async def create_store_token(
    db: AsyncSession, user_id: str, api_key_plain: str, api_key_id: str
) -> str:
    token = generate_store_token()
    store_token = StoreApiToken(
        user_id=user_id,
        token=token,
        api_key_plain=api_key_plain,
        api_key_id=api_key_id,
    )
    db.add(store_token)
    await db.commit()
    logger.info(f"[whoami] Store token created for user {user_id}")
    return token


async def consume_store_token(db: AsyncSession, token: str) -> str | None:
    """Consume a one-time store token. Returns the plain API key or None."""
    result = await db.execute(
        select(StoreApiToken).where(
            StoreApiToken.token == token,
            StoreApiToken.is_used.is_(False),
        )
    )
    store_token = result.scalar_one_or_none()
    if not store_token:
        return None

    # Check expiration (created_at may be naive from SQLite, treat as UTC)
    created_at = store_token.created_at
    if created_at.tzinfo is None:
        created_at = created_at.replace(tzinfo=timezone.utc)
    expire_time = created_at + timedelta(minutes=TOKEN_EXPIRE_MINUTES)
    now = datetime.now(timezone.utc)
    if now > expire_time:
        store_token.is_used = True
        await db.commit()
        logger.info(f"[whoami] Store token expired: {token[:8]}...")
        return None

    # Mark as used
    store_token.is_used = True
    await db.commit()
    logger.info(f"[whoami] Store token consumed for user {store_token.user_id}")
    return store_token.api_key_plain


async def is_store_token_consumed(db: AsyncSession, api_key_id: str, user_id: str) -> bool:
    """Check if the store token associated with an api_key_id has been consumed."""
    result = await db.execute(
        select(StoreApiToken).where(
            StoreApiToken.api_key_id == api_key_id,
            StoreApiToken.user_id == user_id,
        )
    )
    store_token = result.scalar_one_or_none()
    if not store_token:
        return False
    return store_token.is_used
