import logging

from fastapi import APIRouter, Depends, HTTPException, status
from pydantic import BaseModel, Field

from app.deps import get_user_by_apikey
from app.models.user import User
from app.services.cos_service import (
    get_history_profile,
    get_profile,
    list_history,
    rollback_profile,
    save_profile,
)

logger = logging.getLogger("whoami")
router = APIRouter()

# 限制 profile 最大 5000 字符（约 2000 中文字或 2000 英文单词）
PROFILE_MAX_LENGTH = 5000


class UpdateProfileRequest(BaseModel):
    content: str = Field(..., max_length=PROFILE_MAX_LENGTH)


class ProfileResponse(BaseModel):
    content: str | None


class HistoryItem(BaseModel):
    timestamp: str
    size: int
    last_modified: str


@router.get("", response_model=ProfileResponse)
async def get_user_profile(user: User = Depends(get_user_by_apikey)):
    content = get_profile(user.id)
    logger.info("[whoami] api get_profile user=%s found=%s", user.id, content is not None)
    return ProfileResponse(content=content)


@router.post("", status_code=status.HTTP_200_OK)
async def update_user_profile(
    body: UpdateProfileRequest,
    user: User = Depends(get_user_by_apikey),
):
    save_profile(user.id, body.content)
    logger.info("[whoami] api update_profile user=%s len=%d", user.id, len(body.content))
    return {"message": "Profile updated"}


@router.get("/history", response_model=list[HistoryItem])
async def get_profile_history(user: User = Depends(get_user_by_apikey)):
    items = list_history(user.id)
    return [
        HistoryItem(
            timestamp=item["timestamp"],
            size=item["size"],
            last_modified=item["last_modified"],
        )
        for item in items
    ]


@router.post("/rollback/{timestamp}")
async def rollback_to_version(
    timestamp: str,
    user: User = Depends(get_user_by_apikey),
):
    success = rollback_profile(user.id, timestamp)
    if not success:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="History version not found",
        )
    return {"message": f"Profile rolled back to {timestamp}"}
