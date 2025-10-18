# weight-tracker

This project is a simple Go application that provides a small weight-tracking web application.

## Project Structure

```
weight-tracker
├── src
│   ├── main.go
│   └── handlers
│       └── handler.go
├── go.mod
├── go.sum
└── README.md
```

## Getting Started

To run this application, follow these steps:

1. **Clone the repository**:
   ```
   git clone <repository-url>
   cd weight-tracker
   ```

2. **Install dependencies**:
   ```
   go mod tidy
   ```

3. **Run the application**:
   ```
   go run src/main.go
   ```

## Endpoints

- The application exposes endpoints for creating and listing weight entries. Refer to the `src/handlers` package for details.

## Contributing

Feel free to submit issues or pull requests for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.