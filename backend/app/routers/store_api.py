from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.responses import PlainTextResponse
from sqlalchemy.ext.asyncio import AsyncSession

from app.config import settings
from app.database import get_db
from app.services.store_token_service import consume_store_token

router = APIRouter()


@router.get("", response_class=PlainTextResponse)
async def get_store_api(
    token: str,
    db: AsyncSession = Depends(get_db),
) -> str:
    api_key = await consume_store_token(db, token)
    if not api_key:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Token not found, expired, or already used",
        )

    endpoint = settings.frontend_url
    return f"WHOAMI_API_KEY={api_key}\nWHOAMI_ENDPOINT={endpoint}\n"
