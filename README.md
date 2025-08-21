# Dune Security Takehome Assessment

## Live Demo

[https://pretty-imagination-production-3bad.up.railway.app/]

## Project Description

This is a takehome project given to me by Dune Security. It is a dynamic form builder application with ability to see live
analytics of each form that is created. This app also allows for a person to generate a link and share their forms with other people. There is an authentication flow that makes every form unique to one person. You can preview the form, view your form, edit a form, publish them, and save as drafts

## Extra Credit

I am using JWT from the extra credit. I think that should be a requirement for the project, rather than extra credit. Just some small feedback.
I did not have time to implement a DarkMode, but I would approach it with a ThemeProvider and ThemeContext nd utilize that to set DarkMode throughout the app.
I did not have time to add unit tests
I did not have time to add CSV exports, however, i do not think that would be too difficult to implement

## Things To Add

- I could add some way to delete forms. I feel like that should be a requirement, but it was not specified.
- I could add a way to group forms together like how they are with the published and draft tags. That way, you can add different forms to a single group of forms. 

## Tech Stack

- **Frontend**: Next.js 15, React 19, TailwindCSS 4, TypeScript
- **Backend**: Go Fiber API
- **Database**: MongoDB
- **Real-time**: Socket.IO / WebSockets / Go routine
- **Deployment**: Railway

## Getting Started

### Prerequisites

- Node.js (v18+)
- Go (v1.21+)
- MongoDB (local or cloud instance; preferrably local)
- Git

### Local Development Setup

#### 1. Clone the Repository

```bash
git clone https://github.com/1ShoukR/dune-takehome.git
cd dune-takehome
```

#### 2. Backend Setup (Go Fiber)

```bash
# Navigate to server directory
cd server

# Install Go dependencies
go mod download

# Create .env file
cp .env.example .env

# Update .env with your configuration:
# MONGODB_URI=mongodb://localhost:27017/formbuilder
# PORT=8080
# CLIENT_URL=http://localhost:3000

# Run the server
go run cmd/main.go
```

#### 3. Frontend Setup (Next.js)

```bash
# In a new terminal, navigate to client directory
cd client

# Install dependencies
npm install

# Create .env.local file
cp .env.example .env.local

# Update .env.local with:
# NEXT_PUBLIC_API_URL=http://localhost:8080
# NEXT_PUBLIC_WS_URL=ws://localhost:8080/ws


# Run the development server
npm run dev
```

#### 4. MongoDB Setup

```bash
# If using local MongoDB
mongod --dbpath /path/to/your/data

# Or use MongoDB Atlas cloud instance
# Update connection string in backend .env file
```

### Access the Application

- Frontend: <http://localhost:3000>
- Backend API: <http://localhost:8080>
- MongoDB: mongodb://localhost:27017

## Testing Real-Time Analytics

### How to Test Real-Time Features

1. **Setup Multiple Browser Windows**:
   - Open the analytics dashboard in another window
   - Open the public form link in a third window/incognito

2. **Test Real-Time Updates**:
   - Submit a response in the public form
   - Watch the analytics dashboard update immediately

## Project Structure

```
dune-takehome/
├── client/                 # Next.js frontend
│   ├── src/
│   │   ├── app/           # App router pages
│   │   ├── components/    # React components
│   │   ├── hooks/         # Custom React hooks
│   │   └── utils/         # Utility functions
│   └── package.json
├── server/                 # Go Fiber backend
│   ├── cmd/               # Application entry point
│   ├── handlers/          # Request handlers
│   ├── models/            # Data models
│   ├── routes/            # API routes
│   └── go.mod
└── README.md
```

## 💡 Assumptions Made

- I was doing some really quick research while starting this project, because I know the live analytics and web sockets would work correctly, but I immediately thought of "okay what about scale?" I know for this project it might not seem necissary, but I used a go routine to fire off another process to process any analytics broadcasting that needs to be done. Because of that, I have mutexes in place to make sure there is no race condition between any of the connections to the analytics room. I am thinking we could make multiple form owners, in which they can view the form's analytics page like this with ease.

## 🚧 Challenges Faced

- One of the bigger challenges was how I was going to implement the websockets. I was blitzing through this because of the time frame, and to be quite frank, this project is pretty heft for a takehome in 72 hours, so I was prioritizing the important things first. Knowing Gin helped because I have never used Fiber before, but it was Very similar to the Gin framework, so I knew there was some sort of websocket connection module in it. Sure enough, it was. That was challenge 1.
- The next biggest challenge was how I was going to share a room to another person. I needed to figure out how to generate unique share URLs for forms and make them publicly accessible without authentication. The solution was creating a separate public route that bypasses auth middleware and uses the share URL as the identifier instead of the form ID.
- I ended up building custom React hooks to handle form building, validation, and real-time updates. It took some iteration to get it all working, but she is honest work!
- Getting the real-time analytics to work was another challenge. The goroutine approach helped here - form submission returns quickly while analytics calculation and broadcasting happens in the background.
- The actual analytics were more complex than I initially thought. I had to rely heavily on claude for this one.

## 📋 API Endpoints

### Forms

- `GET /api/v1/forms` - List all forms
- `POST /api/v1/forms` - Create a new form
- `GET /api/v1/forms/:id` - Get form by ID
- `PUT /api/v1/forms/:id` - Update form
- `DELETE /api/v1/forms/:id` - Delete form

### Responses

- `POST /api/v1/forms/:id/responses` - Submit form response
- `GET /api/v1/forms/:id/responses` - Get form responses

### Analytics

- `GET /api/v1/forms/:id/analytics` - Get form analytics
- `WS /api/v1/analytics/live` - WebSocket endpoint for real-time updates

## 🔐 Environment Variables

### Backend (.env)

```env
MONGODB_URI=mongodb://localhost:27017/formbuilder
PORT=8080
CLIENT_URL=http://localhost:3000
JWT_SECRET=your-secret-key
```

### Frontend (.env.local)

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080
```

## 📦 Deployment

### Railway Deployment

1. Backend service configured with root directory `/server`
2. Frontend service configured with root directory `/client`
3. MongoDB service added to project
4. Environment variables configured for production URLs

---
Built with ❤️ for Dune Security Take-Home Challenge
