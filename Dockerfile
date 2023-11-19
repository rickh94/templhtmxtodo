FROM oven/bun:latest AS bunbuilder
WORKDIR /app
COPY package.json bun.lockb .
RUN bun install
COPY . .
RUN bun run build

FROM golang:1.21 AS gobuilder

WORKDIR /app
RUN mkdir -p /data/db
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN curl -sSf https://atlasgo.sh | sh -s -- --yes
COPY go.mod go.sum .
RUN go mod download
COPY . .
COPY --from=bunbuilder /app/static/css/main.css ./static/css/main.css
COPY --from=bunbuilder /app/static/js/main.min.js ./static/js/main.min.js
RUN templ generate
RUN CGO_ENABLED=1 GOOS=linux go build -o /templtodo cmd/main.go

CMD touch $DB_PATH && atlas migrate apply --env prod && /templtodo
