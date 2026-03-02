import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.config import settings
from app.database import init_db
from app.routers import auth, cos, keys, profile, web_profile

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(name)s %(levelname)s %(message)s",
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


@app.get("/health")
async def health_check():
    return {"status": "ok"}
