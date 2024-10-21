# Transaction Service

This is a transaction service written in Go. It handles user transactions and balances. This guide will help you set up and run the service using Docker.

## Prerequisites

- Docker installed on your machine
- Docker Compose installed on your machine

## Getting Started

### Clone the Repository

```sh
git clone https://github.com/yourusername/transaction-service.git
cd transaction-service
```

Build and Run the Docker Containers
To build and run the Docker containers, use Docker Compose:

```
docker-compose up --build
```

This will start the service and make it accessible at http://localhost:8080.

## API Endpoints
### Register
Endpoint: ``` POST /register ```

Description: Register a new user.

Request Body:
```
{
  "first_name": "John",
  "last_name": "Doe",
  "phone_number": "1234567890",
  "address": "123 Main St",
  "pin": "123456"
}
```

### Login

Endpoint: ```POST /login ```

Description: Login a user.

Request Body:
``` 
{
  "phone_number": "1234567890",
  "pin": "123456"
} 
```


### Authenticated Routes
For the following routes, you need to include the Authorization header with a valid JWT token obtained from the /login endpoint.

### Topup
Endpoint: ``` POST /topup ```

Description: Top up the user's balance.

Request Body:
```
{
  "amount": 1000
}
```

### Payment
Endpoint: ``` POST /pay ```

Description: Make a payment.

Request Body:
```
{
  "amount": 500,
  "remarks": "Payment for services"
}
```

### Transfer

Endpoint: ``` POST /transfer ```

Description: Transfer money to another user.

Request Body:
```
{
  "to_user_id": "recipient_user_id",
  "amount": 200,
  "remarks": "Transfer to friend"
}
```

### Transaction Report
Endpoint: ``` GET /transactions ```

Description: Get the user's transaction report.

### Update Profile
Endpoint: ``` PUT /profile ```

Description: Update the user's profile.

Request Body:
```
{
  "first_name": "John",
  "last_name": "Doe",
  "address": "123 Main St"
}
```

## Essay Test

for the essay test, you can run on  ``` ./essay ``` folder
