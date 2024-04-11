# Use the official PostgreSQL image from Docker Hub
FROM postgres:latest

# Set the PostgreSQL password (change as needed)
ENV POSTGRES_PASSWORD=mysecretpassword

# Create a new database named "beatbattles"
ENV POSTGRES_DB=beatbattles
