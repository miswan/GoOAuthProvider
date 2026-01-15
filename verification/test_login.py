from playwright.sync_api import sync_playwright

def run():
    with sync_playwright() as p:
        browser = p.chromium.launch()
        page = browser.new_page()
        try:
            page.goto("http://localhost:8081/login")
            # Wait for form to appear
            page.wait_for_selector("form")
            page.fill("input[name='username']", "testuser")
            page.fill("input[name='password']", "password")
            page.screenshot(path="verification/login.png")
            print("Screenshot saved to verification/login.png")
        except Exception as e:
            print(f"Error: {e}")
        finally:
            browser.close()

if __name__ == "__main__":
    run()
