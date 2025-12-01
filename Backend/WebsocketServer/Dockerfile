# Use official Node image
FROM node:18-alpine

# Create app directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install --production

# Copy the rest of the source code
COPY . .

# Expose the port
EXPOSE 8080

# Start the app
CMD ["node", "wsserver.js"]
