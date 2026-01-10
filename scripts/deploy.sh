#!/bin/bash
set -e

REPO_URL=${REPO_URL:-""}
GITHUB_TOKEN=${GITHUB_TOKEN:-""}
GITHUB_ACTOR=${GITHUB_ACTOR:-""}
LARITMO_DATABASE_PASSWORD=${LARITMO_DATABASE_PASSWORD:-""}
LARITMO_AUTH_JWT_SECRET=${LARITMO_AUTH_JWT_SECRET:-""}
MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD:-""}
NEWRELIC_LICENSE_KEY=${NEWRELIC_LICENSE_KEY:-""}


mkdir -p /opt/laritmo
cd /opt/laritmo

curl -o docker-compose.prod.yml https://raw.githubusercontent.com/${REPO_URL}/master/docker-compose.prod.yml

cat > .env << EOF
LARITMO_SERVER_HOST=laritmo.ru
LARITMO_SERVER_USE_TLS=false
LARITMO_DATABASE_PASSWORD=${LARITMO_DATABASE_PASSWORD}
LARITMO_AUTH_JWT_SECRET=${LARITMO_AUTH_JWT_SECRET}
MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
NEWRELIC_ENABLED=true
NEWRELIC_LICENSE_KEY=${NEWRELIC_LICENSE_KEY}
EOF

echo "${GITHUB_TOKEN}" | docker login ghcr.io -u ${GITHUB_ACTOR} --password-stdin

docker pull ghcr.io/createlab/laritmo:latest

docker compose -f docker-compose.prod.yml down
docker compose -f docker-compose.prod.yml up -d

echo "Waiting for database to be ready..."
sleep 10

docker compose -f docker-compose.prod.yml exec -T db mariadb -u laritmo -p"${LARITMO_DATABASE_PASSWORD}" laritmo -e "SELECT 'Database is ready';" || echo "Database check failed, but continuing..."

echo "Applying database migrations..."
docker compose -f docker-compose.prod.yml exec -T app goose -dir /app/migrations mysql "laritmo:${LARITMO_DATABASE_PASSWORD}@tcp(db:3306)/laritmo" up

docker compose -f docker-compose.prod.yml ps

echo "Deployment complete!"

docker image prune -f
