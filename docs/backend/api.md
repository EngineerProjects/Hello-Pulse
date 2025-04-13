# Hello Pulse API Documentation

This document provides comprehensive documentation for all APIs in the Hello Pulse platform. It includes details on authentication, endpoints, request/response formats, and practical examples using cURL.

## Table of Contents

1. [Authentication](#authentication)
2. [User Management](#user-management)
3. [Organization Management](#organization-management)
4. [Project Management](#project-management)
5. [Project Summaries](#project-summaries)
6. [Event Management](#event-management)
7. [File Management](#file-management)

---

## Authentication

All protected endpoints require a valid JWT token which is sent as a cookie. Authentication is handled through the following endpoints.

### Register a New User

Creates a new user account and returns a session token.

- **URL**: `/register`
- **Method**: `POST`
- **Auth Required**: No

**Request Body:**
```json
{
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+33612345678",
  "address": "123 Main St, Paris, France",
  "password": "securePassword123"
}
```

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "User registered successfully",
  "user": {
    "id": "a1b2c3d4-e5f6-7890-abcd-1234567890ab",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe"
  }
}
```

**Error Response:**
- **Code**: 400 Bad Request
- **Content**:
```json
{
  "success": false,
  "error": "Email already exists"
}
```

**Example:**
```bash
curl -X POST http://localhost:5000/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+33612345678",
    "address": "123 Main St, Paris, France",
    "password": "securePassword123"
  }'
```

### Login

Authenticates a user and returns a session token.

- **URL**: `/login`
- **Method**: `POST`
- **Auth Required**: No

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "Login successful",
  "user": {
    "id": "a1b2c3d4-e5f6-7890-abcd-1234567890ab",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe"
  }
}
```

**Error Response:**
- **Code**: 401 Unauthorized
- **Content**:
```json
{
  "success": false,
  "error": "Invalid credentials"
}
```

**Example:**
```bash
curl -X POST http://localhost:5000/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securePassword123"
  }'
```

### Logout

Invalidates the current session token.

- **URL**: `/logout`
- **Method**: `POST`
- **Auth Required**: Yes (Cookie)

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

**Example:**
```bash
curl -X POST http://localhost:5000/logout \
  -b "token=your-session-token"
```

---

## User Management

### Get Current User

Retrieves information about the currently authenticated user.

- **URL**: `/me`
- **Method**: `GET`
- **Auth Required**: Yes (Cookie)

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "user": {
    "id": "a1b2c3d4-e5f6-7890-abcd-1234567890ab",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "phone": "+33612345678",
    "address": "123 Main St, Paris, France",
    "organizationId": "o1p2q3r4-s5t6-7890-uvwx-1234567890ab",
    "role": "Admin",
    "lastActive": "2023-09-15T14:30:45Z",
    "createdAt": "2023-01-15T10:20:30Z",
    "updatedAt": "2023-09-15T14:30:45Z"
  }
}
```

**Error Response:**
- **Code**: 401 Unauthorized
- **Content**:
```json
{
  "success": false,
  "error": "Unauthorized"
}
```

**Example:**
```bash
curl -X GET http://localhost:5000/me \
  -b "token=your-session-token"
```

---

## Organization Management

### Create Organization

Creates a new organization and assigns the current user as its owner.

- **URL**: `/organizations`
- **Method**: `POST`
- **Auth Required**: Yes (Cookie)

**Request Body:**
```json
{
  "name": "Acme Corp"
}
```

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "Organization created successfully",
  "organization": {
    "id": "o1p2q3r4-s5t6-7890-uvwx-1234567890ab",
    "name": "Acme Corp"
  }
}
```

**Error Response:**
- **Code**: 400 Bad Request
- **Content**:
```json
{
  "success": false,
  "error": "User already belongs to an organization"
}
```

**Example:**
```bash
curl -X POST http://localhost:5000/organizations \
  -b "token=your-session-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Acme Corp"
  }'
```

### Join Organization

Adds the current user to an organization using an invite code.

- **URL**: `/organizations/join`
- **Method**: `POST`
- **Auth Required**: Yes (Cookie)

**Request Body:**
```json
{
  "code": "ABCDEF"
}
```

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "Joined organization successfully"
}
```

**Error Response:**
- **Code**: 400 Bad Request
- **Content**:
```json
{
  "success": false,
  "error": "Invalid invite code"
}
```

**Example:**
```bash
curl -X POST http://localhost:5000/organizations/join \
  -b "token=your-session-token" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "ABCDEF"
  }'
```

### Create Invite Code

Creates a new invite code for an organization.

- **URL**: `/organizations/invite-codes`
- **Method**: `POST`
- **Auth Required**: Yes (Cookie)

**Request Body:**
```json
{
  "expirationTimeMs": 86400000
}
```

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "code": "ABCDEF"
}
```

**Error Response:**
- **Code**: 403 Forbidden
- **Content**:
```json
{
  "success": false,
  "error": "Only administrators can create invite codes"
}
```

**Example:**
```bash
curl -X POST http://localhost:5000/organizations/invite-codes \
  -b "token=your-session-token" \
  -H "Content-Type: application/json" \
  -d '{
    "expirationTimeMs": 86400000
  }'
```

### Get Invite Codes

Retrieves all invite codes for the current user's organization.

- **URL**: `/organizations/invite-codes`
- **Method**: `GET`
- **Auth Required**: Yes (Cookie)

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "codes": [
    {
      "id": "i1j2k3l4-m5n6-7890-opqr-1234567890ab",
      "code": "ABCDEF",
      "expirationTimeMs": 1684857600000
    },
    {
      "id": "s1t2u3v4-w5x6-7890-yzab-1234567890ab",
      "code": "GHIJKL",
      "expirationTimeMs": 1684944000000
    }
  ]
}
```

**Error Response:**
- **Code**: 403 Forbidden
- **Content**:
```json
{
  "success": false,
  "error": "Only administrators can view invite codes"
}
```

**Example:**
```bash
curl -X GET http://localhost:5000/organizations/invite-codes \
  -b "token=your-session-token"
```

### Delete Invite Code

Deletes an invite code.

- **URL**: `/organizations/invite-codes`
- **Method**: `DELETE`
- **Auth Required**: Yes (Cookie)

**Request Body:**
```json
{
  "id": "i1j2k3l4-m5n6-7890-opqr-1234567890ab"
}
```

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "Invite code deleted successfully"
}
```

**Error Response:**
- **Code**: 403 Forbidden
- **Content**:
```json
{
  "success": false,
  "error": "Only administrators can delete invite codes"
}
```

**Example:**
```bash
curl -X DELETE http://localhost:5000/organizations/invite-codes \
  -b "token=your-session-token" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "i1j2k3l4-m5n6-7890-opqr-1234567890ab"
  }'
```

---

## Project Management

### Create Project

Creates a new project.

- **URL**: `/projects`
- **Method**: `POST`
- **Auth Required**: Yes (Cookie)

**Request Body:**
```json
{
  "projectName": "New Website",
  "projectDesc": "Company website redesign",
  "parentProjectId": null
}
```

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "Project created successfully",
  "project": {
    "id": "p1q2r3s4-t5u6-7890-vwxy-1234567890ab",
    "name": "New Website",
    "description": "Company website redesign",
    "ownerId": "a1b2c3d4-e5f6-7890-abcd-1234567890ab",
    "organizationId": "o1p2q3r4-s5t6-7890-uvwx-1234567890ab",
    "parentProjectId": null,
    "createdAt": "2023-09-15T15:30:45Z",
    "updatedAt": "2023-09-15T15:30:45Z"
  }
}
```

**Error Response:**
- **Code**: 400 Bad Request
- **Content**:
```json
{
  "success": false,
  "error": "User does not belong to an organization"
}
```

**Example:**
```bash
curl -X POST http://localhost:5000/projects \
  -b "token=your-session-token" \
  -H "Content-Type: application/json" \
  -d '{
    "projectName": "New Website",
    "projectDesc": "Company website redesign",
    "parentProjectId": null
  }'
```

### Get All Projects

Retrieves all root projects for the current user's organization.

- **URL**: `/projects`
- **Method**: `GET`
- **Auth Required**: Yes (Cookie)

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "projects": [
    {
      "id": "p1q2r3s4-t5u6-7890-vwxy-1234567890ab",
      "name": "New Website",
      "description": "Company website redesign",
      "ownerId": "a1b2c3d4-e5f6-7890-abcd-1234567890ab",
      "organizationId": "o1p2q3r4-s5t6-7890-uvwx-1234567890ab",
      "parentProjectId": null,
      "createdAt": "2023-09-15T15:30:45Z",
      "updatedAt": "2023-09-15T15:30:45Z"
    },
    {
      "id": "w1x2y3z4-a5b6-7890-cdef-1234567890ab",
      "name": "Mobile App",
      "description": "Company mobile application",
      "ownerId": "a1b2c3d4-e5f6-7890-abcd-1234567890ab",
      "organizationId": "o1p2q3r4-s5t6-7890-uvwx-1234567890ab",
      "parentProjectId": null,
      "createdAt": "2023-09-10T11:20:35Z",
      "updatedAt": "2023-09-14T09:15:25Z"
    }
  ]
}
```

**Error Response:**
- **Code**: 400 Bad Request
- **Content**:
```json
{
  "success": false,
  "error": "User does not belong to an organization"
}
```

**Example:**
```bash
curl -X GET http://localhost:5000/projects \
  -b "token=your-session-token"
```

### Get Project Details

Retrieves detailed information about a specific project.

- **URL**: `/projects/:id`
- **Method**: `GET`
- **URL Params**: `id=[uuid]` - Project ID
- **Auth Required**: Yes (Cookie)

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "project": {
    "id": "p1q2r3s4-t5u6-7890-vwxy-1234567890ab",
    "name": "New Website",
    "description": "Company website redesign",
    "ownerId": "a1b2c3d4-e5f6-7890-abcd-1234567890ab",
    "organizationId": "o1p2q3r4-s5t6-7890-uvwx-1234567890ab",
    "parentProjectId": null,
    "createdAt": "2023-09-15T15:30:45Z",
    "updatedAt": "2023-09-15T15:30:45Z"
  },
  "childProjects": [
    {
      "id": "c1d2e3f4-g5h6-7890-ijkl-1234567890ab",
      "name": "Frontend Development",
      "description": "Website frontend development",
      "ownerId": "a1b2c3d4-e5f6-7890-abcd-1234567890ab",
      "organizationId": "o1p2q3r4-s5t6-7890-uvwx-1234567890ab",
      "parentProjectId": "p1q2r3s4-t5u6-7890-vwxy-1234567890ab",
      "createdAt": "2023-09-16T10:00:00Z",
      "updatedAt": "2023-09-16T10:00:00Z"
    }
  ],
  "participants": [
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-1234567890ab",
      "firstName": "John",
      "lastName": "Doe",
      "email": "john.doe@example.com"
    },
    {
      "id": "e1f2g3h4-i5j6-7890-klmn-1234567890ab",
      "firstName": "Jane",
      "lastName": "Smith",
      "email": "jane.smith@example.com"
    }
  ]
}
```

**Error Response:**
- **Code**: 404 Not Found
- **Content**:
```json
{
  "success": false,
  "error": "Project not found"
}
```

**Example:**
```bash
curl -X GET http://localhost:5000/projects/p1q2r3s4-t5u6-7890-vwxy-1234567890ab \
  -b "token=your-session-token"
```

### Update Project

Updates a project's details.

- **URL**: `/projects/:id`
- **Method**: `PUT`
- **URL Params**: `id=[uuid]` - Project ID
- **Auth Required**: Yes (Cookie)

**Request Body:**
```json
{
  "projectName": "Updated Website",
  "projectDesc": "Company website redesign with new branding"
}
```

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "Project updated successfully"
}
```

**Error Response:**
- **Code**: 403 Forbidden
- **Content**:
```json
{
  "success": false,
  "error": "Only the project owner can update the project"
}
```

**Example:**
```bash
curl -X PUT http://localhost:5000/projects/p1q2r3s4-t5u6-7890-vwxy-1234567890ab \
  -b "token=your-session-token" \
  -H "Content-Type: application/json" \
  -d '{
    "projectName": "Updated Website",
    "projectDesc": "Company website redesign with new branding"
  }'
```

### Delete Project

Deletes a project and all its child projects.

- **URL**: `/projects/:id`
- **Method**: `DELETE`
- **URL Params**: `id=[uuid]` - Project ID
- **Auth Required**: Yes (Cookie)

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "Project deleted successfully"
}
```

**Error Response:**
- **Code**: 403 Forbidden
- **Content**:
```json
{
  "success": false,
  "error": "Only the project owner can delete the project"
}
```

**Example:**
```bash
curl -X DELETE http://localhost:5000/projects/p1q2r3s4-t5u6-7890-vwxy-1234567890ab \
  -b "token=your-session-token"
```

### Add Project Participant

Adds a user as a participant to a project.

- **URL**: `/projects/add-user`
- **Method**: `POST`
- **Auth Required**: Yes (Cookie)

**Request Body:**
```json
{
  "projectId": "p1q2r3s4-t5u6-7890-vwxy-1234567890ab",
  "userId": "e1f2g3h4-i5j6-7890-klmn-1234567890ab"
}
```

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "Participant added successfully"
}
```

**Error Response:**
- **Code**: 403 Forbidden
- **Content**:
```json
{
  "success": false,
  "error": "Only the project owner can add participants"
}
```

**Example:**
```bash
curl -X POST http://localhost:5000/projects/add-user \
  -b "token=your-session-token" \
  -H "Content-Type: application/json" \
  -d '{
    "projectId": "p1q2r3s4-t5u6-7890-vwxy-1234567890ab",
    "userId": "e1f2g3h4-i5j6-7890-klmn-1234567890ab"
  }'
```

---

## Project Summaries

### Create Summary

Creates a new summary for a project.

- **URL**: `/projects/summaries`
- **Method**: `POST`
- **Auth Required**: Yes (Cookie)

**Request Body:**
```json
{
  "projectId": "p1q2r3s4-t5u6-7890-vwxy-1234567890ab",
  "title": "Initial Design Review",
  "content": "The initial design review went well. Team decided on using a minimalist approach with the company's brand colors."
}
```

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "Summary created successfully",
  "summary": {
    "id": "s1t2u3v4-w5x6-7890-yzab-1234567890ab",
    "title": "Initial Design Review",
    "content": "The initial design review went well. Team decided on using a minimalist approach with the company's brand colors.",
    "projectId": "p1q2r3s4-t5u6-7890-vwxy-1234567890ab",
    "createdBy": "a1b2c3d4-e5f6-7890-abcd-1234567890ab",
    "createdAt": "2023-09-17T14:00:00Z",
    "updatedAt": "2023-09-17T14:00:00Z"
  }
}
```

**Error Response:**
- **Code**: 400 Bad Request
- **Content**:
```json
{
  "success": false,
  "error": "Project not found"
}
```

**Example:**
```bash
curl -X POST http://localhost:5000/projects/summaries \
  -b "token=your-session-token" \
  -H "Content-Type: application/json" \
  -d '{
    "projectId": "p1q2r3s4-t5u6-7890-vwxy-1234567890ab",
    "title": "Initial Design Review",
    "content": "The initial design review went well. Team decided on using a minimalist approach with the company'\''s brand colors."
  }'
```

### Get Project Summaries

Retrieves all summaries for a specific project.

- **URL**: `/projects/:id/summaries`
- **Method**: `GET`
- **URL Params**: `id=[uuid]` - Project ID
- **Auth Required**: Yes (Cookie)

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "summaries": [
    {
      "id": "s1t2u3v4-w5x6-7890-yzab-1234567890ab",
      "title": "Initial Design Review",
      "content": "The initial design review went well. Team decided on using a minimalist approach with the company's brand colors.",
      "projectId": "p1q2r3s4-t5u6-7890-vwxy-1234567890ab",
      "createdBy": "a1b2c3d4-e5f6-7890-abcd-1234567890ab",
      "createdAt": "2023-09-17T14:00:00Z",
      "updatedAt": "2023-09-17T14:00:00Z"
    },
    {
      "id": "c1d2e3f4-g5h6-7890-ijkl-9876543210zy",
      "title": "Development Progress",
      "content": "Frontend development is 50% complete. Backend API endpoints have been defined.",
      "projectId": "p1q2r3s4-t5u6-7890-vwxy-1234567890ab",
      "createdBy": "a1b2c3d4-e5f6-7890-abcd-1234567890ab",
      "createdAt": "2023-09-20T11:30:00Z",
      "updatedAt": "2023-09-20T11:30:00Z"
    }
  ]
}
```

**Error Response:**
- **Code**: 404 Not Found
- **Content**:
```json
{
  "success": false,
  "error": "Project not found"
}
```

**Example:**
```bash
curl -X GET http://localhost:5000/projects/p1q2r3s4-t5u6-7890-vwxy-1234567890ab/summaries \
  -b "token=your-session-token"
```

### Get Summary Details

Retrieves details of a specific summary.

- **URL**: `/projects/summaries/:id`
- **Method**: `GET`
- **URL Params**: `id=[uuid]` - Summary ID
- **Auth Required**: Yes (Cookie)

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "summary": {
    "id": "s1t2u3v4-w5x6-7890-yzab-1234567890ab",
    "title": "Initial Design Review",
    "content": "The initial design review went well. Team decided on using a minimalist approach with the company's brand colors.",
    "projectId": "p1q2r3s4-t5u6-7890-vwxy-1234567890ab",
    "createdBy": "a1b2c3d4-e5f6-7890-abcd-1234567890ab",
    "createdAt": "2023-09-17T14:00:00Z",
    "updatedAt": "2023-09-17T14:00:00Z"
  }
}
```

**Error Response:**
- **Code**: 404 Not Found
- **Content**:
```json
{
  "success": false,
  "error": "Summary not found"
}
```

**Example:**
```bash
curl -X GET http://localhost:5000/projects/summaries/s1t2u3v4-w5x6-7890-yzab-1234567890ab \
  -b "token=your-session-token"
```

### Update Summary

Updates a summary's content.

- **URL**: `/projects/summaries/:id`
- **Method**: `PUT`
- **URL Params**: `id=[uuid]` - Summary ID
- **Auth Required**: Yes (Cookie)

**Request Body:**
```json
{
  "title": "Updated Design Review",
  "content": "The design review went well. Team agreed on a minimalist approach with the company's updated brand colors and typography."
}
```

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "Summary updated successfully"
}
```

**Error Response:**
- **Code**: 400 Bad Request
- **Content**:
```json
{
  "success": false,
  "error": "Only the creator can update this summary"
}
```

**Example:**
```bash
curl -X PUT http://localhost:5000/projects/summaries/s1t2u3v4-w5x6-7890-yzab-1234567890ab \
  -b "token=your-session-token" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Design Review",
    "content": "The design review went well. Team agreed on a minimalist approach with the company'\''s updated brand colors and typography."
  }'
```

### Delete Summary

Deletes a summary.

- **URL**: `/projects/summaries/:id`
- **Method**: `DELETE`
- **URL Params**: `id=[uuid]` - Summary ID
- **Auth Required**: Yes (Cookie)

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "Summary deleted successfully"
}
```

**Error Response:**
- **Code**: 400 Bad Request
- **Content**:
```json
{
  "success": false,
  "error": "Only the creator can delete this summary"
}
```

**Example:**
```bash
curl -X DELETE http://localhost:5000/projects/summaries/s1t2u3v4-w5x6-7890-yzab-1234567890ab \
  -b "token=your-session-token"
```

---

## Event Management

### Create Event

Creates a new event.

- **URL**: `/events`
- **Method**: `POST`
- **Auth Required**: Yes (Cookie)

**Request Body:**
```json
{
  "title": "Project Kickoff Meeting",
  "date": "2023-10-15",
  "startTime": "09:00",
  "endTime": "10:30",
  "userIds": ["e1f2g3h4-i5j6-7890-klmn-1234567890ab"],
  "importance": "high"
}
```

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "Event created successfully",
  "event": {
    "id": "e1v2w3x4-y5z6-7890-abcd-9876543210zy",
    "title": "Project Kickoff Meeting",
    "date": "2023-10-15",
    "startTime": "09:00",
    "endTime": "10:30",
    "importance": "high"
  }
}
```

**Error Response:**
- **Code**: 400 Bad Request
- **Content**:
```json
{
  "success": false,
  "error": "User does not belong to an organization"
}
```

**Example:**
```bash
curl -X POST http://localhost:5000/events \
  -b "token=your-session-token" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Project Kickoff Meeting",
    "date": "2023-10-15",
    "startTime": "09:00",
    "endTime": "10:30",
    "userIds": ["e1f2g3h4-i5j6-7890-klmn-1234567890ab"],
    "importance": "high"
  }'
```

### Get Events

Retrieves all events for the current user.

- **URL**: `/events`
- **Method**: `GET`
- **Auth Required**: Yes (Cookie)

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "events": [
    {
      "id": "e1v2w3x4-y5z6-7890-abcd-9876543210zy",
      "title": "Project Kickoff Meeting",
      "date": "2023-10-15",
      "startTime": "09:00",
      "endTime": "10:30",
      "importance": "high",
      "createdById": "a1b2c3d4-e5f6-7890-abcd-1234567890ab",
      "organizationId": "o1p2q3r4-s5t6-7890-uvwx-1234567890ab"
    },
    {
      "id": "f1g2h3i4-j5k6-7890-lmno-9876543210zy",
      "title": "Weekly Review",
      "date": "2023-10-20",
      "startTime": "14:00",
      "endTime": "15:00",
      "importance": "medium",
      "createdById": "a1b2c3d4-e5f6-7890-abcd-1234567890ab",
      "organizationId": "o1p2q3r4-s5t6-7890-uvwx-1234567890ab"
    }
  ],
  "userId": "a1b2c3d4-e5f6-7890-abcd-1234567890ab"
}
```

**Error Response:**
- **Code**: 401 Unauthorized
- **Content**:
```json
{
  "success": false,
  "error": "Unauthorized"
}
```

**Example:**
```bash
curl -X GET http://localhost:5000/events \
  -b "token=your-session-token"
```

### Delete Event

Deletes an event.

- **URL**: `/events/:id`
- **Method**: `DELETE`
- **URL Params**: `id=[uuid]` - Event ID
- **Auth Required**: Yes (Cookie)

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "Event deleted successfully"
}
```

**Error Response:**
- **Code**: 403 Forbidden
- **Content**:
```json
{
  "success": false,
  "error": "Only the event creator can delete the event"
}
```

**Example:**
```bash
curl -X DELETE http://localhost:5000/events/e1v2w3x4-y5z6-7890-abcd-9876543210zy \
  -b "token=your-session-token"
```

### Add Event Participant

Adds a user to an event.

- **URL**: `/events/add-member`
- **Method**: `POST`
- **Auth Required**: Yes (Cookie)

**Request Body:**
```json
{
  "eventId": "e1v2w3x4-y5z6-7890-abcd-9876543210zy",
  "userId": "g1h2i3j4-k5l6-7890-mnop-1234567890ab"
}
```

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "Participant added successfully"
}
```

**Error Response:**
- **Code**: 403 Forbidden
- **Content**:
```json
{
  "success": false,
  "error": "Only the event creator can add participants"
}
```

**Example:**
```bash
curl -X POST http://localhost:5000/events/add-member \
  -b "token=your-session-token" \
  -H "Content-Type: application/json" \
  -d '{
    "eventId": "e1v2w3x4-y5z6-7890-abcd-9876543210zy",
    "userId": "g1h2i3j4-k5l6-7890-mnop-1234567890ab"
  }'
```

### Remove Event Participant

Removes a user from an event.

- **URL**: `/events/remove-member`
- **Method**: `POST`
- **Auth Required**: Yes (Cookie)

**Request Body:**
```json
{
  "eventId": "e1v2w3x4-y5z6-7890-abcd-9876543210zy",
  "userId": "g1h2i3j4-k5l6-7890-mnop-1234567890ab"
}
```

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "Participant removed successfully"
}
```

**Error Response:**
- **Code**: 403 Forbidden
- **Content**:
```json
{
  "success": false,
  "error": "Only the event creator can remove participants"
}
```

**Example:**
```bash
curl -X POST http://localhost:5000/events/remove-member \
  -b "token=your-session-token" \
  -H "Content-Type: application/json" \
  -d '{
    "eventId": "e1v2w3x4-y5z6-7890-abcd-9876543210zy",
    "userId": "g1h2i3j4-k5l6-7890-mnop-1234567890ab"
  }'
```

### Update Event Title

Updates an event's title.

- **URL**: `/events/:id/update-title`
- **Method**: `POST`
- **URL Params**: `id=[uuid]` - Event ID
- **Auth Required**: Yes (Cookie)

**Request Body:**
```json
{
  "eventId": "e1v2w3x4-y5z6-7890-abcd-9876543210zy",
  "title": "Updated Project Kickoff Meeting"
}
```

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "Event title updated successfully"
}
```

**Error Response:**
- **Code**: 403 Forbidden
- **Content**:
```json
{
  "success": false,
  "error": "Only the event creator can update the event"
}
```

**Example:**
```bash
curl -X POST http://localhost:5000/events/e1v2w3x4-y5z6-7890-abcd-9876543210zy/update-title \
  -b "token=your-session-token" \
  -H "Content-Type: application/json" \
  -d '{
    "eventId": "e1v2w3x4-y5z6-7890-abcd-9876543210zy",
    "title": "Updated Project Kickoff Meeting"
  }'
```

### Get Event Participants

Retrieves all participants of an event.

- **URL**: `/events/:id/participants`
- **Method**: `GET`
- **URL Params**: `id=[uuid]` - Event ID
- **Auth Required**: Yes (Cookie)

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "participants": [
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-1234567890ab",
      "firstName": "John",
      "lastName": "Doe",
      "email": "john.doe@example.com"
    },
    {
      "id": "g1h2i3j4-k5l6-7890-mnop-1234567890ab",
      "firstName": "Alice",
      "lastName": "Johnson",
      "email": "alice.johnson@example.com"
    }
  ]
}
```

**Error Response:**
- **Code**: 404 Not Found
- **Content**:
```json
{
  "success": false,
  "error": "Event not found"
}
```

**Example:**
```bash
curl -X GET http://localhost:5000/events/e1v2w3x4-y5z6-7890-abcd-9876543210zy/participants \
  -b "token=your-session-token"
```

---

## File Management

### Upload File

Uploads a file to the storage system.

- **URL**: `/files`
- **Method**: `POST`
- **Auth Required**: Yes (Cookie)
- **Content-Type**: `multipart/form-data`

**Form Parameters:**
- `file`: The file to upload
- `isPublic`: "true" if the file should be publicly accessible, "false" otherwise

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "File uploaded successfully",
  "file": {
    "id": "f1i2l3e4-5678-90ab-cdef-1234567890ab",
    "fileName": "presentation.pdf",
    "contentType": "application/pdf",
    "size": 2048576,
    "uploadedAt": "2023-09-22T16:45:30Z",
    "isPublic": true
  }
}
```

**Error Response:**
- **Code**: 400 Bad Request
- **Content**:
```json
{
  "success": false,
  "error": "File size exceeds the 100 MB limit"
}
```

**Example:**
```bash
curl -X POST http://localhost:5000/files \
  -b "token=your-session-token" \
  -F "file=@/path/to/presentation.pdf" \
  -F "isPublic=true"
```

### Get User Files

Retrieves files uploaded by the current user.

- **URL**: `/files`
- **Method**: `GET`
- **Query Params**: `includeDeleted=[boolean]` - Whether to include soft-deleted files
- **Auth Required**: Yes (Cookie)

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "files": [
    {
      "id": "f1i2l3e4-5678-90ab-cdef-1234567890ab",
      "fileName": "presentation.pdf",
      "contentType": "application/pdf",
      "size": 2048576,
      "uploadedAt": "2023-09-22T16:45:30Z",
      "isPublic": true,
      "isDeleted": false
    },
    {
      "id": "a1b2c3d4-5678-90ab-cdef-9876543210zy",
      "fileName": "document.docx",
      "contentType": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
      "size": 1048576,
      "uploadedAt": "2023-09-20T14:30:00Z",
      "isPublic": false,
      "isDeleted": false
    }
  ]
}
```

**Error Response:**
- **Code**: 401 Unauthorized
- **Content**:
```json
{
  "success": false,
  "error": "Unauthorized"
}
```

**Example:**
```bash
curl -X GET "http://localhost:5000/files?includeDeleted=false" \
  -b "token=your-session-token"
```

### Get Organization Files

Retrieves files for the current user's organization.

- **URL**: `/files/organization`
- **Method**: `GET`
- **Query Params**: `includeDeleted=[boolean]` - Whether to include soft-deleted files
- **Auth Required**: Yes (Cookie)

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "files": [
    {
      "id": "f1i2l3e4-5678-90ab-cdef-1234567890ab",
      "fileName": "presentation.pdf",
      "contentType": "application/pdf",
      "size": 2048576,
      "uploadedAt": "2023-09-22T16:45:30Z",
      "isPublic": true,
      "isDeleted": false,
      "uploaderId": "a1b2c3d4-e5f6-7890-abcd-1234567890ab"
    },
    {
      "id": "a1b2c3d4-5678-90ab-cdef-9876543210zy",
      "fileName": "document.docx",
      "contentType": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
      "size": 1048576,
      "uploadedAt": "2023-09-20T14:30:00Z",
      "isPublic": false,
      "isDeleted": false,
      "uploaderId": "g1h2i3j4-k5l6-7890-mnop-1234567890ab"
    }
  ]
}
```

**Error Response:**
- **Code**: 400 Bad Request
- **Content**:
```json
{
  "success": false,
  "error": "User does not belong to an organization"
}
```

**Example:**
```bash
curl -X GET "http://localhost:5000/files/organization?includeDeleted=false" \
  -b "token=your-session-token"
```

### Get File Types

Returns a list of supported file types and their extensions.

- **URL**: `/files/types`
- **Method**: `GET`
- **Auth Required**: Yes (Cookie)

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "fileTypes": {
    "documents": [".pdf", ".doc", ".docx", ".txt", ".rtf", ".odt", ".md", ".csv", ".xls", ".xlsx"],
    "images": [".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".bmp"],
    "audio": [".mp3", ".wav", ".ogg", ".flac", ".m4a"],
    "video": [".mp4", ".mov", ".avi", ".mkv", ".webm"],
    "archives": [".zip", ".rar", ".7z", ".tar", ".gz"]
  }
}
```

**Example:**
```bash
curl -X GET http://localhost:5000/files/types \
  -b "token=your-session-token"
```

### Get File by ID

Retrieves details of a specific file.

- **URL**: `/files/:id`
- **Method**: `GET`
- **URL Params**: `id=[uuid]` - File ID
- **Auth Required**: Yes (Cookie)

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "file": {
    "id": "f1i2l3e4-5678-90ab-cdef-1234567890ab",
    "fileName": "presentation.pdf",
    "contentType": "application/pdf",
    "size": 2048576,
    "uploadedAt": "2023-09-22T16:45:30Z",
    "isPublic": true,
    "isDeleted": false,
    "uploaderId": "a1b2c3d4-e5f6-7890-abcd-1234567890ab"
  }
}
```

**Error Response:**
- **Code**: 404 Not Found
- **Content**:
```json
{
  "success": false,
  "error": "File not found"
}
```

**Example:**
```bash
curl -X GET http://localhost:5000/files/f1i2l3e4-5678-90ab-cdef-1234567890ab \
  -b "token=your-session-token"
```

### Get File URL

Generates a presigned URL for accessing a file.

- **URL**: `/files/:id/url`
- **Method**: `GET`
- **URL Params**: `id=[uuid]` - File ID
- **Auth Required**: Yes (Cookie)

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "url": "https://storage.example.com/path/to/file?signature=abc123&expires=1632345678"
}
```

**Error Response:**
- **Code**: 400 Bad Request
- **Content**:
```json
{
  "success": false,
  "error": "Access denied: you don't have permission to access this file"
}
```

**Example:**
```bash
curl -X GET http://localhost:5000/files/f1i2l3e4-5678-90ab-cdef-1234567890ab/url \
  -b "token=your-session-token"
```

### Soft Delete File

Marks a file as deleted without removing it from storage.

- **URL**: `/files/:id`
- **Method**: `DELETE`
- **URL Params**: `id=[uuid]` - File ID
- **Auth Required**: Yes (Cookie)

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "File deleted successfully"
}
```

**Error Response:**
- **Code**: 400 Bad Request
- **Content**:
```json
{
  "success": false,
  "error": "Access denied: only the uploader can delete this file"
}
```

**Example:**
```bash
curl -X DELETE http://localhost:5000/files/f1i2l3e4-5678-90ab-cdef-1234567890ab \
  -b "token=your-session-token"
```

### Restore File

Restores a soft-deleted file.

- **URL**: `/files/:id/restore`
- **Method**: `POST`
- **URL Params**: `id=[uuid]` - File ID
- **Auth Required**: Yes (Cookie)

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "File restored successfully"
}
```

**Error Response:**
- **Code**: 400 Bad Request
- **Content**:
```json
{
  "success": false,
  "error": "File is not deleted"
}
```

**Example:**
```bash
curl -X POST http://localhost:5000/files/f1i2l3e4-5678-90ab-cdef-1234567890ab/restore \
  -b "token=your-session-token"
```

### Batch Delete Files

Soft-deletes multiple files at once.

- **URL**: `/files/batch-delete`
- **Method**: `POST`
- **Auth Required**: Yes (Cookie)

**Request Body:**
```json
{
  "fileIds": [
    "f1i2l3e4-5678-90ab-cdef-1234567890ab",
    "a1b2c3d4-5678-90ab-cdef-9876543210zy"
  ]
}
```

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "All files deleted successfully"
}
```

**Error Response:**
- **Code**: 206 Partial Content
- **Content**:
```json
{
  "success": false,
  "message": "Some files could not be deleted",
  "failedFiles": [
    "a1b2c3d4-5678-90ab-cdef-9876543210zy"
  ]
}
```

**Example:**
```bash
curl -X POST http://localhost:5000/files/batch-delete \
  -b "token=your-session-token" \
  -H "Content-Type: application/json" \
  -d '{
    "fileIds": [
      "f1i2l3e4-5678-90ab-cdef-1234567890ab",
      "a1b2c3d4-5678-90ab-cdef-9876543210zy"
    ]
  }'
```

### Run File Cleanup

Permanently deletes files that were soft-deleted before a threshold (admin only).

- **URL**: `/files/cleanup`
- **Method**: `POST`
- **Query Params**: `days=[integer]` - Number of days since deletion (default: 30)
- **Auth Required**: Yes (Cookie)

**Success Response:**
- **Code**: 200 OK
- **Content**:
```json
{
  "success": true,
  "message": "Cleanup completed successfully"
}
```

**Error Response:**
- **Code**: 403 Forbidden
- **Content**:
```json
{
  "success": false,
  "error": "Only administrators can run cleanup"
}
```

**Example:**
```bash
curl -X POST "http://localhost:5000/files/cleanup?days=30" \
  -b "token=your-session-token"
```

---

This documentation covers all the API endpoints provided by the Hello Pulse platform. For any questions or additional assistance, please contact support at projectsengineer6@gmail.com.