# VAPORTRACE TACTICAL AUDIT REPORT
## CONFIDENTIAL - FOR INTERNAL USE ONLY

| METADATA | VALUE |
| :--- | :--- |
| **AUDIT STATUS** | COMPLETED |
| **MISSION START** | 2026-01-31 10:34:15 |
| **GEN TIME (UTC)** | 2026-01-31 10:34:23 |
| **CLASSIFICATION** | PROPRIETARY / ADVERSARY EMULATION |

---

## 1. EXECUTIVE SUMMARY
This document provides a formal tactical debrief of offensive security operations. Results are mapped to the **MITRE ATT&CK** and **OWASP API 2023** frameworks to facilitate risk-based remediation.

## 2. DETAILED TACTICAL FINDINGS

### PHASE: I. INFIL: 2.1 OpenAPI
| TARGET | MITRE ID | OWASP ID | NIST TAG | STATUS | OBSERVATION |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `https://prod-api.target.com/v1/swagger.json?id=0` | T1595.002 | API9:2023 | ID.RA | **EXPLOITED** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://dev-cluster.target.com/v1/swagger.json?id=100` | T1595.002 | API9:2023 | ID.RA | **VULNERABLE** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://stg-nodes.target.com/v1/swagger.json?id=200` | T1595.002 | API9:2023 | ID.RA | **CRITICAL** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://legacy-v1.target.com/v1/swagger.json?id=300` | T1595.002 | API9:2023 | ID.RA | **ACTIVE** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://edge-gateway.target.com/v1/swagger.json?id=400` | T1595.002 | API9:2023 | ID.RA | **INFO** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://prod-api.target.com/v1/swagger.json?id=500` | T1595.002 | API9:2023 | ID.RA | **EXPLOITED** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://dev-cluster.target.com/v1/swagger.json?id=600` | T1595.002 | API9:2023 | ID.RA | **VULNERABLE** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://stg-nodes.target.com/v1/swagger.json?id=700` | T1595.002 | API9:2023 | ID.RA | **CRITICAL** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://prod-api.target.com/v1/swagger.json?id=0` | T1595.002 | API9:2023 | ID.RA | **EXPLOITED** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://dev-cluster.target.com/v1/swagger.json?id=100` | T1595.002 | API9:2023 | ID.RA | **VULNERABLE** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://stg-nodes.target.com/v1/swagger.json?id=200` | T1595.002 | API9:2023 | ID.RA | **CRITICAL** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://legacy-v1.target.com/v1/swagger.json?id=300` | T1595.002 | API9:2023 | ID.RA | **ACTIVE** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://edge-gateway.target.com/v1/swagger.json?id=400` | T1595.002 | API9:2023 | ID.RA | **INFO** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://prod-api.target.com/v1/swagger.json?id=500` | T1595.002 | API9:2023 | ID.RA | **EXPLOITED** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://dev-cluster.target.com/v1/swagger.json?id=600` | T1595.002 | API9:2023 | ID.RA | **VULNERABLE** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://stg-nodes.target.com/v1/swagger.json?id=700` | T1595.002 | API9:2023 | ID.RA | **CRITICAL** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://prod-api.target.com/v1/swagger.json?id=0` | T1595.002 | API9:2023 | ID.RA | **EXPLOITED** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://dev-cluster.target.com/v1/swagger.json?id=100` | T1595.002 | API9:2023 | ID.RA | **VULNERABLE** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://stg-nodes.target.com/v1/swagger.json?id=200` | T1595.002 | API9:2023 | ID.RA | **CRITICAL** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://legacy-v1.target.com/v1/swagger.json?id=300` | T1595.002 | API9:2023 | ID.RA | **ACTIVE** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://edge-gateway.target.com/v1/swagger.json?id=400` | T1595.002 | API9:2023 | ID.RA | **INFO** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://prod-api.target.com/v1/swagger.json?id=500` | T1595.002 | API9:2023 | ID.RA | **EXPLOITED** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://dev-cluster.target.com/v1/swagger.json?id=600` | T1595.002 | API9:2023 | ID.RA | **VULNERABLE** | Shadow API Discovery: Hidden /internal/debug identified. |
| `https://stg-nodes.target.com/v1/swagger.json?id=700` | T1595.002 | API9:2023 | ID.RA | **CRITICAL** | Shadow API Discovery: Hidden /internal/debug identified. |

### PHASE: I. INFIL: 2.2 JS Mining
| TARGET | MITRE ID | OWASP ID | NIST TAG | STATUS | OBSERVATION |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `https://prod-api.target.commain.bundle.js?id=0` | T1592 | API9:2023 | ID.RA | **EXPLOITED** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://prod-api.target.comvendor.js?id=0` | T1592 | API2:2023 | ID.RA | **EXPLOITED** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://dev-cluster.target.commain.bundle.js?id=100` | T1592 | API9:2023 | ID.RA | **VULNERABLE** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://dev-cluster.target.comvendor.js?id=100` | T1592 | API2:2023 | ID.RA | **VULNERABLE** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://stg-nodes.target.commain.bundle.js?id=200` | T1592 | API9:2023 | ID.RA | **CRITICAL** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://stg-nodes.target.comvendor.js?id=200` | T1592 | API2:2023 | ID.RA | **CRITICAL** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://legacy-v1.target.commain.bundle.js?id=300` | T1592 | API9:2023 | ID.RA | **ACTIVE** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://legacy-v1.target.comvendor.js?id=300` | T1592 | API2:2023 | ID.RA | **ACTIVE** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://edge-gateway.target.commain.bundle.js?id=400` | T1592 | API9:2023 | ID.RA | **INFO** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://edge-gateway.target.comvendor.js?id=400` | T1592 | API2:2023 | ID.RA | **INFO** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://prod-api.target.commain.bundle.js?id=500` | T1592 | API9:2023 | ID.RA | **EXPLOITED** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://prod-api.target.comvendor.js?id=500` | T1592 | API2:2023 | ID.RA | **EXPLOITED** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://dev-cluster.target.commain.bundle.js?id=600` | T1592 | API9:2023 | ID.RA | **VULNERABLE** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://dev-cluster.target.comvendor.js?id=600` | T1592 | API2:2023 | ID.RA | **VULNERABLE** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://stg-nodes.target.commain.bundle.js?id=700` | T1592 | API9:2023 | ID.RA | **CRITICAL** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://stg-nodes.target.comvendor.js?id=700` | T1592 | API2:2023 | ID.RA | **CRITICAL** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://prod-api.target.commain.bundle.js?id=0` | T1592 | API9:2023 | ID.RA | **EXPLOITED** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://prod-api.target.comvendor.js?id=0` | T1592 | API2:2023 | ID.RA | **EXPLOITED** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://dev-cluster.target.commain.bundle.js?id=100` | T1592 | API9:2023 | ID.RA | **VULNERABLE** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://dev-cluster.target.comvendor.js?id=100` | T1592 | API2:2023 | ID.RA | **VULNERABLE** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://stg-nodes.target.commain.bundle.js?id=200` | T1592 | API9:2023 | ID.RA | **CRITICAL** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://stg-nodes.target.comvendor.js?id=200` | T1592 | API2:2023 | ID.RA | **CRITICAL** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://legacy-v1.target.commain.bundle.js?id=300` | T1592 | API9:2023 | ID.RA | **ACTIVE** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://legacy-v1.target.comvendor.js?id=300` | T1592 | API2:2023 | ID.RA | **ACTIVE** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://edge-gateway.target.commain.bundle.js?id=400` | T1592 | API9:2023 | ID.RA | **INFO** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://edge-gateway.target.comvendor.js?id=400` | T1592 | API2:2023 | ID.RA | **INFO** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://prod-api.target.commain.bundle.js?id=500` | T1592 | API9:2023 | ID.RA | **EXPLOITED** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://prod-api.target.comvendor.js?id=500` | T1592 | API2:2023 | ID.RA | **EXPLOITED** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://dev-cluster.target.commain.bundle.js?id=600` | T1592 | API9:2023 | ID.RA | **VULNERABLE** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://dev-cluster.target.comvendor.js?id=600` | T1592 | API2:2023 | ID.RA | **VULNERABLE** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://stg-nodes.target.commain.bundle.js?id=700` | T1592 | API9:2023 | ID.RA | **CRITICAL** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://stg-nodes.target.comvendor.js?id=700` | T1592 | API2:2023 | ID.RA | **CRITICAL** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://prod-api.target.commain.bundle.js?id=0` | T1592 | API9:2023 | ID.RA | **EXPLOITED** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://prod-api.target.comvendor.js?id=0` | T1592 | API2:2023 | ID.RA | **EXPLOITED** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://dev-cluster.target.commain.bundle.js?id=100` | T1592 | API9:2023 | ID.RA | **VULNERABLE** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://dev-cluster.target.comvendor.js?id=100` | T1592 | API2:2023 | ID.RA | **VULNERABLE** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://stg-nodes.target.commain.bundle.js?id=200` | T1592 | API9:2023 | ID.RA | **CRITICAL** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://stg-nodes.target.comvendor.js?id=200` | T1592 | API2:2023 | ID.RA | **CRITICAL** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://legacy-v1.target.commain.bundle.js?id=300` | T1592 | API9:2023 | ID.RA | **ACTIVE** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://legacy-v1.target.comvendor.js?id=300` | T1592 | API2:2023 | ID.RA | **ACTIVE** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://edge-gateway.target.commain.bundle.js?id=400` | T1592 | API9:2023 | ID.RA | **INFO** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://edge-gateway.target.comvendor.js?id=400` | T1592 | API2:2023 | ID.RA | **INFO** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://prod-api.target.commain.bundle.js?id=500` | T1592 | API9:2023 | ID.RA | **EXPLOITED** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://prod-api.target.comvendor.js?id=500` | T1592 | API2:2023 | ID.RA | **EXPLOITED** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://dev-cluster.target.commain.bundle.js?id=600` | T1592 | API9:2023 | ID.RA | **VULNERABLE** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://dev-cluster.target.comvendor.js?id=600` | T1592 | API2:2023 | ID.RA | **VULNERABLE** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |
| `https://stg-nodes.target.commain.bundle.js?id=700` | T1592 | API9:2023 | ID.RA | **CRITICAL** | Hidden Route Extraction: Scraped 14 endpoints from minified source. |
| `https://stg-nodes.target.comvendor.js?id=700` | T1592 | API2:2023 | ID.RA | **CRITICAL** | Credential Leak: Found hardcoded Stripe 'pk_test' key. |

### PHASE: I. INFIL: 3.1 Brute-force
| TARGET | MITRE ID | OWASP ID | NIST TAG | STATUS | OBSERVATION |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `https://prod-api.target.com/api/v0/auth?id=0` | T1589 | - | ID.AM | **EXPLOITED** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://dev-cluster.target.com/api/v0/auth?id=100` | T1589 | - | ID.AM | **VULNERABLE** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://stg-nodes.target.com/api/v0/auth?id=200` | T1589 | - | ID.AM | **CRITICAL** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://legacy-v1.target.com/api/v0/auth?id=300` | T1589 | - | ID.AM | **ACTIVE** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://edge-gateway.target.com/api/v0/auth?id=400` | T1589 | - | ID.AM | **INFO** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://prod-api.target.com/api/v0/auth?id=500` | T1589 | - | ID.AM | **EXPLOITED** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://dev-cluster.target.com/api/v0/auth?id=600` | T1589 | - | ID.AM | **VULNERABLE** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://stg-nodes.target.com/api/v0/auth?id=700` | T1589 | - | ID.AM | **CRITICAL** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://prod-api.target.com/api/v0/auth?id=0` | T1589 | - | ID.AM | **EXPLOITED** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://dev-cluster.target.com/api/v0/auth?id=100` | T1589 | - | ID.AM | **VULNERABLE** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://stg-nodes.target.com/api/v0/auth?id=200` | T1589 | - | ID.AM | **CRITICAL** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://legacy-v1.target.com/api/v0/auth?id=300` | T1589 | - | ID.AM | **ACTIVE** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://edge-gateway.target.com/api/v0/auth?id=400` | T1589 | - | ID.AM | **INFO** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://prod-api.target.com/api/v0/auth?id=500` | T1589 | - | ID.AM | **EXPLOITED** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://dev-cluster.target.com/api/v0/auth?id=600` | T1589 | - | ID.AM | **VULNERABLE** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://stg-nodes.target.com/api/v0/auth?id=700` | T1589 | - | ID.AM | **CRITICAL** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://prod-api.target.com/api/v0/auth?id=0` | T1589 | - | ID.AM | **EXPLOITED** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://dev-cluster.target.com/api/v0/auth?id=100` | T1589 | - | ID.AM | **VULNERABLE** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://stg-nodes.target.com/api/v0/auth?id=200` | T1589 | - | ID.AM | **CRITICAL** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://legacy-v1.target.com/api/v0/auth?id=300` | T1589 | - | ID.AM | **ACTIVE** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://edge-gateway.target.com/api/v0/auth?id=400` | T1589 | - | ID.AM | **INFO** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://prod-api.target.com/api/v0/auth?id=500` | T1589 | - | ID.AM | **EXPLOITED** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://dev-cluster.target.com/api/v0/auth?id=600` | T1589 | - | ID.AM | **VULNERABLE** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |
| `https://stg-nodes.target.com/api/v0/auth?id=700` | T1589 | - | ID.AM | **CRITICAL** | Legacy Version ID: Deprecated auth route accessible via version fuzzing. |

### PHASE: II. EXPLOIT: 4.1 BOLA
| TARGET | MITRE ID | OWASP ID | NIST TAG | STATUS | OBSERVATION |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `https://prod-api.target.com/api/orders/5001?id=0` | T1548 | API1:2023 | PR.AC | **EXPLOITED** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://dev-cluster.target.com/api/orders/5001?id=100` | T1548 | API1:2023 | PR.AC | **VULNERABLE** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://stg-nodes.target.com/api/orders/5001?id=200` | T1548 | API1:2023 | PR.AC | **CRITICAL** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://legacy-v1.target.com/api/orders/5001?id=300` | T1548 | API1:2023 | PR.AC | **ACTIVE** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://edge-gateway.target.com/api/orders/5001?id=400` | T1548 | API1:2023 | PR.AC | **INFO** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://prod-api.target.com/api/orders/5001?id=500` | T1548 | API1:2023 | PR.AC | **EXPLOITED** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://dev-cluster.target.com/api/orders/5001?id=600` | T1548 | API1:2023 | PR.AC | **VULNERABLE** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://stg-nodes.target.com/api/orders/5001?id=700` | T1548 | API1:2023 | PR.AC | **CRITICAL** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://prod-api.target.com/api/orders/5001?id=0` | T1548 | API1:2023 | PR.AC | **EXPLOITED** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://dev-cluster.target.com/api/orders/5001?id=100` | T1548 | API1:2023 | PR.AC | **VULNERABLE** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://stg-nodes.target.com/api/orders/5001?id=200` | T1548 | API1:2023 | PR.AC | **CRITICAL** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://legacy-v1.target.com/api/orders/5001?id=300` | T1548 | API1:2023 | PR.AC | **ACTIVE** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://edge-gateway.target.com/api/orders/5001?id=400` | T1548 | API1:2023 | PR.AC | **INFO** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://prod-api.target.com/api/orders/5001?id=500` | T1548 | API1:2023 | PR.AC | **EXPLOITED** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://dev-cluster.target.com/api/orders/5001?id=600` | T1548 | API1:2023 | PR.AC | **VULNERABLE** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://stg-nodes.target.com/api/orders/5001?id=700` | T1548 | API1:2023 | PR.AC | **CRITICAL** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://prod-api.target.com/api/orders/5001?id=0` | T1548 | API1:2023 | PR.AC | **EXPLOITED** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://dev-cluster.target.com/api/orders/5001?id=100` | T1548 | API1:2023 | PR.AC | **VULNERABLE** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://stg-nodes.target.com/api/orders/5001?id=200` | T1548 | API1:2023 | PR.AC | **CRITICAL** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://legacy-v1.target.com/api/orders/5001?id=300` | T1548 | API1:2023 | PR.AC | **ACTIVE** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://edge-gateway.target.com/api/orders/5001?id=400` | T1548 | API1:2023 | PR.AC | **INFO** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://prod-api.target.com/api/orders/5001?id=500` | T1548 | API1:2023 | PR.AC | **EXPLOITED** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://dev-cluster.target.com/api/orders/5001?id=600` | T1548 | API1:2023 | PR.AC | **VULNERABLE** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |
| `https://stg-nodes.target.com/api/orders/5001?id=700` | T1548 | API1:2023 | PR.AC | **CRITICAL** | Unauthorized Data Access: Accessed Order 5001 (User B) as User A. |

### PHASE: II. EXPLOIT: 5.1 BFLA
| TARGET | MITRE ID | OWASP ID | NIST TAG | STATUS | OBSERVATION |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `https://prod-api.target.com/api/system/reboot?id=0` | T1548.002 | API5:2023 | PR.AC | **EXPLOITED** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://dev-cluster.target.com/api/system/reboot?id=100` | T1548.002 | API5:2023 | PR.AC | **VULNERABLE** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://stg-nodes.target.com/api/system/reboot?id=200` | T1548.002 | API5:2023 | PR.AC | **CRITICAL** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://legacy-v1.target.com/api/system/reboot?id=300` | T1548.002 | API5:2023 | PR.AC | **ACTIVE** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://edge-gateway.target.com/api/system/reboot?id=400` | T1548.002 | API5:2023 | PR.AC | **INFO** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://prod-api.target.com/api/system/reboot?id=500` | T1548.002 | API5:2023 | PR.AC | **EXPLOITED** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://dev-cluster.target.com/api/system/reboot?id=600` | T1548.002 | API5:2023 | PR.AC | **VULNERABLE** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://stg-nodes.target.com/api/system/reboot?id=700` | T1548.002 | API5:2023 | PR.AC | **CRITICAL** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://prod-api.target.com/api/system/reboot?id=0` | T1548.002 | API5:2023 | PR.AC | **EXPLOITED** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://dev-cluster.target.com/api/system/reboot?id=100` | T1548.002 | API5:2023 | PR.AC | **VULNERABLE** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://stg-nodes.target.com/api/system/reboot?id=200` | T1548.002 | API5:2023 | PR.AC | **CRITICAL** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://legacy-v1.target.com/api/system/reboot?id=300` | T1548.002 | API5:2023 | PR.AC | **ACTIVE** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://edge-gateway.target.com/api/system/reboot?id=400` | T1548.002 | API5:2023 | PR.AC | **INFO** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://prod-api.target.com/api/system/reboot?id=500` | T1548.002 | API5:2023 | PR.AC | **EXPLOITED** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://dev-cluster.target.com/api/system/reboot?id=600` | T1548.002 | API5:2023 | PR.AC | **VULNERABLE** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://stg-nodes.target.com/api/system/reboot?id=700` | T1548.002 | API5:2023 | PR.AC | **CRITICAL** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://prod-api.target.com/api/system/reboot?id=0` | T1548.002 | API5:2023 | PR.AC | **EXPLOITED** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://dev-cluster.target.com/api/system/reboot?id=100` | T1548.002 | API5:2023 | PR.AC | **VULNERABLE** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://stg-nodes.target.com/api/system/reboot?id=200` | T1548.002 | API5:2023 | PR.AC | **CRITICAL** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://legacy-v1.target.com/api/system/reboot?id=300` | T1548.002 | API5:2023 | PR.AC | **ACTIVE** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://edge-gateway.target.com/api/system/reboot?id=400` | T1548.002 | API5:2023 | PR.AC | **INFO** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://prod-api.target.com/api/system/reboot?id=500` | T1548.002 | API5:2023 | PR.AC | **EXPLOITED** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://dev-cluster.target.com/api/system/reboot?id=600` | T1548.002 | API5:2023 | PR.AC | **VULNERABLE** | Administrative Escalation: Standard user triggered restricted system action. |
| `https://stg-nodes.target.com/api/system/reboot?id=700` | T1548.002 | API5:2023 | PR.AC | **CRITICAL** | Administrative Escalation: Standard user triggered restricted system action. |

### PHASE: II. EXPLOIT: 5.2 BOPLA
| TARGET | MITRE ID | OWASP ID | NIST TAG | STATUS | OBSERVATION |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `https://prod-api.target.com/api/v2/profile?id=0` | T1496 | API6:2023 | PR.DS | **EXPLOITED** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://dev-cluster.target.com/api/v2/profile?id=100` | T1496 | API6:2023 | PR.DS | **VULNERABLE** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://stg-nodes.target.com/api/v2/profile?id=200` | T1496 | API6:2023 | PR.DS | **CRITICAL** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://legacy-v1.target.com/api/v2/profile?id=300` | T1496 | API6:2023 | PR.DS | **ACTIVE** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://edge-gateway.target.com/api/v2/profile?id=400` | T1496 | API6:2023 | PR.DS | **INFO** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://prod-api.target.com/api/v2/profile?id=500` | T1496 | API6:2023 | PR.DS | **EXPLOITED** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://dev-cluster.target.com/api/v2/profile?id=600` | T1496 | API6:2023 | PR.DS | **VULNERABLE** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://stg-nodes.target.com/api/v2/profile?id=700` | T1496 | API6:2023 | PR.DS | **CRITICAL** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://prod-api.target.com/api/v2/profile?id=0` | T1496 | API6:2023 | PR.DS | **EXPLOITED** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://dev-cluster.target.com/api/v2/profile?id=100` | T1496 | API6:2023 | PR.DS | **VULNERABLE** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://stg-nodes.target.com/api/v2/profile?id=200` | T1496 | API6:2023 | PR.DS | **CRITICAL** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://legacy-v1.target.com/api/v2/profile?id=300` | T1496 | API6:2023 | PR.DS | **ACTIVE** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://edge-gateway.target.com/api/v2/profile?id=400` | T1496 | API6:2023 | PR.DS | **INFO** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://prod-api.target.com/api/v2/profile?id=500` | T1496 | API6:2023 | PR.DS | **EXPLOITED** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://dev-cluster.target.com/api/v2/profile?id=600` | T1496 | API6:2023 | PR.DS | **VULNERABLE** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://stg-nodes.target.com/api/v2/profile?id=700` | T1496 | API6:2023 | PR.DS | **CRITICAL** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://prod-api.target.com/api/v2/profile?id=0` | T1496 | API6:2023 | PR.DS | **EXPLOITED** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://dev-cluster.target.com/api/v2/profile?id=100` | T1496 | API6:2023 | PR.DS | **VULNERABLE** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://stg-nodes.target.com/api/v2/profile?id=200` | T1496 | API6:2023 | PR.DS | **CRITICAL** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://legacy-v1.target.com/api/v2/profile?id=300` | T1496 | API6:2023 | PR.DS | **ACTIVE** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://edge-gateway.target.com/api/v2/profile?id=400` | T1496 | API6:2023 | PR.DS | **INFO** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://prod-api.target.com/api/v2/profile?id=500` | T1496 | API6:2023 | PR.DS | **EXPLOITED** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://dev-cluster.target.com/api/v2/profile?id=600` | T1496 | API6:2023 | PR.DS | **VULNERABLE** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |
| `https://stg-nodes.target.com/api/v2/profile?id=700` | T1496 | API6:2023 | PR.DS | **CRITICAL** | Internal State Injection: Injected 'tier: platinum' via mass assignment. |

### PHASE: II. EXPLOIT: 6.1 JWT
| TARGET | MITRE ID | OWASP ID | NIST TAG | STATUS | OBSERVATION |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `https://prod-api.target.comX-Auth-Token?id=0` | T1606 | API2:2023 | PR.AC | **EXPLOITED** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://dev-cluster.target.comX-Auth-Token?id=100` | T1606 | API2:2023 | PR.AC | **VULNERABLE** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://stg-nodes.target.comX-Auth-Token?id=200` | T1606 | API2:2023 | PR.AC | **CRITICAL** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://legacy-v1.target.comX-Auth-Token?id=300` | T1606 | API2:2023 | PR.AC | **ACTIVE** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://edge-gateway.target.comX-Auth-Token?id=400` | T1606 | API2:2023 | PR.AC | **INFO** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://prod-api.target.comX-Auth-Token?id=500` | T1606 | API2:2023 | PR.AC | **EXPLOITED** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://dev-cluster.target.comX-Auth-Token?id=600` | T1606 | API2:2023 | PR.AC | **VULNERABLE** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://stg-nodes.target.comX-Auth-Token?id=700` | T1606 | API2:2023 | PR.AC | **CRITICAL** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://prod-api.target.comX-Auth-Token?id=0` | T1606 | API2:2023 | PR.AC | **EXPLOITED** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://dev-cluster.target.comX-Auth-Token?id=100` | T1606 | API2:2023 | PR.AC | **VULNERABLE** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://stg-nodes.target.comX-Auth-Token?id=200` | T1606 | API2:2023 | PR.AC | **CRITICAL** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://legacy-v1.target.comX-Auth-Token?id=300` | T1606 | API2:2023 | PR.AC | **ACTIVE** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://edge-gateway.target.comX-Auth-Token?id=400` | T1606 | API2:2023 | PR.AC | **INFO** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://prod-api.target.comX-Auth-Token?id=500` | T1606 | API2:2023 | PR.AC | **EXPLOITED** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://dev-cluster.target.comX-Auth-Token?id=600` | T1606 | API2:2023 | PR.AC | **VULNERABLE** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://stg-nodes.target.comX-Auth-Token?id=700` | T1606 | API2:2023 | PR.AC | **CRITICAL** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://prod-api.target.comX-Auth-Token?id=0` | T1606 | API2:2023 | PR.AC | **EXPLOITED** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://dev-cluster.target.comX-Auth-Token?id=100` | T1606 | API2:2023 | PR.AC | **VULNERABLE** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://stg-nodes.target.comX-Auth-Token?id=200` | T1606 | API2:2023 | PR.AC | **CRITICAL** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://legacy-v1.target.comX-Auth-Token?id=300` | T1606 | API2:2023 | PR.AC | **ACTIVE** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://edge-gateway.target.comX-Auth-Token?id=400` | T1606 | API2:2023 | PR.AC | **INFO** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://prod-api.target.comX-Auth-Token?id=500` | T1606 | API2:2023 | PR.AC | **EXPLOITED** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://dev-cluster.target.comX-Auth-Token?id=600` | T1606 | API2:2023 | PR.AC | **VULNERABLE** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |
| `https://stg-nodes.target.comX-Auth-Token?id=700` | T1606 | API2:2023 | PR.AC | **CRITICAL** | Identity Spoofing: Successfully forged admin token using 'none' algorithm. |

### PHASE: III. EXPAND: 7.1 SSRF
| TARGET | MITRE ID | OWASP ID | NIST TAG | STATUS | OBSERVATION |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `https://prod-api.target.com169.254.169.254?id=0` | T1046 | API7:2023 | DE.CM | **EXPLOITED** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://dev-cluster.target.com169.254.169.254?id=100` | T1046 | API7:2023 | DE.CM | **VULNERABLE** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://stg-nodes.target.com169.254.169.254?id=200` | T1046 | API7:2023 | DE.CM | **CRITICAL** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://legacy-v1.target.com169.254.169.254?id=300` | T1046 | API7:2023 | DE.CM | **ACTIVE** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://edge-gateway.target.com169.254.169.254?id=400` | T1046 | API7:2023 | DE.CM | **INFO** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://prod-api.target.com169.254.169.254?id=500` | T1046 | API7:2023 | DE.CM | **EXPLOITED** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://dev-cluster.target.com169.254.169.254?id=600` | T1046 | API7:2023 | DE.CM | **VULNERABLE** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://stg-nodes.target.com169.254.169.254?id=700` | T1046 | API7:2023 | DE.CM | **CRITICAL** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://prod-api.target.com169.254.169.254?id=0` | T1046 | API7:2023 | DE.CM | **EXPLOITED** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://dev-cluster.target.com169.254.169.254?id=100` | T1046 | API7:2023 | DE.CM | **VULNERABLE** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://stg-nodes.target.com169.254.169.254?id=200` | T1046 | API7:2023 | DE.CM | **CRITICAL** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://legacy-v1.target.com169.254.169.254?id=300` | T1046 | API7:2023 | DE.CM | **ACTIVE** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://edge-gateway.target.com169.254.169.254?id=400` | T1046 | API7:2023 | DE.CM | **INFO** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://prod-api.target.com169.254.169.254?id=500` | T1046 | API7:2023 | DE.CM | **EXPLOITED** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://dev-cluster.target.com169.254.169.254?id=600` | T1046 | API7:2023 | DE.CM | **VULNERABLE** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://stg-nodes.target.com169.254.169.254?id=700` | T1046 | API7:2023 | DE.CM | **CRITICAL** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://prod-api.target.com169.254.169.254?id=0` | T1046 | API7:2023 | DE.CM | **EXPLOITED** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://dev-cluster.target.com169.254.169.254?id=100` | T1046 | API7:2023 | DE.CM | **VULNERABLE** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://stg-nodes.target.com169.254.169.254?id=200` | T1046 | API7:2023 | DE.CM | **CRITICAL** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://legacy-v1.target.com169.254.169.254?id=300` | T1046 | API7:2023 | DE.CM | **ACTIVE** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://edge-gateway.target.com169.254.169.254?id=400` | T1046 | API7:2023 | DE.CM | **INFO** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://prod-api.target.com169.254.169.254?id=500` | T1046 | API7:2023 | DE.CM | **EXPLOITED** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://dev-cluster.target.com169.254.169.254?id=600` | T1046 | API7:2023 | DE.CM | **VULNERABLE** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |
| `https://stg-nodes.target.com169.254.169.254?id=700` | T1046 | API7:2023 | DE.CM | **CRITICAL** | Cloud IAM Role Theft: Exfiltrated AWS credentials from metadata service. |

### PHASE: III. EXPAND: 8.1 DoS
| TARGET | MITRE ID | OWASP ID | NIST TAG | STATUS | OBSERVATION |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `https://prod-api.target.com/api/reports/all?id=0` | T1499 | API4:2023 | RS.AN | **EXPLOITED** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://dev-cluster.target.com/api/reports/all?id=100` | T1499 | API4:2023 | RS.AN | **VULNERABLE** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://stg-nodes.target.com/api/reports/all?id=200` | T1499 | API4:2023 | RS.AN | **CRITICAL** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://legacy-v1.target.com/api/reports/all?id=300` | T1499 | API4:2023 | RS.AN | **ACTIVE** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://edge-gateway.target.com/api/reports/all?id=400` | T1499 | API4:2023 | RS.AN | **INFO** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://prod-api.target.com/api/reports/all?id=500` | T1499 | API4:2023 | RS.AN | **EXPLOITED** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://dev-cluster.target.com/api/reports/all?id=600` | T1499 | API4:2023 | RS.AN | **VULNERABLE** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://stg-nodes.target.com/api/reports/all?id=700` | T1499 | API4:2023 | RS.AN | **CRITICAL** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://prod-api.target.com/api/reports/all?id=0` | T1499 | API4:2023 | RS.AN | **EXPLOITED** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://dev-cluster.target.com/api/reports/all?id=100` | T1499 | API4:2023 | RS.AN | **VULNERABLE** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://stg-nodes.target.com/api/reports/all?id=200` | T1499 | API4:2023 | RS.AN | **CRITICAL** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://legacy-v1.target.com/api/reports/all?id=300` | T1499 | API4:2023 | RS.AN | **ACTIVE** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://edge-gateway.target.com/api/reports/all?id=400` | T1499 | API4:2023 | RS.AN | **INFO** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://prod-api.target.com/api/reports/all?id=500` | T1499 | API4:2023 | RS.AN | **EXPLOITED** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://dev-cluster.target.com/api/reports/all?id=600` | T1499 | API4:2023 | RS.AN | **VULNERABLE** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://stg-nodes.target.com/api/reports/all?id=700` | T1499 | API4:2023 | RS.AN | **CRITICAL** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://prod-api.target.com/api/reports/all?id=0` | T1499 | API4:2023 | RS.AN | **EXPLOITED** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://dev-cluster.target.com/api/reports/all?id=100` | T1499 | API4:2023 | RS.AN | **VULNERABLE** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://stg-nodes.target.com/api/reports/all?id=200` | T1499 | API4:2023 | RS.AN | **CRITICAL** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://legacy-v1.target.com/api/reports/all?id=300` | T1499 | API4:2023 | RS.AN | **ACTIVE** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://edge-gateway.target.com/api/reports/all?id=400` | T1499 | API4:2023 | RS.AN | **INFO** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://prod-api.target.com/api/reports/all?id=500` | T1499 | API4:2023 | RS.AN | **EXPLOITED** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://dev-cluster.target.com/api/reports/all?id=600` | T1499 | API4:2023 | RS.AN | **VULNERABLE** | Backend Service Crash: Resource exhaustion via nested JSON payload. |
| `https://stg-nodes.target.com/api/reports/all?id=700` | T1499 | API4:2023 | RS.AN | **CRITICAL** | Backend Service Crash: Resource exhaustion via nested JSON payload. |

### PHASE: III. EXPAND: 9.1 Persist
| TARGET | MITRE ID | OWASP ID | NIST TAG | STATUS | OBSERVATION |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `https://prod-api.target.comMission Database?id=0` | T1560 | - | PR.DS | **EXPLOITED** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://dev-cluster.target.comMission Database?id=100` | T1560 | - | PR.DS | **VULNERABLE** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://stg-nodes.target.comMission Database?id=200` | T1560 | - | PR.DS | **CRITICAL** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://legacy-v1.target.comMission Database?id=300` | T1560 | - | PR.DS | **ACTIVE** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://edge-gateway.target.comMission Database?id=400` | T1560 | - | PR.DS | **INFO** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://prod-api.target.comMission Database?id=500` | T1560 | - | PR.DS | **EXPLOITED** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://dev-cluster.target.comMission Database?id=600` | T1560 | - | PR.DS | **VULNERABLE** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://stg-nodes.target.comMission Database?id=700` | T1560 | - | PR.DS | **CRITICAL** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://prod-api.target.comMission Database?id=0` | T1560 | - | PR.DS | **EXPLOITED** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://dev-cluster.target.comMission Database?id=100` | T1560 | - | PR.DS | **VULNERABLE** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://stg-nodes.target.comMission Database?id=200` | T1560 | - | PR.DS | **CRITICAL** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://legacy-v1.target.comMission Database?id=300` | T1560 | - | PR.DS | **ACTIVE** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://edge-gateway.target.comMission Database?id=400` | T1560 | - | PR.DS | **INFO** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://prod-api.target.comMission Database?id=500` | T1560 | - | PR.DS | **EXPLOITED** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://dev-cluster.target.comMission Database?id=600` | T1560 | - | PR.DS | **VULNERABLE** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://stg-nodes.target.comMission Database?id=700` | T1560 | - | PR.DS | **CRITICAL** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://prod-api.target.comMission Database?id=0` | T1560 | - | PR.DS | **EXPLOITED** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://dev-cluster.target.comMission Database?id=100` | T1560 | - | PR.DS | **VULNERABLE** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://stg-nodes.target.comMission Database?id=200` | T1560 | - | PR.DS | **CRITICAL** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://legacy-v1.target.comMission Database?id=300` | T1560 | - | PR.DS | **ACTIVE** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://edge-gateway.target.comMission Database?id=400` | T1560 | - | PR.DS | **INFO** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://prod-api.target.comMission Database?id=500` | T1560 | - | PR.DS | **EXPLOITED** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://dev-cluster.target.comMission Database?id=600` | T1560 | - | PR.DS | **VULNERABLE** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |
| `https://stg-nodes.target.comMission Database?id=700` | T1560 | - | PR.DS | **CRITICAL** | Audit Trail Integrity: Findings persisted with NIST framework tagging. |

### PHASE: IV. OBFUSC: 11.1 Proxy
| TARGET | MITRE ID | OWASP ID | NIST TAG | STATUS | OBSERVATION |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `https://prod-api.target.com127.0.0.1:8080?id=0` | T1090 | - | PR.PT | **EXPLOITED** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://dev-cluster.target.com127.0.0.1:8080?id=100` | T1090 | - | PR.PT | **VULNERABLE** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://stg-nodes.target.com127.0.0.1:8080?id=200` | T1090 | - | PR.PT | **CRITICAL** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://legacy-v1.target.com127.0.0.1:8080?id=300` | T1090 | - | PR.PT | **ACTIVE** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://edge-gateway.target.com127.0.0.1:8080?id=400` | T1090 | - | PR.PT | **INFO** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://prod-api.target.com127.0.0.1:8080?id=500` | T1090 | - | PR.PT | **EXPLOITED** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://dev-cluster.target.com127.0.0.1:8080?id=600` | T1090 | - | PR.PT | **VULNERABLE** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://stg-nodes.target.com127.0.0.1:8080?id=700` | T1090 | - | PR.PT | **CRITICAL** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://prod-api.target.com127.0.0.1:8080?id=0` | T1090 | - | PR.PT | **EXPLOITED** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://dev-cluster.target.com127.0.0.1:8080?id=100` | T1090 | - | PR.PT | **VULNERABLE** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://stg-nodes.target.com127.0.0.1:8080?id=200` | T1090 | - | PR.PT | **CRITICAL** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://legacy-v1.target.com127.0.0.1:8080?id=300` | T1090 | - | PR.PT | **ACTIVE** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://edge-gateway.target.com127.0.0.1:8080?id=400` | T1090 | - | PR.PT | **INFO** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://prod-api.target.com127.0.0.1:8080?id=500` | T1090 | - | PR.PT | **EXPLOITED** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://dev-cluster.target.com127.0.0.1:8080?id=600` | T1090 | - | PR.PT | **VULNERABLE** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://stg-nodes.target.com127.0.0.1:8080?id=700` | T1090 | - | PR.PT | **CRITICAL** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://prod-api.target.com127.0.0.1:8080?id=0` | T1090 | - | PR.PT | **EXPLOITED** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://dev-cluster.target.com127.0.0.1:8080?id=100` | T1090 | - | PR.PT | **VULNERABLE** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://stg-nodes.target.com127.0.0.1:8080?id=200` | T1090 | - | PR.PT | **CRITICAL** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://legacy-v1.target.com127.0.0.1:8080?id=300` | T1090 | - | PR.PT | **ACTIVE** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://edge-gateway.target.com127.0.0.1:8080?id=400` | T1090 | - | PR.PT | **INFO** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://prod-api.target.com127.0.0.1:8080?id=500` | T1090 | - | PR.PT | **EXPLOITED** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://dev-cluster.target.com127.0.0.1:8080?id=600` | T1090 | - | PR.PT | **VULNERABLE** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |
| `https://stg-nodes.target.com127.0.0.1:8080?id=700` | T1090 | - | PR.PT | **CRITICAL** | Origin IP Masking: Tactical traffic successfully proxied through Burp. |

### PHASE: IV. OBFUSC: 11.2 Rotation
| TARGET | MITRE ID | OWASP ID | NIST TAG | STATUS | OBSERVATION |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `https://prod-api.target.comProxyPool-Alpha?id=0` | T1090.003 | - | PR.PT | **EXPLOITED** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://dev-cluster.target.comProxyPool-Alpha?id=100` | T1090.003 | - | PR.PT | **VULNERABLE** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://stg-nodes.target.comProxyPool-Alpha?id=200` | T1090.003 | - | PR.PT | **CRITICAL** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://legacy-v1.target.comProxyPool-Alpha?id=300` | T1090.003 | - | PR.PT | **ACTIVE** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://edge-gateway.target.comProxyPool-Alpha?id=400` | T1090.003 | - | PR.PT | **INFO** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://prod-api.target.comProxyPool-Alpha?id=500` | T1090.003 | - | PR.PT | **EXPLOITED** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://dev-cluster.target.comProxyPool-Alpha?id=600` | T1090.003 | - | PR.PT | **VULNERABLE** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://stg-nodes.target.comProxyPool-Alpha?id=700` | T1090.003 | - | PR.PT | **CRITICAL** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://prod-api.target.comProxyPool-Alpha?id=0` | T1090.003 | - | PR.PT | **EXPLOITED** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://dev-cluster.target.comProxyPool-Alpha?id=100` | T1090.003 | - | PR.PT | **VULNERABLE** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://stg-nodes.target.comProxyPool-Alpha?id=200` | T1090.003 | - | PR.PT | **CRITICAL** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://legacy-v1.target.comProxyPool-Alpha?id=300` | T1090.003 | - | PR.PT | **ACTIVE** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://edge-gateway.target.comProxyPool-Alpha?id=400` | T1090.003 | - | PR.PT | **INFO** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://prod-api.target.comProxyPool-Alpha?id=500` | T1090.003 | - | PR.PT | **EXPLOITED** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://dev-cluster.target.comProxyPool-Alpha?id=600` | T1090.003 | - | PR.PT | **VULNERABLE** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://stg-nodes.target.comProxyPool-Alpha?id=700` | T1090.003 | - | PR.PT | **CRITICAL** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://prod-api.target.comProxyPool-Alpha?id=0` | T1090.003 | - | PR.PT | **EXPLOITED** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://dev-cluster.target.comProxyPool-Alpha?id=100` | T1090.003 | - | PR.PT | **VULNERABLE** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://stg-nodes.target.comProxyPool-Alpha?id=200` | T1090.003 | - | PR.PT | **CRITICAL** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://legacy-v1.target.comProxyPool-Alpha?id=300` | T1090.003 | - | PR.PT | **ACTIVE** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://edge-gateway.target.comProxyPool-Alpha?id=400` | T1090.003 | - | PR.PT | **INFO** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://prod-api.target.comProxyPool-Alpha?id=500` | T1090.003 | - | PR.PT | **EXPLOITED** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://dev-cluster.target.comProxyPool-Alpha?id=600` | T1090.003 | - | PR.PT | **VULNERABLE** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |
| `https://stg-nodes.target.comProxyPool-Alpha?id=700` | T1090.003 | - | PR.PT | **CRITICAL** | Rate-Limit Bypass: Egress IP rotated 15 times during session. |

### PHASE: IV. OBFUSC: 12.1 Evasion
| TARGET | MITRE ID | OWASP ID | NIST TAG | STATUS | OBSERVATION |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `https://prod-api.target.comCloudflare WAF?id=0` | T1562.001 | - | PR.PT | **EXPLOITED** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://dev-cluster.target.comCloudflare WAF?id=100` | T1562.001 | - | PR.PT | **VULNERABLE** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://stg-nodes.target.comCloudflare WAF?id=200` | T1562.001 | - | PR.PT | **CRITICAL** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://legacy-v1.target.comCloudflare WAF?id=300` | T1562.001 | - | PR.PT | **ACTIVE** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://edge-gateway.target.comCloudflare WAF?id=400` | T1562.001 | - | PR.PT | **INFO** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://prod-api.target.comCloudflare WAF?id=500` | T1562.001 | - | PR.PT | **EXPLOITED** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://dev-cluster.target.comCloudflare WAF?id=600` | T1562.001 | - | PR.PT | **VULNERABLE** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://stg-nodes.target.comCloudflare WAF?id=700` | T1562.001 | - | PR.PT | **CRITICAL** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://prod-api.target.comCloudflare WAF?id=0` | T1562.001 | - | PR.PT | **EXPLOITED** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://dev-cluster.target.comCloudflare WAF?id=100` | T1562.001 | - | PR.PT | **VULNERABLE** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://stg-nodes.target.comCloudflare WAF?id=200` | T1562.001 | - | PR.PT | **CRITICAL** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://legacy-v1.target.comCloudflare WAF?id=300` | T1562.001 | - | PR.PT | **ACTIVE** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://edge-gateway.target.comCloudflare WAF?id=400` | T1562.001 | - | PR.PT | **INFO** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://prod-api.target.comCloudflare WAF?id=500` | T1562.001 | - | PR.PT | **EXPLOITED** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://dev-cluster.target.comCloudflare WAF?id=600` | T1562.001 | - | PR.PT | **VULNERABLE** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://stg-nodes.target.comCloudflare WAF?id=700` | T1562.001 | - | PR.PT | **CRITICAL** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://prod-api.target.comCloudflare WAF?id=0` | T1562.001 | - | PR.PT | **EXPLOITED** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://dev-cluster.target.comCloudflare WAF?id=100` | T1562.001 | - | PR.PT | **VULNERABLE** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://stg-nodes.target.comCloudflare WAF?id=200` | T1562.001 | - | PR.PT | **CRITICAL** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://legacy-v1.target.comCloudflare WAF?id=300` | T1562.001 | - | PR.PT | **ACTIVE** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://edge-gateway.target.comCloudflare WAF?id=400` | T1562.001 | - | PR.PT | **INFO** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://prod-api.target.comCloudflare WAF?id=500` | T1562.001 | - | PR.PT | **EXPLOITED** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://dev-cluster.target.comCloudflare WAF?id=600` | T1562.001 | - | PR.PT | **VULNERABLE** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |
| `https://stg-nodes.target.comCloudflare WAF?id=700` | T1562.001 | - | PR.PT | **CRITICAL** | WAF Signature Evasion: Randomized JA3 fingerprints and headers. |

### PHASE: PHASE 9.9: EXHAUSTION
| TARGET | MITRE ID | OWASP ID | NIST TAG | STATUS | OBSERVATION |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `https://api.target.com/feed` | T1499 | API4:2023 | RS.AN | **VULNERABLE** | DoS: limit=1000000 (Latency: 5s) |

### PHASE: PHASE II: DISCOVERY
| TARGET | MITRE ID | OWASP ID | NIST TAG | STATUS | OBSERVATION |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `https://api.dummy.com/v1` | T1595 | API9:2023 | ID.AM | **INFO** | Swagger Documentation Exposed |
| `https://api.target.com/v1/swagger.json` | T1595 | API9:2023 | ID.AM | **INFO** | Swagger Documentation Exposed |
| `https://api.target.com/v2/api-docs` | T1595 | API9:2023 | ID.AM | **INFO** | OpenAPI v3 Spec Found |
| `https://api.target.com/.env` | T1595 | API9:2023 | ID.AM | **INFO** | Environment File (403 Forbidden) |
| `https://api.target.com` | T1592 | API8:2023 | PR.IP | **WEAK CONFIG** | Missing Header: HSTS |
| `https://api.target.com` | T1592 | API8:2023 | PR.IP | **WEAK CONFIG** | CORS: * (Wildcard) |

### PHASE: PHASE III: AUTH LOGIC
| TARGET | MITRE ID | OWASP ID | NIST TAG | STATUS | OBSERVATION |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `https://api.dummy.com/users/102` | T1548 | API1:2023 | DE.AE | **EXPLOITED** | BOLA ID Swap Success: Accessed User 102 |
| `https://api.target.com/users/102` | T1548 | API1:2023 | DE.AE | **EXPLOITED** | BOLA ID Swap: Accessed User 102 |
| `https://api.target.com/users/103` | T1548 | API1:2023 | DE.AE | **EXPLOITED** | BOLA ID Swap: Accessed User 103 |
| `https://api.target.com/users/admin` | T1548 | API1:2023 | DE.AE | **INFO** | BOLA ID Swap: Failed |
| `https://api.target.com/admin/user` | T1548.002 | API5:2023 | DE.CM | **VULNERABLE** | BFLA: DELETE Method Allowed |
| `https://api.target.com/admin/settings` | T1548.002 | API5:2023 | DE.CM | **VULNERABLE** | BFLA: POST Method Allowed |

### PHASE: PHASE IV: INJECTION
| TARGET | MITRE ID | OWASP ID | NIST TAG | STATUS | OBSERVATION |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `https://api.dummy.com/hooks` | T1071 | API7:2023 | DE.CM | **CRITICAL** | SSRF Internal Access: 169.254.169.254 |
| `https://api.target.com/hooks` | T1071 | API7:2023 | DE.CM | **CRITICAL** | SSRF Internal: 169.254.169.254 |
| `https://api.target.com/callback` | T1071 | API7:2023 | DE.CM | **CRITICAL** | SSRF Internal: 127.0.0.1 |
| `https://api.target.com/img` | T1071 | API7:2023 | DE.CM | **VULNERABLE** | SSRF: Open Redirect to evil.com |
| `https://api.target.com/profile` | T1538 | API3:2023 | PR.AC | **EXPLOITED** | BOPLA: 'is_admin' Injected |
| `https://api.target.com/order` | T1538 | API3:2023 | PR.AC | **EXPLOITED** | BOPLA: 'discount' Injected |

### PHASE: PHASE VIII: EXFILTRATION
| TARGET | MITRE ID | OWASP ID | NIST TAG | STATUS | OBSERVATION |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `https://api.dummy.com/debug` | T1552 | API2:2023 | PR.AC | **VULNERABLE** | Leaked AWS_KEY: AKIA........ |
| `https://api.target.com/debug` | T1552 | API2:2023 | PR.AC | **VULNERABLE** | Leaked AWS_KEY: AKIA........ |
| `https://api.target.com/logs` | T1552 | API2:2023 | PR.AC | **VULNERABLE** | Leaked JWT Token in Body |

### PHASE: V. COMPL: 13.1 Debrief
| TARGET | MITRE ID | OWASP ID | NIST TAG | STATUS | OBSERVATION |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `https://prod-api.target.commission_logs.md?id=0` | T1020 | - | PR.DS | **EXPLOITED** | Evidence Packaging: Automated Markdown report generated. |
| `https://dev-cluster.target.commission_logs.md?id=100` | T1020 | - | PR.DS | **VULNERABLE** | Evidence Packaging: Automated Markdown report generated. |
| `https://stg-nodes.target.commission_logs.md?id=200` | T1020 | - | PR.DS | **CRITICAL** | Evidence Packaging: Automated Markdown report generated. |
| `https://legacy-v1.target.commission_logs.md?id=300` | T1020 | - | PR.DS | **ACTIVE** | Evidence Packaging: Automated Markdown report generated. |
| `https://edge-gateway.target.commission_logs.md?id=400` | T1020 | - | PR.DS | **INFO** | Evidence Packaging: Automated Markdown report generated. |
| `https://prod-api.target.commission_logs.md?id=500` | T1020 | - | PR.DS | **EXPLOITED** | Evidence Packaging: Automated Markdown report generated. |
| `https://dev-cluster.target.commission_logs.md?id=600` | T1020 | - | PR.DS | **VULNERABLE** | Evidence Packaging: Automated Markdown report generated. |
| `https://stg-nodes.target.commission_logs.md?id=700` | T1020 | - | PR.DS | **CRITICAL** | Evidence Packaging: Automated Markdown report generated. |
| `https://prod-api.target.commission_logs.md?id=0` | T1020 | - | PR.DS | **EXPLOITED** | Evidence Packaging: Automated Markdown report generated. |
| `https://dev-cluster.target.commission_logs.md?id=100` | T1020 | - | PR.DS | **VULNERABLE** | Evidence Packaging: Automated Markdown report generated. |
| `https://stg-nodes.target.commission_logs.md?id=200` | T1020 | - | PR.DS | **CRITICAL** | Evidence Packaging: Automated Markdown report generated. |
| `https://legacy-v1.target.commission_logs.md?id=300` | T1020 | - | PR.DS | **ACTIVE** | Evidence Packaging: Automated Markdown report generated. |
| `https://edge-gateway.target.commission_logs.md?id=400` | T1020 | - | PR.DS | **INFO** | Evidence Packaging: Automated Markdown report generated. |
| `https://prod-api.target.commission_logs.md?id=500` | T1020 | - | PR.DS | **EXPLOITED** | Evidence Packaging: Automated Markdown report generated. |
| `https://dev-cluster.target.commission_logs.md?id=600` | T1020 | - | PR.DS | **VULNERABLE** | Evidence Packaging: Automated Markdown report generated. |
| `https://stg-nodes.target.commission_logs.md?id=700` | T1020 | - | PR.DS | **CRITICAL** | Evidence Packaging: Automated Markdown report generated. |
| `https://prod-api.target.commission_logs.md?id=0` | T1020 | - | PR.DS | **EXPLOITED** | Evidence Packaging: Automated Markdown report generated. |
| `https://dev-cluster.target.commission_logs.md?id=100` | T1020 | - | PR.DS | **VULNERABLE** | Evidence Packaging: Automated Markdown report generated. |
| `https://stg-nodes.target.commission_logs.md?id=200` | T1020 | - | PR.DS | **CRITICAL** | Evidence Packaging: Automated Markdown report generated. |
| `https://legacy-v1.target.commission_logs.md?id=300` | T1020 | - | PR.DS | **ACTIVE** | Evidence Packaging: Automated Markdown report generated. |
| `https://edge-gateway.target.commission_logs.md?id=400` | T1020 | - | PR.DS | **INFO** | Evidence Packaging: Automated Markdown report generated. |
| `https://prod-api.target.commission_logs.md?id=500` | T1020 | - | PR.DS | **EXPLOITED** | Evidence Packaging: Automated Markdown report generated. |
| `https://dev-cluster.target.commission_logs.md?id=600` | T1020 | - | PR.DS | **VULNERABLE** | Evidence Packaging: Automated Markdown report generated. |
| `https://stg-nodes.target.commission_logs.md?id=700` | T1020 | - | PR.DS | **CRITICAL** | Evidence Packaging: Automated Markdown report generated. |

---
## 3. FRAMEWORK ALIGNMENT

| VAPORTRACE COMPONENT | MITRE TACTIC | OWASP TOP 10 | DEFENSIVE CONTEXT |
| :--- | :--- | :--- | :--- |
| I. INFIL | Reconnaissance | API9:2023 | Inventory Management |
| II. EXPLOIT | Priv Escalation | API1 / API2 / API5 | Identity Assurance |
| III. EXPAND | Discovery / Impact | API4 / API7 | Resource Hardening |
| IV. OBFUSC | Defense Evasion | N/A | Stealth Analysis |

---
## 4. REMEDIATION STRATEGY

### Immediate Hardening
1. **Inventory Control:** Remove all exposed Swagger documentation identified in Phase I.
2. **Auth Validation:** Patch JWT algorithm processing to deny 'none' or 'weak' signing keys.
3. **Rate Limiting:** Implement circuit-breaking on high-resource endpoints to mitigate DoS findings.

---
**UNAUTHORIZED DISCLOSURE OF THIS REPORT IS PROHIBITED**
