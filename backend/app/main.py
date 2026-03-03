import logging
from contextlib import asynccontextmanager
from logging.handlers import TimedRotatingFileHandler
from pathlib import Path

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.config import settings
from app.database import init_db
from app.routers import auth, cos, keys, profile, store_api, web_profile

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


@asynccontextmanager
async def lifespan(app: FastAPI):
    logger.info("[whoami] Starting up, initializing database...")
    await init_db()
    logger.info("[whoami] Database initialized")
    yield
    logger.info("[whoami] Shutting down")


app = FastAPI(
    title="whoami API",
    description="Cross-AI user identity profile sync service",
    version="0.1.0",
    lifespan=lifespan,
    root_path="/api",
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
