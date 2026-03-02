from fastapi import APIRouter, Depends

from app.deps import get_user_by_apikey
from app.models.user import User
from app.services.sts_service import get_sts_credential

router = APIRouter()


@router.get("/sts")
async def get_cos_sts(user: User = Depends(get_user_by_apikey)):
    credential = get_sts_credential(user.id)
    return credential
