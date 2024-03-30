# Multi-Tenant Go Server Example

This example showcases a Go server with subdomain routing for a multi-tenant setup. It employs PostgreSQL for data storage and utilizes row-level security to isolate tenant data. The server offers a RESTful API for tenant management, supporting both "pool" and "silo" partitioning strategies for data organization.

For more information on how PostgreSQL's Row-Level Security can be used for multi-tenant data isolation, you can refer to the following resources:

- [PostgreSQL Row-Level Security Documentation](https://www.postgresql.org/docs/current/ddl-rowsecurity.html)
- [Multi-tenant data isolation with PostgreSQL Row-Level Security](https://aws.amazon.com/blogs/database/multi-tenant-data-isolation-with-postgresql-row-level-security/)

## Installation

1. Clone the repository: `git clone https://github.com/alvinchoong/go-multi-tenant-server.git`
2. Install dependencies: `go mod download`

## Configuration

1. Set the `TENANT_DB` environment variable with the mapping of `tenant slugs` to db. if unspecified for a tenant, it will default to `pool`.

## Usage

1. Run the server: `make up seed server-run-app`
2. Access the API at <http://lvh.me:8080>

## Tenant

`special.lvh.me:8080` is already set up using the "silo" partitioning model.

For other tenants that should use the "pool" partitioning model follow the steps below

1. Each tenant should have a unique subdomain (e.g., `tenant1.lvh.me:8080`, `tenant2.lvh.me:8080`).
2. Use the `lvh.me:8080/api/tenants` endpoint to onboard new tenants and associate them with their subdomains.

## API Endpoints

- `/api/tenants`: Create tenant (GET, POST)
- `/api/users`: Manage tenant users (GET, POST)
- `/api/users/{id}`: Manage tenant users (GET, DELETE)
