import requests
import time
import sys
import os.path
from multiprocessing import Pool
from tqdm import tqdm
import pickle
import datetime
import os

CACHE_FILE = os.path.join(os.path.expanduser("~"), "Documents", f"{datetime.datetime.now().strftime('%Y-%m-%d')}_marketdata.txt")

def load_cached_data():
    if os.path.exists(CACHE_FILE):
        with open(CACHE_FILE, 'rb') as f:
            return pickle.load(f)
    return None

def cache_data(data):
    with open(CACHE_FILE, 'wb') as f:
        pickle.dump(data, f)

def get_volume(item_id):
    url = f"https://esi.evetech.net/latest/markets/10000002/history/{item_id}/"
    params = {
        'datasource': 'tranquility'
    }
    retries = 0
    while retries < 1:
        response = requests.get(url, params=params)
        if response.status_code == 200:
            volume = sum([datum['volume'] for datum in response.json()])
            return (item_id, volume)
        else:
            retries += 1
            print(f"Error: Unable to retrieve volume data for item ID {item_id} (attempt {retries} of 1)")
            time.sleep(5)
    if retries == 1:
        print(f"Error: Unable to retrieve volume data for item ID {item_id} (exceeded maximum number of retries)")
        return None

def retrieve_item_ids():
    # Set up the GET request parameters
    url = 'https://esi.evetech.net/latest/universe/types/'
    params = {
        'datasource': 'tranquility',
        'page': '1'
    }

    # Retrieve a list of all item IDs in the EVE universe
    item_ids = []
    retries = 0
    while retries < 3:
        response = requests.get(url, params=params, timeout=30)
        if response.status_code != 200:
            retries += 1
            print(f"Error: Unable to retrieve item list (attempt {retries} of 3)")
            time.sleep(5)
        else:
            retries = 0
            item_data = response.json()
            if len(item_data) == 0:
                break
            item_ids.extend(item_data)
            print(f"Retrieved {len(item_ids)} item IDs so far...")
            sys.stdout.flush()
            params['page'] = str(int(params['page']) + 1)
    if retries == 3:
        print("Error: Unable to retrieve item list (exceeded maximum number of retries)")

    return item_ids

def retrieve_volumes(item_ids):
    # Retrieve volume data for each item using multiple processes
    pool = Pool()
    chunk_size = 1000  # adjust as needed
    volumes = []
    with tqdm(total=len(item_ids), desc='Processing volumes') as progress_bar:
        for chunk in [item_ids[i:i + chunk_size] for i in range(0, len(item_ids), chunk_size)]:
            results = pool.imap_unordered(get_volume, chunk)
            with tqdm(total=len(chunk), leave=False, desc='Processing chunk') as chunk_bar:
                for volume in results:
                    if volume is not None:
                        volumes.append(volume)
                    chunk_bar.update(1)
            progress_bar.update(len(chunk))

    volumes.sort(key=lambda x: x[1], reverse=True)

    return volumes[:1000]

def retrieve_market_data(item_ids, use_cache=True):
    # Check if cached data exists and return it if requested
    if use_cache:
        cached_data = load_cached_data()
        if cached_data is not None:
            return cached_data

    # Retrieve market data for the top 1000 items by volume
    url = 'https://api.evemarketer.com/ec/marketstat/json'
    params = {
        'type_id': ','.join([str(item_id) for item_id, volume in item_ids]),
        'usesystem': '30000142',
        'region_limit': '10000002',
        'hours': '720'
    }
    response = requests.get(url, params=params)

    # Check the response status code
    if response.status_code != 200:
        print('Error: Unable to retrieve market data')
    else:
        market_data = response.json()
        results = []
        for item_data in market_data:
            item_id = int(item_data['buy']['forQuery']['types'][0])
            volume = next(volume for (item_id_, volume) in item_ids if item_id_ == item_id)
            result = {
                'type_id': item_id,
                'volume': volume,
                'buy_price': item_data['buy']['max'],
                'sell_price': item_data['sell']['min'],
                'profit_margin': (item_data['sell']['min'] - item_data['buy']['max']) / item_data['sell']['min'] * 100
            }
            results.append(result)

        results.sort(key=lambda x: x['profit_margin'], reverse=True)

        # Cache the retrieved data
        cache_data(results)

        return results[:1000]


def main():
    # Check if cached data exists and retrieve it if present
    cached_data = load_cached_data()
    if cached_data is not None:
        market_data = cached_data
    else:
        item_ids = retrieve_item_ids()
        volumes = retrieve_volumes(item_ids)
        market_data = retrieve_market_data(volumes)

    # Print the top 10 items by profit margin
    print('Top 10 items by profit margin:')
    for i, item in enumerate(market_data[:10]):
        print(f"{i + 1}. {item['type_id']}: {item['profit_margin']:.2f}%")

    # Save the market data to a file
    file_name = f"{time.strftime('%Y-%m-%d')}_marketdata.txt"
    file_path = os.path.join(os.path.expanduser("~"), "Documents", file_name)
    save_market_data(market_data, file_path)


if __name__ == '__main__':
    main()

