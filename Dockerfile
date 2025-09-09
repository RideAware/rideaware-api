FROM python:3.10-slim AS builder

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential gcc \
 && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .

RUN python -m pip install --upgrade pip && \
    pip wheel --no-deps -r requirements.txt -w /wheels && \
    pip wheel --no-deps gunicorn -w /wheels

FROM python:3.10-slim AS runtime

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1 \
    PORT=5000 \
    WSGI_MODULE=server:app \
    GUNICORN_WORKERS=2 \
    GUNICORN_THREADS=4 \
    GUNICORN_TIMEOUT=60 \
    FLASK_APP=server.py

WORKDIR /app

RUN groupadd -g 10001 app && useradd -m -u 10001 -g app app

COPY --from=builder /wheels /wheels
RUN pip install --no-cache-dir /wheels/* && rm -rf /wheels

# Install python-dotenv if not already in requirements.txt
RUN pip install python-dotenv

USER app

COPY --chown=app:app . .

# Copy .env file specifically
COPY --chown=app:app .env .env

EXPOSE 5000

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD python -c "import os,socket; s=socket.socket(); s.settimeout(2); s.connect(('127.0.0.1', int(os.getenv('PORT', '5000')))); s.close()"

CMD ["sh", "-c", "exec gunicorn $WSGI_MODULE --bind=0.0.0.0:$PORT --workers=$GUNICORN_WORKERS --threads=$GUNICORN_THREADS --timeout=$GUNICORN_TIMEOUT --access-logfile=- --error-logfile=- --keep-alive=5"]