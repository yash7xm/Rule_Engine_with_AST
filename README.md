# Rule Engine - Backend

### Overview

The backend of the Rule Engine is developed using **Golang** and uses **PostgreSQL** as the database for storing rules and metadata. The backend provides APIs for rule creation, combination, and evaluation using an **Abstract Syntax Tree (AST).**

### Tech Stack

1. **Backend:** Golang
2. **Database:** PostgreSQL
3. **Deployment:** Render

### API Endpoints

1. `POST /create_rule`: Create a new rule from a rule string.
2. `POST /combine_rules`: Combine multiple rules.
3. `POST /evaluate`: Evaluate a rule against user attributes.

## Setup

### Prerequisites

1. **Golang:** Ensure Go is installed (v1.18 or higher).
2. **PostgreSQL:** Install PostgreSQL (v12 or higher).

## Steps to Run Locally

1. Clone the repository:

```
git clone https://github.com/yash7xm/Rule_Engine_with_AST.git
cd Rule_Engine_with_AST
```

2. Create a database:
   `CREATE DATABASE Rule_Engine;`

3. Create environment file:

```
DATABASE_URL = "postgres://postgres:yourpassword@localhost:5432/Rule_Engine?sslmode=disable"
```

4. Install Go dependencies:
   `go mod tidy`

5.Run the server:
`cd cmd/server && go run main.go`

## Run the Go tests using:

`cd test && go test`
