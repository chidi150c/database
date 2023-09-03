# Use an official SQLite runtime as a parent image
FROM sqlite:latest

# Set the working directory
WORKDIR /app

# Copy your SQLite database file to the container (if needed)
COPY sqlite.db /app

# Specify a command to run when the container starts
CMD ["sqlite3", "sqlite.db"]
