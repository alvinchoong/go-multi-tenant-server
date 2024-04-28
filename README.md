# Multi-Tenant Go Server Example

This example showcases a Go server with subdomain routing for a multi-tenant setup. It employs PostgreSQL for data storage and utilizes row-level security to isolate tenant data. The server offers a RESTful API for tenant management, supporting both "pool" and "silo" partitioning strategies for data organization.

For more information on how PostgreSQL's Row-Level Security can be used for multi-tenant data isolation, you can refer to the following resources:

- [PostgreSQL Row-Level Security Documentation](https://www.postgresql.org/docs/current/ddl-rowsecurity.html)
- [Multi-tenant data isolation with PostgreSQL Row-Level Security](https://aws.amazon.com/blogs/database/multi-tenant-data-isolation-with-postgresql-row-level-security/)

## Installation

1. Clone the repository: `git clone https://github.com/alvinchoong/go-multi-tenant-server.git`
2. Install dependencies: `go mod download`

## Configuration

1. Set the `DATABASE_SILO_RW_URLS` environment variable with the mapping of `user slugs` to db. if unspecified for a user, it will default using the `pooled` db.

## Usage

1. Run the server: `make up seed run`
2. Access the API at <http://lvh.me:8080>

**Note:** `lvh.me` resolves to 127.0.0.1, useful for local development to simulate subdomains without modifying the hosts file.

## User

`special.lvh.me:8080` is already set up using the "silo" partitioning model.

For other users that should use the "pool" partitioning model follow the steps below

1. Each user should have a unique subdomain (e.g., `user1.lvh.me:8080`, `user2.lvh.me:8080`).
2. Use the `lvh.me:8080/api/users` endpoint to onboard new users and associate them with their subdomains.

## API Endpoints

- `/api/users`: Create user (GET, POST)
- `/api/todos`: Manage user todos (GET, POST)
- `/api/todos/{id}`: Manage user todos (GET, DELETE)

## Notes on Row-Level Security

- **Superusers:** Always bypass row-level security. Identified by `pg_roles.rolsuper = true`

- **BYPASSRLS Roles:** Always bypass row-level security. Identified by `pg_roles.rolbypassrls = true`

- **Table Owners:** By default, bypass row-level security. Can be enforced using `ALTER TABLE ... FORCE ROW LEVEL SECURITY`

For more information, see [PostgreSQL documentation](https://www.postgresql.org/docs/current/ddl-rowsecurity.html#DDL-ROWSECURITY).
