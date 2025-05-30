**This is a personal project where I explore building a recommendation system using Go, gRPC, and PostgreSQL — keeping things clean, modular, and practical.**

## 1. How to run?

- Navigate to the project root folder:
  ```bash
  cd ../recommendation
  ```

- Generate proto DTOs: 
    ```bash
    make dto
    ```
  

- Setup databases:
    ```bash
    make dbs
  ```
  

## 2. Database

**I keep it simple, but structured — PostgreSQL is my choice for reliable data handling.**

### Initial Data Setup

- SQL initialization files are located at:  
  `/services/$(SERVICE_NAME)/external`

- To populate the database with sample data, you can run these SQL files directly using your preferred PostgreSQL client.

- **Note:** In practice, I use proper database migrations instead of manually applying raw SQL files. Migrations help maintain consistency and version control across environments.


## 3. API's documenting

- `services/core/swagger`