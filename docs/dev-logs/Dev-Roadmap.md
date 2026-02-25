## 🚀 Strategic Roadmap

### **Part I: The Hardened Core & Intelligence [STABLE]**

| Phase | Sub-Phase | Focus / Technical Deliverable | Status |
| --- | --- | --- | --- |
| **Sprint 1: Foundation** |
|  | 1.1 | Cobra CLI Engine: Subcommand-based architecture (map, scan, auth). | ✅ DONE |
|  | 1.2 | Interactive Shell UI: Advanced REPL with readline auto-completion. | ✅ DONE |
|  | 1.3 | The Burp Bridge: Industrial-strength HTTP client with native proxy support. | ✅ DONE |
|  | 1.4 | SSL/TLS Hardening: Automatic bypass of self-signed certs for proxies. | ✅ DONE |
|  | 1.5 | Global Config: Persistent flag management for headers and authentication. | ✅ DONE |
| **Sprint 2: Recon** |
|  | 2.1 | Spec Ingestion: Automated parsing of Swagger (v2) and OpenAPI (v3). | ✅ DONE |
|  | 2.2 | JS Route Scraper: Regex-based endpoint extraction from JS bundles. | ✅ DONE |
|  | 2.3 | Version Walker: Identification of deprecated versions (/v1/ vs /v2/). | ✅ DONE |
|  | 2.4 | Parameter Miner: Automatic identification of hidden query params/headers. | ✅ DONE |
| **Sprint 3: Auth Logic** |
|  | 3.1 | BOLA Prober (API1): Tactical ID-swapping engine with session stores. | ✅ DONE |
|  | 3.2 | BOPLA/Mass Assignment (API3): Fuzzing bodies for hidden properties. | ✅ DONE |
|  | 3.3 | BFLA Module (API5): Hierarchical access testing via method manipulation. | ✅ DONE |
| **Sprint 4: Injection** |
|  | 4.1 | Resource Exhaustion (API4): Probing pagination and payload limits. | ✅ DONE |
|  | 4.2 | SSRF Tracker (API7): Detecting OOB callbacks via URL-parameter injection. | ✅ DONE |
|  | 4.3 | Security Misconfig (API8): Automated CORS and Security Header audit. | ✅ DONE |
|  | 4.4 | Integration Probe (API10): Unsafe consumption in webhooks/3rd party. | ✅ DONE |
| **Sprint 5: Intel** |
|  | 5.1 | SQLite Persistence: Local-first mission database for session continuity. | ✅ DONE |
|  | 5.2 | Async Log Worker: Non-blocking background commitments of findings. | ✅ DONE |
|  | 5.3 | Classified Reporting: NIST-aligned Markdown/PDF debrief generator. | ✅ DONE |
|  | 5.4 | Database Management: Built-in init_db and reset_db control. | ✅ DONE |
| **Sprint 6: Evasion** |
|  | 6.1 | Header Randomization: Rotating User-Agents and JA3 fingerprints. | ✅ DONE |
|  | 6.2 | IP Rotation: Integration with proxy-chains and Tor. | ✅ DONE |
|  | 6.3 | Timing Attacks: Implementing jitter and "Sleepy Probes" for NHPP. | ✅ DONE |
| **Sprint 7: Flow & Logic** |
|  | 7.1 | Flow Engine Implementation: Command suite, recording, and replay. | ✅ DONE |
|  | 7.2 | State-Machine Mapping: Logical order enforcement & out-of-order testing. | ✅ DONE |
|  | 7.3 | Race Condition Engine: Multi-threaded "Turbo Intruder" probes. | ✅ DONE |
| **Sprint 8: Post-Exfil** |
|  | 8.1 | Discovery Vault: Real-time regex scanning of all responses for secrets. | ✅ DONE |
|  | 8.2 | Cloud Pivot Engine: Interception of IMDS (169.254.169.254) requests. | ✅ DONE |
|  | 8.3 | Ghost-Weaver Agent: OIDC interception and encrypted exfiltration. | ✅ DONE |
|  | 8.4 | NHPP Evasion: Masking data as "Deprecated Dependency" system logs. | ✅ DONE |
|  | 8.5 | OOB Validation: Automated validation for leaked tokens/infrastructure. | ✅ DONE |
| **Sprint 9: Hardening** |
|  | 9.1 | Report Engine: Refactored NIST generator with Vault integration. | ✅ DONE |
|  | 9.1.1 | Tactical UI: Integrated spinners and real-time feedback tables. | ✅ DONE |
|  | 9.2 | Surgical BOLA: Response Diffing engine to eliminate False Positives. | ✅ DONE |
|  | 9.3 | Concurrency Engine: High-speed channel-based worker pools. | ✅ DONE |
|  | 9.4 | Environment Sensing: Burp/ZAP detection with X-Header signaling. | ✅ DONE |
|  | 9.5 | Discovery-to-Engine: Automating map-to-scan handover pipeline. | ✅ DONE |
|  | 9.6 | Universal Proxy: Refactored SafeDo with multi-module mirroring. | ✅ DONE |
|  | 9.7 | BOLA Concurrency: Multi-threaded mass scanner upgrade. | ✅ DONE |
|  | 9.8 | Industrialized BOPLA: Concurrent JSON property fuzzing. | ✅ DONE |
|  | 9.9 | Industrialized BFLA: Method Matrix worker pool (Verb-Tampering). | ✅ DONE |
|  | 9.10 | Universal Concurrency: GenericExecutor standardization. | ✅ DONE |
|  | 9.11 | Ghost Masquerade: Process renaming to kworker_system_auth. | ✅ DONE |
|  | 9.13 | Refactor: Framework-Tagged DB (OWASP/MITRE/NIST) Integration | ✅ DONE |

### **Part II: The Hydra TUI & Autonomous Systems [STABLE]**

| Phase | Sub-Phase | Focus / Technical Deliverable | Status |
| --- | --- | --- | --- |
| **Sprint 10: Hydra** |
|  | 10.1 | Universal Target Function (Global Context) | ✅ DONE |
|  | 10.2 | Project Mosaic: The Hydra-TUI Dashboard | ✅ DONE |
|  | 10.2.1 | Terminal Multi-Pane (Quadrants + F-Tabs Switcher) | ✅ DONE |
|  | 10.2.2 | Legacy Shell Fallback (CLI Flag Logic) | ✅ DONE |
|  | 10.3 | Contextual Aggregator & Information Gathering | ✅ DONE |
|  | 10.4 | Tactical Interceptor (F2 Modal Manipulation) | ✅ DONE |
|  | 10.5 | AI Base Integration (Heuristic Brain) | ✅ DONE |
|  | 10.6 | AI Payload Generation & Autonomous Fuzzing | ✅ DONE |

### **Part III: The Future Evolution [ACTIVE-DEV]**

| Phase | Sub-Phase | Focus / Technical Deliverable | Status |
| --- | --- | --- | --- |
| **Sprint 11: Autonomy** |  
|  | **11.1** | **Dynamic Dependency Injection (DDI)** | ✅ DONE |
|  | 11.2 | State-Machine driven payload selection | ❌ ACTIVE |
|  | 11.3 | Autonomous lateral movement within API subnets | ❌ [NEW] |
| **Sprint 12: Evasion V2** |
|  |  | 12.1 | Deep Traffic Shaping: Mimicking legitimate API traffic |
|  | 12.2 | Encrypted OOB: Secure exfiltration via custom protocols | ❌ [NEW] |
|  | 12.3 | Behavioral Jitter: Randomized inter-packet timing | ❌ [NEW] |
| **Sprint 13: The Hive** |
|  |  | 13.1 | Hybrid C2 Architecture: gRPC Control Plane |
|  | 13.2 | RESTful Management API for the Hive Master | ❌ [NEW] |
|  | 13.3 | VaporTrace Console: Web-based Mission Dashboard | ❌ [NEW] |
| **Sprint 14: Pivot** |  
|  | 14.1 | Cross-Tenant Leakage: Exploiting shared infrastructure |
|  | 14.2 | K8s Escape: API-to-Cluster orchestration pivoting | ❌ [NEW] |
|  | 14.3 | Serverless Poisoning: Attacking Lambda/Cloud-Function logic | ❌ [NEW] |
| **Sprint 15: Mastery** | 
|  | 15.1 | Post-Quantum Cryptography for NHPP |
|  | 15.2 | Multi-Agent Swarm Logic (Coordinated BOLA) | ❌ [NEW] |
| **Sprint 15: BlueTeam** | Autonomous Heuristic Remediation: 
|  | 16.1 | Development of a "Blue-Team Mirror" that uses the AI Brain to suggest specific code-level middleware fixes for discovered BOLA/BFLA vulnerabilities. ❌ [NEW] |