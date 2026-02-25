# VAPORTRACE PENETRATION TEST REPORT
**CONFIDENTIAL - INTERNAL USE ONLY**

| META | VALUE |
| :--- | :--- |
| **DATE** | 2026-02-04 |
| **MISSION START** | 2026-02-04 08:41:56 |
| **CLASSIFICATION** | PROPRIETARY / ADVERSARY EMULATION |
| **ENGINE VERSION** | VaporTrace v3.1 (Tactical Suite) |

---

## 1. EXECUTIVE SUMMARY

### 1.1 Risk Overview
VaporTrace Tactical Suite performed an automated adversarial emulation against the target infrastructure. This section provides a high-level overview of the security posture based on discovered vulnerabilities.

**OVERALL RISK RATING:** MODERATE

| METRIC | VALUE |
| :--- | :--- |
| **Total Findings** | 50 |
| **Unique Targets** | 12 |
| **Average CVSS** | 0.4 / 10.0 |

### 1.2 Vulnerability Distribution
Breakdown of findings by severity (CVSS v3.1):

- **CRITICAL (9.0+):** 0  (░░░░░░░░░░)
- **HIGH (7.0-8.9):**     0  (░░░░░░░░░░)
- **MEDIUM (4.0-6.9):**   4  (░░░░░░░░░░)
- **LOW (0.0-3.9):**      46  (█████████░)

---

## 2. REMEDIATION PRIORITY TRACKER
The following table prioritizes vulnerabilities requiring immediate attention. **Sorted by Severity (CVSS Descending).**

| SEVERITY | CVSS | VULNERABILITY (OWASP) | CVE ID | AFFECTED TARGET | ACTION |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 🟡 | 5.4 | API8 | CVE-2024-AUDIT | `https://petstore.swagger.io/v2/swagger.json` | Remediate < 30 Days |
| 🟡 | 5.4 | API8 | CVE-2024-AUDIT | `https://petstore.swagger.io/v2/swagger.json` | Remediate < 30 Days |
| 🟡 | 5.4 | API8 | CVE-2024-AUDIT | `https://petstore.swagger.io/v2/swagger.json` | Remediate < 30 Days |
| 🟡 | 5.4 | API8 | CVE-2024-AUDIT | `https://petstore.swagger.io/v2/swagger.json` | Remediate < 30 Days |

---

## 3. TECHNICAL FINDINGS (DEEP DIVE)
Detailed evidence logs for engineering teams. **Sorted Chronologically (Execution Order).**

### [INFO] API9:2023 Improper Inventory Management on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:30Z
- **Vector/Command:** `map`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** Swagger/OpenAPI Documentation Found (Deep Parse)

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1595.002` | Reconnaissance |
| **NIST CSF v2.0** | `ID.AM-07` | Control Mapping |
| **CVE / CVSS** | `N/A` | **0.0** (Severity Score) |

---
### [EXPLOITED]  on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:30Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** Leaked EMAIL: apiteam@swagger.io

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` |  |
| **NIST CSF v2.0** | `` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [VULNERABLE] API3:2023 on https://petstore.swagger.io?admin=true
- **Timestamp:** 2026-02-04T11:43:41Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io?admin=true`
- **Details:** Potential Hidden Parameter: admin

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1596` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [VULNERABLE] API3:2023 on https://petstore.swagger.io?internal=true
- **Timestamp:** 2026-02-04T11:43:41Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io?internal=true`
- **Details:** Potential Hidden Parameter: internal

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1596` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [VULNERABLE] API3:2023 on https://petstore.swagger.io?config=true
- **Timestamp:** 2026-02-04T11:43:41Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io?config=true`
- **Details:** Potential Hidden Parameter: config

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1596` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [VULNERABLE] API3:2023 on https://petstore.swagger.io?role=true
- **Timestamp:** 2026-02-04T11:43:42Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io?role=true`
- **Details:** Potential Hidden Parameter: role

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1596` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [VULNERABLE] API3:2023 on https://petstore.swagger.io/v2/swagger.json?debug=true
- **Timestamp:** 2026-02-04T11:43:50Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json?debug=true`
- **Details:** Potential Hidden Parameter: debug

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1596` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [EXPLOITED]  on https://petstore.swagger.io/v2/swagger.json?debug=true
- **Timestamp:** 2026-02-04T11:43:50Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json?debug=true`
- **Details:** Leaked EMAIL: apiteam@swagger.io

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` |  |
| **NIST CSF v2.0** | `` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [VULNERABLE] API3:2023 on https://petstore.swagger.io/v2/swagger.json?admin=true
- **Timestamp:** 2026-02-04T11:43:51Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json?admin=true`
- **Details:** Potential Hidden Parameter: admin

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1596` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [EXPLOITED]  on https://petstore.swagger.io/v2/swagger.json?admin=true
- **Timestamp:** 2026-02-04T11:43:51Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json?admin=true`
- **Details:** Leaked EMAIL: apiteam@swagger.io

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` |  |
| **NIST CSF v2.0** | `` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [VULNERABLE] API3:2023 on https://petstore.swagger.io/v2/swagger.json?test=true
- **Timestamp:** 2026-02-04T11:43:52Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json?test=true`
- **Details:** Potential Hidden Parameter: test

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1596` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [EXPLOITED]  on https://petstore.swagger.io/v2/swagger.json?test=true
- **Timestamp:** 2026-02-04T11:43:52Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json?test=true`
- **Details:** Leaked EMAIL: apiteam@swagger.io

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` |  |
| **NIST CSF v2.0** | `` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [VULNERABLE] API3:2023 on https://petstore.swagger.io/v2/swagger.json?dev=true
- **Timestamp:** 2026-02-04T11:43:52Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json?dev=true`
- **Details:** Potential Hidden Parameter: dev

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1596` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [EXPLOITED]  on https://petstore.swagger.io/v2/swagger.json?dev=true
- **Timestamp:** 2026-02-04T11:43:52Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json?dev=true`
- **Details:** Leaked EMAIL: apiteam@swagger.io

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` |  |
| **NIST CSF v2.0** | `` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [VULNERABLE] API3:2023 on https://petstore.swagger.io/v2/swagger.json?internal=true
- **Timestamp:** 2026-02-04T11:43:53Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json?internal=true`
- **Details:** Potential Hidden Parameter: internal

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1596` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [EXPLOITED]  on https://petstore.swagger.io/v2/swagger.json?internal=true
- **Timestamp:** 2026-02-04T11:43:53Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json?internal=true`
- **Details:** Leaked EMAIL: apiteam@swagger.io

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` |  |
| **NIST CSF v2.0** | `` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [VULNERABLE] API3:2023 on https://petstore.swagger.io/v2/swagger.json?config=true
- **Timestamp:** 2026-02-04T11:43:54Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json?config=true`
- **Details:** Potential Hidden Parameter: config

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1596` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [EXPLOITED]  on https://petstore.swagger.io/v2/swagger.json?config=true
- **Timestamp:** 2026-02-04T11:43:54Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json?config=true`
- **Details:** Leaked EMAIL: apiteam@swagger.io

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` |  |
| **NIST CSF v2.0** | `` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [VULNERABLE] API3:2023 on https://petstore.swagger.io/v2/swagger.json?role=true
- **Timestamp:** 2026-02-04T11:43:54Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json?role=true`
- **Details:** Potential Hidden Parameter: role

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1596` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [EXPLOITED]  on https://petstore.swagger.io/v2/swagger.json?role=true
- **Timestamp:** 2026-02-04T11:43:54Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json?role=true`
- **Details:** Leaked EMAIL: apiteam@swagger.io

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` |  |
| **NIST CSF v2.0** | `` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (Relative): /v2

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (Relative): /pet

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (Relative): /pet/findByStatus

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (Relative): /pet/findByTags

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (Relative): /store/inventory

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (Relative): /store/order

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (Relative): /user/createWithList

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (Relative): /user/login

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (Relative): /user/logout

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (Relative): /user/createWithArray

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (Relative): /user

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (Absolute): http://swagger.io/irc/

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (Absolute): http://swagger.io/terms/

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (Absolute): http://www.apache.org/licenses/LICENSE-2

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (Absolute): https://petstore.swagger.io/oauth/authorize

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (SPA-Route): definitions/ApiResponse

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (SPA-Route): definitions/Pet

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (SPA-Route): definitions/Order

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (SPA-Route): definitions/User

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (SPA-Route): definitions/Category

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [INFO] Unknown on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:43:59Z
- **Vector/Command:** `scrape`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** JS Discovery (SPA-Route): definitions/Tag

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [VULNERABLE] API8:2023 Security Misconfiguration on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:44:23Z
- **Vector/Command:** `audit`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** Weak CORS Policy: *

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1562.001` | Defense Evasion |
| **NIST CSF v2.0** | `PR.PS-01` | Control Mapping |
| **CVE / CVSS** | `CVE-2024-AUDIT` | **5.4** (Severity Score) |

---
### [WEAK CONFIG] API8:2023 Security Misconfiguration on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:44:23Z
- **Vector/Command:** `audit`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** Missing Header: Strict-Transport-Security

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1562.001` | Defense Evasion |
| **NIST CSF v2.0** | `PR.PS-01` | Control Mapping |
| **CVE / CVSS** | `CVE-2024-AUDIT` | **5.4** (Severity Score) |

---
### [WEAK CONFIG] API8:2023 Security Misconfiguration on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:44:23Z
- **Vector/Command:** `audit`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** Missing Header: Content-Security-Policy

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1562.001` | Defense Evasion |
| **NIST CSF v2.0** | `PR.PS-01` | Control Mapping |
| **CVE / CVSS** | `CVE-2024-AUDIT` | **5.4** (Severity Score) |

---
### [WEAK CONFIG] API8:2023 Security Misconfiguration on https://petstore.swagger.io/v2/swagger.json
- **Timestamp:** 2026-02-04T11:44:23Z
- **Vector/Command:** `audit`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json`
- **Details:** Missing Header: X-Content-Type-Options

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1562.001` | Defense Evasion |
| **NIST CSF v2.0** | `PR.PS-01` | Control Mapping |
| **CVE / CVSS** | `CVE-2024-AUDIT` | **5.4** (Severity Score) |

---
### [VULNERABLE] API3:2023 on https://petstore.swagger.io/v2/swagger.json?test=true
- **Timestamp:** 2026-02-04T11:44:53Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json?test=true`
- **Details:** Potential Hidden Parameter: test

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1596` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [VULNERABLE] API3:2023 on https://petstore.swagger.io/v2/swagger.json?dev=true
- **Timestamp:** 2026-02-04T11:44:54Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json?dev=true`
- **Details:** Potential Hidden Parameter: dev

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1596` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [VULNERABLE] API3:2023 on https://petstore.swagger.io/v2/swagger.json?internal=true
- **Timestamp:** 2026-02-04T11:44:54Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json?internal=true`
- **Details:** Potential Hidden Parameter: internal

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1596` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [VULNERABLE] API3:2023 on https://petstore.swagger.io/v2/swagger.json?config=true
- **Timestamp:** 2026-02-04T11:44:55Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json?config=true`
- **Details:** Potential Hidden Parameter: config

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1596` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
### [VULNERABLE] API3:2023 on https://petstore.swagger.io/v2/swagger.json?role=true
- **Timestamp:** 2026-02-04T11:44:58Z
- **Vector/Command:** `manual`
- **Target URL:** `https://petstore.swagger.io/v2/swagger.json?role=true`
- **Details:** Potential Hidden Parameter: role

**Compliance Mapping:**
| Framework | ID / Control | Description / Tactic |
| :--- | :--- | :--- |
| **MITRE ATT&CK** | `T1596` | Untriaged |
| **NIST CSF v2.0** | `N/A` | Control Mapping |
| **CVE / CVSS** | `-` | **0.0** (Severity Score) |

---
## 4. METHODOLOGY & FRAMEWORK ALIGNMENT

This assessment was conducted using the **VaporTrace Tactical Engine**, adhering to standard Adversary Emulation protocols.

### 4.1 Framework Reference
- **MITRE ATT&CK:** Used to classify adversary tactics and techniques (T-Codes).
- **OWASP API Security Top 10 (2023):** Primary standard for API vulnerability classification.
- **NIST CSF v2.0:** Used for mapping findings to defensive controls (Identify, Protect, Detect, Respond, Recover).
- **CVSS v3.1:** Common Vulnerability Scoring System for severity quantification.

**End of Report**
