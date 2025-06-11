FROM python:3.9-slim
WORKDIR /app
COPY main.py db.json ./
RUN pip install fastapi uvicorn
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]