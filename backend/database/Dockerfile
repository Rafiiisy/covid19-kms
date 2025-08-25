# =============================================================================
# COVID-19 KMS PostgreSQL Database Container
# =============================================================================

FROM postgres:15-alpine

# Set environment variables
ENV POSTGRES_DB=covid19_kms
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=password
ENV POSTGRES_PORT=5432

# Create directory for initialization scripts
RUN mkdir -p /docker-entrypoint-initdb.d

# Copy initialization SQL scripts
COPY init-scripts/ /docker-entrypoint-initdb.d/

# Copy custom PostgreSQL configuration (if needed)
# COPY postgresql.conf /etc/postgresql/postgresql.conf

# Expose PostgreSQL port
EXPOSE 5432

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD pg_isready -U $POSTGRES_USER -d $POSTGRES_DB || exit 1

# Labels for better organization
LABEL maintainer="COVID-19 KMS Team"
LABEL description="PostgreSQL database for COVID-19 Knowledge Management System"
LABEL version="1.0"

# Optional: Set custom data directory
# VOLUME ["/var/lib/postgresql/data"]

# The postgres image already has an entrypoint, so we don't need to specify one

