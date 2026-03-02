import logging
from datetime import datetime, timezone

from qcloud_cos import CosConfig, CosS3Client

from app.config import settings

logger = logging.getLogger("whoami")


def _get_cos_client() -> CosS3Client:
    config = CosConfig(
        Region=settings.cos_region,
        SecretId=settings.cos_secret_id,
        SecretKey=settings.cos_secret_key,
    )
    return CosS3Client(config)


def _profile_key(user_id: str) -> str:
    return f"profiles/{user_id}/current.md"


def _history_key(user_id: str, timestamp: str) -> str:
    return f"profiles/{user_id}/history/{timestamp}.md"


def get_profile(user_id: str) -> str | None:
    client = _get_cos_client()
    key = _profile_key(user_id)
    try:
        response = client.get_object(Bucket=settings.cos_bucket, Key=key)
        return response["Body"].get_raw_stream().read().decode("utf-8")
    except Exception as e:
        if "NoSuchKey" in str(e):
            return None
        logger.error(f"[whoami] COS get_profile error: {e}")
        raise


def save_profile(user_id: str, content: str) -> None:
    client = _get_cos_client()
    key = _profile_key(user_id)

    # Backup current to history before overwriting
    _backup_to_history(client, user_id)

    client.put_object(
        Bucket=settings.cos_bucket,
        Body=content.encode("utf-8"),
        Key=key,
        ContentType="text/markdown; charset=utf-8",
    )
    logger.info(f"[whoami] Profile saved for user {user_id}")

    # Clean up old history versions (keep latest 3)
    _cleanup_history(client, user_id)


def _backup_to_history(client: CosS3Client, user_id: str) -> None:
    key = _profile_key(user_id)
    try:
        client.get_object(Bucket=settings.cos_bucket, Key=key)
    except Exception:
        return  # No current profile to backup

    timestamp = datetime.now(timezone.utc).strftime("%Y%m%dT%H%M%SZ")
    history_key = _history_key(user_id, timestamp)

    client.copy_object(
        Bucket=settings.cos_bucket,
        Key=history_key,
        CopySource={
            "Bucket": settings.cos_bucket,
            "Key": key,
            "Region": settings.cos_region,
        },
    )
    logger.info(f"[whoami] Profile backed up to {history_key}")


def _cleanup_history(client: CosS3Client, user_id: str) -> None:
    prefix = f"profiles/{user_id}/history/"
    response = client.list_objects(Bucket=settings.cos_bucket, Prefix=prefix)
    contents = response.get("Contents", [])

    if len(contents) <= 3:
        return

    # Sort by key (timestamp-based naming ensures chronological order)
    contents.sort(key=lambda x: x["Key"])
    to_delete = contents[: len(contents) - 3]

    for item in to_delete:
        client.delete_object(Bucket=settings.cos_bucket, Key=item["Key"])
        logger.info(f"[whoami] Deleted old history: {item['Key']}")


def list_history(user_id: str) -> list[dict]:
    client = _get_cos_client()
    prefix = f"profiles/{user_id}/history/"
    response = client.list_objects(Bucket=settings.cos_bucket, Prefix=prefix)
    contents = response.get("Contents", [])

    results = []
    for item in contents:
        key = item["Key"]
        # Extract timestamp from key: profiles/{user_id}/history/{timestamp}.md
        ts = key.split("/")[-1].replace(".md", "")
        results.append({
            "timestamp": ts,
            "key": key,
            "size": item.get("Size", 0),
            "last_modified": item.get("LastModified", ""),
        })

    results.sort(key=lambda x: x["timestamp"], reverse=True)
    return results


def get_history_profile(user_id: str, timestamp: str) -> str | None:
    client = _get_cos_client()
    key = _history_key(user_id, timestamp)
    try:
        response = client.get_object(Bucket=settings.cos_bucket, Key=key)
        return response["Body"].get_raw_stream().read().decode("utf-8")
    except Exception as e:
        if "NoSuchKey" in str(e):
            return None
        logger.error(f"[whoami] COS get_history_profile error: {e}")
        raise


def rollback_profile(user_id: str, timestamp: str) -> bool:
    client = _get_cos_client()
    history_key = _history_key(user_id, timestamp)

    # Check history version exists
    try:
        client.get_object(Bucket=settings.cos_bucket, Key=history_key)
    except Exception:
        return False

    # Backup current before rollback
    _backup_to_history(client, user_id)

    # Copy history version to current
    current_key = _profile_key(user_id)
    client.copy_object(
        Bucket=settings.cos_bucket,
        Key=current_key,
        CopySource={
            "Bucket": settings.cos_bucket,
            "Key": history_key,
            "Region": settings.cos_region,
        },
    )
    logger.info(f"[whoami] Profile rolled back to {timestamp} for user {user_id}")

    _cleanup_history(client, user_id)
    return True
