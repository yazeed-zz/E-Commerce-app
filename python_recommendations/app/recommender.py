from app.database import get_connection

def get_popular_products(limit=10):
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("""
        SELECT p.id, p.name, p.price, p.rating_avg,
               COUNT(oi.product_id) as purchase_count
        FROM products p
        LEFT JOIN order_items oi ON p.id = oi.product_id
        GROUP BY p.id, p.name, p.price, p.rating_avg
        ORDER BY purchase_count DESC, p.rating_avg DESC
        LIMIT %s
    """, (limit,))
    rows = cur.fetchall()
    cur.close()
    conn.close()
    return [
        {
            "id": str(r[0]),
            "name": r[1],
            "price": float(r[2]),
            "rating_avg": float(r[3]) if r[3] else 0,
            "purchase_count": r[4]
        } for r in rows
    ]

def get_similar_products(product_id, limit=6):
    conn = get_connection()
    cur = conn.cursor()
    # Get category of current product
    cur.execute("SELECT category_id FROM products WHERE id = %s", (product_id,))
    row = cur.fetchone()
    if not row:
        cur.close()
        conn.close()
        return []
    category_id = row[0]

    # Get products in same category
    cur.execute("""
        SELECT id, name, price, rating_avg
        FROM products
        WHERE category_id = %s AND id != %s
        ORDER BY rating_avg DESC
        LIMIT %s
    """, (category_id, product_id, limit))
    rows = cur.fetchall()
    cur.close()
    conn.close()
    return [
        {
            "id": str(r[0]),
            "name": r[1],
            "price": float(r[2]),
            "rating_avg": float(r[3]) if r[3] else 0
        } for r in rows
    ]

def get_customers_also_bought(product_id, limit=6):
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("""
        SELECT DISTINCT p.id, p.name, p.price, p.rating_avg,
               COUNT(*) as freq
        FROM order_items oi1
        JOIN order_items oi2 ON oi1.order_id = oi2.order_id
        JOIN products p ON oi2.product_id = p.id
        WHERE oi1.product_id = %s AND oi2.product_id != %s
        GROUP BY p.id, p.name, p.price, p.rating_avg
        ORDER BY freq DESC
        LIMIT %s
    """, (product_id, product_id, limit))
    rows = cur.fetchall()
    cur.close()
    conn.close()
    return [
        {
            "id": str(r[0]),
            "name": r[1],
            "price": float(r[2]),
            "rating_avg": float(r[3]) if r[3] else 0,
            "frequency": r[4]
        } for r in rows
    ]

def get_personalized_recommendations(user_id, limit=10):
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("""
        SELECT DISTINCT p.id, p.name, p.price, p.rating_avg,
               COUNT(*) as score
        FROM orders o
        JOIN order_items oi ON o.id = oi.order_id
        JOIN products p2 ON oi.product_id = p2.id
        JOIN products p ON p.category_id = p2.category_id
        WHERE o.user_id = %s
        AND p.id NOT IN (
            SELECT oi2.product_id FROM orders o2
            JOIN order_items oi2 ON o2.id = oi2.order_id
            WHERE o2.user_id = %s
        )
        GROUP BY p.id, p.name, p.price, p.rating_avg
        ORDER BY score DESC
        LIMIT %s
    """, (user_id, user_id, limit))
    rows = cur.fetchall()
    cur.close()
    conn.close()
    return [
        {
            "id": str(r[0]),
            "name": r[1],
            "price": float(r[2]),
            "rating_avg": float(r[3]) if r[3] else 0
        } for r in rows
    ]