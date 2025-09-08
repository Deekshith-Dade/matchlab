FROM python:3.12-slim

RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

RUN curl -LsSf https://astral.sh/uv/install.sh | sh
ENV PATH="/root/.local/bin:${PATH}"

RUN uv --version

WORKDIR /app
COPY sim ./sim

RUN cd sim && uv sync --frozen || true

CMD cd sim && uv run main.py
