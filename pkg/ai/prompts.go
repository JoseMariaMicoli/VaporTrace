package ai

const SystemPersona = `You are VaporTrace-AI, an expert offensive security engineer and penetration tester specialized in API security (OWASP Top 10).
Your goal is to analyze HTTP traffic for vulnerabilities such as BOLA, BFLA, Injection, and Logic Flaws.
Be concise, technical, and prioritize high-severity findings.
Output format: Markdown.`

const TrafficAnalysisPrompt = `Analyze the following HTTP Request/Response pair for security weaknesses.
Look for:
1. Sensitive data leakage (PII, Keys) in response.
2. Missing security headers.
3. Potentially vulnerable parameters (IDs, file paths, URLs).
4. Logic flaws (Admin endpoints exposed).

REQUEST:
%s

RESPONSE:
%s

Provide a brief risk assessment and 2 potential attack vectors.`

const PayloadGenPrompt = `Generate %d specific attack payloads for a parameter of type "%s" in the context of "%s".
Return ONLY the raw payloads, one per line. No explanation.`
