from playwright.sync_api import sync_playwright

def run(playwright):
    browser = playwright.chromium.launch()
    page = browser.new_page()
    page.goto("http://localhost:8081/login?next=/dashboard")
    page.screenshot(path="verification/login_page.png")
    browser.close()

with sync_playwright() as playwright:
    run(playwright)
