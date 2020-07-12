# gin-go
RBAC scaffolding based on Gin + GROM + Casbin + Dig
# Project Structure
```bash
.
├── cmd
│   └── server            # Main
├── docs                  # Documents
├── internal
│   └── app               # Your Application
│       ├── api           # api
│       ├── config        # Config's Structure
│       ├── context       #
│       ├── errors        # Errors
│       ├── middleware    # Gin's Middleware
│       ├── models        # Model gorm mongodb
│       ├── repositories  # Repository DB
│       ├── schema        # Sechemas
│       ├── services      # Business Logic Layer
│       │   └── impl      # BLL Implement
│       └── test          # Test Cases
└── pkg                   # Common Packages
    ├── app               # Extend Gin
    ├── jwt               # JWT Auth
    ├── logger            # Logers
    └── utils             # Utilities
