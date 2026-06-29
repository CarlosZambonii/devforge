"""
Servico de inferencia: carrega o modelo treinado e expoe /predict.
"""
import re
import pickle
from urllib.parse import urlparse
from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI(title="DevForge ML Service")

with open("model/url_classifier.pkl", "rb") as f:
    model = pickle.load(f)


def extract_features(url: str) -> list:
    parsed = urlparse(url if "://" in url else "http://" + url)
    host = parsed.netloc
    return [
        len(url), len(host), url.count("."), url.count("-"),
        url.count("@"), url.count("/"),
        1 if re.search(r"\d+\.\d+\.\d+\.\d+", host) else 0,
        1 if parsed.scheme == "https" else 0,
        len(re.findall(r"\d", url)),
        1 if any(w in url.lower() for w in
            ["login", "secure", "account", "verify", "update", "bank"]) else 0,
        len(host.split(".")),
    ]


class PredictRequest(BaseModel):
    url: str


class PredictResponse(BaseModel):
    url: str
    risk: str
    score: float


@app.get("/health")
def health():
    return {"status": "ok"}


@app.post("/predict", response_model=PredictResponse)
def predict(req: PredictRequest):
    features = [extract_features(req.url)]
    proba = model.predict_proba(features)[0][1]  # prob de ser suspeita
    risk = "high" if proba >= 0.7 else "medium" if proba >= 0.4 else "low"
    return PredictResponse(url=req.url, risk=risk, score=round(float(proba), 3))
