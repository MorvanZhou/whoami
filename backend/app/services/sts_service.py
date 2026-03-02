import json
import logging

from sts.sts import Sts

from app.config import settings

logger = logging.getLogger("whoami")


def get_sts_credential(user_id: str) -> dict:
    sts_config = {
        "url": "https://sts.tencentcloudapi.com/",
        "domain": "sts.tencentcloudapi.com",
        "duration_seconds": 1800,
        "secret_id": settings.cos_secret_id,
        "secret_key": settings.cos_secret_key,
        "bucket": settings.cos_bucket,
        "region": settings.cos_region,
        "allow_prefixes": [f"profiles/{user_id}/*"],
        "allow_actions": [
            "name/cos:GetObject",
            "name/cos:PutObject",
            "name/cos:HeadObject",
            "name/cos:GetBucket",
        ],
        "policy": json.dumps({
            "version": "2.0",
            "statement": [
                {
                    "action": [
                        "cos:GetObject",
                        "cos:PutObject",
                        "cos:HeadObject",
                        "cos:GetBucket",
                    ],
                    "effect": "allow",
                    "resource": [
                        f"qcs::cos:{settings.cos_region}:uid/{settings.cos_bucket.split('-')[-1]}:"
                        f"{settings.cos_bucket}/profiles/{user_id}/*"
                    ],
                }
            ],
        }),
    }

    try:
        sts = Sts(sts_config)
        response = sts.get_credential()
        logger.info(f"[whoami] STS credential generated for user {user_id}")
        return {
            "credentials": response["credentials"],
            "expiredTime": response["expiredTime"],
            "startTime": response.get("startTime", ""),
            "bucket": settings.cos_bucket,
            "region": settings.cos_region,
            "profileKey": f"profiles/{user_id}/current.md",
            "historyPrefix": f"profiles/{user_id}/history/",
        }
    except Exception as e:
        logger.error(f"[whoami] STS generation error: {e}")
        raise
