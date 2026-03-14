import hashlib
import logging
import secrets
import time

from authlib.integrations.httpx_client import AsyncOAuth2Client
from cachetools import TTLCache
from fastapi import APIRouter, Depends, Form, HTTPException, Request, Response, status
from fastapi.responses import RedirectResponse
from jose import jwt as jose_jwt
from sqlalchemy.ext.asyncio import AsyncSession

from app.config import settings
from app.database import get_db
from app.deps import get_current_user
from app.limiter import limiter
from app.models.user import User
from app.services.auth_service import create_jwt_token, find_or_create_user

logger = logging.getLogger("whoami")
router = APIRouter()

# In-memory state store for OAuth CSRF protection (TTL 10 min, max 10k entries)
_oauth_states: TTLCache[str, str] = TTLCache(maxsize=10000, ttl=600)

# PKCE code verifiers for Twitter OAuth 2.0
_pkce_verifiers: TTLCache[str, str] = TTLCache(maxsize=10000, ttl=600)

# Whitelist of allowed redirect paths after OAuth callback
ALLOWED_REDIRECTS = {"/dashboard", "/settings", "/keys"}

GITHUB_AUTHORIZE_URL = "https://github.com/login/oauth/authorize"
GITHUB_TOKEN_URL = "https://github.com/login/oauth/access_token"
GITHUB_USER_URL = "https://api.github.com/user"
GITHUB_EMAILS_URL = "https://api.github.com/user/emails"

GOOGLE_AUTHORIZE_URL = "https://accounts.google.com/o/oauth2/v2/auth"
GOOGLE_TOKEN_URL = "https://oauth2.googleapis.com/token"
GOOGLE_USERINFO_URL = "https://www.googleapis.com/oauth2/v2/userinfo"

MICROSOFT_AUTHORIZE_URL = "https://login.microsoftonline.com/common/oauth2/v2.0/authorize"
MICROSOFT_TOKEN_URL = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
MICROSOFT_USERINFO_URL = "https://graph.microsoft.com/v1.0/me"

APPLE_AUTHORIZE_URL = "https://appleid.apple.com/auth/authorize"
APPLE_TOKEN_URL = "https://appleid.apple.com/auth/token"

TWITTER_AUTHORIZE_URL = "https://twitter.com/i/oauth2/authorize"
TWITTER_TOKEN_URL = "https://api.twitter.com/2/oauth2/token"
TWITTER_USERINFO_URL = "https://api.twitter.com/2/users/me"


@router.get("/github/login")
@limiter.limit("5/minute")
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

    return _build_callback_response(create_jwt_token(user.id), redirect_target)


@router.get("/google/login")
@limiter.limit("5/minute")
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

    return _build_callback_response(create_jwt_token(user.id), redirect_target)


def _build_callback_response(jwt_token: str, redirect_target: str) -> RedirectResponse:
    if redirect_target and redirect_target not in ALLOWED_REDIRECTS:
        redirect_target = ""
    redirect_url = f"{settings.frontend_url}/auth/callback"
    if redirect_target:
        redirect_url += f"?redirect={redirect_target}"
    response = RedirectResponse(url=redirect_url)
    response.set_cookie(
        key="access_token",
        value=jwt_token,
        httponly=True,
        secure=True,
        samesite="lax",
        max_age=settings.jwt_expire_days * 86400,
        path="/",
    )
    return response


# --------------- Microsoft ---------------

@router.get("/microsoft/login")
@limiter.limit("5/minute")
async def microsoft_login(request: Request, redirect: str | None = None):
    state = secrets.token_urlsafe(32)
    _oauth_states[state] = redirect or ""

    client = AsyncOAuth2Client(
        client_id=settings.microsoft_client_id,
        client_secret=settings.microsoft_client_secret,
        redirect_uri=settings.microsoft_redirect_uri,
    )
    url, _ = client.create_authorization_url(
        MICROSOFT_AUTHORIZE_URL,
        state=state,
        scope="openid email profile User.Read",
    )
    return RedirectResponse(url=url)


@router.get("/microsoft/callback")
async def microsoft_callback(
    code: str,
    state: str,
    db: AsyncSession = Depends(get_db),
):
    if state not in _oauth_states:
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="Invalid state")

    redirect_target = _oauth_states.pop(state)

    async with AsyncOAuth2Client(
        client_id=settings.microsoft_client_id,
        client_secret=settings.microsoft_client_secret,
        redirect_uri=settings.microsoft_redirect_uri,
    ) as client:
        token = await client.fetch_token(MICROSOFT_TOKEN_URL, code=code)
        client.token = token

        resp = await client.get(MICROSOFT_USERINFO_URL)
        user_data = resp.json()

    user = await find_or_create_user(
        db=db,
        provider="microsoft",
        provider_id=str(user_data["id"]),
        email=user_data.get("mail") or user_data.get("userPrincipalName"),
        name=user_data.get("displayName"),
        avatar_url=None,
    )

    return _build_callback_response(create_jwt_token(user.id), redirect_target)


# --------------- Apple ---------------

def _generate_apple_client_secret() -> str:
    """Generate a short-lived JWT as Apple OAuth client_secret."""
    now = int(time.time())
    headers = {"kid": settings.apple_key_id, "alg": "ES256"}
    payload = {
        "iss": settings.apple_team_id,
        "iat": now,
        "exp": now + 86400 * 180,
        "aud": "https://appleid.apple.com",
        "sub": settings.apple_client_id,
    }
    return jose_jwt.encode(payload, settings.apple_private_key, algorithm="ES256", headers=headers)


@router.get("/apple/login")
@limiter.limit("5/minute")
async def apple_login(request: Request, redirect: str | None = None):
    state = secrets.token_urlsafe(32)
    _oauth_states[state] = redirect or ""

    client = AsyncOAuth2Client(
        client_id=settings.apple_client_id,
        redirect_uri=settings.apple_redirect_uri,
    )
    url, _ = client.create_authorization_url(
        APPLE_AUTHORIZE_URL,
        state=state,
        scope="name email",
        response_mode="form_post",
    )
    return RedirectResponse(url=url)


@router.post("/apple/callback")
async def apple_callback(
    code: str = Form(...),
    state: str = Form(...),
    db: AsyncSession = Depends(get_db),
):
    if state not in _oauth_states:
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="Invalid state")

    redirect_target = _oauth_states.pop(state)

    client_secret = _generate_apple_client_secret()

    async with AsyncOAuth2Client(
        client_id=settings.apple_client_id,
        client_secret=client_secret,
        redirect_uri=settings.apple_redirect_uri,
    ) as client:
        token = await client.fetch_token(
            APPLE_TOKEN_URL,
            code=code,
            grant_type="authorization_code",
        )

    id_token = token.get("id_token", "")
    claims = jose_jwt.get_unverified_claims(id_token)

    apple_user_id = claims.get("sub", "")
    email = claims.get("email")

    user = await find_or_create_user(
        db=db,
        provider="apple",
        provider_id=apple_user_id,
        email=email,
        name=None,
        avatar_url=None,
    )

    return _build_callback_response(create_jwt_token(user.id), redirect_target)


# --------------- Twitter/X ---------------

def _generate_pkce_pair() -> tuple[str, str]:
    """Generate PKCE code_verifier and S256 code_challenge."""
    import base64

    verifier = secrets.token_urlsafe(64)[:128]
    digest = hashlib.sha256(verifier.encode("ascii")).digest()
    challenge = base64.urlsafe_b64encode(digest).rstrip(b"=").decode("ascii")
    return verifier, challenge


@router.get("/twitter/login")
@limiter.limit("5/minute")
async def twitter_login(request: Request, redirect: str | None = None):
    state = secrets.token_urlsafe(32)
    _oauth_states[state] = redirect or ""

    verifier, challenge = _generate_pkce_pair()
    _pkce_verifiers[state] = verifier

    client = AsyncOAuth2Client(
        client_id=settings.twitter_client_id,
        client_secret=settings.twitter_client_secret,
        redirect_uri=settings.twitter_redirect_uri,
    )
    url, _ = client.create_authorization_url(
        TWITTER_AUTHORIZE_URL,
        state=state,
        scope="tweet.read users.read offline.access",
        code_challenge=challenge,
        code_challenge_method="S256",
    )
    return RedirectResponse(url=url)


@router.get("/twitter/callback")
async def twitter_callback(
    code: str,
    state: str,
    db: AsyncSession = Depends(get_db),
):
    if state not in _oauth_states:
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="Invalid state")

    redirect_target = _oauth_states.pop(state)
    verifier = _pkce_verifiers.pop(state, "")

    async with AsyncOAuth2Client(
        client_id=settings.twitter_client_id,
        client_secret=settings.twitter_client_secret,
        redirect_uri=settings.twitter_redirect_uri,
        token_endpoint_auth_method="client_secret_basic",
    ) as client:
        token = await client.fetch_token(
            TWITTER_TOKEN_URL,
            code=code,
            code_verifier=verifier,
        )
        client.token = token

        resp = await client.get(
            TWITTER_USERINFO_URL,
            params={"user.fields": "profile_image_url,name,username"},
        )
        user_data = resp.json().get("data", {})

    user = await find_or_create_user(
        db=db,
        provider="twitter",
        provider_id=str(user_data.get("id", "")),
        email=None,
        name=user_data.get("name"),
        avatar_url=user_data.get("profile_image_url"),
    )

    return _build_callback_response(create_jwt_token(user.id), redirect_target)


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
