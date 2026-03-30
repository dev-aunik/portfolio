# syntax=docker/dockerfile:1

# ─── Development stage ───────────────────────────────────────────────────────
FROM node:20-alpine AS development
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
EXPOSE 4200
CMD ["npx", "ng", "serve", "--host", "0.0.0.0", "--poll", "2000"]

# ─── Build stage ─────────────────────────────────────────────────────────────
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --omit=dev
COPY . .
RUN npx ng build --configuration=production

# ─── Production stage (nginx) ────────────────────────────────────────────────
FROM nginx:1.27-alpine AS production
# Copy built files
COPY --from=builder /app/dist/frontend/browser /usr/share/nginx/html
# Custom nginx config for Angular SPA routing
RUN printf 'server {\n\
    listen 80;\n\
    root /usr/share/nginx/html;\n\
    index index.html;\n\
    gzip on;\n\
    gzip_types text/css application/javascript application/json;\n\
    location / { try_files $uri $uri/ /index.html; }\n\
    location /api { proxy_pass http://backend:8080; }\n\
    add_header X-Frame-Options DENY;\n\
    add_header X-Content-Type-Options nosniff;\n\
    }\n' > /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
