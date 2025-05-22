 # Build frontend
 FROM node:18-alpine AS frontend-build
 WORKDIR /app/frontend
 COPY frontend/frontend-app/package.json ./
 RUN npm install
 COPY frontend/frontend-app ./
 RUN npm run build

 # Build backend
 FROM golang:1.19-alpine AS backend-build
 WORKDIR /app
 COPY go.mod go.sum ./
 # Disable Go module proxy to avoid download errors and install git for VCS operations
 ENV GOPROXY=direct
 RUN apk add --no-cache git && \
     go mod download
 COPY . ./
 RUN go build -o qrapp main.go

 # Final stage
 FROM alpine:3.17
 RUN apk add --no-cache ca-certificates
 WORKDIR /app
 COPY --from=backend-build /app/qrapp ./
 COPY --from=frontend-build /app/frontend/dist ./dist
 EXPOSE 8080
 CMD ["./qrapp"]