# Stage 1: Build the frontend assets
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package.json and pnpm-lock.yaml to leverage Docker cache
COPY package.json pnpm-lock.yaml ./

# Install pnpm and dependencies
RUN npm install -g pnpm && pnpm install

# Copy the rest of the application source code
COPY . .

# Build the application
RUN pnpm build

# Stage 2: Serve the static files with a lightweight server
FROM nginx:stable-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]