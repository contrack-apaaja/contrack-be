# Contrack BE

Gin-based backend (MVC) with Supabase HTTP client stub.

Environment
- SUPABASE_URL - your Supabase project URL (optional for hello endpoint)
- SUPABASE_KEY - your Supabase anon/service role key (optional)
- PORT - server port (default 8080)

Run locally

1. Set env vars (optional for supabase):

	On Windows (cmd.exe):

	set SUPABASE_URL=https://your-project.supabase.co
	set SUPABASE_KEY=your-key
	set PORT=8080

2. Tidy and build:

	go mod tidy
	go build

3. Run:

	.\contrack-be.exe

Project layout (recommended)

cmd/           - entrypoints (cmd/server)
configs/       - configuration files (optional)
internal/      - application code
	config/      - config loading
	controllers/ - HTTP controllers
	router/      - route registration
	services/    - services (supabase, db, etc.)
	models/      - domain models
	repository/  - data access layer
	middlewares/ - HTTP middlewares
docs/          - documentation
scripts/       - helper scripts

Endpoints
- GET /api/hello - returns a JSON message and whether supabase env is configured

Run (dev)

go run ./cmd/server

