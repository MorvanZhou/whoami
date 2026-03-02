import logging
from datetime import datetime, timedelta, timezone

from jose import JWTError, jwt
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from app.config import settings
from app.models.user import User

logger = logging.getLogger("whoami")


def create_jwt_token(user_id: str) -> str:
    expire = datetime.now(timezone.utc) + timedelta(days=settings.jwt_expire_days)
    payload = {"sub": user_id, "exp": expire}
    return jwt.encode(payload, settings.jwt_secret, algorithm=settings.jwt_algorithm)


def decode_jwt_token(token: str) -> str | None:
    try:
        payload = jwt.decode(token, settings.jwt_secret, algorithms=[settings.jwt_algorithm])
        return payload.get("sub")
    except JWTError:
        return None


async def find_or_create_user(
    db: AsyncSession,
    provider: str,
    provider_id: str,
    email: str | None = None,
    name: str | None = None,
    avatar_url: str | None = None,
) -> User:
    result = await db.execute(
        select(User).where(User.provider == provider, User.provider_id == provider_id)
    )
    user = result.scalar_one_or_none()

    if user:
        # Update user info on each login
        if email:
            user.email = email
        if name:
            user.name = name
        if avatar_url:
            user.avatar_url = avatar_url
        await db.commit()
        await db.refresh(user)
        logger.info(f"[whoami] Existing user logged in: {user.id} via {provider}")
        return user

    user = User(
        provider=provider,
        provider_id=str(provider_id),
        email=email,
        name=name,
        avatar_url=avatar_url,
    )
    db.add(user)
    await db.commit()
    await db.refresh(user)
    logger.info(f"[whoami] New user created: {user.id} via {provider}")
    return user
