# Intelligent Goods Distribution System (Microservices)

This project is an implementation of a diploma thesis on the topic **"Intelligent Goods Distribution System"**. The system is built on a microservices architecture (Java Spring Boot + Go) and solves the task of warehouse stock balancing (Load Balancing) using the **Weighted Worst-Fit Decreasing (WFD)** algorithm.

## 📋 Requirements

- Docker
- Docker Compose

## 🚀 1. Starting the System

For a clean start (with removal of old data), execute:

```bash
# Stop containers and clear data volumes (DB)
docker-compose down -v

# Build and start the system in background mode
docker-compose up --build -d
```

⏳ **Wait 15-30 seconds** after startup for the database to initialize and Liquibase to create tables.

## 🛠️ 2. Seed Data Preparation

Since the database starts empty, you need to create a test scenario: 2 warehouses, 1 product type, and 1 incoming supply.

Execute the following commands in the terminal:

**Step 1: Update admin password** (password will be 'password'):

```bash
docker exec -it warehouse-db psql -U user -d warehouse_db -c "UPDATE users SET password_hash = '\$2a\$10\$dXJ3SW6G7P50lGmMkkmwe.20cQQubK3.HZWzG3YB1tlRy.fqvM/BG' WHERE username = 'admin';"
```

> **Windows PowerShell users:** If the above doesn't work, use single quotes:
> ```powershell
> docker exec -it warehouse-db psql -U user -d warehouse_db -c 'UPDATE users SET password_hash = ''$2a$10$dXJ3SW6G7P50lGmMkkmwe.20cQQubK3.HZWzG3YB1tlRy.fqvM/BG'' WHERE username = ''admin'';'
> ```

**Step 2: Add LOGISTICIAN role to admin:**

```bash
docker exec -it warehouse-db psql -U user -d warehouse_db -c "INSERT INTO user_roles (user_id, role_id) SELECT u.id, r.id FROM users u, roles r WHERE u.username = 'admin' AND r.name = 'ROLE_LOGISTICIAN' ON CONFLICT DO NOTHING;"
```

**Step 3: Create warehouses:**

```bash
docker exec -it warehouse-db psql -U user -d warehouse_db -c "INSERT INTO warehouses (id, total_capacity, created_by) VALUES (1, 1000.0, 'setup'), (2, 500.0, 'setup') ON CONFLICT DO NOTHING;"
```

**Step 4: Create product:**

```bash
docker exec -it warehouse-db psql -U user -d warehouse_db -c "INSERT INTO products (id, volume_m3, created_by) VALUES (100, 2.0, 'setup') ON CONFLICT DO NOTHING;"
```

**Step 5: Create supply:**

```bash
docker exec -it warehouse-db psql -U user -d warehouse_db -c "INSERT INTO supplies (id, warehouse_id, status, arrival_date, created_by) VALUES (555, 1, 'RECEIVED', NOW(), 'setup') ON CONFLICT DO NOTHING;"
docker exec -it warehouse-db psql -U user -d warehouse_db -c "INSERT INTO supply_items (supply_id, product_id, quantity) VALUES (555, 100, 50) ON CONFLICT DO NOTHING;"
```

**Verify the password was updated correctly:**

```bash
docker exec -it warehouse-db psql -U user -d warehouse_db -c "SELECT username, password_hash FROM users WHERE username = 'admin';"
```

The `password_hash` should start with `$2a$10$dXJ3SW6G7P50lGmMkkmwe...`

## 🔑 3. Authorization (Getting a Token)

The system is protected by JWT. To manage it, you need to log in as an Administrator.

**Request:**

```bash
curl -X POST http://localhost:8080/api/auth/login \
     -H "Content-Type: application/json" \
     -d "{\"username\": \"admin\", \"password\": \"password\"}"
```

**Response:**

Copy the token value (`eyJ...`) from the response.

```json
{"token":"eyJhbGciOiJIUzI1NiJ9..."}
```

## ⚡ 4. Trigger Distribution Calculation

Initiate the Go algorithm to distribute supply #555.

*(Replace `YOUR_TOKEN` with the copied string)*

**Request:**

```bash
curl -X POST http://localhost:8080/api/v1/distribution/calculate \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -d "{\"supplyId\": 555}"
```

**Expected response:** `202 Accepted`

## 📊 5. Checking Results

### Method A: View Logs

See how services communicate with each other:

```bash
docker-compose logs -f --tail=100
```

Look for messages: `calculation completed`, `moves_generated`, `Successfully created shipments`.

### Method B: Database Analytics (Recommended)

Check where the system decided to send the goods. This query will show a detailed movement plan:

```bash
docker exec -it warehouse-db psql -U user -d warehouse_db -c "
SELECT 
    s.id as shipment_id, 
    s.source_id as from_wh, 
    s.destination_id as to_wh, 
    s.status, 
    si.quantity as qty, 
    p.volume_m3
FROM shipments s 
JOIN shipment_items si ON s.id = si.shipment_id
JOIN products p ON si.product_id = p.id;"
```

**Expected result (Ideal balance):**

The system should leave ~33 units at warehouse #1 (1000 m³) and send ~17 units to warehouse #2 (500 m³) to equalize their occupancy.

## 🧪 6. Running Scientific Experiment (Benchmark)

To confirm the effectiveness of the WFD algorithm compared to First-Fit (for the diploma thesis), run a benchmark test inside the Go container:

```bash
docker exec -it distribution-engine go test -v ./internal/algorithm/benchmark_test.go ./internal/algorithm/models.go ./internal/algorithm/packer.go
```

You will see a comparison table showing the **Coefficient of Variation (CV)**.

| CV Value | Balance Quality |
|----------|-----------------|
| CV ≈ 0.9 | Poor balance (First-Fit) |
| CV ≈ 0.0 | Ideal balance (WFD) |

## 🔌 API Endpoints (Reference)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/auth/login` | Get JWT token | No |
| POST | `/api/auth/register` | User registration | No |
| POST | `/api/v1/distribution/calculate` | Run distribution algorithm | Yes (Admin/Logistician) |
| GET | `/api/admin/users` | List of users | Yes (Admin) |

