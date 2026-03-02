import logging
import secrets

from authlib.integrations.httpx_client import AsyncOAuth2Client
from fastapi import APIRouter, Depends, HTTPException, Request, Response, status
from fastapi.responses import RedirectResponse
from sqlalchemy.ext.asyncio import AsyncSession

from app.config import settings
from app.database import get_db
from app.deps import get_current_user
from app.models.user import User
from app.services.auth_service import create_jwt_token, find_or_create_user

logger = logging.getLogger("whoami")
router = APIRouter()

# In-memory state store for OAuth CSRF protection
_oauth_states: dict[str, str] = {}

GITHUB_AUTHORIZE_URL = "https://github.com/login/oauth/authorize"
GITHUB_TOKEN_URL = "https://github.com/login/oauth/access_token"
GITHUB_USER_URL = "https://api.github.com/user"
GITHUB_EMAILS_URL = "https://api.github.com/user/emails"

GOOGLE_AUTHORIZE_URL = "https://accounts.google.com/o/oauth2/v2/auth"
GOOGLE_TOKEN_URL = "https://oauth2.googleapis.com/token"
GOOGLE_USERINFO_URL = "https://www.googleapis.com/oauth2/v2/userinfo"


@router.get("/github/login")
async def github_login(request: Request, redirect: str | None = None):
    state = secrets.token_urlsafe(32)
    # Encode redirect target in state
    state_value = redirect or ""
    _oauth_states[state] = state_value

    client = AsyncOAuth2Client(
        client_id=settings.github_client_id,
        client_secret=settings.github_client_secret,
        redirect_uri=settings.github_redirect_uri,
    )
    url, _ = client.create_authorization_url(GITHUB_AUTHORIZE_URL, state=state, scope="user:email")
    return RedirectResponse(url=url)


@router.get("/github/callback")
async def github_callback(
    code: str,
    state: str,
    db: AsyncSession = Depends(get_db),
):
    if state not in _oauth_states:
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="Invalid state")

    redirect_target = _oauth_states.pop(state)

    async with AsyncOAuth2Client(
        client_id=settings.github_client_id,
        client_secret=settings.github_client_secret,
        redirect_uri=settings.github_redirect_uri,
    ) as client:
        token = await client.fetch_token(GITHUB_TOKEN_URL, code=code)
        client.token = token

        resp = await client.get(GITHUB_USER_URL)
        user_data = resp.json()

        # Get primary email
        email = user_data.get("email")
        if not email:
            emails_resp = await client.get(GITHUB_EMAILS_URL)
            emails = emails_resp.json()
            primary = next((e for e in emails if e.get("primary")), None)
            email = primary["email"] if primary else None

    user = await find_or_create_user(
        db=db,
        provider="github",
        provider_id=str(user_data["id"]),
        email=email,
        name=user_data.get("name") or user_data.get("login"),
        avatar_url=user_data.get("avatar_url"),
    )

    jwt_token = create_jwt_token(user.id)
    redirect_url = f"{settings.frontend_url}/auth/callback?token={jwt_token}"
    if redirect_target:
        redirect_url += f"&redirect={redirect_target}"

    return RedirectResponse(url=redirect_url)


@router.get("/google/login")
async def google_login(request: Request, redirect: str | None = None):
    state = secrets.token_urlsafe(32)
    state_value = redirect or ""
    _oauth_states[state] = state_value

    client = AsyncOAuth2Client(
        client_id=settings.google_client_id,
        client_secret=settings.google_client_secret,
        redirect_uri=settings.google_redirect_uri,
    )
    url, _ = client.create_authorization_url(
        GOOGLE_AUTHORIZE_URL,
        state=state,
        scope="openid email profile",
        access_type="offline",
    )
    return RedirectResponse(url=url)


@router.get("/google/callback")
async def google_callback(
    code: str,
    state: str,
    db: AsyncSession = Depends(get_db),
):
    if state not in _oauth_states:
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="Invalid state")

    redirect_target = _oauth_states.pop(state)

    async with AsyncOAuth2Client(
        client_id=settings.google_client_id,
        client_secret=settings.google_client_secret,
        redirect_uri=settings.google_redirect_uri,
    ) as client:
        token = await client.fetch_token(GOOGLE_TOKEN_URL, code=code)
        client.token = token

        resp = await client.get(GOOGLE_USERINFO_URL)
        user_data = resp.json()

    user = await find_or_create_user(
        db=db,
        provider="google",
        provider_id=str(user_data["id"]),
        email=user_data.get("email"),
        name=user_data.get("name"),
        avatar_url=user_data.get("picture"),
    )

    jwt_token = create_jwt_token(user.id)
    redirect_url = f"{settings.frontend_url}/auth/callback?token={jwt_token}"
    if redirect_target:
        redirect_url += f"&redirect={redirect_target}"

    return RedirectResponse(url=redirect_url)


@router.get("/me")
async def get_me(user: User = Depends(get_current_user)):
    return {
        "id": user.id,
        "email": user.email,
        "name": user.name,
        "avatar_url": user.avatar_url,
        "provider": user.provider,
        "created_at": user.created_at.isoformat() if user.created_at else None,
    }


@router.post("/logout")
async def logout(response: Response):
    response.delete_cookie("access_token")
    return {"message": "Logged out"}
