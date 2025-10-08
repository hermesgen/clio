module github.com/hermesgen/clio

go 1.24.0

toolchain go1.24.7

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/securecookie v1.1.2
	github.com/hermesgen/hm v0.1.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/mattn/go-sqlite3 v1.14.32
	github.com/stretchr/testify v1.11.1
	github.com/yuin/goldmark v1.7.13
	gopkg.in/yaml.v2 v2.4.0
)

// Local development - remove after hm is published
replace github.com/hermesgen/hm => ../hm

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gertd/go-pluralize v0.2.1 // indirect
	github.com/go-chi/chi/v5 v5.2.3 // indirect
	github.com/gorilla/csrf v1.7.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.42.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
