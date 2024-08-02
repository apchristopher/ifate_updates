# api/vercel_app.py
def handler(request):
    return {
        "statusCode": 200,
        "body": json.dumps({"message": "Hello, World!"}),
        "headers": {
            "Content-Type": "application/json"
        }
    }
