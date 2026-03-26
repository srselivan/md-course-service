# Course Service

Clean Architecture Go Project by Evrone

## Project Structure

```
course-service/
├── cmd/                           # Executable applications
│   └── app/
│       └── main.go               # Application entry point
│
├── internal/                       # Private application code
│   ├── app/
│   │   └── app.go                # Application initialization
│   │
│   ├── config/
│   │   └── config.go             # Configuration management
│   │
│   ├── entity/                   # Domain entities (business models)
│   │   └── course.go             # Course domain entity
│   │
│   ├── usecase/                  # Business logic / Use cases
│   │   └── course/
│   │       ├── get_course.go      # Get course use case
│   │       ├── list_courses.go    # List courses use case
│   │       └── interface.go       # Repository interfaces
│   │
│   ├── repository/               # Data access layer (adapters)
│   │   └── course/
│   │       └── postgres.go        # PostgreSQL implementation
│   │
│   ├── delivery/                 # HTTP handlers (controllers)
│   │   └── http/
│   │       └── handler/
│   │           └── course.go      # HTTP handlers for courses
│   │
│   └── middleware/               # HTTP middleware
│       └── middleware.go          # Logger, recovery, etc.
│
├── pkg/                           # Shared utilities
│   ├── logger/                    # Logger package
│   │   └── logger.go
│   │
│   └── server/                    # HTTP server wrapper
│       └── server.go
│
├── migrations/                    # Database migrations
├── go.mod                         # Go modules file
└── README.md                      # This file

```

## Architecture Layers

### 1. **Entity Layer** (`internal/entity/`)
- Contains domain models and business rules
- Independent from any frameworks or external libraries

### 2. **Use Case Layer** (`internal/usecase/`)
- Contains business logic
- Implements use cases that orchestrate the flow of data
- Depends on entities and repository interfaces

### 3. **Interface/Adapter Layer** (`internal/repository/` and `internal/delivery/`)
- **Repository**: Implements data access interfaces
- **Delivery/Handler**: Implements HTTP handlers for API endpoints

### 4. **Framework Layer** (`cmd/`)
- Contains the main application entry point
- Wires up all dependencies
- Minimal business logic

## Key Features

✅ Clean Architecture separation of concerns
✅ Dependency Inversion Principle
✅ Testability through interfaces
✅ Easy to extend and maintain
✅ Framework/library agnostic

## Getting Started

### Prerequisites
- Go 1.21 or higher

### Installation

1. Clone the repository:
```bash
git clone https://github.com/new-e-learning/course-service.git
cd courses-service
```

2. Initialize dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run cmd/app/main.go
```

## Environment Variables

```
APP_NAME=course-service
APP_VERSION=1.0.0
APP_ENV=development
SERVER_HOST=localhost
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=course
```

## API Endpoints

### Courses

- `GET /courses` - List all courses
- `GET /courses/{id}` - Get course by ID
- `POST /courses` - Create a new course
- `PUT /courses/{id}` - Update a course
- `DELETE /courses/{id}` - Delete a course

## Development

### Adding a New Feature

1. Create domain entity in `internal/entity/`
2. Define repository interface in `internal/usecase/`
3. Create use case in `internal/usecase/`
4. Implement repository in `internal/repository/`
5. Create HTTP handler in `internal/delivery/http/handler/`

### Testing

```bash
go test ./...
```

## License

MIT

## Suggestions for a good README

Every project is different, so consider which of these sections apply to yours. The sections used in the template are suggestions for most open source projects. Also keep in mind that while a README can be too long and detailed, too long is better than too short. If you think your README is too long, consider utilizing another form of documentation rather than cutting out information.

## Name
Choose a self-explaining name for your project.

## Description
Let people know what your project can do specifically. Provide context and add a link to any reference visitors might be unfamiliar with. A list of Features or a Background subsection can also be added here. If there are alternatives to your project, this is a good place to list differentiating factors.

## Badges
On some READMEs, you may see small images that convey metadata, such as whether or not all the tests are passing for the project. You can use Shields to add some to your README. Many services also have instructions for adding a badge.

## Visuals
Depending on what you are making, it can be a good idea to include screenshots or even a video (you'll frequently see GIFs rather than actual videos). Tools like ttygif can help, but check out Asciinema for a more sophisticated method.

## Installation
Within a particular ecosystem, there may be a common way of installing things, such as using Yarn, NuGet, or Homebrew. However, consider the possibility that whoever is reading your README is a novice and would like more guidance. Listing specific steps helps remove ambiguity and gets people to using your project as quickly as possible. If it only runs in a specific context like a particular programming language version or operating system or has dependencies that have to be installed manually, also add a Requirements subsection.

## Usage
Use examples liberally, and show the expected output if you can. It's helpful to have inline the smallest example of usage that you can demonstrate, while providing links to more sophisticated examples if they are too long to reasonably include in the README.

## Support
Tell people where they can go to for help. It can be any combination of an issue tracker, a chat room, an email address, etc.

## Roadmap
If you have ideas for releases in the future, it is a good idea to list them in the README.

## Contributing
State if you are open to contributions and what your requirements are for accepting them.

For people who want to make changes to your project, it's helpful to have some documentation on how to get started. Perhaps there is a script that they should run or some environment variables that they need to set. Make these steps explicit. These instructions could also be useful to your future self.

You can also document commands to lint the code or run tests. These steps help to ensure high code quality and reduce the likelihood that the changes inadvertently break something. Having instructions for running tests is especially helpful if it requires external setup, such as starting a Selenium server for testing in a browser.

## Authors and acknowledgment
Show your appreciation to those who have contributed to the project.

## License
For open source projects, say how it is licensed.

## Project status
If you have run out of energy or time for your project, put a note at the top of the README saying that development has slowed down or stopped completely. Someone may choose to fork your project or volunteer to step in as a maintainer or owner, allowing your project to keep going. You can also make an explicit request for maintainers.
