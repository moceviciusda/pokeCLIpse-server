ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go build -v -o /run-app .

# Node.js layer to install and run the Prisma tunnel
FROM node:23-slim as node
WORKDIR /app
RUN npm install @prisma/ppg-tunnel

FROM debian:bookworm

RUN apt update && apt install -y ca-certificates curl nodejs npm

COPY --from=builder /run-app /usr/local/bin/
COPY --from=builder /go/bin/goose /usr/local/bin/
COPY --from=node /app/node_modules /app/node_modules
COPY ./sql/schema /migrations

# Create a startup script that sets up the psql tunnel and then runs the app
RUN echo '#!/bin/bash\n\
# Start Prisma tunnel in the background\n\
echo "Starting Prisma tunnel..."\n\
npx @prisma/ppg-tunnel --host 127.0.0.1 --port 52604 &\n\
\n\
# Wait for the tunnel to be ready\n\
echo "Waiting for Prisma tunnel to be ready..."\n\
sleep 5\n\
\n\
# Set the DATABASE_URL for the application to use the tunnel\n\
export APP_DATABASE_URL="postgres://postgres:postgres@localhost:52604?sslmode=disable"\n\
\n\
# Run database migrations using the tunnel connection\n\
echo "Running database migrations..."\n\
goose -dir /migrations postgres "$APP_DATABASE_URL" up\n\
\n\
# Start the application\n\
echo "Starting application..."\n\
exec run-app\n' > /usr/local/bin/start.sh

RUN chmod +x /usr/local/bin/start.sh

CMD ["/usr/local/bin/start.sh"]