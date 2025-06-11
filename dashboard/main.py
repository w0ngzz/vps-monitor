from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse
import json, time, os

app = FastAPI()

@app.post("/report")
async def report(req: Request):
    data = await req.json()
    if data.get('token') != "testtoken123":
        return JSONResponse(status_code=403, content={"error": "unauthorized"})

    data['timestamp'] = int(time.time())
    db = {}
    if os.path.exists("db.json"):
        db = json.load(open("db.json"))
    db[data['host_id']] = data
    with open("db.json", "w") as f:
        json.dump(db, f)
    return {"status": "ok"}

@app.get("/status")
def get_status():
    if not os.path.exists("db.json"):
        return {}
    return json.load(open("db.json"))