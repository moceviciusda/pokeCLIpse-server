ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN go build -v -o /run-app .


FROM debian:bookworm

RUN apt update && apt install -y ca-certificates

COPY --from=builder /run-app /usr/local/bin/

COPY --from=builder /go/bin/goose /usr/local/bin/
COPY ./sql/schema /migrations

RUN echo '#!/bin/bash\necho "Running database migrations..."\ngoose -dir /migrations postgres "$DATABASE_URL" up\necho "Starting application..."\nexec run-app' > /usr/local/bin/start.sh
RUN chmod +x /usr/local/bin/start.sh

CMD ["/usr/local/bin/start.sh"]