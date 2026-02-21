import os
import subprocess
import tempfile
import threading
import uuid

from flask import Flask, request, jsonify
from minio import Minio

app = Flask(__name__)

SYNC_API_KEY = os.getenv('SYNC_API_KEY', '')
MINIO_ENDPOINT = os.getenv('MINIO_ENDPOINT', 'minio:9000')
MINIO_ACCESS_KEY = os.getenv('MINIO_ACCESS_KEY', 'minioadmin')
MINIO_SECRET_KEY = os.getenv('MINIO_SECRET_KEY', 'minioadmin')
MINIO_BUCKET_VIDEOS = os.getenv('MINIO_BUCKET_VIDEOS', 'videos')
MINIO_USE_SSL = os.getenv('MINIO_USE_SSL', 'false').lower() == 'true'

# In-memory job store: job_id -> {"status": "processing"|"completed"|"failed", "output_url": str|None}
jobs = {}
jobs_lock = threading.Lock()


def get_minio_client():
    return Minio(
        MINIO_ENDPOINT,
        access_key=MINIO_ACCESS_KEY,
        secret_key=MINIO_SECRET_KEY,
        secure=MINIO_USE_SSL,
    )


def run_lipsync(job_id: str, avatar_data: bytes, audio_data: bytes, avatar_filename: str, audio_filename: str):
    """
    Принимает аватар и аудио из тела запроса (bytes). Пишет во временные файлы,
    собирает видео через ffmpeg, загружает результат в MinIO.
    """
    try:
        with jobs_lock:
            jobs[job_id] = {"status": "processing", "output_url": None}

        client = get_minio_client()

        with tempfile.TemporaryDirectory() as tmpdir:
            avatar_path = os.path.join(tmpdir, avatar_filename or "avatar.jpg")
            audio_path = os.path.join(tmpdir, audio_filename or "audio.mp3")
            output_path = os.path.join(tmpdir, "output.mp4")

            with open(avatar_path, "wb") as f:
                f.write(avatar_data)
            with open(audio_path, "wb") as f:
                f.write(audio_data)

            if not avatar_data or not audio_data:
                raise ValueError("Avatar or audio is empty")

            # Create video: static image + audio using ffmpeg
            cmd = [
                "ffmpeg", "-y",
                "-loop", "1", "-i", avatar_path,
                "-i", audio_path,
                "-c:v", "libx264", "-tune", "stillimage",
                "-c:a", "aac", "-shortest",
                "-pix_fmt", "yuv420p",
                output_path,
            ]
            result = subprocess.run(cmd, capture_output=True, text=True, timeout=300)
            if result.returncode != 0:
                raise RuntimeError(f"ffmpeg failed: {result.stderr[:500]}")

            if not os.path.exists(output_path) or os.path.getsize(output_path) == 0:
                raise RuntimeError("ffmpeg produced no output")

            # Загружаем результат в MinIO
            object_name = f"lipsync/{job_id}.mp4"
            with open(output_path, "rb") as f:
                client.put_object(
                    MINIO_BUCKET_VIDEOS,
                    object_name,
                    f,
                    length=os.path.getsize(output_path),
                    content_type="video/mp4",
                )

            output_url = f"/{MINIO_BUCKET_VIDEOS}/{object_name}"

            with jobs_lock:
                jobs[job_id] = {"status": "completed", "output_url": output_url}

    except Exception as e:
        with jobs_lock:
            jobs[job_id] = {"status": "failed", "output_url": None}
        app.logger.exception("lipsync job %s failed: %s", job_id, e)


@app.route('/health', methods=['GET'])
def health():
    return jsonify({"status": "ok", "service": "lipsync"}), 200


@app.route('/sync', methods=['POST'])
def lipsync():
    """
    Start a lip sync job. Expects multipart/form-data: avatar (image file), audio (audio file).
    Returns job_id and status. Poll GET /sync/<job_id> for completion and output_url.
    """
    try:
        if "avatar" not in request.files or "audio" not in request.files:
            return jsonify({"error": "avatar and audio files are required in request body"}), 400

        avatar_file = request.files["avatar"]
        audio_file = request.files["audio"]

        if avatar_file.filename == "" or audio_file.filename == "":
            return jsonify({"error": "avatar and audio must be non-empty files"}), 400

        avatar_data = avatar_file.read()
        audio_data = audio_file.read()

        if not avatar_data or not audio_data:
            return jsonify({"error": "avatar and audio content cannot be empty"}), 400

        job_id = f"job_{uuid.uuid4().hex[:12]}"

        thread = threading.Thread(
            target=run_lipsync,
            args=(job_id, avatar_data, audio_data, avatar_file.filename or "avatar.jpg", audio_file.filename or "audio.mp3"),
            daemon=True,
        )
        thread.start()

        return jsonify({
            "job_id": job_id,
            "status": "processing",
            "message": "Lip sync job started",
        }), 202

    except Exception as e:
        return jsonify({"error": str(e)}), 500


@app.route('/sync/<job_id>', methods=['GET'])
def get_job_status(job_id):
    """Return current job status and output_url when completed."""
    with jobs_lock:
        job = jobs.get(job_id)

    if job is None:
        return jsonify({"error": "Job not found", "job_id": job_id}), 404

    return jsonify({
        "job_id": job_id,
        "status": job["status"],
        "output_url": job.get("output_url") or "",
    }), 200


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)
