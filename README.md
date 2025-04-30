# InstaPoll

InstaPoll is a real-time polling application designed for teams and communities to quickly create, share, and participate in polls with instant results.

## Core Features

* **Easy Poll Creation:** Intuitive interface for setting up polls quickly.
* **Ranked Choice Voting:** Allows users to rank preferences, providing more nuanced results than standard multiple-choice voting.
* **Real-Time Results:** Votes are tallied and displayed dynamically (planned feature).
* **User Accounts:** (Planned) Associate polls with creators and manage user profiles.

## Tech Stack

* **Frontend:** React with TypeScript (or similar modern JS framework)
* **Backend:** Go with the Gin framework
* **Database:** MongoDB
* **Infrastructure (Primary Target: AWS):**
    * AWS ECS for container orchestration
    * AWS S3 + CloudFront for static frontend hosting
    * AWS Route 53 for DNS management
    * AWS ACM for SSL/TLS certificates
    * Terraform for Infrastructure as Code (IaC)
* **CI/CD:** GitHub Actions
* **Containerization:** Docker

## Architecture

The application is designed with a cloud-native approach, leveraging containerization for scalability and resilience. Backend services handle poll management and voting logic, communicating with the frontend via APIs.

**For detailed architecture diagrams and component descriptions, please see the [Architecture Overview document](./docs/architecture.md).**

## Prerequisites

Before you begin, ensure you have the following installed:

* Go (version 1.21 or later)
* Node.js (version 18 or later)
* Docker & Docker Compose
* AWS CLI (configured with appropriate credentials, if deploying to AWS)
* Terraform (if managing infrastructure via IaC)

## Getting Started

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/jessepye/InstaPoll.git
    cd instapoll
    ```

2.  **Set up Cloud Infrastructure (Optional - for AWS deployment):**
    * Navigate to the Terraform directory: `cd terraform`
    * Initialize Terraform: `terraform init`
    * Review the plan: `terraform plan`
    * Apply the configuration: `terraform apply`
    * *Note: Ensure your AWS credentials are configured correctly.*

3.  **Install Dependencies:**
    * **Frontend:**
        ```bash
        cd frontend
        npm install
        cd ..
        ```
    * **Backend:**
        ```bash
        cd backend
        go mod download
        cd ..
        ```

## Development

To run the application locally for development:

1.  **Configure Environment:**
    * Ensure you have a running MongoDB instance accessible.
    * Set necessary environment variables for the backend (e.g., database connection string, port). Often done via a `.env` file or system variables.

2.  **Start the Backend Server:**
    ```bash
    cd backend
    # Example: Set env vars if needed, then run
    # export MONGODB_URI="mongodb://localhost:27017"
    go run main.go
    ```
    The backend will typically be available at `http://localhost:8080` (or the configured port).

3.  **Start the Frontend Development Server:**
    ```bash
    cd frontend
    npm run dev # Or your specific command to start the dev server
    ```
    The frontend will usually be available at `http://localhost:3000` (or another port specified by your framework).

## Testing

* **Frontend Tests:**
    ```bash
    cd frontend
    npm test # Or your specific test command
    ```
* **Backend Tests:**
    ```bash
    cd backend
    go test ./...
    ```

## Deployment

Deployment to AWS is automated via GitHub Actions (example flow):

1.  Pushing changes to the `main` branch triggers the workflow.
2.  The workflow typically performs these steps:
    * Runs linters and tests.
    * Builds Docker containers for the backend service(s).
    * Pushes container images to a registry (e.g., AWS ECR).
    * Builds the frontend application.
    * Uploads static frontend assets to AWS S3.
    * Updates the AWS ECS service to deploy the new backend container.
    * Invalidates the CloudFront cache (if necessary).
    * *(Note: Requires `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, and potentially other secrets configured in GitHub Actions.)*

## Security Considerations

* Utilize AWS IAM roles and least-privilege permissions.
* Secure S3 buckets with appropriate policies.
* Configure security headers via CloudFront.
* Enforce HTTPS/TLS using ACM certificates.
* Implement robust input validation on the backend.
* Protect MongoDB with authentication and network rules.
* Consider rate limiting for public API endpoints.

## Contributing

This is currently a personal project, but suggestions and feedback are welcome. Please feel free to open an issue to discuss potential changes or report bugs.

## Author

* **Jesse Pye** - [GitHub](https://github.com/jessepye)
