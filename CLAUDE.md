# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **job scraping and interview tracking application** with a Go backend and Next.js frontend.

## Architecture

```
myjob_interview/
├── cmd/app/          # Go backend entry point (main + wire DI)
├── internel/         # Core business logic
│   ├── client/       # Job board scrapers (jobsdb, jobthai)
│   ├── handler/      # HTTP handlers (Gin)
│   ├── model/        # Data models
│   ├── repository/   # MongoDB data access
│   ├── usecase/      # Business logic (scraping, skill analysis)
│   ├── provider/     # Provider registration pattern
│   └── db/           # Database connections
├── front/            # Next.js 16 frontend
└── workspace/        # Volume mount for MongoDB data
```

## Tech Stack

- **Backend**: Go 1.25.4, Gin, Wire (DI), MongoDB driver, Ollama (for AI skill analysis)
- **Frontend**: Next.js 16, React 19, Tailwind CSS 4
- **Database**: MongoDB
- **Containerization**: Docker Compose

## Key Commands

```bash
# Build and run with Docker Compose
docker-compose up --build

# Run backend locally (requires MongoDB on localhost:27017)
cd /home/jakkrit/Desktop/projects/myjob_interview
go run ./cmd/app/

# Run frontend
cd /home/jakkrit/Desktop/projects/myjob_interview/front
npm run dev

# Backend linting
golangci-lint run ./...

# Frontend linting
cd front && npm run lint

# Generate Wire DI code (after modifying cmd/app/wire.go)
go run github.com/google/wire/cmd/wire ./cmd/app/
```

## API Endpoints

- `GET /api/v1/job` - Get all jobs
- `PUT /api/v1/job/:id/status` - Update job status

## Scraping Schedule

Jobs are scraped every 30 minutes from:
- **jobsdb**: Thai Golang jobs in Bangkok
- **jobthai**: Thai Golang jobs

## AI Skill Analysis

The scraper uses Ollama with model `scb10x/typhoon2.5-qwen3-4b:latest` to analyze job descriptions and extract skills into: languages, frameworks, tools, databases, hardSkills, softSkills.

## Environment Variables (.env)

```
GIN_MODE="debug"
MONGO_URI="mongodb://mongoadmin:mongoadmin@127.0.0.1:27017"
MONGO_DB_NAME="myjob_interview"
APP_PORT="8077"
APP_HOST="0.0.0.0"
```

## Frontend Configuration

- API base URL: `http://localhost:8077/api/v1`
- Next.js runs on port 3000
- Tailwind CSS v4 with PostCSS
