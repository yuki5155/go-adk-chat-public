# Vue.js Development Dockerfile
FROM node:22.15.0

# Set working directory
WORKDIR /usr/src/app

# Set environment variables for development
ENV NODE_ENV=development
ENV TZ=Asia/Tokyo

# Configure timezone
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Install system dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        git \
        curl \
        ca-certificates \
        && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install global tools for development
RUN npm install -g @vue/cli@latest create-vue@latest vite

# Create a non-root user
RUN groupadd -g 1001 nodejs && \
    useradd -u 1001 -g nodejs -s /bin/bash -m nodejs

# Change ownership to the non-root user
RUN chown -R nodejs:nodejs /usr/src/app

# Copy entrypoint script
COPY --chown=nodejs:nodejs docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# Expose port for Vite dev server
EXPOSE 5173

# Default command for development
ENTRYPOINT ["/bin/bash"]
CMD ["-c", "chown -R nodejs:nodejs /usr/src/app/node_modules 2>/dev/null || true && su nodejs -c 'npm install && npm run dev'"]
