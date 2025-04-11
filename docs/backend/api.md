# Hello Pulse Backend API Documentation

This document provides a comprehensive overview of the Hello Pulse backend API, data models, and architecture to guide frontend development, including practical usage examples with curl.

## API Structure and Endpoints with Usage Examples

### Authentication Endpoints

| Method | Endpoint | Description | Authentication Required |
|--------|----------|-------------|------------------------|
| POST | `/register` | Register a new user | No |
| POST | `/login` | Authenticate user and get session token | No |
| POST | `/logout` | Invalidate current session | No |
| GET | `/me` | Get current authenticated user details | Yes |

#### Authentication API Usage Examples

**Register a new user:**
```bash
curl -X POST http://localhost:5000/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "5551234567",
    "address": "123 Main St",
    "password": "securepassword"
  }'
```

**Register Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "user": {
    "id": "user-uuid",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe"
  }
}
```

**Login with email and password:**
```bash
curl -X POST http://localhost:5000/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }' \
  -c cookies.txt
```

**Login Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "user": {
    "id": "user-uuid",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe"
  }
}
```
> Note: The JWT token is stored in an HTTP-only cookie, which curl saves to the cookies.txt file.

**Get current user details:**
```bash
curl -X GET http://localhost:5000/me \
  -b cookies.txt
```

**Me Response:**
```json
{
  "success": true,
  "user": {
    "id": "user-uuid",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "phone": "5551234567",
    "address": "123 Main St",
    "role": "User",
    "organization": {
      "id": "org-uuid",
      "name": "My Organization"
    },
    "createdAt": "2023-01-01T00:00:00Z"
  }
}
```

**Logout (invalidate session):**
```bash
curl -X POST http://localhost:5000/logout \
  -b cookies.txt
```

**Logout Response:**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

### Organization Endpoints

| Method | Endpoint | Description | Authentication Required |
|--------|----------|-------------|------------------------|
| POST | `/organizations` | Create a new organization | Yes |
| POST | `/organizations/join` | Join an organization using invite code | Yes |
| GET | `/organizations/invite-codes` | Get invite codes for user's organization | Yes |
| POST | `/organizations/invite-codes` | Create a new invite code | Yes |
| DELETE | `/organizations/invite-codes` | Delete an invite code | Yes |

#### Organization API Usage Examples

**Create a new organization:**
```bash
curl -X POST http://localhost:5000/organizations \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "name": "My Organization"
  }'
```

**Create Organization Response:**
```json
{
  "success": true,
  "message": "Organization created successfully",
  "organization": {
    "id": "org-uuid",
    "name": "My Organization"
  }
}
```

**Join an organization using invite code:**
```bash
curl -X POST http://localhost:5000/organizations/join \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "code": "ABC123"
  }'
```

**Join Organization Response:**
```json
{
  "success": true,
  "message": "Joined organization successfully"
}
```

**List invite codes for organization:**
```bash
curl -X GET http://localhost:5000/organizations/invite-codes \
  -b cookies.txt
```

**List Invite Codes Response:**
```json
{
  "success": true,
  "codes": [
    {
      "id": "invite-uuid-1",
      "code": "ABC123",
      "expirationTimeMs": 1672531200000
    },
    {
      "id": "invite-uuid-2",
      "code": "XYZ789",
      "expirationTimeMs": 1672617600000
    }
  ]
}
```

**Create a new invite code:**
```bash
curl -X POST http://localhost:5000/organizations/invite-codes \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "expirationTimeMs": 1672531200000
  }'
```

**Create Invite Code Response:**
```json
{
  "success": true,
  "code": "ABC123"
}
```

**Delete an invite code:**
```bash
curl -X DELETE http://localhost:5000/organizations/invite-codes \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "id": "invite-uuid-1"
  }'
```

**Delete Invite Code Response:**
```json
{
  "success": true,
  "message": "Invite code deleted successfully"
}
```

### Project Endpoints

| Method | Endpoint | Description | Authentication Required |
|--------|----------|-------------|------------------------|
| POST | `/projects` | Create a new project | Yes |
| GET | `/projects` | Get all (root) projects for user's organization | Yes |
| GET | `/projects/:id` | Get project details by ID | Yes |
| PUT | `/projects/:id` | Update project details | Yes |
| DELETE | `/projects/:id` | Delete a project | Yes |
| POST | `/projects/add-user` | Add a user to a project | Yes |

#### Project API Usage Examples

**Create a new project:**
```bash
curl -X POST http://localhost:5000/projects \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "projectName": "New Project",
    "projectDesc": "Project description",
    "parentProjectId": null
  }'
```

**Create Project Response:**
```json
{
  "success": true,
  "message": "Project created successfully",
  "project": {
    "id": "project-uuid",
    "projectName": "New Project",
    "projectDesc": "Project description",
    "ownerId": "user-uuid",
    "organizationId": "org-uuid",
    "parentProjectId": null,
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  }
}
```

**Create a sub-project:**
```bash
curl -X POST http://localhost:5000/projects \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "projectName": "Sub Project",
    "projectDesc": "A child project",
    "parentProjectId": "parent-project-uuid"
  }'
```

**Get all root projects:**
```bash
curl -X GET http://localhost:5000/projects \
  -b cookies.txt
```

**Get All Projects Response:**
```json
{
  "success": true,
  "projects": [
    {
      "projectId": "project-uuid-1",
      "projectName": "Project 1",
      "projectDesc": "Description 1",
      "ownerId": "user-uuid",
      "organizationId": "org-uuid",
      "parentProjectId": null,
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z",
      "participants": [
        {
          "userId": "user-uuid",
          "firstName": "John",
          "lastName": "Doe",
          "email": "john@example.com"
        }
      ]
    },
    {
      "projectId": "project-uuid-2",
      "projectName": "Project 2",
      "projectDesc": "Description 2",
      "ownerId": "user-uuid",
      "organizationId": "org-uuid",
      "parentProjectId": null,
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z",
      "participants": [
        {
          "userId": "user-uuid",
          "firstName": "John",
          "lastName": "Doe",
          "email": "john@example.com"
        }
      ]
    }
  ]
}
```

**Get a specific project:**
```bash
curl -X GET http://localhost:5000/projects/project-uuid-1 \
  -b cookies.txt
```

**Get Project Response:**
```json
{
  "success": true,
  "project": {
    "projectId": "project-uuid-1",
    "projectName": "Project Name",
    "projectDesc": "Description",
    "ownerId": "user-uuid",
    "organizationId": "org-uuid",
    "parentProjectId": null,
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  },
  "childProjects": [
    {
      "projectId": "child-project-uuid",
      "projectName": "Child Project",
      "projectDesc": "Child Description",
      "ownerId": "user-uuid",
      "organizationId": "org-uuid",
      "parentProjectId": "project-uuid-1",
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z"
    }
  ],
  "participants": [
    {
      "userId": "user-uuid",
      "firstName": "John",
      "lastName": "Doe",
      "email": "john@example.com"
    }
  ]
}
```

**Update a project:**
```bash
curl -X PUT http://localhost:5000/projects/project-uuid-1 \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "projectName": "Updated Project Name",
    "projectDesc": "Updated description"
  }'
```

**Update Project Response:**
```json
{
  "success": true,
  "message": "Project updated successfully"
}
```

**Delete a project:**
```bash
curl -X DELETE http://localhost:5000/projects/project-uuid-1 \
  -b cookies.txt
```

**Delete Project Response:**
```json
{
  "success": true,
  "message": "Project deleted successfully"
}
```

**Add a user to a project:**
```bash
curl -X POST http://localhost:5000/projects/add-user \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "projectId": "project-uuid-1",
    "userId": "another-user-uuid"
  }'
```

**Add User Response:**
```json
{
  "success": true,
  "message": "Participant added successfully"
}
```

### Event Endpoints

| Method | Endpoint | Description | Authentication Required |
|--------|----------|-------------|------------------------|
| POST | `/events` | Create a new event | Yes |
| GET | `/events` | Get all events for user | Yes |
| DELETE | `/events/:id` | Delete an event | Yes |
| POST | `/events/add-member` | Add a user to an event | Yes |
| POST | `/events/remove-member` | Remove a user from an event | Yes |
| POST | `/events/:id/update-title` | Update an event's title | Yes |
| GET | `/events/:id/participants` | Get participants of an event | Yes |
| GET | `/events/fetch-users` | Get users from the organization for event creation | Yes |

#### Event API Usage Examples

**Create a new event:**
```bash
curl -X POST http://localhost:5000/events \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "title": "Team Meeting",
    "date": "2023-05-20",
    "startTime": "14:00",
    "endTime": "15:00",
    "userIds": ["user-uuid-1", "user-uuid-2"],
    "importance": "important"
  }'
```

**Create Event Response:**
```json
{
  "success": true,
  "message": "Event created successfully",
  "event": {
    "id": "event-uuid",
    "title": "Team Meeting",
    "date": "2023-05-20",
    "startTime": "14:00",
    "endTime": "15:00",
    "importance": "important"
  }
}
```

**Get all events for a user:**
```bash
curl -X GET http://localhost:5000/events \
  -b cookies.txt
```

**Get Events Response:**
```json
{
  "success": true,
  "events": [
    {
      "id": "event-uuid-1",
      "title": "Team Meeting",
      "date": "2023-05-20",
      "startTime": "14:00",
      "endTime": "15:00",
      "importance": "important",
      "createdById": "user-uuid",
      "organizationId": "org-uuid",
      "createdBy": {
        "userId": "user-uuid",
        "firstName": "John",
        "lastName": "Doe",
        "email": "john@example.com"
      },
      "users": [
        {
          "userId": "user-uuid-1",
          "firstName": "Jane",
          "lastName": "Smith",
          "email": "jane@example.com"
        }
      ]
    },
    {
      "id": "event-uuid-2",
      "title": "Project Review",
      "date": "2023-05-22",
      "startTime": "10:00",
      "endTime": "11:00",
      "importance": "very important",
      "createdById": "user-uuid",
      "organizationId": "org-uuid",
      "createdBy": {
        "userId": "user-uuid",
        "firstName": "John",
        "lastName": "Doe",
        "email": "john@example.com"
      },
      "users": []
    }
  ],
  "userId": "user-uuid"
}
```

**Delete an event:**
```bash
curl -X DELETE http://localhost:5000/events/event-uuid-1 \
  -b cookies.txt
```

**Delete Event Response:**
```json
{
  "success": true,
  "message": "Event deleted successfully"
}
```

**Add a user to an event:**
```bash
curl -X POST http://localhost:5000/events/add-member \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "eventId": "event-uuid-1",
    "userId": "user-uuid-2"
  }'
```

**Add Member Response:**
```json
{
  "success": true,
  "message": "Participant added successfully"
}
```

**Remove a user from an event:**
```bash
curl -X POST http://localhost:5000/events/remove-member \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "eventId": "event-uuid-1",
    "userId": "user-uuid-2"
  }'
```

**Remove Member Response:**
```json
{
  "success": true,
  "message": "Participant removed successfully"
}
```

**Update an event title:**
```bash
curl -X POST http://localhost:5000/events/event-uuid-1/update-title \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "eventId": "event-uuid-1",
    "title": "Updated Team Meeting"
  }'
```

**Update Title Response:**
```json
{
  "success": true,
  "message": "Event title updated successfully"
}
```

**Get event participants:**
```bash
curl -X GET http://localhost:5000/events/event-uuid-1/participants \
  -b cookies.txt
```

**Get Participants Response:**
```json
{
  "success": true,
  "participants": [
    {
      "id": "user-uuid-1",
      "firstName": "Jane",
      "lastName": "Smith",
      "email": "jane@example.com"
    },
    {
      "id": "user-uuid-2",
      "firstName": "Bob",
      "lastName": "Johnson",
      "email": "bob@example.com"
    }
  ]
}
```

**Fetch organization users for event creation:**
```bash
curl -X GET http://localhost:5000/events/fetch-users \
  -b cookies.txt
```

**Fetch Users Response:**
```json
{
  "success": true,
  "users": [
    {
      "userId": "user-uuid-1",
      "firstName": "Jane",
      "lastName": "Smith",
      "email": "jane@example.com"
    },
    {
      "userId": "user-uuid-2",
      "firstName": "Bob",
      "lastName": "Johnson",
      "email": "bob@example.com"
    }
  ]
}
```

## Data Models and Relationships

### User Model

```typescript
interface User {
  id: string;            // UUID
  firstName: string;
  lastName: string;
  email: string;         // Unique
  phone: string;
  address: string;
  organizationId: string | null;
  role: string;          // "Admin" or "User"
  lastActive: string;    // ISO timestamp
  createdAt: string;     // ISO timestamp
  updatedAt: string;     // ISO timestamp
}
```

### Organization Model

```typescript
interface Organization {
  id: string;            // UUID
  name: string;          // Unique
  ownerId: string;       // UUID of user who created
  createdAt: string;     // ISO timestamp
  updatedAt: string;     // ISO timestamp
}
```

### Project Model

```typescript
interface Project {
  id: string;            // UUID
  name: string;
  description: string;
  ownerId: string;       // UUID of user who created
  organizationId: string; // UUID of organization
  parentProjectId: string | null; // UUID of parent project (hierarchical)
  createdAt: string;     // ISO timestamp
  updatedAt: string;     // ISO timestamp
  
  // Returned in responses but not in model
  participants?: User[];
  childProjects?: Project[];
}
```

### Event Model

```typescript
interface Event {
  id: string;            // UUID
  title: string;
  date: string;          // ISO date
  startTime: string;     // ISO time
  endTime: string;       // ISO time
  organizationId: string; // UUID of organization
  createdById: string;   // UUID of creator user
  importance: string;    // "not important", "important", "very important"
  createdAt: string;     // ISO timestamp
  updatedAt: string;     // ISO timestamp
  
  // Returned in responses but not in model
  participants?: User[];
  createdBy?: User;
}
```

### Session Model

```typescript
interface Session {
  id: string;            // UUID
  userId: string;        // UUID of related user
  token: string;         // Unique session token
  expiresAt: string;     // ISO timestamp
  createdAt: string;     // ISO timestamp
}
```

### File Model

```typescript
interface File {
  id: string;            // UUID
  fileName: string;
  bucketName: string;
  objectName: string;
  contentType: string;
  uploadedAt: string;    // ISO timestamp
  isDeleted: boolean;
  deletedAt: string | null; // ISO timestamp
  uploaderId: string;    // UUID of user who uploaded
  organizationId: string; // UUID of organization
  isPublic: boolean;
}
```

### InviteCode Model

```typescript
interface InviteCode {
  id: string;            // UUID
  value: string;         // The actual code
  organizationId: string; // UUID of organization
  expirationTime: string; // ISO timestamp
  createdAt: string;     // ISO timestamp
  updatedAt: string;     // ISO timestamp
}
```

### Key Relationships

- **User** belongs to an **Organization** (optional)
- **Organization** has many **Users**
- **Project** belongs to an **Organization**
- **Project** can have a parent **Project** (hierarchical structure)
- **Project** has many **Users** as participants (many-to-many)
- **Event** belongs to an **Organization**
- **Event** is created by a **User**
- **Event** has many **Users** as participants (many-to-many)
- **File** is uploaded by a **User**
- **File** belongs to an **Organization**
- **InviteCode** belongs to an **Organization**
- **Session** belongs to a **User**

## Application Architecture

### Architecture Overview

Hello Pulse follows a clean architecture approach with domain-driven design principles:

```
┌─────────────────────────────────────────┐
│                                         │
│                API Layer                │
│  (Handlers, Routes, Request/Response)   │
│                                         │
├─────────────────────────────────────────┤
│                                         │
│               Service Layer             │
│          (Business Logic)               │
│                                         │
├─────────────────────────────────────────┤
│                                         │
│             Repository Layer            │
│             (Data Access)               │
│                                         │
├─────────────────────────────────────────┤
│                                         │
│               Domain Layer              │
│               (Models)                  │
│                                         │
└─────────────────────────────────────────┘
```

### Layer Structure

1. **Domain Models Layer**
   - Contains business entities (User, Organization, Project, Event, etc.)
   - Defines core business rules and validation

2. **Repository Layer**
   - Abstracts data access operations
   - Provides CRUD operations for domain models
   - Handles database interactions

3. **Service Layer**
   - Contains business logic
   - Orchestrates operations across multiple repositories
   - Implements domain-specific workflows
   - Handles authorization rules

4. **API Handler Layer**
   - Manages HTTP requests/responses
   - Handles request validation
   - Translates between HTTP and domain objects
   - Returns appropriate status codes and responses

### Authentication Flow

1. User submits credentials (email/password)
2. Server authenticates and creates a session with JWT token
3. Token is sent to client as an HTTP-only cookie
4. Client includes cookie in subsequent requests
5. Auth middleware validates token and attaches user to request context
6. Protected routes check for authenticated user in context

### Infrastructure Components

- **Web Server**: Go with Gin framework
- **Database**: PostgreSQL with pgvector extension
- **File Storage**: MinIO (S3-compatible object storage)
- **API Gateway**: Nginx as reverse proxy
- **Development**: Docker and Docker Compose

## Error Handling

### Common Error Responses

```json
{
  "success": false,
  "error": "Error message"
}
```

### Common Error Messages

- Authentication: "Unauthorized", "Invalid credentials"
- Organization: "User already belongs to an organization", "Invalid invite code"
- Projects: "Project not found", "Access denied", "Only the project owner can update the project"
- Events: "Event not found", "Only the event creator can update the event"

### HTTP Status Codes

- 200: Success
- 400: Bad Request (invalid input)
- 401: Unauthorized (authentication required)
- 403: Forbidden (insufficient permissions)
- 404: Not Found (resource doesn't exist)
- 500: Internal Server Error

## Testing with curl

For convenience, here's a complete workflow for testing the API with curl:

### Setup Environment Variables

```bash
# Base URL
BASE_URL=http://localhost:5000

# Create a file to store cookies
touch cookies.txt
```

### User Registration and Authentication

```bash
# Register a new user
curl -X POST $BASE_URL/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "5551234567",
    "address": "123 Main St",
    "password": "securepassword"
  }'

# Login to get authentication cookie
curl -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }' \
  -c cookies.txt

# Verify current user
curl -X GET $BASE_URL/me \
  -b cookies.txt
```

### Organization Management

```bash
# Create an organization
curl -X POST $BASE_URL/organizations \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "name": "My Organization"
  }'

# Create an invite code
curl -X POST $BASE_URL/organizations/invite-codes \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "expirationTimeMs": 1672531200000
  }'

# View invite codes
curl -X GET $BASE_URL/organizations/invite-codes \
  -b cookies.txt
```

### Project Management

```bash
# Create a project
curl -X POST $BASE_URL/projects \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "projectName": "New Project",
    "projectDesc": "Project description",
    "parentProjectId": null
  }'

# List all projects
curl -X GET $BASE_URL/projects \
  -b cookies.txt

# Add a user to a project (replace with actual IDs)
curl -X POST $BASE_URL/projects/add-user \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "projectId": "project-uuid",
    "userId": "another-user-uuid"
  }'
```

### Event Management

```bash
# Create an event
curl -X POST $BASE_URL/events \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "title": "Team Meeting",
    "date": "2023-05-20",
    "startTime": "14:00",
    "endTime": "15:00",
    "userIds": [],
    "importance": "important"
  }'

# List all events
curl -X GET $BASE_URL/events \
  -b cookies.txt
```

### Logout

```bash
# Logout
curl -X POST $BASE_URL/logout \
  -b cookies.txt
```