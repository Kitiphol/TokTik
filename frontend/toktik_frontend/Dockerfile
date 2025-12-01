# # ---------- 1. Build Stage ----------
# FROM node:18-alpine AS builder

# # Set working directory
# WORKDIR /app

# # Copy dependency files
# COPY package.json package-lock.json* ./

# # Install dependencies
# RUN npm ci

# # Copy all source files
# COPY . .

# # Build the Next.js app
# RUN npm run build

# # ---------- 2. Production Runner ----------
# FROM node:18-alpine AS runner

# WORKDIR /app

# ENV NODE_ENV=production

# # Copy only necessary files from builder
# COPY --from=builder /app/.next ./.next
# COPY --from=builder /app/public ./public
# COPY --from=builder /app/package.json ./package.json
# COPY --from=builder /app/node_modules ./node_modules

# # Expose the port the app runs on
# EXPOSE 3000

# # Start Next.js app
# CMD ["npm", "start"]


# ---------- 1. Build Stage ----------
FROM node:18-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm ci

COPY . .
RUN npm run build

# ---------- 2. Production Runner ----------
FROM node:18-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production

# Copy standalone output
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/static ./.next/static

EXPOSE 3000

CMD ["node", "server.js"]
