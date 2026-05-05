// Package secrets provides secret interpolation for envchain.
//
// Values in .env files may reference external secrets using the syntax:
//
//	${<provider>:<key>}
//
// For example:
//
//	DATABASE_URL=postgres://user:${static:DB_PASS}@localhost/mydb
//	API_TOKEN=${env:MY_API_TOKEN}
//
// Built-in providers:
//
//   - env    — reads from the current process environment (EnvProvider)
//   - static — a fixed in-memory map, useful for testing (StaticProvider)
//
// Custom providers can be added by implementing the Provider interface and
// registering them with NewResolver.
package secrets
