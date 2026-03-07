# AI Chat Backend
Backend server for AI chat application with MongoDB, JWT authentication, and RESTful API.

## Technology Stack
- **Go 1.21+** (Programming Language)
- **Gin** (Web Framework)
- **MongoDB** (Database)
- **JWT** (Authentication)
- **bcrypt** (Password Hashing)
- **CORS** (Cross-Origin Resource Sharing)

## API Endpoints

### Authentication
- `POST /api/auth/register`  
  Register a new user (username, email, password)
- `POST /api/auth/login`  
  Authenticate user and return JWT token
- `GET /api/profile`  
  Get current user profile (JWT-protected)

### Chats
- `POST /api/chats/`  
  Create a new chat (optional title)
- `POST /api/chats/message`  
  Add a message to a chat (role: user/assistant, content)
- `GET /api/chats/user/:userID`  
  List all chats for a user
- `GET /api/chats/pinned/:userID`  
  List only pinned chats for a user
- `GET /api/chats/:chatID`  
  Get full chat details including messages
- `PUT /api/chats/:chatID`  
  Update chat metadata (title, last message preview, etc.)
- `PUT /api/chats/:chatID/pin`  
  Toggle chat pin status
- `DELETE /api/chats/:chatID`  
  Delete chat and its messages


This backend serves as a lightweight, secure API layer between the Next.js frontend and local AI services (Ollama + ComfyUI).