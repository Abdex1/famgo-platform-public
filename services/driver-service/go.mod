module famgo/driver-service

go 1.25.0

require (
	github.com/Abdex1/FamGo-platform/shared v0.0.0
	github.com/go-chi/chi/v5 v5.0.11
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/jmoiron/sqlx v1.3.5
	github.com/lib/pq v1.12.3
	github.com/shopspring/decimal v1.4.0
	gorm.io/datatypes v1.2.7
	gorm.io/gorm v1.31.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/crypto v0.51.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	gorm.io/driver/mysql v1.5.6 // indirect
	gorm.io/driver/postgres v1.5.4 // indirect
)

replace github.com/Abdex1/FamGo-platform/shared => ../../shared
