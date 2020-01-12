import time

from selenium import webdriver


def main(d):
    while True:
        d.get("http://frontend:8080/")
        print("Polled frontend")

        time.sleep(2)


if __name__ == '__main__':
    time.sleep(30)
    with webdriver.Remote(
            command_executor="http://selenium:4444/wd/hub",
            options=webdriver.ChromeOptions()
    ) as driver:
        main(driver)
