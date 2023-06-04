# 爬虫
import re
from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support.ui import Select
from selenium.webdriver.edge.options import Options
# 导入环境变量
import os

# 选项
options = Options()
options.add_argument('--headless')
options.add_argument('--disable-gpu')
options.add_argument('--no-sandbox')
options.add_argument('--disable-dev-shm-usage')

# 获取用户的主目录以及项目位置
PATH = os.path.dirname(os.path.abspath(__file__)) + "/../"

# 通过爬虫读取公文通的信息
url = r'https://www1.szu.edu.cn/board/'
# 定位到公文通的首页（先登录）
browser = webdriver.Edge(options=options)
browser.get(url)
username = browser.find_element('id', 'username')
username.send_keys("账号")
passwd = browser.find_element('id', 'password')
passwd.send_keys("密码", Keys.RETURN)
for department in ['党政办公室', '教务部', '招生办公室', '研究生院', '科学技术部']:
    browser.get(url + 'infolist.asp?')
    # 定位并设置年份为2022年
    year = browser.find_element('name', 'dayy')
    year_list = Select(year)
    year_list.select_by_visible_text('2022年')
    # 定位并设置不同的发文单位
    dpt = browser.find_element('name', 'from_username')
    dpt_list = Select(dpt)
    dpt_list.select_by_value(department)
    # 搜索
    search_box = browser.find_element('name', 'searchb1')
    search_box.click()
    # 导入最近40篇文章
    table = browser.find_element(
        'xpath', '/html/body/table/tbody/tr[2]/td/table/tbody/tr[3]/td/table/tbody/tr[3]/td')
    table_html = table.get_attribute('innerHTML')
    article_url = re.compile(
        r'<a target="_blank" class="fontcolor3" href="(?P<url>.*?)">', re.S).finditer(table_html)
    for i, article in enumerate(article_url):
        if i == 40:
            break
        id = article.group('url')
        a_url = url + id
        browser.get(a_url)
        content = browser.find_element(
            'xpath', '/html/body/table/tbody/tr[2]/td/table/tbody/tr[3]/td/table/tbody/tr[2]/td')

        # 写入文件 /data/departmentName/title.txt
        file_path = os.path.join(PATH, "pages", f"{department}_{i}.txt")
        if os.path.exists(file_path):
            print(f"File {file_path} already exists, skipping...")
        else:
            with open(file_path, 'w') as f:
                f.write(content.text)
browser.quit()
