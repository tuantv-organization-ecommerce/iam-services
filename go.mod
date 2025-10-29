module github.com/tvttt/iam-services

go 1.24

require (
	github.com/casbin/casbin/v2 v2.82.0
	github.com/casbin/gorm-adapter/v3 v3.20.0
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/stretchr/testify v1.8.4
	github.com/tvttt/gokits v0.0.0
	go.uber.org/zap v1.26.0
	golang.org/x/crypto v0.40.0
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.31.0
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.25.5
)

// Use local gokits for development
replace github.com/tvttt/gokits => ../gokits

// Exclude incompatible versions for Go 1.19
exclude google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f

exclude google.golang.org/protobuf v1.36.6

require (
	github.com/casbin/govaluate v1.1.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/glebarez/go-sqlite v1.20.3 // indirect
	github.com/glebarez/sqlite v1.7.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/microsoft/go-mssqldb v0.17.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230126093431-47fa9a501578 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
	gorm.io/driver/mysql v1.4.1 // indirect
	gorm.io/driver/sqlserver v1.4.1 // indirect
	gorm.io/plugin/dbresolver v1.3.0 // indirect
	modernc.org/libc v1.22.2 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.5.0 // indirect
	modernc.org/sqlite v1.20.3 // indirect
)
