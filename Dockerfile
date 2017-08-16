FROM golang:1.8

# Create the directory where the application will reside
RUN mkdir /app
RUN mkdir /data
RUN mkdir /data/mj

# Copy the application files (needed for production)
ADD conf_local.toml /app/conf_local.toml
ADD conf_dev.toml /app/conf_dev.toml
ADD conf_pro.toml /app/conf_pro.toml
ADD qnmahjong /app/qnmahjong

# Set the working directory to the app directory
WORKDIR /app

# Expose the application on port 8080.
# This should be the same as in the app.conf file
EXPOSE 5001
EXPOSE 5002

# Set the entry point of the container to the application executables
ENTRYPOINT ["/app/qnmahjong"]