# Schnur

Schnur is a small Go service that analyzes strings. It exposes HTTP endpoints to store, analyze, query, and delete strings. The service computes properties such as length, palindrome check, unique character count, word count, SHA256 hash, and a character frequency map.

## Features

- Analyze and store strings via POST `/strings`
- Retrieve analyzed string by value via GET `/strings/{value}`
- List and filter stored strings via GET `/strings` with query parameters
- Search using simple natural language queries via GET `/strings/filter-by-natural-language?query=...`
- Delete a stored string via DELETE `/strings/{value}`

## Prerequisites

- Go toolchain (the module declares `go 1.24.4`; Go 1.20+ should work)
- git (to clone the repo)

## Dependencies

Dependencies are managed via Go modules (`go.mod`). The main third-party dependency used for local development is:

- `github.com/joho/godotenv` — optional helper to load a `.env` file

To download modules manually run:

```bash
go mod download
```

## Environment variables

The server reads the following environment variable:

- `PORT` — (optional) TCP port the server listens on. Defaults to `8080` when not set.

For local development you can create a `.env` file in the project root with:

```env
PORT=8080
```

`main.go` uses `godotenv.Load()` so `.env` will be loaded automatically when present.

## Run locally

1. Clone the repository and change into the project directory:

```bash
git clone <repo-url>
cd Schnur
```

2. (Optional) Download dependencies:

```bash
go mod download
```

3. Build and run:

```bash
go build -o schnur
./schnur
```

Or run directly with:

```bash
go run main.go
```

The server will print the port it's using and start listening for HTTP requests.

## API (quick examples)

All endpoints use JSON for request/response unless noted otherwise.

- POST /strings

	Request body example:

	```json
	{ "value": "hello world" }
	```

	Response: `201 Created` with the stored string object and computed properties.

- GET /strings/{value}

	Retrieve the stored analysis for the exact string value. Example: `GET /strings/hello%20world`.

- GET /strings

	List stored strings with optional query filters:

	- `is_palindrome=true|false`
	- `min_length=<n>`
	- `max_length=<n>`
	- `word_count=<n>`
	- `contains_character=<c>` (single character)

	Example: `GET /strings?min_length=5&contains_character=a`

- GET /strings/filter-by-natural-language?query=...

	Supports a few simple English-like queries such as:

	- `strings longer than 5`
	- `single word palindromic strings`
	- `strings containing the letter a`

- DELETE /strings/{value}

	Remove a stored string by value.

## Testing

Run the test suite:

```bash
go test ./...
```

If tests fail, check the failing test output and open the corresponding `_test.go` files under `cmd/`.

## Troubleshooting

- If the server fails to start because the port is already in use, change the `PORT` variable.
- If environment variables in `.env` are not loaded, ensure the file is in the project root and formatted as `KEY=VALUE` pairs.

## Contributing

Contributions are welcome. Please open an issue or submit a pull request. Include tests for new behavior.

## License

No license file is included in this repository. Add a license if you plan to publish it.

---

