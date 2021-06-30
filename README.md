Prometheus ODBC Exporter
---

Universal Prometheus exporter for all SQL DB accessed by ODBC.

- <https://github.com/ncabatoff/dbms_exporter>
- <https://github.com/alexbrainman/odbc>

ODBC Drivers:
- [MySQL](https://dev.mysql.com/downloads/connector/odbc/)
- [MariaDB](https://downloads.mariadb.org/connector-odbc/)
- [PostgreSQL](https://odbc.postgresql.org)
- [Microsoft SQL Server](https://docs.microsoft.com/en-us/sql/connect/odbc/download-odbc-driver-for-sql-server)
- [ClickHouse](https://github.com/ClickHouse/clickhouse-odbc)
- [Exasol](https://docs.exasol.com/connect_exasol/drivers/odbc/odbc_linux.htm)

TODO:
- Graceful termination

Notepads
```bash
export CGO_CFLAGS="-g -O2 -I/opt/local/include"
export LIBRARY_PATH="$LIBRARY_PATH:/opt/local/lib"
```
