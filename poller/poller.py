import random
import time

from selenium import webdriver
from selenium.common.exceptions import TimeoutException
from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.wait import WebDriverWait


def main(d):
    while True:
        d.get("http://frontend:8080/")

        if random.random() > 0.7:
            try:
                wait = WebDriverWait(driver, 10)
                wait.until(EC.presence_of_element_located((By.CSS_SELECTOR, "table#all-data")))
                rows = d.find_elements_by_css_selector("table#all-data tr")
                row = rows[random.randint(0, len(rows) - 1)]

                td = row.find_elements_by_tag_name("td")[-1]
                td.find_element_by_tag_name("i").click()

                print("Upvoted some document")
            except TimeoutException as e:
                print("Timeout: {}".format(e))
            except Exception as e:
                print("Some exception: {}".format(e))

        print("Polled frontend", flush=True)

        time.sleep(2)


if __name__ == '__main__':
    time.sleep(30)
    with webdriver.Remote(
            command_executor="http://selenium:4444/wd/hub",
            options=webdriver.ChromeOptions()
    ) as driver:
        main(driver)
