# MLP API - Postman Collection

This directory contains Postman collection and environment files for testing the MLP API.

## Files

- `MLP.postman_collection.json` - Complete API collection with all endpoints
- `MLP.postman_environment.json` - Environment variables

## Import into Postman

### Option 1: Using Postman Desktop App

1. Open Postman
2. Click **Import** button (top left)
3. Select **File** tab
4. Choose `MLP.postman_collection.json`
5. Click **Import**

### Option 2: Using Postman Web

1. Go to [Postman Web](https://web.postman.co/)
2. Click **Import** in the workspace
3. Drag and drop the collection file
4. Click **Import**

## Import Environment

1. Click **Import** button
2. Select `MLP.postman_environment.json`
3. Click **Import**
4. Select the environment from the dropdown (top right)

## Collection Structure

### 1. Health
- **GET** Health Check - Test if API is running

### 2. Authentication
- **POST** Register - Create a new user account
- **POST** Login - Authenticate and receive JWT token

### 3. Lectures (Authenticated)
- **POST** Create Lecture - Create a new lecture
- **GET** Get All User Lectures - Retrieve all lectures for authenticated user
- **GET** Get Lecture by ID - Retrieve specific lecture

### 4. Audios (Authenticated)
- **GET** Get Audio by ID - Retrieve specific audio
- **GET** Get Audios by Lecture ID - Retrieve all audios for a lecture

### 5. Videos (Authenticated)
- **GET** Get Video by ID - Retrieve specific video
- **GET** Get Videos by Audio ID - Retrieve all videos for an audio

## Quick Start Guide

### Step 1: Start the API Server

```bash
# Using Docker
make docker-up

# Or locally
make run
```

### Step 2: Test Basic Connectivity

1. Open the **Health** folder
2. Run **Health Check** request
3. You should receive: `{"status":"ok"}`

### Step 3: Create an Account

1. Open the **Authentication** folder
2. Run **Register** request
3. The response will contain your user details
4. The `user_id` is automatically saved to the environment

### Step 4: Login

1. Run **Login** request with the same credentials
2. The JWT token is automatically saved to environment variable `jwt_token`
3. All subsequent authenticated requests will use this token

### Step 5: Create a Lecture

1. Open the **Lectures** folder
2. Run **Create Lecture** request
3. The `lecture_id` is automatically saved to the environment

### Step 6: Explore Other Endpoints

All authenticated endpoints automatically use the JWT token from the environment.

## Environment Variables

The collection uses these environment variables:

| Variable | Description | Auto-populated |
|----------|-------------|----------------|
| `base_url` | API base URL | No (default: http://localhost:8080) |
| `jwt_token` | JWT authentication token | Yes (from login) |
| `user_id` | Current user ID | Yes (from register/login) |
| `lecture_id` | Last created lecture ID | Yes (from create lecture) |
| `audio_id` | Last audio ID | Yes (from get audios) |
| `video_id` | Last video ID | No |

## Authentication

All endpoints under **Lectures**, **Audios**, and **Videos** require authentication.

The JWT token is automatically included in the `Authorization` header as:
```
Authorization: Bearer {{jwt_token}}
```

### Manual Token Setup

If you need to manually set a token:

1. Click on the environment (top right)
2. Find `jwt_token` variable
3. Paste your token in the **Current Value** field
4. Click **Save**

## Example Request Bodies

### Register/Login

```json
{
    "email": "test@example.com",
    "password": "password123"
}
```

### Create Lecture

```json
{
    "topic": "Introduction to Machine Learning",
    "content": "Machine learning is a subset of artificial intelligence..."
}
```

## Testing Workflow

### Complete Flow Test

1. **Health Check** → Verify API is running
2. **Register** → Create new user account
3. **Login** → Get JWT token (auto-saved)
4. **Create Lecture** → Create a new lecture (auto-saves lecture_id)
5. **Get All User Lectures** → View all your lectures
6. **Get Lecture by ID** → View specific lecture details
7. **Get Audios by Lecture ID** → View audios for the lecture
8. **Get Videos by Audio ID** → View videos for an audio

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
    "error": "failed to create lecture"
}
```

## Scripts

The collection includes automatic scripts that:

1. **Extract JWT token** from login response
2. **Save user_id** from register/login
3. **Save lecture_id** from create lecture
4. **Save audio_id** from get audios

These scripts run automatically - no manual intervention needed!

## Environment

The collection uses **MLP Local Environment** which works for both local development and Docker (since Docker exposes port 8080 to localhost).

## Tips

1. **Run in sequence**: Execute requests in order for the first time
2. **Check environment**: Verify variables are populated after each request
3. **Token expiration**: If you get 401 errors, re-run the Login request
4. **CORS**: If using from browser, ensure CORS is configured (already done in the backend)

## Troubleshooting

### "Could not get response"
- Check if the server is running: `curl http://localhost:8080/api/v1/health`
- Verify port 8080 is not in use by another service

### "401 Unauthorized"
- Run the Login request again to get a fresh token
- Check if JWT_EXPIRATION in .env is not too short

### "404 Not Found"
- Verify the endpoint URL matches your API routes
- Check if the resource ID exists (user_id, lecture_id, etc.)

### "500 Internal Server Error"
- Check server logs: `make docker-logs` or check terminal output
- Verify database connection is working
- Ensure MinIO is running

## Advanced Usage

### Collection Runner

1. Click **Runner** button (top toolbar)
2. Select **MLP API** collection
3. Choose the environment
4. Click **Run MLP API**

This will execute all requests in sequence.

### Export Results

After running requests:
1. Click **Runner**
2. Click on a run
3. Click **Export Results**

## Support

For issues or questions:
- Check the main [README.md](../README.md)
- Review API documentation
- Check server logs for errors

## Update Collection

When API endpoints change:
1. Re-import the updated collection file
2. Postman will ask if you want to replace or merge
3. Choose **Replace** to update all endpoints
