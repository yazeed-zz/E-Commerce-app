from fastapi import FastAPI, HTTPException
from app.recommender import (
    get_popular_products,
    get_similar_products,
    get_customers_also_bought,
    get_personalized_recommendations
)

app = FastAPI(title="Recommendation Service")

@app.get("/")
def health_check():
    return {"status": "Recommendation service is running ✅"}

@app.get("/recommendations/popular")
def popular_products(limit: int = 10):
    return {"products": get_popular_products(limit)}

@app.get("/recommendations/similar/{product_id}")
def similar_products(product_id: str, limit: int = 6):
    return {"products": get_similar_products(product_id, limit)}

@app.get("/recommendations/also-bought/{product_id}")
def also_bought(product_id: str, limit: int = 6):
    return {"products": get_customers_also_bought(product_id, limit)}

@app.get("/recommendations/personalized/{user_id}")
def personalized(user_id: str, limit: int = 10):
    return {"products": get_personalized_recommendations(user_id, limit)}