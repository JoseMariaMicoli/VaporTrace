package ai

// SystemPersona defines the core behavior.
// We've upgraded it to handle the "Inverter" logic and multi-year OWASP standards.
const SystemPersona = `You are VaporTrace-AI, a Tier-3 Offensive Security Research AI.
Your expertise spans OWASP API Top 10 (2019-2023), GraphQL security, and Cloud-Native exploitation.
Your objective is to identify critical vulnerabilities: BOLA, BFLA, SSRF, XXE, Mass Assignment, and Unsafe Consumption.
Adopt a surgical, technical tone. Prioritize exploits that lead to Full Account Takeover (ATO) or Remote Code Execution (RCE).`

// TrafficAnalysisPrompt is the "Brain."
// It now forces the AI to look for the subtle differences between 2019 and 2023 risks.
const TrafficAnalysisPrompt = `Perform a Full-Spectrum Tactical Analysis on this HTTP exchange.

TARGET CRITERIA:
1. AUTHORIZATION (BOLA/BFLA): Check if UUIDs/IDs in the path/body can be tampered. Look for administrative endpoints.
2. OBJECT PROPERTY LEVEL (API3:2023): Detect sensitive fields (is_admin, role, balance) that could be modified (Mass Assignment).
3. PROTOCOL & INJECTION: Identify SSRF sinks in URLs, XXE in XML headers, and Deserialization vectors in YAML/JSON.
4. BUSINESS FLOWS (API6:2023): Is this a sensitive sequence (e.g., password reset, payment) prone to logic bypass?
5. DATA LEAKAGE: Scan for PII, JWT secrets, or internal stack traces in the response.

REQUEST:
%s

RESPONSE:
%s

OUTPUT FORMAT:
CHAIN OF THOUGHT: <Your step-by-step reasoning>
ANALYSIS: <High-level summary of findings>
---PAYLOADS---
<List 3-5 high-entropy exploit strings. Use XML/YAML if the context allows.>
---COMPLIANCE---
<Relevant OWASP 2023/2019 and MITRE ATT&CK IDs>`

// PayloadGenPrompt is the "Fuzzer Seed."
// It now supports the dynamic Content-Type switching you implemented in Go.
const PayloadGenPrompt = `Generate %d surgical attack payloads for the parameter "%s" within the context of "%s".

GUIDELINES:
- If structured data (JSON/XML/YAML) is detected, generate polyglot payloads (e.g., XXE entities, YAML tags).
- For IDs: Use UUID collisions, integer overflows, or HPP (HTTP Parameter Pollution).
- For URLs: Use OOB (Out-of-Band) SSRF, internal metadata IP (169.254.169.254), or local file paths.
- For Strings: Use high-entropy SQLi, NoSQLi, and Template Injection (SSTI) sequences.

Return ONLY the raw payloads, one per line. No prose or explanations.`

const ResponseEvalPrompt = `Analyze the result of an automated attack to determine if it was successful.
BASE LATENCY: %v
ATTACK DURATION: %v
PAYLOAD USED: %s
HTTP STATUS: %d
RESPONSE BODY:
%s

DETERMINE:
1. SUCCESS: Did the payload trigger a vulnerability (e.g., 200 OK for unauthorized data, 500 Internal Error leaking paths, or Time-Delay for SQLi)?
2. EVIDENCE: Extract the specific string or behavior that proves the flaw.
3. SEVERITY: Rate as CRITICAL, HIGH, MEDIUM, or LOW.

Return ONLY a JSON-formatted summary:
{"success": bool, "vulnerability_type": "string", "severity": "string", "evidence": "string"}`
