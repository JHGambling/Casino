# Casino/casino-app/Dockerfile
FROM node:18 AS build

# Set working directory
WORKDIR /app

# Copy package files and install dependencies
COPY package.json svelte.config.js playwright.config.ts ./
COPY src ./src
COPY static ./static

RUN npm install
RUN npm run build

# Use Node.js to serve the built app
FROM node:18 AS production

# Set working directory
WORKDIR /app

# Copy built files and dependencies
COPY --from=build /app/build ./build
COPY package.json ./

RUN npm install --production

# Expose port 3000
EXPOSE 3000

# Start the app
CMD ["node", "build"]
