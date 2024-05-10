# Multi-Tenant Go Server Example

This project provides an example of a Go multi-tenant server setup with subdomain routing, demonstrating how to isolate tenant data using PostgreSQL Row-Level Security (RLS).

## Getting Started

Follow these steps to clone the repository, install dependencies, and run the server:

1. **Clone the repository**:

    ```sh
    git clone https://github.com/alvinchoong/go-multi-tenant-server.git
    ```

2. **Install dependencies**:

    ```sh
    go mod download
    ```

3. **Start the server**:

    ```sh
    make up run
    ```

4. **Access the API** at one of the following URLs to simulate multi-tenancy:
    - <http://user-1.lvh.me:8080>
    - <http://user-2.lvh.me:8080>
    - <http://user-3.lvh.me:8080>

> **Note:** `lvh.me` resolves to `127.0.0.1`, allowing local development with subdomain simulation without modifying the hosts file.

## Postman Collection

This collection provides a set of pre-configured API requests for testing and exploring the multi-tenant server's functionality.

1. **Import the Postman Collection**:
   - In Postman, click "Import" and select "Link."
   - Paste this URL:
   `https://github.com/alvinchoong/go-multi-tenant-server/blob/main/docs/multi-tenant.postman_collection.json`

2. **Update Collection Variables**:
   - Go to the imported collection's "Variables" tab in Postman.

3. **Explore the API**:
   - Test and explore the API endpoints provided in the collection.

## How It Works

1. **Row-Level Security**:

    PostgreSQL's Row-Level Security (RLS) ensures each tenant can access only their own data. RLS policies restrict row access by comparing each row's tenant identifier with the current session's tenant:

    ```sql
    CREATE POLICY user_isolation_policy ON users
      USING (slug = current_setting('app.current_user'));

    CREATE POLICY todo_isolation_policy ON todos
      USING (user_slug = current_setting('app.current_user'));
    ```

2. **Subdomain Routing**:

    The server uses subdomains to identify tenants. The middleware extracts the subdomain (tenant identifier) from the request and stores it in the context for subsequent access:

    ```go
    // extractTenantMiddleware extracts the subdomain (tenant identifier) and adds it to the request context
    func extractTenantMiddleware(host string) func(next http.Handler) http.Handler {
        return func(next http.Handler) http.Handler {
            fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                ctx := r.Context()

                // Extract the subdomain (tenant identifier) by removing the main domain
                subdomain := strings.TrimSuffix(r.Host, "."+host)
                if subdomain != "" && subdomain != host {
                    // Store the subdomain in the request context for future access
                    ctx = context.WithValue(ctx, TenantCtxKey, subdomain)
                }

                // Serve the request with the modified context
                next.ServeHTTP(w, r.WithContext(ctx))
            })

            return http.HandlerFunc(fn)
        }
    }
    ```

3. **Tenant Context Setup via `BeforeAcquire` Hook**:

    The `pgx` PostgreSQL driver provides a `BeforeAcquire` hook to customize connection setup before acquisition from the pool. This hook extracts the tenant's identifier from the request context and sets it in the database session to enforce tenant-specific access:

    ```go
    // DB hook before acquiring a connection
    beforeAcquire := func(ctx context.Context, conn *pgx.Conn) bool {
        // Extract the tenant identifier from the request context
        if s := router.TenantFromCtx(ctx); s != "" {
            // Set the tenant for the current database session
            rows, err := conn.Query(ctx, "SELECT set_config('app.current_user', $1, false)", s)
            if err != nil {
                // Log the error and discard the connection
                slog.Error("beforeAcquire conn.Query", slog.Any("err", err))
                return false
            }
            rows.Close()
        }
        return true
    }
    ```

## Row-Level Security Notes

- **Superusers**:
  Superusers have unrestricted access to the database and can bypass all RLS policies, identified by `pg_roles.rolsuper = true`.

- **BYPASSRLS Roles**:
  Roles that are allowed to bypass RLS are identified by `pg_roles.rolbypassrls = true`.

- **Table Owners**:
  Table owners bypass RLS by default but can enforce it with `ALTER TABLE ... FORCE ROW LEVEL SECURITY`.

For further details, refer to the official [PostgreSQL documentation](https://www.postgresql.org/docs/current/ddl-rowsecurity.html#DDL-ROWSECURITY).

## Resources

- [PostgreSQL Row-Level Security Documentation](https://www.postgresql.org/docs/current/ddl-rowsecurity.html)
- [Multi-tenant data isolation with PostgreSQL Row-Level Security](https://aws.amazon.com/blogs/database/multi-tenant-data-isolation-with-postgresql-row-level-security/)
