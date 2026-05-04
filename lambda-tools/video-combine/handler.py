import json
import os
import subprocess
import tempfile
import uuid

import boto3
import requests

S3_BUCKET = os.environ["TOOLS_S3_BUCKET"]
PRESIGNED_URL_EXPIRY = int(os.environ.get("PRESIGNED_URL_EXPIRY", "3600"))
S3_ENDPOINT_URL = os.environ.get("AWS_ENDPOINT_URL")           # internal: container → MinIO
S3_PUBLIC_ENDPOINT_URL = os.environ.get("S3_PUBLIC_ENDPOINT_URL")   # external: browser → MinIO
# Rewrite incoming presigned URL hosts to the internal endpoint (local dev only)
S3_URL_REWRITE_FROM = os.environ.get("S3_URL_REWRITE_FROM")   # e.g. http://localhost:9000
S3_URL_REWRITE_TO = os.environ.get("S3_URL_REWRITE_TO")       # e.g. http://go-google-auth-minio:9000

s3 = boto3.client("s3", endpoint_url=S3_ENDPOINT_URL)
s3_public = boto3.client("s3", endpoint_url=S3_PUBLIC_ENDPOINT_URL) if S3_PUBLIC_ENDPOINT_URL else s3


def _rewrite_url(url: str) -> str:
    if S3_URL_REWRITE_FROM and S3_URL_REWRITE_TO:
        return url.replace(S3_URL_REWRITE_FROM, S3_URL_REWRITE_TO)
    return url


def _parse_body(event: dict) -> dict:
    """Support both API Gateway proxy format and direct RIE invocation."""
    if "body" in event:
        raw = event["body"]
        return json.loads(raw) if isinstance(raw, str) else raw
    return event


def _download(url: str, dest: str) -> None:
    with requests.get(url, stream=True, timeout=60) as r:
        r.raise_for_status()
        with open(dest, "wb") as f:
            for chunk in r.iter_content(chunk_size=8192):
                f.write(chunk)


def lambda_handler(event, context):
    try:
        body = _parse_body(event)
        video1_url = body.get("video1_url", "")
        video2_url = body.get("video2_url", "")

        if not video1_url or not video2_url:
            return _response(400, {"error": "video1_url and video2_url are required"})

        with tempfile.TemporaryDirectory() as tmp:
            v1 = os.path.join(tmp, "video1.mp4")
            v2 = os.path.join(tmp, "video2.mp4")
            concat_list = os.path.join(tmp, "list.txt")
            output = os.path.join(tmp, "output.mp4")

            _download(_rewrite_url(video1_url), v1)
            _download(_rewrite_url(video2_url), v2)

            with open(concat_list, "w") as f:
                f.write(f"file '{v1}'\nfile '{v2}'\n")

            subprocess.run(
                [
                    "ffmpeg", "-y",
                    "-f", "concat", "-safe", "0",
                    "-i", concat_list,
                    "-c", "copy",
                    output,
                ],
                check=True,
                capture_output=True,
            )

            key = f"combined/{uuid.uuid4()}.mp4"
            s3.upload_file(output, S3_BUCKET, key, ExtraArgs={"ContentType": "video/mp4"})

        presigned_url = s3_public.generate_presigned_url(
            "get_object",
            Params={"Bucket": S3_BUCKET, "Key": key},
            ExpiresIn=PRESIGNED_URL_EXPIRY,
        )

        return _response(200, {"output_url": presigned_url})

    except subprocess.CalledProcessError as e:
        return _response(500, {"error": "ffmpeg failed", "detail": e.stderr.decode()})
    except Exception as e:
        return _response(500, {"error": str(e)})


def _response(status: int, body: dict) -> dict:
    return {
        "statusCode": status,
        "headers": {"Content-Type": "application/json"},
        "body": json.dumps(body),
    }
