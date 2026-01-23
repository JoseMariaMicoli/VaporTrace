from http.server import BaseHTTPRequestHandler, HTTPServer

class MockMetadataHandler(BaseHTTPRequestHandler):
    def do_PUT(self):
        # Simulate AWS IMDSv2 Token Acquisition
        if self.path == "/latest/api/token":
            self.send_response(200)
            self.send_header("X-aws-ec2-metadata-token-ttl-seconds", "21600")
            self.end_headers()
            self.wfile.write(b"vtrace-mock-token-777")
            print("[PIVOT TEST] Issued IMDSv2 Token")

    def do_GET(self):
        # Simulate AWS Credential Extraction
        if "/latest/meta-data/iam/security-credentials/" in self.path:
            token = self.headers.get("X-aws-ec2-metadata-token")
            if token == "vtrace-mock-token-777":
                self.send_response(200)
                self.end_headers()
                self.wfile.write(b'{"AccessKeyId": "AKIA-MOCK", "SecretAccessKey": "SECRET-VTRACE"}')
                print("[PIVOT TEST] Disclosed Mock IAM Credentials")
            else:
                self.send_response(401)
                self.end_headers()
        else:
            self.send_response(404)
            self.end_headers()

if __name__ == "__main__":
    server = HTTPServer(('127.0.0.1', 8080), MockMetadataHandler)
    print("Mock Cloud Metadata Service running on http://127.0.0.1:8080...")
    server.serve_forever()