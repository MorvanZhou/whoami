import logging

from sqlalchemy import inspect, text
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine

from app.config import settings
from app.models import Base

logger = logging.getLogger("whoami")

engine = create_async_engine(settings.database_url, echo=False)
async_session = async_sessionmaker(engine, class_=AsyncSession, expire_on_commit=False)


def _migrate_api_keys(conn) -> None:
    inspector = inspect(conn)
    if "api_keys" in inspector.get_table_names():
        columns = [col["name"] for col in inspector.get_columns("api_keys")]
        if "key_suffix" not in columns:
            conn.execute(text("ALTER TABLE api_keys ADD COLUMN key_suffix VARCHAR DEFAULT '' NOT NULL"))
            logger.info("[whoami] Migration: added key_suffix column to api_keys")


async def init_db() -> None:
    async with engine.begin() as conn:
        await conn.run_sync(_migrate_api_keys)
        await conn.run_sync(Base.metadata.create_all)


async def get_db():
    async with async_session() as session:
        yield session
