# 📘 VaporTrace v3.1-Hydra | Tactical Manual

## 1. 🧠 Neuro Brain Initialization (Setup)
Before launching attacks, you must configure the Hybrid Brain. We use **Groq** as the Primary (Fast/Free) and **Ollama** as the local fallback.

### **Step A: Start Local Fallback**
1.  Open a separate terminal.
2.  Run Ollama:
    ```bash
    ollama serve
    ```
3.  Ensure the model is loaded (do this once):
    ```bash
    ollama pull mistral
    ```

### **Step B: Configure Primary Brain (Groq)**
1.  Start VaporTrace: `./VaporTrace`
2.  In the CLI input, type:
    ```bash
    neuro config groq gsk_YOUR_GROQ_API_KEY_HERE llama-3.1-8b-instant
    ```
    *(Note: Using `llama-3.1-8b-instant` is recommended for maximum speed).*

### **Step C: Activation & Verification**
1.  Enable the engine:
    ```bash
    neuro on
    ```
2.  Test connectivity and latency:
    ```bash
    test-neuro
    ```
    *   **Success:** Logs show `[green]NEURO ONLINE:[-] ...`
    *   **Fail/429:** Logs show `[red]Primary Error... Switching to Fallback` (Uses Ollama).

---

## 2. 🛡️ The Tactical Interceptor (Man-in-the-Middle)
The Interceptor allows you to pause, modify, and inject AI payloads into HTTP requests *before* they leave your machine.

### **Step A: Engage Interceptor**
*   **Hotkey:** Press **`Ctrl + I`**
*   **Visual Indicator:** The bottom status bar will turn **RED** and display:
    `Ctrl+I: INTERCEPTING (ACTIVE)`

### **Step B: Trigger Traffic**
Run a command that generates HTTP traffic.
*   Example: `map -u http://target-api.com` or `test-bola`.
*   **Result:** The **Interceptor Modal** pops up immediately. The Logic thread pauses.

### **Step C: Manual Manipulation**
Inside the modal:
*   **TAB Key:** Switch between Method, URL, Headers, and Body fields.
*   **Editing:** Type directly to modify headers (e.g., change `User-Agent` or add `admin=true`).

### **Step D: Traffic Actions**
| Action | Hotkey | Description |
| :--- | :--- | :--- |
| **FORWARD** | `Ctrl + F` | Sends the modified request to the target. Resumes the logic thread. |
| **DROP** | `Ctrl + D` | Cancels the request entirely. It never hits the network. |
| **SYNC VAULT** | `Ctrl + S` | Saves the current request snapshot to the `Loot Database` (F3) for later reporting, without sending it yet. |

---

## 3. ⚡ Neuro-Kinetic Features (AI inside Interceptor)
While the Interceptor is open, you can use the AI to generate attacks for you on the fly.

### **Feature: Neuro-Brute (`Ctrl + B`)**
**Scenario:** You see a JSON body like `{"user_id": 101, "role": "guest"}` inside the Interceptor.
1.  **Focus:** Ensure your cursor is in the **Body** text area.
2.  **Trigger:** Press **`Ctrl + B`**.
3.  **What happens:**
    *   The engine sends the body to Groq (or Ollama).
    *   It generates 5 aggressive mutations (e.g., SQLi, Mass Assignment).
    *   **Output:** Switch to the **Neural Tab (F6)** to see the generated payloads. You can then copy/paste them back into the Interceptor manually (Safety precaution).

### **Feature: Neuro-Inverter (`Ctrl + N`)**
**Scenario:** You suspect logic flaws (BOLA) but don't want to fuzz manually.
1.  **Trigger:** Press **`Ctrl + N`** inside the Interceptor.
2.  **What happens:**
    *   Toggles a global flag `NeuroInverterActive`.
    *   Any request forwarded while this is ON will automatically undergo "Logic Inversion" attempts (e.g., swapping `GET` to `DELETE` or swapping User IDs) in the background.

---

## 4. 🔍 Deep Traffic Analysis (Snapshotting)
This allows you to analyze traffic *after* it has happened (Post-Mortem).

1.  **Navigate:** Switch to **Traffic View (F4)** using the F-keys.
2.  **Select:** Ensure you can see a Request/Response pair you are interested in.
3.  **Trigger:** Press **`Ctrl + A`**.
4.  **Process:**
    *   The HTTP Snapshot is sent to the Neuro Engine.
    *   The AI performs "Chain of Thought" reasoning.
    *   It maps findings to **MITRE ATT&CK** and **OWASP**.
5.  **Results:**
    *   Switch to **Neural View (F6)**.
    *   Read the AI's analysis.
    *   **Auto-Attack:** If the AI finds high-probability exploits in the snapshot, it will *automatically* fire 3 fuzzing packets at the target in the background.

---

## 🔢 Summary of Hotkeys (v3.1-Hydra)

| Key Combination | Scope | Function |
| :--- | :--- | :--- |
| **`Ctrl + I`** | Global | **Toggle Interceptor** (On/Off) |
| **`Ctrl + F`** | Modal | **Forward** packet to network |
| **`Ctrl + D`** | Modal | **Drop** packet |
| **`Ctrl + B`** | Modal | **Neuro Brute:** Gen payloads for current field |
| **`Ctrl + S`** | Modal | **Sync:** Save to Loot DB |
| **`Ctrl + A`** | F4 Tab | **Analyze:** Send snapshot to AI Brain |
| **`F1 - F6`** | Global | Switch Tabs (Logs, Map, Loot, Traffic, Context, **Neural**) |

---
## Neuro Notes:
### 1. The neuro-gen 
How it works:
It allows you to manually generate payloads without sending them, using the AI.
Usage: neuro-gen "search parameter sql injection" 5
Result: It generates 5 specific payloads for that context and prints them to the Neural Tab (F6).
(I have included the fix for usage in the code section below).
### 2. Tactical Hotkeys Explained
These shortcuts are context-sensitive. Here is exactly what they do:
Inside Interceptor (Ctrl + I)
#### Ctrl + B (Neuro Brute):
Function: Takes the text currently in the Body field of the interceptor and sends it to the AI.
Goal: The AI generates 5-10 high-entropy mutations (SQLi, XSS, JSON injection) based on that specific body.
Output: Results appear in the Neural Tab (F6). You can then copy/paste them back into the interceptor to fire them.
#### Ctrl + S (Sync to Vault):
Function: "Snapshot." It saves the current Request/Response details to your Loot Database (F3) and the internal SQLite DB immediately.
Goal: Useful if you see an interesting packet but aren't ready to exploit it yet. It bookmarks it for the Report.
#### Ctrl + N (Neuro Invert):
Function: Toggles a logic switch called "Inverter Mode."
Goal: If active, the system automatically tries to "Invert" logic flow on the next forwarded packet (e.g., swapping POST to GET, or swapping user_id values) without you typing it manually.
Inside Traffic View (F4)
#### Ctrl + A (Analyze):
Function: Takes the currently selected Request/Response pair in the traffic window.
Goal: Sends the entire snapshot to the AI (Groq/Ollama).
Output: The AI performs a "Deep Analysis," looking for BOLA, BFLA, and Sensitive Data Exposure, and prints a report to Tab 6.

##### Ctrl + H Shows Keybindings

---




Here is the comprehensive, step-by-step usage guide for the **VaporTrace Phase III** architecture. This guide covers compilation, initialization, telemetry gathering, and the new **AI Strategic Planner (HITL)** workflow.

---

# 🛡️ VaporTrace Phase III - Operator Manual

## 1. Installation & Initialization

Before starting a mission, ensure the binary is compiled with the latest patches.

### Step 1: Compile the Suite
Ensure you are in the root directory of the project.
```bash
go mod tidy
go build -o VaporTrace main.go
```

### Step 2: Launch the Dashboard
Run the tool in TUI (Tactical User Interface) mode.
```bash
./VaporTrace
```

### Step 3: Initialize Persistence
VaporTrace uses an SQLite backend to track findings and map them to MITRE/OWASP.
1.  Type `init_db` in the command bar.
2.  **Verify:** Check the **LOGS (F1)** tab for `[green]Database Persistence Initialized`.

---

## 2. Telemetry Gathering (Feeding the Brain)

The **Strategic Brain** relies on data. It will not generate a plan if the system is empty. You must populate the **Sensors** first.

### Step 4: Set Global Scope
Define the target to ensure all modules share the same context.
```text
target https://httpbin.org
```

### Step 5: Populate Discovery Map (Sensor F2)
The AI needs to know the API topology (endpoints) to recommend attacks.
*   **Option A (Swagger):**
    ```text
    map -u https://httpbin.org/spec.json
    ```
*   **Option B (JS Scraping):**
    ```text
    scrape https://target.com/assets/app.bundle.js
    ```
*   **Verification:** Press **F2** to see the populated Endpoint Map.

### Step 6: Populate Loot/Auth (Sensor F3)
The AI increases confidence scores if it detects valid credentials (e.g., swapping a found JWT to test BOLA).
*   **Manual Injection:**
    ```text
    auth attacker eyJhbGciOiJIUzI1Ni...
    ```
*   **Auto-Capture:** Browse the target using the VaporTrace Proxy. Any captured tokens will automatically appear in **LOOT (F3)**.

---

## 3. The Strategic Planner (HITL Workflow)

This is the core Phase III workflow. You will use the AI to generate an attack plan, refine it, and execute it.

### Step 7: Analyze Telemetry
Trigger the Neuro-Engine to correlate Map + Loot + Traffic.
```text
analyze
```
*   **Behavior:** The screen will automatically switch to **PLAN (F5)**.
*   **Visual:** You will see the **Strategic Action Buffer** table populated with "PENDING" actions (e.g., BOLA candidates, BFLA candidates).

### Step 8: Human-in-the-Loop (Refinement)
Review the proposed table in Tab F5. The AI might guess IDs incorrectly or suggest unsafe actions.

*   **Edit a Payload:**
    If Action #1 is a BOLA attack using `ID: 1337`, but you know the target ID is `999`:
    ```text
    edit 1 ID: 999
    ```
    *Result:* The table updates the Payload column immediately.

*   **Drop an Action:**
    If Action #2 is a destructive DELETE method you don't want to run:
    ```text
    drop 2
    ```
    *Result:* The status changes to `[red]DROPPED`.

### Step 9: Commit the Plan
Once the buffer is refined, authorize the engine to fire.
```text
commit
```
*   **Behavior:**
    1.  The engine iterates through the buffer.
    2.  It skips `DROPPED` actions.
    3.  It executes `PENDING` actions asynchronously.
    4.  Status updates to `[green]EXECUTED`.

### Step 10: Review Results
*   Switch to **LOGS (F1)** to see the raw output of the attacks.
*   Switch to **NEURO (F6)** if you used the `remediate` command.

---

## 4. Remediation & Reporting

### Step 11: Generate Fixes (Sprint 16 Preview)
If the BOLA attack (Action #1) was successful, ask the Blue Team module for a fix.
```text
remediate BOLA
```
*   **Result:** Check **NEURO (F6)** for a Golang/Node.js code snippet to fix the vulnerability.

### Step 12: Mission Debrief
Generate the final artifact for the client.
```text
report
```
*   **Result:** A timestamped Markdown file (e.g., `VAPORTRACE_PEN_TEST_2026...md`) is created in the `reports/` folder.

---

## 🔒 Advanced: Interceptor Mode

To manually manipulate traffic before it hits the sensors:
1.  Press `Ctrl + I` to toggle the Interceptor **ON**.
2.  Trigger an action (e.g., `exhaust https://httpbin.org/get limit`).
3.  The UI will freeze and show a **Modal**.
4.  **Edit** the Body/Headers.
5.  Press `Ctrl + F` to Forward or `Ctrl + D` to Drop.



1. SYSTEM ARCHITECTURE DIAGRAM
The VaporTrace Strategic Brain implements a centralized telemetry aggregation model. The pkg/engine acts as the cortex, pulling state from three sensory inputs located in pkg/logic and pkg/db, processing them through a Heuristic State Machine, and outputting actionable vectors to the UI.
code
Ascii
[ F2: TOPOLOGY SENSOR ]      [ F3: CREDENTIAL SENSOR ]      [ F4: TRAFFIC SENSOR ]
      (pkg/logic/store.go)         (pkg/logic/loot.go)            (pkg/logic/network.go)
             |                            |                              |
             | GetEndpoints()             | GetLootSummary()             | GetTrafficHistory()
             |                            |                              |
             +-------------+--------------+--------------+---------------+
                           |                             |
                           v                             v
                  +-----------------------------------------+
                  |      PKG/ENGINE: STRATEGIC CORTEX       |
                  | 1. Aggregates data into 'TelemetryState'|
                  | 2. Runs Heuristic State Machine         |
                  | 3. Generates 'TacticalAction' objects   |
                  +-----------------------------------------+
                                       |
                                       v
                  +-----------------------------------------+
                  |      TAB F5: STRATEGIC ACTION BUFFER    |
                  | [ID] [TYPE] [TARGET] [CONFIDENCE]       |
                  |  1.  BFLA   /admin   HIGH (403->200)    |
                  |  2.  SSRF   /hook    CRITICAL (AWS Key) |
                  +-----------------------------------------+
                                       ^
                                       | (HITL Loop)
                               [ HUMAN OPERATOR ] 