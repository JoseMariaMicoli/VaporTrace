/*
Copyright (c) 2026 José María Micoli
Licensed under the Business Source License 1.1
Change Date: 2033-02-17
Change License: Apache-2.0

You may:
✔ Study
✔ Modify
✔ Use for internal security testing

You may NOT:
✘ Offer as a commercial service
✘ Sell derived competing products
*/

package discovery

// Top100Params - High-probability query parameters for mining
var Top100Params = []string{
	"id", "user", "username", "password", "passwd", "admin", "debug", "test",
	"url", "redirect", "next", "callback", "oauth", "token", "auth", "key",
	"secret", "api_key", "session", "view", "limit", "offset", "page", "q",
	"search", "query", "filter", "file", "filename", "path", "folder", "dir",
	"cmd", "exec", "command", "email", "mail", "role", "permissions", "group",
	"order", "sort", "lang", "locale", "version", "v", "source", "config",
	"env", "details", "info", "action", "method", "grant", "access", "code",
	"client_id", "client_secret", "state", "nonce", "response_type", "scope",
	"profile", "image", "avatar", "upload", "download", "log", "report", "date",
	"start", "end", "type", "format", "json", "xml", "rss", "api", "uid", "uuid",
	"csrf", "xsrf", "_token", "payment", "card", "amount", "price", "currency",
}

// Top100Paths - Common administrative and API paths
var Top100Paths = []string{
	"admin", "administrator", "login", "register", "signin", "signup", "auth",
	"api", "v1", "v2", "api/v1", "api/v2", "dashboard", "console", "control",
	"manage", "manager", "internal", "intranet", "private", "test", "testing",
	"debug", "metrics", "health", "healthz", "status", "info", "version",
	"docs", "documentation", "swagger", "openapi", "api-docs", "swagger.json",
	"config", "conf", "settings", "setup", "install", "update", "upload",
	"files", "images", "assets", "static", "public", "resources", "css", "js",
	"users", "accounts", "profiles", "members", "groups", "roles", "sessions",
	"logs", "audit", "reports", "stats", "search", "query", "graphql", "graph",
	"payment", "billing", "invoice", "orders", "cart", "shop", "store",
	"webhook", "callback", "notifications", "events", "secrets", "keys", "tokens",
	"backup", "bak", "db", "database", "sql", "dump", "server-status", ".env",
}
