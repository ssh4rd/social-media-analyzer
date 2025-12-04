# Architecture Overview

## Project Summary

**Social Media Analyzer** is a Go web application designed to analyze VK (VKontakte) social media groups. It fetches group data and posts from the VK API, stores them in PostgreSQL, and provides analytics visualizations through a web interface.

---

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      WEB CLIENT (Browser)                   │
│                   (HTML/CSS/JavaScript)                     │
└──────────────────────────┬──────────────────────────────────┘
                           │
                  ┌────────▼────────┐
                  │   HTTP Server   │
                  │   (Port 3000)   │
                  └────────┬────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
    ┌───▼───┐      ┌──────▼──────┐     ┌────▼─────┐
    │Router │      │ Controllers │     │ Static   │
    │       │      │             │     │ Files    │
    └───┬───┘      └──────┬──────┘     └──────────┘
        │                 │
        │        ┌────────▼────────┐
        │        │    Services     │
        │        │                 │
        │        └────────┬────────┘
        │                 │
        │        ┌────────▼────────┐
        │        │   Models &      │
        │        │   Database      │
        │        │   (PostgreSQL)  │
        │        └─────────────────┘
        │
        └─────────────────────────────────────┐
                                              │
                         ┌────────────────────▼────┐
                         │   VK API Integration    │
                         │   (External Service)    │
                         └─────────────────────────┘
```

---

## Layered Architecture

The application follows a **4-layer Clean Architecture** pattern:

### 1. **Transport Layer** (HTTP Interface)
- **Location**: `internal/transport/http/`
- **Responsibility**: Handles HTTP requests and responses
- **Components**:
  - **Router** (`router/router.go`): Custom lightweight HTTP router with trie-based path matching
    - Supports path parameters (`:param`)
    - Supports catch-all routes (`*wildcard`)
    - Middleware support
    - Method-based routing (GET, POST)
  - **Controllers** (`controller/`):
    - `MainController`: Handles main page rendering with group analytics
    - `GroupController`: Handles group addition and post fetching

### 2. **Service Layer** (Business Logic)
- **Location**: `internal/service/`
- **Responsibility**: Core business logic and data processing
- **Components**:
  - **VKService**: External API integration
    - Parses VK group links to extract group names
    - Fetches group information from VK API
    - Fetches wall posts from groups
    - Handles VK API responses and error handling
  
  - **AnalyticsService**: Data analysis and calculations
    - Calculates group statistics (subscribers, likes, comments, etc.)
    - Computes average and maximum values from posts
    - Prepares chart data (dependence of likes/comments on subscribers)
  
  - **TemplateDataService**: Template data preparation
    - Converts analytics data to template-friendly format
    - Prepares chart data in JSON format
    - Orchestrates data preparation for rendering

### 3. **Data Layer** (Database & Models)
- **Location**: `internal/models/`, `internal/db/`
- **Responsibility**: Data persistence and modeling
- **Components**:
  - **Models** (`models/`):
    - `Group`: Represents a VK group
      - Fields: ID, Domain, Subscribers, ParsedAt
      - Relationships: One-to-Many with Posts
    - `Post`: Represents a wall post from a group
      - Fields: ID, GroupID, Date, Text, Views, Reactions, Likes, Comments
      - Relationships: Many-to-One with Group
  
  - **Database** (`db/db.go`):
    - PostgreSQL connection management
    - Uses GORM ORM for database operations
    - Connection pooling and error handling
  
  - **Migrations** (`db/migrations/`):
    - Schema initialization using GORM AutoMigrate
    - Creates Group and Post tables with constraints

### 4. **Configuration & Utilities**
- **Location**: `internal/config/`, `internal/util/`
- **Configuration** (`config/config.go`):
  - Loads environment variables
  - Provides server, database, and VK API configuration
  - DSN construction for database connection
  
- **Utilities** (`util/`):
  - Helper functions (date conversion, etc.)

---

## Data Flow

### Flow 1: Viewing Main Page with Analytics

```
User Request (GET /)
    ↓
Router matches GET / route
    ↓
MainController.GetMainPage()
    ↓
TemplateDataService.PrepareGroupsForTemplate()
    ├─ Calls AnalyticsService.CalculateGroupStats()
    │  ├─ Fetches all Groups from database
    │  └─ For each group, calculates stats from Posts
    └─ Formats data for template
    ↓
TemplateDataService.PrepareChartDataForTemplate()
    ├─ Calls AnalyticsService.CalculateChartData()
    │  └─ Prepares chart data (subscribers vs likes/comments)
    └─ Returns JSON-formatted data
    ↓
PageData {Groups, ChartData} is passed to template
    ↓
Template renders HTML with:
    - Group statistics table (paginated, 20 rows per page)
    - Line charts showing dependency of likes/comments on subscribers
    ↓
Browser receives HTML + JavaScript
```

### Flow 2: Adding a New Group

```
User submits group link (POST /api/groups)
    ↓
GroupController.AddGroup()
    ├─ Validates request
    ├─ VKService.ParseGroupFromLink()
    │  ├─ Extracts screen_name from link
    │  └─ Fetches group info from VK API
    ├─ Check if group exists in database
    │  ├─ If exists: Update group data
    │  └─ If not: Create new group
    └─ Returns HTTP 201 with GroupID
    ↓
Asynchronous post fetching (goroutine)
    ├─ GroupController.fetchAndSaveWallPosts()
    ├─ Delete all existing posts for this group
    ├─ VKService.GetWallPosts() fetches latest 100 posts
    └─ Save converted posts to database
    ↓
Response sent immediately to client
(User sees success message while posts are being fetched)
```

### Flow 3: Data Update on Re-parse

```
User submits same group link again (POST /api/groups)
    ↓
GroupController.AddGroup()
    ├─ Query database: WHERE domain = ?
    ├─ If found (existing group):
    │  └─ UPDATE group SET subscribers=..., parsed_at=now()
    └─ If not found:
       └─ INSERT new group
    ↓
Fetch and save posts
    ├─ DELETE FROM posts WHERE group_id = ?
    └─ INSERT fresh posts from VK API
    ↓
Database now contains updated group with latest posts
(No duplicates)
```

---

## Class/Structure Relationships

```
┌─────────────────────────────────────────────────────────────┐
│                       CONTROLLERS                           │
├─────────────────────────────────────────────────────────────┤
│ MainController                 GroupController              │
│ ├─ templateDataService        ├─ db (*gorm.DB)            │
│ └─ GetMainPage()              ├─ vkService                │
│                               ├─ AddGroup()               │
│                               └─ fetchAndSaveWallPosts()  │
└────────┬──────────────────────────────────┬────────────────┘
         │                                  │
         │                                  ▼
         │                        ┌──────────────────┐
         │                        │   VKService      │
         │                        ├──────────────────┤
         │                        │ - accessToken    │
         │                        │ - apiVersion     │
         │                        │ - httpClient     │
         │                        ├──────────────────┤
         │                        │ + ParseGroup..() │
         │                        │ + GetGroupInfo() │
         │                        │ + GetWallPosts() │
         │                        └──────────────────┘
         │
         ▼
┌──────────────────────────────────────┐
│    TemplateDataService               │
├──────────────────────────────────────┤
│ - analyticsService                   │
├──────────────────────────────────────┤
│ + PrepareGroupsForTemplate()         │
│ + PrepareChartDataForTemplate()      │
└────────────┬─────────────────────────┘
             │
             ▼
┌──────────────────────────────────────┐
│    AnalyticsService                  │
├──────────────────────────────────────┤
│ - db (*gorm.DB)                      │
├──────────────────────────────────────┤
│ + CalculateGroupStats()              │
│ + CalculateChartData()               │
│ - calculateGroupStat()               │
└────────────┬─────────────────────────┘
             │
             ▼
┌──────────────────────────────────────┐
│    DATABASE LAYER (GORM)             │
├──────────────────────────────────────┤
│ Models:                              │
│  ├─ Group                           │
│  │  ├─ ID (PK)                      │
│  │  ├─ Domain (Unique Index)        │
│  │  ├─ Subscribers                  │
│  │  ├─ ParsedAt                     │
│  │  └─ Posts (1:N relation)         │
│  │                                  │
│  └─ Post                            │
│     ├─ ID (PK)                      │
│     ├─ GroupID (FK)                 │
│     ├─ Date                         │
│     ├─ Text                         │
│     ├─ Views, Reactions, Likes      │
│     ├─ Comments                     │
│     └─ Group (N:1 relation)         │
└──────────────────────────────────────┘
         │
         ▼
    PostgreSQL Database
```

---

## Key Design Patterns

### 1. **Dependency Injection**
- Services are injected into controllers and other services
- Example: `MainController` receives `TemplateDataService` instead of creating it

### 2. **Separation of Concerns**
- **Controllers**: HTTP handling only
- **Services**: Business logic only
- **Models**: Data representation only
- **Database**: Persistence only

### 3. **Asynchronous Processing**
- Group post fetching happens in a goroutine
- Client gets immediate response while data loads in background

### 4. **Repository Pattern** (Not extensively used yet)
- Located in `internal/repo/` for potential data access abstraction
- Currently minimal usage; services interact directly with GORM

### 5. **Template Function Map**
- Custom `json` template function for serializing Go values to JSON in templates

---

## Data Models in Detail

### Group Model
```go
type Group struct {
    ID          uint      // Primary key
    Domain      string    // Unique group identifier (screen_name)
    Subscribers int       // Number of subscribers
    ParsedAt    time.Time // When the group was last parsed
    Posts       []Post    // Related posts (1:N)
}
```

**Database Constraints**:
- Primary Key: `ID`
- Unique Index: `Domain` (prevents duplicate groups)

### Post Model
```go
type Post struct {
    ID        uint   // Primary key
    GroupID   uint   // Foreign key to Group
    Date      string // Publication date
    Group     Group  // Related group (N:1)
    Views     int    // Post views
    Reactions int    // Post reactions
    Likes     int    // Post likes count
    Text      string // Post content
    Comments  int    // Comments count
}
```

### Statistics Models
```go
type GroupStats struct {
    ID                 uint      // Group ID
    Domain             string    // Group domain
    Subscribers        int       // Subscriber count
    ParsedAt           string    // Formatted date
    TotalPosts         int       // Number of posts analyzed
    TotalLikes         int       // Sum of all likes
    AvgLikesPerPost    float64   // Average likes per post
    MaxLikesPerPost    int       // Maximum likes on single post
    AvgCommentsPerPost float64   // Average comments per post
    PostsLastWeek      int       // Posts from last week
}

type ChartData struct {
    Subscribers []int     // Array of subscriber counts
    AvgLikes    []float64 // Array of avg likes values
    AvgComments []float64 // Array of avg comments values
}
```

---

## HTTP Routes

```
GET  /                    → MainController.GetMainPage()
                            Returns main page with analytics
                            
POST /api/groups          → GroupController.AddGroup()
                            Request: { "link": "https://vk.com/groupname" }
                            Response: { "message": "...", "group_id": 123 }
                            
GET  /static/*            → Static file server
                            CSS, JavaScript, images
```

---

## Database Schema

```sql
-- Groups table
CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    domain TEXT UNIQUE NOT NULL,
    subscribers INT DEFAULT 0,
    parsed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Posts table
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    group_id INT NOT NULL REFERENCES groups(id),
    date TEXT NOT NULL,
    text TEXT NOT NULL,
    views INT NOT NULL,
    reactions INT NOT NULL,
    likes INT NOT NULL,
    comments INT NOT NULL
);
```

---

## Frontend Architecture

### HTML Template (`web/templates/main.html`)
- Server-side rendering with Go templates
- Displays group statistics in paginated table
- Two line charts for analytics visualization
- Bootstrap 5 for styling

### JavaScript Files (`web/static/js/`)
1. **tables.js**:
   - Pagination logic (20 rows per page)
   - Table data extraction and rendering
   - Dynamic table updates

2. **charts.js**:
   - Chart.js integration
   - Renders line charts for:
     - Dependency of average likes on subscribers
     - Dependency of average comments on subscribers

3. **groups.js**:
   - Group addition form handling
   - AJAX requests to POST /api/groups
   - Loading spinner and error/success alerts

---

## Configuration

Environment variables (loaded in `internal/config/config.go`):

```env
# Server
PORT=3000

# Database
DB_HOST=postgres          # Container service name in Docker
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=social-media-analyzer

# VK API
VK_ACCESS_TOKEN=your_token_here
VK_API_VERSION=5.131
```

---

## Deployment Architecture

```
┌─────────────────────────────────────┐
│      Docker Compose (Local)         │
├─────────────────────────────────────┤
│                                     │
│  ┌─────────────────────────────┐   │
│  │  app Container              │   │
│  │  (Go application)           │   │
│  │  Port 3000 (exposed)        │   │
│  │  - Web server               │   │
│  │  - Migrations run on start  │   │
│  │  - Hot reload via volume    │   │
│  └────────┬────────────────────┘   │
│           │                        │
│           │ (TCP connection)       │
│           ▼                        │
│  ┌─────────────────────────────┐   │
│  │  postgres Container         │   │
│  │  (PostgreSQL 16 Alpine)     │   │
│  │  Port 5432 (exposed)        │   │
│  │  - Database service         │   │
│  │  - Persistent volume        │   │
│  │  - Health checks enabled    │   │
│  └─────────────────────────────┘   │
│                                     │
└─────────────────────────────────────┘
         │              │
         │              ▼
    (localhost)   PostgreSQL Volume
      :3000       (postgres_data)
      :5432
```

---

## Summary of Key Components

| Component | Purpose | Language |
|-----------|---------|----------|
| Router | HTTP routing with parameters | Go |
| Controllers | HTTP request handling | Go |
| VKService | External API integration | Go |
| AnalyticsService | Data analysis and calculations | Go |
| TemplateDataService | Data formatting for templates | Go |
| Models | Data structure definitions | Go |
| Database | PostgreSQL via GORM | PostgreSQL + Go |
| Templates | HTML rendering | Go Templates + HTML |
| Frontend | User interface | HTML/CSS/JavaScript |

---

## Data Processing Pipeline

```
Raw VK API Response
    ↓
VKService (JSON parsing)
    ↓
Go Structs (VKGroupInfo, VKWallPost, etc.)
    ↓
Database Models (Group, Post)
    ↓
Stored in PostgreSQL
    ↓
AnalyticsService (calculations)
    ↓
GroupStats / ChartData
    ↓
TemplateDataService (formatting)
    ↓
TemplateGroupData / ChartDataForTemplate
    ↓
JSON/HTML in template
    ↓
Browser display
```

---

## Notes

- **Update Logic**: When the same group is re-parsed, the group data is updated (not duplicated) and all old posts are deleted and replaced with new ones
- **Unique Constraint**: Domain field has a unique index to prevent duplicate group entries
- **Async Processing**: Post fetching happens asynchronously to avoid blocking HTTP responses
- **Custom Router**: The application implements its own HTTP router instead of using a framework like Gin or Echo
- **Template Safety**: The `json` template function safely serializes Go values for JavaScript use
