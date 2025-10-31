Zymm Project - Backend

Written in golang v1.22.2

Docker
------

This repo includes a multi-stage `Dockerfile` and `docker-compose.yml` to build and run the backend.

Local build (requires Docker):

```powershell
# build the image
docker build -t zymm-backend:latest .

# run with the existing .env file
docker run --env-file .env -p 5000:5000 zymm-backend:latest
```

Using docker-compose (recommended for dev):

```powershell
# build and start services (backend + mssql)
docker-compose up --build

# stop
docker-compose down
```

Deploy to a VPS (simple approach)
---------------------------------

1. Copy the repo or build artifacts to the VPS (or push the image to a registry):

```powershell
# build locally then push to registry (example)
docker build -t your-registry/zymm-backend:latest .
docker push your-registry/zymm-backend:latest
```

2. On the VPS, install Docker and docker-compose, then pull or build the image and run:

```powershell
# using git on the VPS
git clone <repo> && cd zymm-backend
docker-compose up -d --build
```

Notes:
- The application expects MSSQL. The compose file includes an MSSQL service and maps port 1433. For production, use a managed DB or a separate DB host and update `.env` accordingly.
- Keep your `.env` out of source control for secrets. The provided `.env` in this repo is for local dev only.

