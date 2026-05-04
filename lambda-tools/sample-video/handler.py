import json
import os
import subprocess
import tempfile
import uuid

import boto3

S3_BUCKET = os.environ["TOOLS_S3_BUCKET"]
PRESIGNED_URL_EXPIRY = int(os.environ.get("PRESIGNED_URL_EXPIRY", "86400"))
S3_ENDPOINT_URL = os.environ.get("AWS_ENDPOINT_URL")
S3_PUBLIC_ENDPOINT_URL = os.environ.get("S3_PUBLIC_ENDPOINT_URL")

s3 = boto3.client("s3", endpoint_url=S3_ENDPOINT_URL)
s3_public = boto3.client("s3", endpoint_url=S3_PUBLIC_ENDPOINT_URL) if S3_PUBLIC_ENDPOINT_URL else s3

SAMPLES = [
    {"color": "blue", "frequency": 440},
    {"color": "red",  "frequency": 880},
]


def _parse_body(event: dict) -> dict:
    if "body" in event:
        raw = event["body"]
        return json.loads(raw) if isinstance(raw, str) else (raw or {})
    return event


def _generate_video(color: str, frequency: int, duration: int, dest: str) -> None:
    subprocess.run(
        [
            "ffmpeg", "-y",
            "-f", "lavfi", "-i", f"color=c={color}:size=320x240:duration={duration}:rate=30",
            "-f", "lavfi", "-i", f"sine=frequency={frequency}:duration={duration}",
            "-c:v", "libx264", "-c:a", "aac",
            dest,
        ],
        check=True,
        capture_output=True,
    )


def _upload_and_sign(local_path: str, key: str) -> str:
    s3.upload_file(local_path, S3_BUCKET, key, ExtraArgs={"ContentType": "video/mp4"})
    return s3_public.generate_presigned_url(
        "get_object",
        Params={"Bucket": S3_BUCKET, "Key": key},
        ExpiresIn=PRESIGNED_URL_EXPIRY,
    )


def lambda_handler(event, context):
    try:
        body = _parse_body(event)
        duration = int(body.get("duration", 3))
        if duration < 1 or duration > 30:
            return _response(400, {"error": "duration must be between 1 and 30"})

        run_id = uuid.uuid4().hex[:8]
        urls = {}

        with tempfile.TemporaryDirectory() as tmp:
            for sample in SAMPLES:
                color = sample["color"]
                dest = f"{tmp}/{color}.mp4"
                key = f"samples/{run_id}/{color}.mp4"
                _generate_video(color, sample["frequency"], duration, dest)
                urls[f"video_{color}_url"] = _upload_and_sign(dest, key)

        return _response(200, {
            "video1_url": urls["video_blue_url"],
            "video2_url": urls["video_red_url"],
            "duration": duration,
        })

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
