# InstaPoll Project Plan

This document outlines the planned milestones for developing the InstaPoll application.

## Milestone 1: Core Infrastructure & Basic Functionality

**Goal:** Set up the project foundation, implement minimal functional polls without user accounts, and establish basic deployment and local development tooling.

* [X] Set up repository structure and documentation framework (`README.md`, `docs/`)
* [ ] Implement basic backend API for poll creation & retrieval (Go/Gin)
    * [X] Define Poll Data Structure for MongoDB
        * [X] Create poll model/struct (`models/poll.go`) with `bson` tags
        * [X] Define required fields (title, options, timestamps, etc.)
        * [X] Set up validation rules (`poll.Validate()`)
    * [ ] Set Up API Routes
        * [ ] Create endpoint for poll creation (POST `/api/polls`)
        * [ ] Create endpoint for retrieving polls (GET `/api/polls`)
        * [ ] Create endpoint for getting single poll (GET `/api/polls/:id`)
        * [ ] Add basic error handling middleware
    * [ ] Implement Database Layer (MongoDB)
        * [ ] Set up MongoDB connection (using Go driver) in `main.go`
        * [ ] Define and create MongoDB collections for polls (implicitly done via operations)
        * [ ] Implement CRUD operations for polls using MongoDB in handlers
        * [ ] Add basic database error handling in handlers
    * [X] Add Input Validation
        * [X] Validate poll creation payload
        * [ ] Sanitize user input
        * [X] Return appropriate error messages
    * [ ] Add Basic Security
        * [ ] Implement basic rate limiting
        * [ ] Add CORS configuration
        * [ ] Set up basic request logging
* [ ] Set up Local Development Environment
    * [ ] **Create `docker-compose.yml` in backend directory to manage local MongoDB for testing**
* [ ] Create simple frontend for creating and viewing polls (React/TS)
    * [ ] Create Basic UI Components
        * [ ] Build poll creation form
        * [ ] Create poll display component
        * [ ] Add loading states
        * [ ] Implement error handling UI
    * [ ] Set Up API Integration
        * [ ] Create API service layer
        * [ ] Implement poll creation function
        * [ ] Add poll fetching functions
        * [ ] Handle API errors
    * [ ] Implement Basic State Management
        * [ ] Set up state for polls list and individual polls
        * [ ] Add loading states
        * [ ] Handle error states
    * [ ] Add Basic Styling
        * [ ] Create responsive layout
        * [ ] Style poll creation form & display
* [ ] Implement initial deployment pipeline (Basic CI/CD)
    * [ ] Set up basic CI/CD pipeline (GitHub Actions) to build Go binary and Docker image
    * [ ] Add step to build static frontend assets
    * [ ] Add initial manual deployment steps for backend container (e.g., to local Docker)

## Milestone 2: Authentication & User Management

**Goal:** Enable user accounts, associate polls with creators, and implement basic security features.

* [ ] Define User schema for MongoDB
* [ ] Implement user registration and login functionality
    * [ ] Backend API endpoints for signup/login
    * [ ] Password hashing and storage
    * [ ] Implement user CRUD operations in MongoDB
* [ ] Add user profile functionality (basic view/edit)
* [ ] Associate polls with creators (add `userId` or `creatorId` to Poll schema)
* [ ] Implement basic permissions model (e.g., only creator can edit/delete poll)
* [ ] Add session management (e.g., JWT tokens) and security features
    * [ ] Implement token generation and validation
    * [ ] Secure relevant API endpoints

## Milestone 3: Enhanced Poll Features

**Goal:** Make polls more useful and flexible.

* [ ] Support for different question types (e.g., multiple choice, ranked choice voting)
* [ ] Add poll expiration options (set `ExpiresAt` field)
* [ ] Implement poll sharing functionality (generate unique shareable link)
* [ ] Create basic results visualization on the frontend (e.g., bar chart for votes)
* [ ] Add simple analytics for poll creators (e.g., view vote counts)

## Milestone 4: Real-Time Capabilities & UI Polish

**Goal:** Create a dynamic, responsive experience with improved aesthetics.

* [ ] Implement WebSocket for real-time updates
    * [ ] Add WebSocket support to the Go backend service for broadcasting vote updates
    * [ ] Frontend subscribes to updates for viewed polls
* [ ] Refine UI/UX with a consistent design system
    * [ ] Implement using a chosen component library (e.g., Material UI, Chakra UI, Shadcn/ui)
* [ ] Enhance mobile responsiveness across all views
* [ ] Implement subtle animations for voting and results display
* [ ] Create a polished landing page highlighting features

## Milestone 5: Production Readiness & Scalability

**Goal:** Prepare the application for production deployment with a focus on scalability, monitoring, and robustness using AWS infrastructure.

* [ ] Configure Production Infrastructure (Terraform)
    * [ ] Define ECS Cluster, Task Definitions, Service
    * [ ] Configure Application Load Balancer (ALB) rules and health checks
    * [ ] Set up S3 bucket for static hosting
    * [ ] Create CloudFront distribution with ACM certificate
    * [ ] Configure Route 53 DNS records
    * [ ] Provision managed MongoDB instance (Atlas/DocumentDB) or configure self-hosted replica set
* [ ] Implement Caching Layer
    * [ ] Implement caching layer (e.g., using Redis/AWS ElastiCache) for frequently accessed poll data or results
* [ ] Optimize Performance
    * [ ] Optimize MongoDB queries and indexing strategies
    * [ ] Analyze frontend bundle size and optimize loading
* [ ] Implement Scalability Features
    * [ ] Configure ECS Auto Scaling based on metrics (CPU/Memory)
* [ ] Set up Monitoring & Logging
    * [ ] Set up comprehensive performance monitoring using AWS CloudWatch Metrics and Logs
    * [ ] Configure alarms for key metrics (errors, latency, resource utilization)
* [ ] Finalize CI/CD Pipeline
    * [ ] Configure full CI/CD pipeline in GitHub Actions for automated testing and deployment to AWS ECS & S3/CloudFront

