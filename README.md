# AI Chat Backend

Backend server for AI chat application with MongoDB, JWT authentication, and RESTful API.

## Technology Stack

- **Go 1.21** (Programming Language)
- **Gin** (Web Framework)
- **MongoDB** (Database)
- **JWT** (Authentication)
- **bcrypt** (Password Hashing)
- **CORS** (Cross-Origin Resource Sharing)

## API Endpoints

### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `GET /api/profile` - Get user profile (protected)

### Chats
- `POST /api/chats/` - Create new chat
- `POST /api/chats/message` - Add message to chat
- `GET /api/chats/user/:userID` - Get user chats
- `GET /api/chats/pinned/:userID` - Get pinned chats
- `GET /api/chats/:chatID` - Get chat details
- `PUT /api/chats/:chatID` - Update chat
- `PUT /api/chats/:chatID/pin` - Toggle pin chat
- `DELETE /api/chats/:chatID` - Delete chat

## Integration with External Services

### Stable Diffusion
The backend integrates with [AbdBarho/stable-diffusion-webui-docker](https://github.com/AbdBarho/stable-diffusion-webui-docker) for image generation:
- Communicates with Stable Diffusion API on port `7860`
- Receives image generation requests from frontend
- Stores generated images as base64 strings in chat messages
- Compatible with AUTOMATIC1111 web UI running in Docker