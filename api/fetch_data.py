import requests
from bs4 import BeautifulSoup
import pandas as pd
import json
import io
from typing import Dict
from fastapi import FastAPI
from starlette.responses import StreamingResponse

app = FastAPI()

def fetch_and_parse_data(base_url: str, download_url: str, headers: Dict[str, str]) -> pd.DataFrame:
    session = requests.Session()
    try:
        page = session.get(base_url).content
        standard_ids = [
            json.loads(div['data-standard']).get('id')
            for div in BeautifulSoup(page, 'lxml').find_all('div', attrs={'data-standard': True})
        ]
        csv_content = session.post(download_url, headers=headers, json=standard_ids).content.decode('utf-8-sig')
        return pd.read_csv(io.StringIO(csv_content), skiprows=1)
    finally:
        session.close()

@app.get("/fetch-data")
def fetch_data():
    base_url = 'https://www.instituteforapprenticeships.org/apprenticeship-standards/'
    download_url = f'{base_url}download'
    headers = {'Content-Type': 'application/json', 'Accept': 'application/json'}
    
    # Fetch and parse the data
    data = fetch_and_parse_data(base_url, download_url, headers)
    
    # Convert DataFrame to CSV and stream it
    csv_buffer = io.StringIO()
    data.to_csv(csv_buffer, index=False)
    csv_buffer.seek(0)
    
    return StreamingResponse(
        csv_buffer,
        media_type='text/csv',
        headers={'Content-Disposition': 'attachment; filename=apprenticeship_data.csv'}
    )
