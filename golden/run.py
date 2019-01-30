# python3 -m pip install requests-html tqdm
from requests_html import HTMLSession
from tqdm import tqdm 

session = HTMLSession()
images = {}
links = open('links.txt','r').read().split('\n')
for i,line in tqdm(enumerate(links),total=len(links)):
  line = line.strip()
  if len(line) == 0:
    continue
  r = session.get(line)
  for img in r.html.find('img'):
    if 'src' not in img.attrs:
      continue
    if '/swatches/' in img.attrs['src'] and '-l.jpg' in img.attrs['src']:
      images[img.attrs['src']] = ""

print("\n".join(list(images.keys())))
