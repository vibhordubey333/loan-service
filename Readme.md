# Loan Service API

A RESTful API service for managing loans with state transitions and investment tracking.

## Architecture

This project follows Clean Architecture principles with the following layers:

1. Domain Layer (`internal/domain`)
    - Contains business entities and rules
    - Defines core types and interfaces
    - Independent of external frameworks

2. Repository Layer (`internal/repository`)
    - Handles data persistence
    - Implements database operations
    - Uses PostgreSQL for storage

3. Service Layer (`internal/service`)
    - Implements business logic
    - Manages state transitions
    - Handles email notifications and PDF generation

4. Handler Layer (`internal/handler`)
    - HTTP request handling
    - Input validation
    - Response formatting

## API Endpoints

### Create Loan
```http
POST /api/v1/loans
```

Request body:
```json
{
  "borrower_id_number": "string",
  "principal_amount": number,
  "rate": number,
  "roi": number
}
```

### Get Loan
```http
GET /api/v1/loans/{id}
```

### Approve Loan
```http
POST /api/v1/loans/{id}/approve
```

Request body:
```json
{
  "field_validator_id": "string",
  "proof_image_url": "string"
}
```

### Invest in Loan
```http
POST /api/v1/loans/{id}/invest
```

Request body:
```json
{
  "investor_id": "uuid",
  "amount": number
}
```

### Disburse Loan
```http
POST /api/v1/loans/{id}/disburse
```

Request body: 

```json
{
"field_officer_id": "string",
"signed_agreement_url": "string"
}
```

## Loan States

1. PROPOSED
   - Initial state when loan is created
   - Contains basic loan information

2. APPROVED
   - Requires field validator verification
   - Includes proof of borrower visit
   - Cannot revert to PROPOSED state

3. INVESTED
   - Achieved when total investments equal principal
   - Triggers agreement letter generation
   - Sends notifications to investors

4. DISBURSED
   - Final state when loan is given to borrower
   - Requires signed agreement
   - Records disbursement details
