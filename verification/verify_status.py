from playwright.sync_api import sync_playwright

def verify():
    with sync_playwright() as p:
        browser = p.chromium.launch()
        page = browser.new_page()
        try:
            page.goto("http://localhost:8081")
            page.wait_for_selector("h1:text('System Operational')") # Ensure it loads
            page.screenshot(path="verification/status_page.png")
            print("Screenshot taken")
        except Exception as e:
            print(f"Error: {e}")
        finally:
            browser.close()

if __name__ == "__main__":
    verify()
