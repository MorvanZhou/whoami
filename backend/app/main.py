import asyncio
import logging
from contextlib import asynccontextmanager
from logging.handlers import TimedRotatingFileHandler
from pathlib import Path

from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from slowapi.errors import RateLimitExceeded

from app.config import settings
from app.database import async_session, init_db
from app.limiter import limiter
from app.routers import auth, cos, keys, profile, store_api, web_profile
from app.services.store_token_service import cleanup_expired_tokens

CLEANUP_INTERVAL_SECONDS = 3600  # 1 hour

# Ensure logs directory exists
LOG_DIR = Path(__file__).resolve().parent.parent / "logs"
LOG_DIR.mkdir(exist_ok=True)

# Formatter
log_formatter = logging.Formatter("%(asctime)s %(name)s %(levelname)s %(message)s")

# Console handler
console_handler = logging.StreamHandler()
console_handler.setFormatter(log_formatter)

# File handler: rotate daily, keep 30 days
file_handler = TimedRotatingFileHandler(
    filename=LOG_DIR / "whoami.log",
    when="midnight",
    interval=1,
    backupCount=30,
    encoding="utf-8",
)
file_handler.setFormatter(log_formatter)

# Root logger
logging.basicConfig(
    level=logging.INFO,
    handlers=[console_handler, file_handler],
)
logger = logging.getLogger("whoami")

# Suppress noisy COS SDK logs (e.g. NoSuchKey ERROR on new users)
logging.getLogger("qcloud_cos").setLevel(logging.WARNING)


async def _periodic_cleanup() -> None:
    """Background task: clean up expired store_api_tokens every hour."""
    while True:
        await asyncio.sleep(CLEANUP_INTERVAL_SECONDS)
        try:
            async with async_session() as db:
                deleted = await cleanup_expired_tokens(db)
                if deleted:
                    logger.info(f"[whoami] Cleaned up {deleted} expired store tokens")
        except Exception as e:
            logger.error(f"[whoami] Token cleanup error: {e}")


@asynccontextmanager
async def lifespan(app: FastAPI):
    logger.info("[whoami] Starting up, initializing database...")
    await init_db()
    logger.info("[whoami] Database initialized")
    cleanup_task = asyncio.create_task(_periodic_cleanup())
    yield
    cleanup_task.cancel()
    logger.info("[whoami] Shutting down")


app = FastAPI(
    title="whoami API",
    description="Cross-AI user identity profile sync service",
    version="0.1.0",
    lifespan=lifespan,
    root_path="/api",
    # Disable docs in production
    docs_url="/docs" if settings.env == "dev" else None,
    redoc_url="/redoc" if settings.env == "dev" else None,
    openapi_url="/openapi.json" if settings.env == "dev" else None,
)

# --- Rate limiter ---
app.state.limiter = limiter


@app.exception_handler(RateLimitExceeded)
async def rate_limit_handler(request: Request, exc: RateLimitExceeded):
    return JSONResponse(
        status_code=429,
        content={"detail": "Rate limit exceeded. Please try again later."},
    )

app.add_middleware(
    CORSMiddleware,
    allow_origins=[settings.frontend_url],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(auth.router, prefix="/auth", tags=["auth"])
app.include_router(keys.router, prefix="/keys", tags=["keys"])
app.include_router(profile.router, prefix="/profile", tags=["profile"])
app.include_router(web_profile.router, prefix="/web/profile", tags=["web-profile"])
app.include_router(cos.router, prefix="/cos", tags=["cos"])
app.include_router(store_api.router, prefix="/storeapi", tags=["storeapi"])


@app.get("/health")
async def health_check():
    return {"status": "ok"}
