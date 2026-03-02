from sqlalchemy.orm import DeclarativeBase


class Base(DeclarativeBase):
    pass


from app.models.user import User  # noqa: E402, F401
from app.models.api_key import ApiKey  # noqa: E402, F401
