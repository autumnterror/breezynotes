## Welcome to BreezyNotes. The place where fantasy becomes reality!
#### Development by [Breezy Innovation RZN](https://about.breezynotes.ru) 
![BREEZYNOTES](https://i.ibb.co/PvRh0KvX/favicon.png)
### Technology stack:
![MongoDB](https://img.shields.io/badge/MongoDB-%234ea94b.svg?style=for-the-badge&logo=mongodb&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
### Frontend repository:
[GitHub](https://github.com/DaniilaRyadinsky/breezy)
# Development plan

## Common goal
To develop a minimally viable version of the notes service with a microservice architecture:
- Authorization and user management (PostgreSQL)
- Working with notes, tags and blocks (MongoDB)
- Interaction between services via gRPC
- HTTP API via Echo (Gateway)
- Prepare the base for further expansion of the block system

---

# Monthly plan (4 weeks)

## Week 1 — Project Skeleton and Protocols
**Purpose:** assembly and launch of the system framework.

### Tasks
- Create a monorepository and a structure `/services/{auth,blocknote,gateway}`
- Configure docker-compose with PostgreSQL, MongoDB and Redis
- Common packages `/pkg/{proto,auth,mw,errs,config,logger}`
- Implement health-check for all services
- Create proto files for Auth and BlockNote services
- Implement the Gateway (Echo) skeleton without business logic

**Checkpoint:** All services are lifted via `docker compose up', Gateway gives 200 to `/healthz'.

---

## Week 2 — Auth-service (PostgreSQL) + JWT
**Purpose:** registration, authorization, refresh stream via cookies.

### Tasks
- Create tables for `users` and `sessions' (migrations)
- Implement:
- `createUser`
  - `Auth`
  - `Refresh`
  - `GetUser`
  - `UpdateProfile`
- Add middleware auto-refresh to Gateway
- To test the flow: `reg → login → me → refresh → logout`

**Checkpoint:** curl allows you to go through the full registration and authorization scenario.

---

## Week 3 — BlockNote-service (MongoDB) + basic CRUD
**Purpose:** CRUD notes, tags, blocks and a simple ACL.

### Tasks
- Collections of `notes`, `blocks', `tags`
- Methods:
- `CreateNote`, `GetNote', `GetAllNotes`
  - `NoteToTrash`, `NoteFromTrash`, `GetNotesFromTrash`
  - `CreateBlock`, `GetAllBlocksInNote`, `ReplaceBlockOrder`, `DeleteBlock`
  - `CreateTag`, `AddTagToNote`, `GetNotesByTag`
- Implement `OpBlock` and the `BlockPlugin` plug-in interface
- Prototype of the `Text` and `Image` blocks

**Checkpoint:** You can create a note, add a block, change the order, and restore it from the trash.

---

## Week 4 — Sharing, Cache, and Finalizing
**Goal:** Complete the functionality and prepare the MVP release.
### Tasks
- `ShareNote` and `ChangeUserRole`
- Checking roles during access
- Mini cache in Redis (optional)
- Metrics, logs, documentation
- Mini search by title and tags

**Checkpoint:** The end-to-end script is completely error-free.

---
# Architecture by services

### Auth-service
- Database: PostgreSQL  
- Main methods: `Auth`, `Refresh`, `createUser`, `getUser`, `UpdateProfile`  
- Stores users and refresh sessions

### BlockNote-service
- DATABASE: MongoDB  
- Main entities: `Note`, `Block', `Tag`  
- ACL via `author`, `editors', `readers`  
- Support for block plugins (`BlockPlugin')

### Gateway
- Framework: Echo  
- REST API → gRPC calls  
- Middleware: gzip, recover, request-id, auto-refresh

---

# Final checkpoint (MVP is ready)
- Registration and login via cookies  
- Creating notes, tags, and blocks  
- Working with the shopping cart  
-  Sharing notes  
- Auto-refresh JWT  
- Documentation and Postman collection
