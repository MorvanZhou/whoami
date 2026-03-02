import hashlib
import logging
import secrets
from datetime import datetime, timezone

from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from app.config import settings
from app.models.api_key import ApiKey

logger = logging.getLogger("whoami")


def generate_api_key() -> str:
    random_part = secrets.token_hex(24)
    return f"{settings.api_key_prefix}{random_part}"


def hash_api_key(key: str) -> str:
    return hashlib.sha256(key.encode()).hexdigest()


def get_key_prefix(key: str) -> str:
    return key[:8]


def get_key_suffix(key: str) -> str:
    return key[-4:]


async def create_api_key(
    db: AsyncSession, user_id: str, label: str | None = None
) -> tuple[str, ApiKey]:
    plain_key = generate_api_key()
    key_hash = hash_api_key(plain_key)
    key_prefix = get_key_prefix(plain_key)
    key_suffix = get_key_suffix(plain_key)

    api_key = ApiKey(
        user_id=user_id,
        key_hash=key_hash,
        key_prefix=key_prefix,
        key_suffix=key_suffix,
        label=label,
    )
    db.add(api_key)
    await db.commit()
    await db.refresh(api_key)

    logger.info(f"[whoami] API key created for user {user_id}, prefix: {key_prefix}")
    return plain_key, api_key


async def list_api_keys(db: AsyncSession, user_id: str) -> list[ApiKey]:
    result = await db.execute(
        select(ApiKey).where(ApiKey.user_id == user_id, ApiKey.is_active.is_(True))
    )
    return list(result.scalars().all())


async def revoke_api_key(db: AsyncSession, key_id: str, user_id: str) -> bool:
    result = await db.execute(
        select(ApiKey).where(ApiKey.id == key_id, ApiKey.user_id == user_id)
    )
    api_key = result.scalar_one_or_none()
    if not api_key:
        return False

    await db.delete(api_key)
    await db.commit()
    logger.info(f"[whoami] API key deleted: {key_id}")
    return True


async def verify_api_key(db: AsyncSession, plain_key: str) -> ApiKey | None:
    key_hash = hash_api_key(plain_key)
    result = await db.execute(
        select(ApiKey).where(ApiKey.key_hash == key_hash, ApiKey.is_active.is_(True))
    )
    api_key = result.scalar_one_or_none()
    if api_key:
        api_key.last_used = datetime.now(timezone.utc)
        await db.commit()
    return api_key
