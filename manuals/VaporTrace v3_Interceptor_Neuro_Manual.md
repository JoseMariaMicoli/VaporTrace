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