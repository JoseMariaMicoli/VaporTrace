#!/usr/bin/env python3
import subprocess
import sys
import os
import shutil
import argparse
from collections import Counter
from datetime import datetime

# ==========================================
# CONFIGURATION MAP (The "Plan")
# ==========================================
# Mapped from the provided 139 commits.
PROJECT_CONFIG = [
    {
        "phase_name": "Phase 1 - Foundation & Early Discovery (Sprint 1-5)",
        "features": [
            {
                "name": "feat/foundation-setup",
                "commits": [
                    "5d2df74",  # Initial commit
                    "291af5e",  # complete phase 1 - foundation, cobra cli, and burp bridge
                    "4a2d8ca",  # Complete phase 2.0 - swagger ingestion and automated shadow route probing
                ]
            },
            {
                "name": "feat/discovery-phase2-3",
                "commits": [
                    "5f98f7e",  # Phase2 Completed phase 2 - discovery, shadow probing, and parameter mining
                    "dd5ba75",  # Updated roadmap to include CLI UI checkpoint and BOLA stability of phase3
                    "f81deba",  # Finalize roadmap and add detailed interactive shell usage
                    "95932e7",  # Phase3.2 BOPLA Implemented
                    "1fb0433",  # Completed Phase3.3 BFLA Implemented
                ]
            },
            {
                "name": "feat/api-resource-exhaustion",
                "commits": [
                    "9d6dbfa",  # Changed Banner UI and added MITRE MAPPING and IR TEMPLATE to the README
                    "8fbb60e",  # Implemented phase4.1 Resource Exhaustion (API4)
                    "8ef59da",  # Implemented phase4.1 Resource Exhaustion (API4)
                    "2bfd086",  # Implemented Phase4.2 SSRF Tracker API7
                    "b008392",  # Implemented Phase 4.3 Security Misconfiguration API8
                    "5bce872",  # Implemented Phase4.4  Integration Probe (API10)
                    "36268dc",  # Implemented Phase4.4  Integration Probe (API10)
                    "4ad5231",  # Added Phase5 to the readme
                    "4502195",  # Phase 5.1 completed
                    "01d5bc5",  # Phase 5.2 completed
                    "8144b74",  # PHASE 5 Successfully Completed
                ]
            }
        ],
        "qa_fixes": [
            {
                "name": "fix/readme-and-burp-bridge",
                "commits": [
                    "0c682f6",  # Corrected the readme
                    "0da300e",  # Readme corrected
                    "6b81e7e",  # Readme corrected part2
                    "4870748",  # SplashScreen in readme fixed
                    "3f6eeb1",  # Fixed readme
                    "acb0a74",  # BUG: BurpSuite Proxy Integration FIXED
                ]
            }
        ],
        "release_tag": "v0.5.0"
    },
    {
        "phase_name": "Phase 2 - Roadmap Expansion & Core Stability (Sprint 6-9)",
        "features": [
            {
                "name": "feat/phase9-roadmap-execution",
                "commits": [
                    "4baa75f",  # Phase 9.1 Completed
                    "e299067",  # Phase 9.2 completed
                    "c1be227",  # Phase 9.3 Completed
                    "d314c5e",  # Phase 9.4 Completed
                    "f8672d4",  # Updated README to-do tasks checklist
                    "0d0a81c",  # Updated README to-do tasks checklist 2
                    "1334629",  # Updated README to-do tasks checklist 3
                    "948b0ff",  # Phase 9.5 Completed
                    "48e33dd",  # Phase 9.5 Completed 2
                    "d91b885",  # Phase 9.6 completed
                    "741153e",  # Phase 9.6 completed 2
                    "6e2d076",  # Phase 9.7 completed
                    "2741015",  # Phase 9.8 completed
                    "7b148ea",  # Phase 9.9 completed
                    "dc16fe0",  # Phase 9.9 completed
                    "1d53e9b",  # Phase 9.10 completed
                ]
            },
            {
                "name": "feat/phase6-7-8-consolidation",
                "commits": [
                    "f9bd53f",  # Changed Phase 9 to stable now ACTIVE Phase 6
                    "c650d2f",  # Changed Phase 9 to stable now ACTIVE Phase 6
                    "c04379d",  # Phase 6 completed
                    "012687c",  # Phase 7.1 completed
                    "a20f7eb",  # Phase 7.1 completed
                    "94eeff6",  # Phase 7.2 completed
                    "324f2e7",  # Phase 7.3 completed
                    "6aef642",  # Phase 7.3 completed
                    "d02ce1f",  # Phase 7.3 completed update readme
                    "2f66ba9",  # Phase 8.1 completed
                    "9fe4fd1",  # Phase 8.1 completed
                    "0651687",  # Phase 8.2 completed
                    "754bd10",  # Added to the roadmap phase 9.3 OOB
                    "25caf79",  # Added to the roadmap phase 9.3 OOB
                    "50ae80d",  # Added to the roadmap phase 9.3 OOB
                    "219cdee",  # Phase 8.3.1 completed
                    "142e294",  # Phase 8.3.1 completed
                    "bdabe58",  # Phase 8 DONE Project STATUS STABLE
                ]
            }
        ],
        "qa_fixes": [],  # No specific fixes logged in this block
        "release_tag": "v0.9.0"
    },
    {
        "phase_name": "Phase 3 - Report Engine & TUI Dashboard (Sprint 10)",
        "features": [
            {
                "name": "feat/report-engine-ui",
                "commits": [
                    "2e78ea0",  # Updated Readme
                    "27e8b2b",  # Updated Readme and report engine
                    "70dfbc4",  # Removed new UI
                    "65d2ee5",  # Phase 10.1 GLOBAL TARGER DONE
                    "9d5babe",  # Updated the ROADMAP
                    "8b68245",  # Updated the ROADMAP
                    "d1d7053",  # Updated the ROADMAP
                    "b6f5c0c",  # Phases 10.2.1 and 10.2.2 DONE
                    "99e9608",  # Updated RoadMap
                    "9414e15",  # Phase 9.13 and Phase 10.2.2 DONE
                    "f8c9bb6",  # Updated MITRE ATT&CK Mapping
                    "50813bd",  # Updated SEED_DB
                    "38ece9f",  # Updated REPORT ENGINE
                    "99d4089",  # Changed report location and somo adjustments are applied
                    "86c5fa3",  # Report Engine and TUI Dashboard DONE
                    "1e3d797",  # Updated Report Gen and Seed_DB
                    "1ce161b",  # Phase 10.3 DONE - Phase 10.4 ACTIVE implemented and testing functionalities
                    "40c1589",  # Phase 10.3 DONE - Phase 10.4 ACTIVE implemented and testing functionalities
                    "63a421d",  # Phase 10.4 DONE updated README
                    "0efc4fc",  # Phase 10.5 and Phase 10.6 DONE IN TESTING
                ]
            },
            {
                "name": "feat/docs-install-guide",
                "commits": [
                    "c09b06f",  # Updated README
                    "4d7a6cb",  # Updated README report examples
                    "a338c15",  # Updated installation guide
                    "a7601fb",  # Updated installation guide
                    "a4965df",  # Updated Description
                    "ad125a7",  # Updated readme and added grok
                    "ecc8617",  # Updated readme
                    "1c5b02d",  # Updated readme
                ]
            }
        ],
        "qa_fixes": [
            {
                "name": "fix/ui-reports-and-bugs",
                "commits": [
                    "6f2b7e1",  # REPORT COMPILATION BUG FIXED
                    "310f094",  # BUG FIX: MASSBOLA UI BREAK FIXED
                    "b86d24f",  # QA Testing Bug Fixes
                    "bcd9fb5",  # Improved many functionalities and BUG FIXES APPLIED
                ]
            }
        ],
        "release_tag": "v1.0.0"
    },
    {
        "phase_name": "Phase 4 - Documentation Overhaul & Regression Fixes (Sprint 11)",
        "features": [
            {
                "name": "feat/manuals-and-help",
                "commits": [
                    "99be263",  # Added manuals and md2pdf.py
                    "fc13b91",  # Added report tab in dev
                    "e39648a",  # Added text editor syntax highlight and editing functionality
                    "9d42e0d",  # Added show help modal
                    "a30929e",  # Updated README
                    "7894b21",  # Created full documentation (manuals and dev-logs) and complete 11.1 from Sprint 11
                ]
            }
        ],
        "qa_fixes": [
            {
                "name": "fix/readme-markdown-regression",
                "commits": [
                    "2491f8b",  # Bug Fixed, improvements and hardened functionalities
                    "9f8a14f",  # Fixed README BANNER
                    "832ac41",  # Fixed README
                    "6bc2be5",  # Fixed README commands markdown
                    "6b4bb56",  # Fixed README commands markdown
                    "7d78dc4",  # Fixed README mermaid
                    "abdf3a1",  # Fixed README mermaid
                    "c510631",  # Fixed README mermaid
                    "d50d006",  # Fixed README mermaid
                    "c0467c1",  # Fixed README mermaid
                    "5c6e6f3",  # Fixed README mermaid
                    "046804e",  # Untracked .obsidian directory
                    "ddb35c1",  # Untracked .obsidian directory
                    "87157ff",  # Fix TUI SSRF cascade breaking bug
                ]
            }
        ],
        "release_tag": "v2.0.0"
    },
    {
        "phase_name": "Phase 5 - Advanced Offense & Tiers 1-3 (Sprint 12-19)",
        "features": [
            {
                "name": "feat/sprint12-oob-stealth",
                "commits": [
                    "946e720",  # Sprint 12 Add OOB Exfil, stochastic jitter, plus 22 userAgents rotating but with matching
                    "dc4d420",  # Phase12 DONE under testing and Updated Documentation and Dev-Log
                    "6abd375",  # Updated README
                    "7a271dc",  # Updated README
                    "d9da855",  # Updated README and DOCS
                    "072fb62",  # Stealth functionality implemented
                    "49b808e",  # Updated Dev-Logs
                    "4efae5d",  # Updated Documentation tree
                ]
            },
            {
                "name": "feat/tier-upgrades",
                "commits": [
                    "faf67ff",  # Sprint 17 under development
                    "1f95000",  # Sprint 18 Spider/Crawler implemented need to add usage / help and manual
                    "6ef128f",  # Sprint 18 DONE Tier1 and Tier2 Upgrade Done
                    "f9327d8",  # Sprint 19 Tier 3 Part1 of 3 Done
                    "4b28ece",  # Sprint 19 Tier 3 Part2 of 3 Done
                    "523e374",  # Sprint 19 DONE Tier 3 Part3 of 3 Done
                ]
            }
        ],
        "qa_fixes": [
            {
                "name": "fix/neuro-pipeline-race-conditions",
                "commits": [
                    "9c91950",  # BUG Fixes applied still remain some fixes to apply in the next commit, under testing
                    "b999c8b",  # Fixed double / bug in BOLA, BFLA and BOPLA (pipeline)
                    "e23cef3",  # Fixed Neuro docs, nil point and race conditions in neuro engine
                ]
            }
        ],
        "release_tag": "v3.0.0"
    },
    {
        "phase_name": "Phase 6 - Tier 4 Autonomy & Final Polish",
        "features": [
            {
                "name": "feat/tier4-autonomy",
                "commits": [
                    "7b20d45",  # Tier-4 Part 2 of 3 DONE
                    "bf35c7d",  # Tier-4 DONE
                    "b8af0ab",  # Updated README
                    "a6a32d0",  # Updated README
                    "1cc5cc2",  # feat(tier4): harden autonomous logic and optimize TUI visibility
                    "515c8de",  # JA3/JA4 TLS Fingerprint evasion improved but not working properly. Under testing
                ]
            },
            {
                "name": "feat/licensing-compliance",
                "commits": [
                    "281db80",  # Update README license notice
                    "0e9ab28",  # Update LICENSE and source/docs headers
                    "329b402",  # Version 3.1.1 published
                ]
            }
        ],
        "qa_fixes": [
             {
                "name": "fix/license-badges",
                "commits": [
                    "6b65819",  # Corrected the License Badge in the Readme in compliance with actual licensing
                    "b7a84c5",  # Corrected the License Badge in the Readme in compliance with actual licensing
                ]
            }
        ],
        "release_tag": "v3.1.1"
    }
]


# ==========================================
# CORE UTILITIES
# ==========================================

RUN_ID = datetime.now().strftime("%Y%m%d-%H%M%S")
LOG_FILE = f"git_reconstruction_{RUN_ID}.log"
BACKUP_REFS = {}
DRY_RUN = False
USE_SIGNING = True
FORCE_TAGS = False
ALLOW_UNTRACKED = False
SHELVED_ROOT = f"/tmp/git_recon_shelved_{RUN_ID}"
SHELVED_PATHS = []


def log_event(message):
    print(message)
    with open(LOG_FILE, "a", encoding="utf-8") as f:
        f.write(message + "\n")


def print_restore_instructions():
    if not BACKUP_REFS:
        return
    log_event("\n[!] Recovery commands (if you need rollback):")
    for branch_name, backup_ref in BACKUP_REFS.items():
        log_event(f"    git branch -f {branch_name} {backup_ref}")
    log_event("    git checkout main")


def fail_critical(message):
    log_event(message)
    print_restore_instructions()
    log_event(f"[!] Execution log: {LOG_FILE}")
    sys.exit(1)


def run_git(args, env=None, description=None):
    """
    Executes a git command using subprocess.
    - inherits stdin/stdout/stderr for GPG interactivity.
    - checks for errors and halts execution if non-zero exit.
    """
    if description:
        log_event(f"\n[+] Action: {description}")
    log_event(f"    Command: git {' '.join(args)}")

    if DRY_RUN:
        return

    run_git_once(args, env=env, description=description)


def run_git_once(args, env=None, description=None):
    """Single-attempt git command runner that raises on failure."""
    final_env = os.environ.copy()
    if env:
        final_env.update(env)

    subprocess.run(
        ["git"] + args,
        check=True,
        env=final_env
    )


def cherry_pick_with_fallback(commit_hash, env_dates, description):
    """
    Cherry-pick with automatic conflict fallback:
    1) normal cherry-pick
    2) abort and retry with -X theirs
    """
    base_args = ["cherry-pick"] + signing_flag_or_empty() + ["-x", commit_hash]

    if DRY_RUN:
        run_git(base_args, env=env_dates, description=description)
        return

    def cherry_pick_in_progress():
        return os.path.exists(".git/CHERRY_PICK_HEAD")

    def cherry_pick_has_no_changes():
        idx_clean = subprocess.run(
            ["git", "diff", "--cached", "--quiet"],
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL
        ).returncode == 0
        wt_clean = subprocess.run(
            ["git", "diff", "--quiet"],
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL
        ).returncode == 0
        return idx_clean and wt_clean

    def maybe_skip_empty_cherry_pick():
        if cherry_pick_in_progress() and cherry_pick_has_no_changes():
            log_event(f"[~] Empty cherry-pick for {commit_hash}; skipping.")
            subprocess.run(["git", "cherry-pick", "--skip"], check=True)
            return True
        return False

    def try_force_theirs_continue():
        """
        For conflict classes not handled by -X theirs (e.g. modify/delete),
        force resolve by checking out theirs for all paths, staging, and continuing.
        """
        if not cherry_pick_in_progress():
            return False
        log_event(f"[~] Attempting forced conflict resolution for {commit_hash} via checkout --theirs + continue")
        subprocess.run(["git", "checkout", "--theirs", "--", "."], check=False)
        subprocess.run(["git", "add", "-A"], check=False)
        try:
            subprocess.run(["git", "cherry-pick", "--continue"], check=True)
            return True
        except subprocess.CalledProcessError:
            if maybe_skip_empty_cherry_pick():
                return True
            return False

    try:
        run_git_once(base_args, env=env_dates, description=description)
        return
    except subprocess.CalledProcessError:
        if maybe_skip_empty_cherry_pick():
            return
        log_event(f"[~] Conflict/failed cherry-pick for {commit_hash}; retrying with -X theirs")
        subprocess.run(["git", "cherry-pick", "--abort"], check=False)
        retry_args = ["cherry-pick"] + signing_flag_or_empty() + ["-X", "theirs", "-x", commit_hash]
        try:
            run_git_once(retry_args, env=env_dates, description=f"{description} (retry -X theirs)")
            return
        except subprocess.CalledProcessError as e:
            if maybe_skip_empty_cherry_pick():
                return
            if try_force_theirs_continue():
                return
            fail_critical(
                f"\n[!] CRITICAL ERROR during: {description}\n"
                f"[!] Command failed with exit code {e.returncode} even after -X theirs retry.\n"
                "[!] Please resolve manually and re-run."
            )


def run_git_capture(args):
    try:
        result = subprocess.run(
            ["git"] + args,
            capture_output=True,
            text=True,
            check=True
        )
        return result.stdout.strip()
    except subprocess.CalledProcessError:
        fail_critical(f"[!] Failed command: git {' '.join(args)}")
        return ""


def get_untracked_paths():
    result = subprocess.run(
        ["git", "status", "--porcelain"],
        capture_output=True,
        text=True,
        check=True
    )
    paths = []
    for line in result.stdout.splitlines():
        if line.startswith("?? "):
            paths.append(line[3:])
    return paths


def get_commit_touched_paths(commit_hash):
    output = run_git_capture(["show", "--pretty=format:", "--name-only", commit_hash])
    return {line.strip() for line in output.splitlines() if line.strip()}


def shelve_conflicting_untracked(commit_hash):
    """
    Moves untracked files that would collide with a cherry-pick into /tmp.
    This preserves operator files while allowing deterministic replay.
    """
    untracked = set(get_untracked_paths())
    if not untracked:
        return

    touched = get_commit_touched_paths(commit_hash)
    conflicts = sorted(untracked.intersection(touched))
    if not conflicts:
        return

    os.makedirs(SHELVED_ROOT, exist_ok=True)
    for rel_path in conflicts:
        if not os.path.lexists(rel_path):
            continue
        dest_path = os.path.join(SHELVED_ROOT, rel_path)
        os.makedirs(os.path.dirname(dest_path), exist_ok=True)
        shutil.move(rel_path, dest_path)
        SHELVED_PATHS.append(rel_path)
        log_event(f"[~] Shelved conflicting untracked path: {rel_path} -> {dest_path}")


def branch_exists(branch_name):
    result = subprocess.run(
        ["git", "show-ref", "--verify", "--quiet", f"refs/heads/{branch_name}"]
    )
    return result.returncode == 0


def tag_exists(tag_name):
    result = subprocess.run(
        ["git", "show-ref", "--verify", "--quiet", f"refs/tags/{tag_name}"]
    )
    return result.returncode == 0


def commit_exists(commit_hash):
    result = subprocess.run(
        ["git", "cat-file", "-e", f"{commit_hash}^{{commit}}"],
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL
    )
    return result.returncode == 0


def commit_already_in_head(commit_hash):
    """Returns True if commit is already reachable from HEAD."""
    result = subprocess.run(
        ["git", "merge-base", "--is-ancestor", commit_hash, "HEAD"],
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL
    )
    return result.returncode == 0


def get_commit_date(commit_hash):
    """
    Retrieves the strictly formatted ISO-8601 Author Date from a specific commit.
    Used to spoof the Committer Date during cherry-picks.
    """
    return run_git_capture(["show", "-s", "--format=%aI", commit_hash])


def get_current_branches():
    """Returns a list of local branch names."""
    output = run_git_capture(["branch", "--format=%(refname:short)"])
    return output.splitlines() if output else []


def is_working_directory_clean(allow_untracked=False):
    """Checks if there are uncommitted changes."""
    result = subprocess.run(
        ["git", "status", "--porcelain"],
        capture_output=True,
        text=True
    )
    lines = [line for line in result.stdout.splitlines() if line.strip()]
    if allow_untracked:
        lines = [line for line in lines if not line.startswith("?? ")]
    return len(lines) == 0


def parse_args():
    parser = argparse.ArgumentParser(description="Forensic GitFlow reconstruction")
    parser.add_argument("--dry-run", action="store_true", help="Validate and print commands without changing git refs")
    parser.add_argument("--yes", action="store_true", help="Skip destructive operation confirmation prompt")
    parser.add_argument("--no-sign", action="store_true", help="Disable GPG signing (-S / -s)")
    parser.add_argument("--force-tags", action="store_true", help="Delete existing release tags before creating them")
    parser.add_argument("--allow-untracked", action="store_true", help="Allow untracked files in working tree")
    parser.add_argument("--log-file", default=LOG_FILE, help="Path to execution log file")
    return parser.parse_args()


def preflight_validate():
    all_commits = []
    all_tags = []

    for phase in PROJECT_CONFIG:
        all_tags.append(phase["release_tag"])
        for feat in phase["features"]:
            all_commits.extend(feat["commits"])
        for fix in phase["qa_fixes"]:
            all_commits.extend(fix["commits"])

    duplicate_hashes = [h for h, count in Counter(all_commits).items() if count > 1]
    if duplicate_hashes:
        fail_critical(f"[!] Duplicate commit hashes found in PROJECT_CONFIG: {duplicate_hashes}")

    invalid_hashes = [h for h in all_commits if len(h) != 7 or any(c not in "0123456789abcdef" for c in h)]
    if invalid_hashes:
        fail_critical(f"[!] Invalid short hash format found: {invalid_hashes}")

    missing_hashes = [h for h in all_commits if not commit_exists(h)]
    if missing_hashes:
        fail_critical(f"[!] Missing commits in repository history: {missing_hashes}")

    duplicate_tags = [t for t, count in Counter(all_tags).items() if count > 1]
    if duplicate_tags:
        fail_critical(f"[!] Duplicate release tags in PROJECT_CONFIG: {duplicate_tags}")

    colliding_tags = [t for t in all_tags if tag_exists(t)]
    if colliding_tags and not FORCE_TAGS:
        fail_critical(
            f"[!] These tags already exist: {colliding_tags}\n"
            "[!] Re-run with --force-tags if you explicitly want to replace them."
        )

    if USE_SIGNING and shutil.which("gpg") is None:
        fail_critical("[!] GPG signing requested but 'gpg' binary was not found. Use --no-sign if needed.")

    log_event("[+] Preflight validation passed.")


def maybe_confirm_or_exit():
    if DRY_RUN or os.environ.get("GIT_RECON_AUTO_YES") == "1":
        return

    log_event(
        "\n[!] This script will rewrite local branches (main/dev/QA).\n"
        "[!] Backups are created first, but this operation is destructive."
    )
    answer = input("Type RECONSTRUCT to proceed: ").strip()
    if answer != "RECONSTRUCT":
        fail_critical("[!] Confirmation failed. Aborted by operator.")


def create_branch_backups(existing_branches):
    log_event(">>> Performing backups...")
    for branch_name in ("main", "dev", "QA"):
        if branch_name not in existing_branches:
            continue

        immutable_backup = f"{branch_name}-backup-{RUN_ID}"
        mutable_backup = f"{branch_name}-backup"

        if branch_exists(immutable_backup):
            fail_critical(f"[!] Backup branch collision: {immutable_backup} already exists.")

        run_git(
            ["branch", immutable_backup, branch_name],
            description=f"Creating immutable backup {immutable_backup} from {branch_name}"
        )
        run_git(
            ["branch", "-f", mutable_backup, immutable_backup],
            description=f"Updating {mutable_backup} pointer"
        )
        BACKUP_REFS[branch_name] = immutable_backup

    if not BACKUP_REFS:
        fail_critical("[!] No source branches found to back up (expected at least 'main').")


def signing_flag_or_empty():
    return ["-S"] if USE_SIGNING else []


def unique_work_branch(base_name):
    return f"{base_name}-{RUN_ID}"


# ==========================================
# MAIN EXECUTION FLOW
# ==========================================

def main():
    global DRY_RUN, USE_SIGNING, FORCE_TAGS, LOG_FILE, ALLOW_UNTRACKED

    args = parse_args()
    DRY_RUN = args.dry_run
    USE_SIGNING = not args.no_sign
    FORCE_TAGS = args.force_tags
    ALLOW_UNTRACKED = args.allow_untracked
    LOG_FILE = args.log_file

    # Run cleanliness guard before any logging side effects inside the repo.
    if not is_working_directory_clean(allow_untracked=ALLOW_UNTRACKED):
        print("[!] Error: Working directory is not clean. Commit or stash changes first.")
        print(f"[!] Execution log: {LOG_FILE}")
        sys.exit(1)

    log_event(">>> STARTING OFFENSIVE GIT RECONSTRUCTION <<<")
    log_event(f"[+] Run ID: {RUN_ID}")
    log_event(f"[+] Dry run mode: {DRY_RUN}")
    log_event(f"[+] Signing enabled: {USE_SIGNING}")
    log_event(f"[+] Force tags: {FORCE_TAGS}")
    log_event(f"[+] Allow untracked: {ALLOW_UNTRACKED}")

    preflight_validate()

    if not args.yes:
        maybe_confirm_or_exit()

    existing_branches = get_current_branches()
    create_branch_backups(existing_branches)

    # Optional tag cleanup if user explicitly requested overwrite.
    if FORCE_TAGS:
        for phase in PROJECT_CONFIG:
            tag_name = phase["release_tag"]
            if tag_exists(tag_name):
                run_git(["tag", "-d", tag_name], description=f"Removing existing tag {tag_name}")

    # Gets the very first commit reachable from HEAD
    root_output = run_git_capture(["rev-list", "--max-parents=0", "HEAD"])
    root_lines = [line for line in root_output.splitlines() if line.strip()]
    if not root_lines:
        fail_critical("[!] Failed to find root commit.")
    root_commit = root_lines[-1]
    log_event(f">>> Root commit identified: {root_commit}")

    run_git(["checkout", root_commit], description="Checking out root commit (Detached HEAD)")

    run_git(["checkout", "-B", "main", root_commit], description="Resetting main to root")
    run_git(["checkout", "-B", "dev", root_commit], description="Resetting dev to root")
    run_git(["checkout", "-B", "QA", root_commit], description="Resetting QA to root")
    run_git(["checkout", "dev"], description="Switching to dev context")

    for phase in PROJECT_CONFIG:
        phase_name = phase["phase_name"]
        log_event("\n========================================")
        log_event(f" PROCESSING PHASE: {phase_name}")
        log_event("========================================")

        for feat in phase["features"]:
            feat_branch = unique_work_branch(feat["name"])

            run_git(["checkout", "dev"], description="Aligning dev")
            run_git(["checkout", "-b", feat_branch], description=f"Created {feat_branch}")

            for commit_hash in feat["commits"]:
                if commit_already_in_head(commit_hash):
                    log_event(f"[~] Skipping {commit_hash}: already present in HEAD for {feat_branch}")
                    continue
                shelve_conflicting_untracked(commit_hash)
                orig_date = get_commit_date(commit_hash)
                env_dates = {
                    "GIT_AUTHOR_DATE": orig_date,
                    "GIT_COMMITTER_DATE": orig_date
                }
                cherry_pick_with_fallback(
                    commit_hash,
                    env_dates,
                    description=f"Cherry-picking {commit_hash} into {feat_branch}"
                )

            run_git(["checkout", "dev"], description="Switching to dev")
            merge_args = ["merge", "--no-ff"] + signing_flag_or_empty() + [feat_branch]
            run_git(merge_args, description=f"Merging {feat_branch} into dev")
            run_git(["branch", "-d", feat_branch], description=f"Cleaning up {feat_branch}")

        run_git(["checkout", "QA"], description="Switching to QA")
        qa_merge_args = ["merge", "--no-ff"] + signing_flag_or_empty() + ["dev"]
        run_git(qa_merge_args, description=f"Promoting {phase_name} dev code to QA")

        if phase["qa_fixes"]:
            log_event(f">>> Applying QA Fixes for {phase_name}")
            for fix in phase["qa_fixes"]:
                fix_base = fix["name"] if fix["name"].startswith("fix/") else f"fix/{fix['name']}"
                fix_branch = unique_work_branch(fix_base)

                run_git(["checkout", "QA"], description="Aligning QA")
                run_git(["checkout", "-b", fix_branch], description=f"Created {fix_branch}")

                for commit_hash in fix["commits"]:
                    if commit_already_in_head(commit_hash):
                        log_event(f"[~] Skipping {commit_hash}: already present in HEAD for {fix_branch}")
                        continue
                    shelve_conflicting_untracked(commit_hash)
                    orig_date = get_commit_date(commit_hash)
                    env_dates = {
                        "GIT_AUTHOR_DATE": orig_date,
                        "GIT_COMMITTER_DATE": orig_date
                    }
                    cherry_pick_with_fallback(
                        commit_hash,
                        env_dates,
                        description=f"Cherry-picking fix {commit_hash}"
                    )

                run_git(["checkout", "QA"], description="Switching to QA")
                fix_merge_args = ["merge", "--no-ff"] + signing_flag_or_empty() + [fix_branch]
                run_git(fix_merge_args, description=f"Merging fix {fix_branch} into QA")
                run_git(["branch", "-d", fix_branch], description="Cleaning up fix branch")

            run_git(["checkout", "dev"], description="Switching to dev for Sync")
            dev_sync_args = ["merge", "--no-ff"] + signing_flag_or_empty() + ["QA"]
            run_git(dev_sync_args, description=f"Syncing QA fixes back to dev")

        run_git(["checkout", "main"], description="Switching to main")
        main_merge_args = ["merge", "--no-ff"] + signing_flag_or_empty() + ["QA"]
        run_git(main_merge_args, description=f"Releasing {phase_name} to main")

        tag_name = phase["release_tag"]
        if USE_SIGNING:
            tag_args = ["tag", "-s", tag_name, "-m", f"Release {tag_name}"]
            tag_desc = f"Signed Tag: {tag_name}"
        else:
            tag_args = ["tag", "-a", tag_name, "-m", f"Release {tag_name}"]
            tag_desc = f"Annotated Tag: {tag_name}"
        run_git(tag_args, description=tag_desc)

    log_event("\n>>> RECONSTRUCTION COMPLETE <<<")
    log_event("Verify history with: git log --graph --oneline --all")
    if SHELVED_PATHS:
        log_event(f"[!] Shelved untracked paths were moved to: {SHELVED_ROOT}")
        log_event("[!] Review and manually restore/delete them as needed after verification.")
    log_event(f"[+] Execution log: {LOG_FILE}")


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        fail_critical("\n[!] User aborted execution.")
