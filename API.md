# MLP API Documentation

Complete API documentation for the MLP lecture generation system.

## Overview

The MLP API provides a complete workflow for:
1. **Generating lectures** using Gemini AI
2. **Converting text to speech** using VoiceRSS
3. **Creating lip-synced videos** using avatar images

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

All authenticated endpoints require a JWT token in the `Authorization` header:

```
Authorization: Bearer <token>
```

Get the token by logging in via `/auth/login`.

---

## Complete Workflow

### Step 1: Register & Login

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "created_at": "2026-02-22T00:00:00Z"
  }
}
```

### Step 2: Generate Lecture with AI

```http
POST /api/v1/lectures
Authorization: Bearer <token>
Content-Type: application/json

{
  "topic": "Introduction to Machine Learning"
}
```

**Response:**
```json
{
  "id": "lecture-uuid",
  "topic": "Introduction to Machine Learning",
  "content": "Machine learning is a subset of artificial intelligence...",
  "status": "completed"
}
```

### Step 3: Create Audio from Lecture

```http
POST /api/v1/audios
Authorization: Bearer <token>
Content-Type: application/json

{
  "lecture_id": "lecture-uuid",
  "language": "en-us",
  "voice": "Amy",
  "rate": 0
}
```

**Response:**
```json
{
  "id": "audio-uuid",
  "lecture_id": "lecture-uuid",
  "url": "/audios/lecture-uuid/audio-uuid.mp3",
  "language": "en-us",
  "voice": "Amy",
  "created_at": "2026-02-22T00:00:00Z"
}
```

### Step 4: Create Video with Avatar

```http
POST /api/v1/videos
Authorization: Bearer <token>
Content-Type: multipart/form-data

audio_id: audio-uuid
avatar: [image file]
```

**Response:**
```json
{
  "id": "video-uuid",
  "audio_id": "audio-uuid",
  "url": "",
  "status": "processing",
  "created_at": "2026-02-22T00:00:00Z"
}
```

### Step 5: Check Video Status

```http
GET /api/v1/videos/video-uuid
Authorization: Bearer <token>
```

**Response (Processing):**
```json
{
  "id": "video-uuid",
  "audio_id": "audio-uuid",
  "url": "",
  "status": "processing",
  "created_at": "2026-02-22T00:00:00Z"
}
```

**Response (Completed):**
```json
{
  "id": "video-uuid",
  "audio_id": "audio-uuid",
  "url": "https://storage.sync.so/output/job_12345.mp4",
  "status": "completed",
  "created_at": "2026-02-22T00:00:00Z"
}
```

---

## API Endpoints

### Authentication

#### POST /auth/register
Register a new user account.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (201):**
```json
{
  "id": "user-uuid",
  "email": "user@example.com",
  "created_at": "2026-02-22T00:00:00Z"
}
```

#### POST /auth/login
Authenticate and receive JWT token.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "user-uuid",
    "email": "user@example.com",
    "created_at": "2026-02-22T00:00:00Z"
  }
}
```

---

### Lectures

#### POST /lectures
Generate a lecture using Gemini AI.

**Headers:**
- `Authorization: Bearer <token>`

**Request:**
```json
{
  "topic": "Introduction to Machine Learning"
}
```

**Response (201):**
```json
{
  "id": "lecture-uuid",
  "topic": "Introduction to Machine Learning",
  "content": "Machine learning is a subset of artificial intelligence that focuses on building systems...",
  "status": "completed"
}
```

#### GET /lectures
Get all lectures for the authenticated user.

**Headers:**
- `Authorization: Bearer <token>`

**Response (200):**
```json
[
  {
    "id": "lecture-uuid",
    "user_id": "user-uuid",
    "topic": "Introduction to Machine Learning",
    "content": "...",
    "status": "completed",
    "created_at": "2026-02-22T00:00:00Z"
  }
]
```

#### GET /lectures/:id
Get a specific lecture by ID.

**Headers:**
- `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "id": "lecture-uuid",
  "user_id": "user-uuid",
  "topic": "Introduction to Machine Learning",
  "content": "...",
  "status": "completed",
  "created_at": "2026-02-22T00:00:00Z"
}
```

---

### Audios

#### POST /audios
Convert lecture text to speech using VoiceRSS.

**Headers:**
- `Authorization: Bearer <token>`

**Request:**
```json
{
  "lecture_id": "lecture-uuid",
  "language": "en-us",
  "voice": "Amy",
  "rate": 0
}
```

**Parameters:**
- `lecture_id` (required): UUID of the lecture
- `language` (optional): Language code (default: "en-us")
  - Options: `en-us`, `en-gb`, `es-es`, `fr-fr`, `de-de`, `it-it`, `pt-br`, `ru-ru`, `zh-cn`, `ja-jp`, `ko-kr`, `ar-sa`
- `voice` (optional): Voice name (depends on language)
- `rate` (optional): Speech rate from -10 (slowest) to 10 (fastest), default: 0

**Response (201):**
```json
{
  "id": "audio-uuid",
  "lecture_id": "lecture-uuid",
  "url": "/audios/lecture-uuid/audio-uuid.mp3",
  "language": "en-us",
  "voice": "Amy",
  "created_at": "2026-02-22T00:00:00Z"
}
```

#### GET /audios/:id
Get a specific audio by ID.

**Headers:**
- `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "id": "audio-uuid",
  "lecture_id": "lecture-uuid",
  "url": "/audios/lecture-uuid/audio-uuid.mp3",
  "language": "en-us",
  "voice": "Amy",
  "created_at": "2026-02-22T00:00:00Z"
}
```

#### GET /audios/lecture/:lecture_id
Get all audios for a specific lecture.

**Headers:**
- `Authorization: Bearer <token>`

**Response (200):**
```json
[
  {
    "id": "audio-uuid",
    "lecture_id": "lecture-uuid",
    "url": "/audios/lecture-uuid/audio-uuid.mp3",
    "language": "en-us",
    "voice": "Amy",
    "created_at": "2026-02-22T00:00:00Z"
  }
]
```

---

### Videos

#### POST /videos
Create a lip-synced video from audio and avatar image.

**Headers:**
- `Authorization: Bearer <token>`
- `Content-Type: multipart/form-data`

**Form Data:**
- `audio_id` (required): UUID of the audio
- `avatar` (required): Avatar image file (JPEG/PNG)

**Response (202):**
```json
{
  "id": "video-uuid",
  "audio_id": "audio-uuid",
  "url": "",
  "status": "processing",
  "created_at": "2026-02-22T00:00:00Z"
}
```

**Status Values:**
- `processing`: Video is being generated
- `completed`: Video is ready with URL
- `failed`: Generation failed
- `timeout`: Generation timed out

#### GET /videos/:id
Get video status and URL.

**Headers:**
- `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "id": "video-uuid",
  "audio_id": "audio-uuid",
  "url": "https://storage.sync.so/output/job_12345.mp4",
  "status": "completed",
  "created_at": "2026-02-22T00:00:00Z"
}
```

#### GET /videos/audio/:audio_id
Get all videos for a specific audio.

**Headers:**
- `Authorization: Bearer <token>`

**Response (200):**
```json
[
  {
    "id": "video-uuid",
    "audio_id": "audio-uuid",
    "url": "https://storage.sync.so/output/job_12345.mp4",
    "status": "completed",
    "created_at": "2026-02-22T00:00:00Z"
  }
]
```

---

## Error Responses

### 400 Bad Request
```json
{
  "error": "invalid request body"
}
```

### 401 Unauthorized
```json
{
  "error": "missing authorization header"
}
```

### 404 Not Found
```json
{
  "error": "lecture not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "failed to generate lecture"
}
```

---

## Storage Structure

### MinIO Buckets

1. **audios**: MP3 files from TTS
   - Path: `/{lecture_id}/{audio_id}.mp3`

2. **videos**: Generated lip-sync videos
   - Path: `/{audio_id}/{video_id}.mp4`

3. **avatars**: Uploaded avatar images
   - Path: `/{audio_id}/{avatar_id}.jpg`

---

## Supported Languages (VoiceRSS)

| Code | Language |
|------|----------|
| en-us | English (United States) |
| en-gb | English (United Kingdom) |
| es-es | Spanish (Spain) |
| fr-fr | French (France) |
| de-de | German (Germany) |
| it-it | Italian (Italy) |
| pt-br | Portuguese (Brazil) |
| ru-ru | Russian (Russia) |
| zh-cn | Chinese (China) |
| ja-jp | Japanese (Japan) |
| ko-kr | Korean (South Korea) |
| ar-sa | Arabic (Saudi Arabia) |

---

## Rate Limits

Currently no rate limits enforced. Production deployment should implement:
- Rate limiting per user
- Request throttling
- Maximum file sizes

---

## Security

- All passwords are hashed using bcrypt
- JWT tokens expire after 24 hours (configurable)
- All authenticated endpoints require valid JWT
- File uploads are validated for type and size

---

## Development

### Testing with cURL

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Generate Lecture
curl -X POST http://localhost:8080/api/v1/lectures \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"topic":"Machine Learning Basics"}'

# Create Audio
curl -X POST http://localhost:8080/api/v1/audios \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"lecture_id":"<uuid>","language":"en-us"}'

# Create Video
curl -X POST http://localhost:8080/api/v1/videos \
  -H "Authorization: Bearer <token>" \
  -F "audio_id=<uuid>" \
  -F "avatar=@/path/to/image.jpg"
```

---

## Next Steps

1. Add pagination to list endpoints
2. Implement webhook notifications for video completion
3. Add video preview thumbnails
4. Support multiple audio formats
5. Add batch processing endpoints
6. Implement WebSocket for real-time status updates
