# 爬虫
import re
from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support.ui import Select
from selenium.webdriver.edge.options import Options
from configparser import ConfigParser
# 导入环境变量
import os

# 选项
options = Options()
# options.add_argument('--headless')
# options.add_argument('--disable-gpu')
# options.add_argument('--no-sandbox')
# options.add_argument('--disable-dev-shm-usage')

# 获取用户的主目录以及项目位置
PATH = os.path.dirname(os.path.abspath(__file__)) + "/../"

# 读取公文通账号密码
config = ConfigParser()
# TODO: 自行创建config.ini文件
config.read(PATH + "python/config.ini")
account = config["Profile"]["account"]
password = config["Profile"]["password"]

# 通过爬虫读取深圳大学
url = r'https://xq40.szu.edu.cn/'
# 定位到公文通的首页（先登录）
browser = webdriver.Edge(options=options)
browser.get(url)

browser.get(url + "/tzgg.htm")
# 找到当前页面的所有文章table
table = browser.find_element('xpath', '/html/body/div[4]/div/div[2]/ul')
table_html = table.get_attribute('innerHTML')
article_url = re.compile(
    r'<a href="(?P<url>.*?)" class="wl">', re.S).finditer(table_html)
# 找到当前页面的跳转按钮
nav = browser.find_element(
    'xpath', '/html/body/div[4]/div/div[2]/div/div/span[1]')
nav_html = nav.get_attribute('innerHTML')
nav_url = re.compile(r'<a href="(?P<url>.*?)">', re.S).finditer(nav_html)
for i, nav in enumerate(nav_url):
    for j, article in enumerate(article_url):
        page_url = article.group('url')
        if re.compile(r'weixin').search(page_url):
            browser.get(page_url)

        else:
            browser.get(page_url)
            if j == 0:
                username = browser.find_element('id', 'username')
                username.send_keys(account)
                passwd = browser.find_element('id', 'password')
                passwd.send_keys(password, Keys.RETURN)

            content = browser.find_element(
                'xpath', '/html/body/table/tbody/tr[2]/td/table/tbody/tr[3]/td/table/tbody/tr[2]/td')

        # 写入文件 /data/departmentName/title.txt
        file_path = os.path.join(PATH, "pages", f"page_{i + j}.txt")
        if os.path.exists(file_path):
            print(f"File {file_path} already exists, skipping...")
        else:
            with open(file_path, 'w') as f:
                f.write(content.text)

    id = nav.group('url')
    browser.get(url + id)
browser.quit()
