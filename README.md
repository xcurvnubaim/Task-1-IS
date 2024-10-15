# Encryption and Load Testing Application

| Team Member                           | NRP         |
|---------------------------------------|-------------|
| Surya Fadli Alamsyah                 | 5025221059  |
| Mochammad Zharif Asyam Marzuqi       | 5025221163  |



## Overview

This application manages user data encryption using AES, RC4, and DES algorithms. Each user has unique encryption keys securely stored in HashiCorp Vault. When a user saves data, the application retrieves the appropriate keys from the vault, encrypts the data, and stores it in the database. The data is decrypted when requested by the user.

For user authentication, the application uses **bcrypt** for password hashing and **JWT** (JSON Web Tokens) for session management. This ensures secure handling of sensitive information.

## Application Flow

1. **User Registration**: 
    - AES, RC4, and DES keys are generated and stored in HashiCorp Vault.
2. **Data Saving**: 
    - The application retrieves the encryption keys from the vault, encrypts the data, and saves it in the database.
3. **Data Retrieval**: 
    - The application retrieves the encryption keys, decrypts the data, and returns it to the user.
4. **Password Hashing**: 
    - Passwords are hashed using bcrypt, ensuring only the user knows their password.

## Load Testing Setup

The application has been tested using **Grafana k6** for load testing. Below is the configuration and results from load testing on 100 virtual users for 20 seconds.

### Scenarios Tested
- Upload and Encrypt with AES
- Upload and Encrypt with RC4
- Upload and Encrypt with DES
- Download and Decrypt with AES
- Download and Decrypt with RC4
- Download and Decrypt with DES

### Load Testing Tool
- **Grafana k6**

### Parameter Testing
- **`http_reqs`**: Total number of HTTP requests made during the test.
- **`http_req_waiting`**: Time spent waiting for the response after sending the request, excluding connection and response times.
- **`iteration_duration`**: Total time taken for each iteration of the VUs (Virtual Users), including setup, execution, and any sleep time.
- **`rps`**: Requests per second, indicating the number of requests sent by the VUs per second during the test.

## Test Results

| Scenario                       | Total Http Requests | Request Per Second | Http Request Duration (avg, min, med, max, p90, p95) | Iteration Duration (avg, min, med, max, p90, p95)  |
|---------------------------------|---------------------|--------------------|-----------------------------------------------------|----------------------------------------------------|
| Upload and Encrypt with AES      | 392                 | 15.18              | avg=4.7s, min=318.51ms, med=2.05s, max=22.72s, p(90)=10.52s, p(95)=21.9s  | avg=5.71s, min=1.31s, med=3.06s, max=23.72s, p(90)=11.52s, p(95)=22.9s  |
| Upload and Encrypt with RC4      | 655                 | 17.13              | avg=2.61s, min=308.74ms, med=838.61ms, max=37.21s, p(90)=5.28s, p(95)=13.05s  | avg=3.61s, min=1.3s, med=1.83s, max=38.22s, p(90)=6.29s, p(95)=14.05s  |
| Upload and Encrypt with DES      | 568                 | 19.54              | avg=2.84s, min=303.93ms, med=963.29ms, max=20.69s, p(90)=6.55s, p(95)=9.25s   | avg=3.85s, min=1.3s, med=1.96s, max=21.69s, p(90)=7.55s, p(95)=10.25s  |
| Download and Decrypt with AES    | 612                 | 16.26              | avg=2.59s, min=183.29ms, med=1.81s, max=29.78s, p(90)=5.39s, p(95)=6.55s     | avg=3.59s, min=1.18s, med=2.81s, max=30.78s, p(90)=6.39s, p(95)=7.56s  |
| Download and Decrypt with RC4    | 576                 | 23.99              | avg=2.68s, min=189.51ms, med=1.76s, max=20.67s, p(90)=5.65s, p(95)=7.57s     | avg=3.68s, min=1.19s, med=2.76s, max=21.67s, p(90)=6.65s, p(95)=8.57s  |
| Download and Decrypt with DES    | 611                 | 26.69              | avg=2.47s, min=180.23ms, med=1.87s, max=11.99s, p(90)=5.17s, p(95)=6.05s     | avg=3.47s, min=1.18s, med=2.87s, max=12.99s, p(90)=6.17s, p(95)=7.05s  |

## How to Run Backend

1. Clone the repository.
2. Setup Environment.
```sh
cp .env.example .env
```
3. Install dependencies.
```sh
go mod download
```
4. Migrate Database
```sh
make migrate
```
5. Build the application
```sh
make build
```
6. Run the application
```sh
make run
```

## How to run frontend
1. Go to frontend folder
```sh
cd task-1-is-fe
```
2. Install dependencies
```sh
npm install
```
3. Start development server
```sh
npm run dev
```