"""
This script scrapes the Kanshudo's 10000 most useful Japanese words and saves them to a separate file
The collections are grouped by usefulness level and packed as 100 per page, needing multiple requests to fully scrape the list.
I intended to use this list for my app. However, Leeds' list contains almost 5000 more useful words.
https://www.kanshudo.com/collections/vocab_usefulness

pip install requests
pip install beautifulsoup4
"""

from bs4 import BeautifulSoup
import requests
import re

def main():
    count = {
        1: 5,
        2: 10,
        3: 15,
        4: 20,
        5: 50
    }
    f = open('kanshudo.txt', 'w', encoding='utf-8')
    for usfulness in range(1, 6):
        for page in range(count[usfulness]):
            if page == 0:
                page_link = '1'
            else:
                page_link = str(page) + '01'
            r = requests.get('https://www.kanshudo.com/collections/vocab_usefulness/UFN-{}-{}'.format(usfulness, page_link))
            soup = BeautifulSoup(r.text, 'html.parser')
            jukugos = soup.find_all('div', class_='jukugo')
            words = [re.search("'.+'", jukugo.a['onclick']).group()[1:-1] for jukugo in jukugos]
            print(len(words))
            for word in words:
                f.write(word + '\n')
    f.close()

if __name__ == '__main__':
    main()