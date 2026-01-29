#!/usr/bin/env python3
from http.server import HTTPServer, BaseHTTPRequestHandler
import json

class IMDSMockHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        # Phase 1: Discovery Trigger
        # This returns the string that loot.go is looking for.
        if self.path == "/":
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b"VaporTrace Discovery: Metadata Service active at 127.0.0.1")
            print("[+] DISCOVERY: Sent trigger string to scanner.")
            return

        # Phase 3: Credential Leak
        if "iam/security-credentials" in self.path:
            self.send_response(200)
            self.send_header('Content-Type', 'application/json')
            self.end_headers()
            mock_creds = {
                "Code": "Success",
                "AccessKeyId": "ASIA-VAPOR-TRACE-2026",
                "SecretAccessKey": "9y$B&E)H@McQfTjW$C&F)J@NcRfUjXn",
                "Token": "IQoJb3JpZ2luX2VjE[...]mock-token-data"
            }
            self.wfile.write(json.dumps(mock_creds).encode())
            print("[!!!] EXFILTRATION: Sent Mock Credentials to VaporTrace.")

    def do_PUT(self):
        # Phase 2: IMDSv2 Token Acquisition
        if "/latest/api/token" in self.path:
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b"mock-v2-token-777")
            print("[+] HANDSHAKE: Issued IMDSv2 Token.")

    # Silence POST errors from the fuzzer to keep logs clean
    def do_POST(self):
        # We send BOTH an AWS Key and the IP to force multiple triggers
        self.send_response(200)
        self.end_headers()
        content = b"CRITICAL_LEAK: IP=127.0.0.1 KEY=AKIAVAPORTRACE2026TEST"
        self.wfile.write(content)
        print(f"[!] TRIGGER DISPATCHED: {content.decode()}")

if __name__ == "__main__":
    print("[*] VaporTrace IMDSv2 Mock Service Live on http://127.0.0.1:80")
    try:
        HTTPServer(('127.0.0.1', 80), IMDSMockHandler).serve_forever()
    except PermissionError:
        print("[!] ERROR: Run with 'sudo' to bind to port 80.")