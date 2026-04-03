# AP2 Assignment 1: Order & Payment Microservices
**Author:** Diasbek Amangeldiev  
**Deadline:** 01.04.2026 23:59

## Architecture Overview
This project implements two microservices following **Clean Architecture** principles (Domain, UseCase, Repository, and Transport layers). 

### Bounded Contexts
1. **Order Context**: Responsible for managing the order lifecycle (Pending, Paid, Failed, Cancelled).
2. **Payment Context**: Responsible for transaction processing and limit enforcement ($1000/100,000 cents).

### Communication
The services communicate synchronously over **REST**. The Order Service acts as a client to the Payment Service.

## Failure Handling & Resilience
- [cite_start]**Timeout**: The Order Service uses a custom `http.Client` with a **2-second timeout**[cite: 101].
- [cite_start]**Availability**: If the Payment Service is down, the Order Service returns a `503 Service Unavailable`[cite: 105].
- **State Consistency**: Orders are initially saved as `Pending`. [cite_start]If a payment call fails or times out, the order is marked as `Failed` to ensure the user knows the transaction did not complete[cite: 106].

## Database Setup
- Each service has its own dedicated database schema.
- SQL migrations are located in the `/migrations` folder of each service.
- **Order DB**: Stores order details.
- **Payment DB**: Stores transaction records.

## How to Run
1. Run PostgreSQL locally on port 5432.
2. Create `order_db` and `payment_db`.
3. Run migrations provided in each service.
4. `go run cmd/order-service/main.go`
5. `go run cmd/payment-service/main.go`

##The Diagram
[ ORDER SERVICE (:8080) ]          [ PAYMENT SERVICE (:8081) ]
-------------------------          ---------------------------
|   Transport (HTTP)    |          |    Transport (HTTP)     |
|          |            |          |           ^             |
|   UseCase (Logic)     | --REST-->|    UseCase (Logic)      |
|          |            | (Timeout)|           |             |
|   Repository (DB)     |          |    Repository (DB)      |
-------------------------          ---------------------------
           |                                   |
    [ Order Database ]                 [ Payment Database ]