# Use an official Node.js runtime as the base image
FROM node:14

# Set the working directory inside the container
WORKDIR /usr/src/frontend

# Copy package.json and package-lock.json to the working directory
COPY package.json .
COPY package-lock.json .

# Install app dependencies
RUN npm install