# InstaPoll

A real-time polling platform built with modern web technologies and AWS cloud services. InstaPoll enables teams and communities to quickly create and participate in polls with instant results.

## Features

- **Cloud-Native Architecture**
  - AWS ECS for container orchestration
  - S3 + CloudFront for static hosting
  - DynamoDB for data persistence
  - Route 53 for DNS management
  - ACM for SSL/TLS

- **Voting System**
  - Single-choice voting
  - Real-time vote processing via WebSocket
  - Vote validation and integrity checks

- **Tech Stack**
  - Frontend: React with TypeScript
  - Backend: Go with Gin framework
  - Infrastructure: Terraform
  - CI/CD: GitHub Actions
  - Container: Docker

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Frontend      │     │  API Gateway    │     │  Poll Service   │
│  (React/TS)     │◄───►│    (AWS)        │◄───►│    (Go)         │
└─────────────────┘     └─────────────────┘     └─────────────────┘
                                                        ▲
                                                        │
┌─────────────────┐     ┌─────────────────┐             │
│  Vote Service   │◄───►│   DynamoDB      │◄────────────┘
│    (Go)         │     │    (AWS)        │
└─────────────────┘     └─────────────────┘
```

## Development

### Prerequisites
- Go 1.21+
- Node.js 18+
- Docker
- AWS CLI
- Terraform

### Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/instapoll.git
   cd instapoll
   ```

2. Set up AWS infrastructure:
   ```bash
   cd terraform
   terraform init
   terraform plan
   terraform apply
   ```

3. Configure GitHub Secrets:
   - `AWS_ACCESS_KEY_ID`
   - `AWS_SECRET_ACCESS_KEY`
   - `CLOUDFRONT_DISTRIBUTION_ID`

4. Install dependencies:
   ```bash
   # Frontend
   cd frontend
   npm install

   # Backend
   cd ../backend
   go mod download
   ```

5. Start development servers:
   ```bash
   # Frontend
   cd frontend
   npm run dev

   # Backend
   cd backend
   go run main.go
   ```

## Testing

```bash
# Frontend tests
cd frontend
npm test

# Backend tests
cd backend
go test ./...
```

## Deployment

The application uses GitHub Actions for automated deployment:

1. Push to main branch
2. GitHub Actions:
   - Runs tests
   - Builds containers
   - Deploys to AWS
   - Updates CloudFront cache

## Security

- AWS IAM roles and policies
- S3 bucket policies
- CloudFront security headers
- HTTPS/TLS with ACM
- Input validation
- Rate limiting

## Contributing

While this is primarily a personal project, suggestions and feedback are welcome.


## Author

Jesse Pye - [GitHub](https://github.com/jessepye)

---

## Features

- **Quick Poll Creation:** Create and share polls with minimal setup.
- **Multiple Voting Styles:** Supports single-choice, ranked choice, and other custom voting mechanisms.
- **Real-Time Updates:** Leverages asynchronous processing for near real-time vote tallying.
- **Distributed Architecture:** Designed with microservices, message queues, and container orchestration for scalability and fault tolerance.
- **User-Friendly Interface:** Modern front-end framework with responsive design and PWA capabilities.

