# -------- Stage 1: Build Frontend --------
FROM node:20-alpine AS frontend-build

WORKDIR /app/todo-vue-js
COPY todo-vue-js/package*.json ./
RUN npm install
COPY todo-vue-js/ .
RUN npm run build

# -------- Stage 2: Build Backend --------
FROM golang:1.25-alpine AS backend-build

WORKDIR /app/todo-api
COPY todo-api/go.mod todo-api/go.sum ./
RUN go mod download
COPY todo-api/ .
RUN go build -o todo-api main.go

# -------- Stage 3: Final Image --------
FROM alpine:latest
RUN apk add --no-cache bash postgresql-client

WORKDIR /app

# copy backend binary
COPY --from=backend-build /app/todo-api/todo-api ./

# copy frontend build
COPY --from=frontend-build /app/todo-vue-js/dist ./frontend/dist

# copy wait script
COPY todo-api/wait-for-postgres.sh ./wait-for-postgres.sh
RUN chmod +x ./wait-for-postgres.sh

EXPOSE 8080

# run backend after waiting for DB
CMD ["./wait-for-postgres.sh", "db", "./todo-api"]