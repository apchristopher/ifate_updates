# vercel_app.py
import requests
from bs4 import BeautifulSoup
import pandas as pd
import json

def fetch_and_parse_data(base_url, download_url, headers):
    session = requests.Session()
    page = session.get(base_url).content
    standard_ids = [json.loads(div['data-standard']).get('id') for div in BeautifulSoup(page, 'lxml').find_all('div', attrs={'data-standard': True})]
    csv_content = session.post(download_url, headers=headers, json=standard_ids).content.decode('utf-8-sig')
    session.close()
    return pd.read_csv(pd.compat.StringIO(csv_content), skiprows=1)

def handler(request):
    base_url = 'https://www.instituteforapprenticeships.org/apprenticeship-standards/'
    download_url = f'{base_url}download'
    headers = {'Content-Type': 'application/json', 'Accept': 'application/json'}

    try:
        data = fetch_and_parse_data(base_url, download_url, headers)
        # Convert DataFrame to JSON
        json_data = data.to_json(orient='records')
        return {
            "statusCode": 200,
            "body": json_data,
            "headers": {
                "Content-Type": "application/json"
            }
        }
    except Exception as e:
        return {
            "statusCode": 500,
            "body": json.dumps({"error": str(e)}),
            "headers": {
                "Content-Type": "application/json"
            }
        }
