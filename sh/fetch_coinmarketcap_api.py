# python3 fetch_coinmarketcap_api.py
import os
import json
import redis
import asyncio
import aiohttp
import MySQLdb

from bs4 import BeautifulSoup


class APIBase:

    def __init__(self, loop):
        self.loop = loop
        self.sem = asyncio.Semaphore(value=10)

    async def get(self, *args, **kwargs):
        html = False
        async with self.sem:
            async with aiohttp.ClientSession(loop=self.loop) as session:
                if kwargs.get('html', False):
                    del kwargs['html']
                    html = True

                async with session.get(*args, **kwargs) as resp:
                    if html:
                        return await resp.read()
                    return await resp.json()


class CoinMarcketCap(APIBase):

    def __init__(self, loop):
        super(CoinMarcketCap, self).__init__(loop)
        self.rs = redis.StrictRedis(host='127.0.0.1')
        self.db = MySQLdb.connect(db='xcoin', passwd='MYSQLzhouligang153', user='root', host='127.0.0.1')
        self.fetch_sa_img()

    def get_coin_range_key(self):
        return 'coins'

    def get_coin_x_key(self, symbol, name):
        return 'coin_%s_%s' % (symbol, name)

    def if_fetch_images(self):
        key = 'if_fetch_image'
        res = self.rs.setnx(key, 1)
        self.rs.expire(key, 24*60*60)
        return res

    def fetch_all_coin_sAn_to_redis(self):
        c = self.db.cursor()
        c.execute('select symbol, name from coin;')
        rows = c.fetchall()
        pairs = ["%s_%s" % (symbol, name) for symbol, name in rows]

        key = 'coins_exists_key'
        self.rs.sadd(key, *pairs)

    def fetch_sa_img(self):
        c = self.db.cursor()
        c.execute('select symbol, name, img_url from coin;')
        rows = c.fetchall()
        d = {'%s_%s' % (symbol, name): img_url for symbol, name, img_url in rows}
        self.sa_img_dict = d
        print(self.sa_img_dict)

    def update_coin_img_url(self, symbol, name, img_url):
        c = self.db.cursor()
        sql = 'update coin set img_url = "{img_url}" where symbol="{symbol}" and name = "{name}"'.format(img_url=img_url, symbol=symbol, name=name)
        rows = c.execute(sql)
        if not rows:
            print(symbol, name, img_url)

        self.db.commit()

    def check_image_already_exsits(self, symbol, name, first_time=False):
        if first_time:
            return False

        key = 'coins_exists_key'
        mems = self.rs.smembers(key)

        if not len(mems):
            self.fetch_all_coin_sAn_to_redis()
            mems = self.rs.smembers(key)

        mems = [i.decode() for i in mems]
        item = '%s_%s' % (symbol, name)
        return item in membs

    def insert_to_db(self, coins):
        for coin in coins:
            try:
                c = self.db.cursor()
                c.execute('insert into coin (name, symbol, img_url) values ("%s", "%s", "%s")' % (coin['name'], coin['symbol'], coin["img_url"]))
                self.db.commit()
            except Exception as e:
                print(e)
                self.db.rollback()

    async def fetch_coins(self):
        url = 'https://api.coinmarketcap.com/v1/ticker/?limit=0'
        coins = await self.get(url)
        print('coin leng: %s' % len(coins))

        key = self.get_coin_range_key()
        self.rs.delete(key)

        coinids = ['%s_%s' % (i['symbol'], i['name']) for i in coins]
        self.rs.rpush(key, *coinids)

        for coin in coins:
            print('iterate coin: %s' % coin['symbol'])
            k = self.get_coin_x_key(coin['symbol'], coin['name'])
            coin.update({
                'percent_change_one_day': coin['percent_change_24h'] or '',
                'percent_change_one_hour': coin['percent_change_1h'] or '',
                'percent_change_one_week': coin['percent_change_7d'] or '',
                'img_url': self.sa_img_dict.get('%s_%s' % (coin['symbol'], coin['name']))
            })
            for i in ('percent_change_24h', 'percent_change_1h', 'percent_change_7d'):
                del coin[i]
            self.rs.set(k, json.dumps(coin))

        self.insert_to_db(coins)

    async def fetch_images(self):
        if not self.if_fetch_images():
            # 一天只去下载一次，后面又做验证，某个币的图片是否需要重新下
            return

        url = 'https://coinmarketcap.com/%s'
        for i in range(1, 20):
            html = await self.get(url % i, html=True)
            soup = BeautifulSoup(html, 'html.parser')
            tables = soup.find_all('table')
            if not len(tables):
                break

            table = tables[0]
            tbody = table.find_all('tbody')[0]
            trs = tbody.find_all('tr')

            for idx, tr in enumerate(trs):
                symbol = tr.find_all('span', 'hidden-xs')[0].text
                tds = tr.find_all('td')
                image = tds[1].find_all('img')[0]
                name = tds[1]['data-sort']

                if self.sa_img_dict.get('%s_%s' % (symbol, name)):
                    # already downloaded image, just continue
                    continue

                img_url = image['src']
                if 'coinmarketcap' not in img_url:
                    img_url = image['data-src']
                    if 'coinmarketcap' not in img_url:
                        img_url = ''
                        print('===========%s===========' % idx)
                        continue

                img_url = img_url.replace('16x16', '64x64')
                self.update_coin_img_url(symbol, name, img_url)
                self.sa_img_dict.update({
                    '%s_%s' % (symbol, name): img_url
                })
                os.system('wget %s' % (img_url))

    async def start(self):
        await self.fetch_images()
        await self.fetch_coins()

if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    cmc = CoinMarcketCap(loop)
    loop.run_until_complete(cmc.start())

