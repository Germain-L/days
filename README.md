# Days - Full-Stack Day Tracker

![Full CI](https://github.com/Germain-L/days/workflows/Full%20CI/badge.svg)
![Backend CI](https://github.com/Germain-L/days/workflows/Backend%20CI/badge.svg)
![Android CI](https://github.com/Germain-L/days/workflows/Android%20CI/badge.svg)

A modern full-stack application for tracking your days with a Go backend and Android frontend.

## 🏗️ Architecture

This project consists of two main components:

### 🖥️ Backend (Go)

- **REST API** server built with Go
- **PostgreSQL** database with migrations
- **JWT Authentication** for secure access
- **Swagger documentation** for API endpoints
- **Comprehensive testing** with coverage reports

### 📱 Frontend (Android)

- **Modern Android app** built with Jetpack Compose
- **Interactive calendar** interface
- **Material 3 design** with dark/light theme support
- **Local data persistence** with export/import features

## 🚀 Quick Start

### Using Task (Recommended)

We use [Task](https://taskfile.dev/) for build automation:

```bash
# Install Task first: https://taskfile.dev/installation/

# View all available tasks
task --list

# Start development environment
task dev                    # Shows setup instructions for both backend and Android

# Run backend development server
task backend:dev            # Starts backend with swagger docs

# Run Android development
task app:dev               # Android development workflow

# Run all tests
task test                  # Tests both backend and Android

# Build everything
task build                 # Builds both backend and Android
```

### Manual Setup

#### Backend Setup

```bash
cd backend

# Install dependencies
go mod download

# Run tests
go test ./...

# Start server
go run cmd/server/main.go
```

#### Android Setup

```bash
cd app

# Build debug APK
./gradlew assembleDebug

# Run tests
./gradlew test

# Install to device
./gradlew installDebug
```

## 🧪 Testing

Our GitHub Actions workflows automatically:

- **Backend CI**: Go tests, vet analysis, and builds
- **Android CI**: Lint checks, unit tests, and APK builds  
- **Full CI**: Combined status check for both components

Local testing:

```bash
# Test backend
task backend:test

# Test Android
task app:test

# Test everything
task test
```

## 📚 Documentation

- **Backend**: Swagger docs available at `http://localhost:8080/swagger/` when running
- **Android**: See [app/README.md](app/README.md) for detailed Android documentation
- **API**: Full REST API documentation generated with Swagger

## 🛠️ Development

### Backend Development

```bash
task backend:dev    # Auto-reloading server with swagger docs
task backend:test   # Run tests with coverage
task backend:lint   # Run go vet and formatting
```

### Android Development  

```bash
task app:dev        # Clean build, install and run
task app:test       # Unit tests and lint
task app:build      # Debug build
```

## 📦 Project Structure

```tree
.
├── .github/workflows/     # GitHub Actions CI/CD
├── backend/              # Go REST API server
│   ├── cmd/server/       # Main application entry
│   ├── internal/         # Private application code
│   ├── db/              # Database queries and migrations
│   └── docs/            # Swagger documentation
├── app/                 # Android application
│   ├── app/src/         # Android source code
│   ├── gradle/          # Gradle wrapper and libs
│   └── build.gradle.kts # Android build configuration
└── Taskfile.yml         # Build automation
```

## 🔧 Technologies

### Backend

- **Go 1.23** - Modern, efficient backend language
- **PostgreSQL** - Reliable database with ACID compliance
- **JWT** - Secure authentication tokens
- **Swagger** - API documentation and testing
- **SQLC** - Type-safe SQL code generation

### Android

- **Kotlin** - Modern Android development language
- **Jetpack Compose** - Declarative UI toolkit
- **Material 3** - Latest Material Design system
- **MVVM Architecture** - Clean separation of concerns
- **Coroutines & Flow** - Reactive programming

### DevOps

- **GitHub Actions** - Automated CI/CD pipelines
- **Task** - Cross-platform task runner
- **Docker** - Containerization for backend services

## 🚦 Status

- ✅ Backend API with full CRUD operations
- ✅ Android app with calendar interface
- ✅ Automated testing and CI/CD
- ✅ Swagger API documentation
- ✅ Local development environment
- 🔄 Database migrations and seeding
- 🔄 Cloud deployment configuration

## 📄 License

This project is developed for educational and personal use.

---

Built with ❤️ using modern full-stack development practices.
