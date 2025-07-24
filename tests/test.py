import requests

# 服务的基础 URL
BASE_URL = "http://192.168.100.208:8080"

# 测试端点
ENDPOINT = "/api/pore_diameter"

# 测试文件路径
CIF_FILE = "tests/hMOF-1.cif"

def test_pore_diameter():
    url = f"{BASE_URL}{ENDPOINT}"
    try:
        # 打开 CIF 文件并发送 POST 请求
        with open(CIF_FILE, "rb") as file:
            files = {"structure_file": file}  # 修改字段名为 structure_file
            response = requests.post(url, files=files)

        # 打印响应状态码和内容
        print(f"Status Code: {response.status_code}")
        print(f"Response: {response.json() if response.headers.get('Content-Type') == 'application/json' else response.text}")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    test_pore_diameter()