from fastapi import APIRouter, Depends, HTTPException, status
from pydantic import BaseModel, Field
from sqlalchemy.ext.asyncio import AsyncSession

from app.database import get_db
from app.deps import get_current_user
from app.models.user import User
from app.config import settings
from app.services.key_service import create_api_key, list_api_keys, revoke_api_key
from app.services.store_token_service import TOKEN_EXPIRE_MINUTES, create_store_token, is_store_token_consumed

router = APIRouter()


class CreateKeyRequest(BaseModel):
    label: str = Field(..., min_length=1, max_length=15)


class ApiKeyResponse(BaseModel):
    id: str
    key_prefix: str
    key_suffix: str
    label: str | None
    created_at: str
    last_used: str | None


class CreateKeyResponse(BaseModel):
    id: str
    label: str
    store_url: str
    expires_in: int


@router.post("", response_model=CreateKeyResponse, status_code=status.HTTP_201_CREATED)
async def create_key(
    body: CreateKeyRequest,
    user: User = Depends(get_current_user),
    db: AsyncSession = Depends(get_db),
):
    plain_key, api_key = await create_api_key(db, user.id, body.label)
    token = await create_store_token(db, user.id, plain_key, api_key.id)
    store_url = f"{settings.frontend_url}/api/storeapi?token={token}"
    return CreateKeyResponse(
        id=api_key.id,
        label=api_key.label or "",
        store_url=store_url,
        expires_in=TOKEN_EXPIRE_MINUTES * 60,
    )


@router.get("", response_model=list[ApiKeyResponse])
async def list_keys(
    user: User = Depends(get_current_user),
    db: AsyncSession = Depends(get_db),
):
    keys = await list_api_keys(db, user.id)
    return [
        ApiKeyResponse(
            id=k.id,
            key_prefix=k.key_prefix,
            key_suffix=k.key_suffix,
            label=k.label,
            created_at=k.created_at.isoformat() if k.created_at else "",
            last_used=k.last_used.isoformat() if k.last_used else None,
        )
        for k in keys
    ]


@router.delete("/{key_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_key(
    key_id: str,
    user: User = Depends(get_current_user),
    db: AsyncSession = Depends(get_db),
):
    success = await revoke_api_key(db, key_id, user.id)
    if not success:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="API key not found")


@router.get("/{key_id}/status")
async def get_key_setup_status(
    key_id: str,
    user: User = Depends(get_current_user),
    db: AsyncSession = Depends(get_db),
):
    """Check if the store token for this key has been consumed (agent completed setup)."""
    consumed = await is_store_token_consumed(db, key_id, user.id)
    return {"configured": consumed}
