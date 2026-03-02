from fastapi import APIRouter, Depends, HTTPException, status
from pydantic import BaseModel
from sqlalchemy.ext.asyncio import AsyncSession

from app.database import get_db
from app.deps import get_current_user
from app.models.user import User
from app.services.key_service import create_api_key, list_api_keys, revoke_api_key

router = APIRouter()


class CreateKeyRequest(BaseModel):
    label: str | None = None


class ApiKeyResponse(BaseModel):
    id: str
    key_prefix: str
    key_suffix: str
    label: str | None
    created_at: str
    last_used: str | None


class CreateKeyResponse(BaseModel):
    id: str
    key: str
    key_prefix: str
    key_suffix: str
    label: str | None
    created_at: str


@router.post("", response_model=CreateKeyResponse, status_code=status.HTTP_201_CREATED)
async def create_key(
    body: CreateKeyRequest,
    user: User = Depends(get_current_user),
    db: AsyncSession = Depends(get_db),
):
    plain_key, api_key = await create_api_key(db, user.id, body.label)
    return CreateKeyResponse(
        id=api_key.id,
        key=plain_key,
        key_prefix=api_key.key_prefix,
        key_suffix=api_key.key_suffix,
        label=api_key.label,
        created_at=api_key.created_at.isoformat() if api_key.created_at else "",
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
