import logging

from fastapi import APIRouter, Depends
from pydantic import BaseModel, Field

from app.deps import get_current_user
from app.models.user import User
from app.services.cos_service import get_profile, save_profile

logger = logging.getLogger("whoami")
router = APIRouter()

PROFILE_MAX_LENGTH = 5000


class UpdateProfileRequest(BaseModel):
    content: str = Field(..., max_length=PROFILE_MAX_LENGTH)


class ProfileResponse(BaseModel):
    content: str | None


@router.get("", response_model=ProfileResponse)
async def get_my_profile(user: User = Depends(get_current_user)):
    content = get_profile(user.id)
    return ProfileResponse(content=content)


@router.post("")
async def update_my_profile(
    body: UpdateProfileRequest,
    user: User = Depends(get_current_user),
):
    save_profile(user.id, body.content)
    return {"message": "Profile updated"}
