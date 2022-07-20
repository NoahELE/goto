set windows-shell := ["pwsh.exe", "-NoLogo", "-Command"]

default: build dev

dev:
    .\goto.exe

build: build-frontend build-backend

build-frontend:
    cd frontend && npm run build

build-backend:
    go build -tags=jsoniter

