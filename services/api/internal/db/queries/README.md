Handwritten SQL queries for the Revenue Leak Detective project should be placed in this directory.
Preferred layout:
- Group related queries by domain/table (e.g., `users.sql`) using sqlc `-- name:` blocks; or
- Use one-file-per-query for large/complex statements.
- File names should reflect purpose or domain.

When adding new queries:
- Add a brief header comment with purpose and expected cardinality.
- Avoid `SELECT *`; list columns explicitly for sqlc stability.
- Add ORDER BY for deterministic reads; paginate for large result sets.
- Use `:execrows` when the affected row count matters.