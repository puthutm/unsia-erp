# Stage 1: Builder
FROM node:18-alpine AS builder
WORKDIR /app

# Copy only package.json (NOT package-lock.json - avoid npm cache bug)
COPY package.json ./

# Install dependencies fresh without lock file
RUN npm install --legacy-peer-deps

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Stage 2: Runner (Production)
FROM node:18-alpine AS runner
WORKDIR /app

ENV NODE_ENV production

RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# Copy build artifacts
COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3001
ENV PORT 3001
ENV HOSTNAME "0.0.0.0"

CMD ["node", "server.js"]
