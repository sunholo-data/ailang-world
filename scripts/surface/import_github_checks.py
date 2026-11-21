#!/usr/bin/env python3
"""Capture one GitHub PR's checks into a local World daemon state directory."""

import argparse
from datetime import datetime, timezone
import json
import os
from pathlib import Path
import re
import subprocess
import sys
import tempfile

SCHEMA = "world/surface-github-checks/v1"
REPO_ID = re.compile(r"[a-z][a-z0-9-]{0,31}")
SLUG = re.compile(r"[a-zA-Z0-9][a-zA-Z0-9-]{0,38}/[a-zA-Z0-9_][a-zA-Z0-9_.-]{0,99}")
HEAD = re.compile(r"[0-9a-f]{40}")
GITHUB_HOST = "github.com"
ORIGIN_PREFIXES = (f"https://{GITHUB_HOST}/", f"git@{GITHUB_HOST}:", f"ssh://git@{GITHUB_HOST}/")
STATUS = {
    "EXPECTED", "ERROR", "FAILURE", "PENDING", "SUCCESS",
    "IN_PROGRESS", "QUEUED", "REQUESTED", "WAITING",
}
CONCLUSION = {
    "ACTION_REQUIRED", "CANCELLED", "FAILURE", "NEUTRAL", "SKIPPED", "STALE",
    "STARTUP_FAILURE", "SUCCESS", "TIMED_OUT",
}


def command(args, timeout):
    return subprocess.run(
        args, stdout=subprocess.PIPE, stderr=subprocess.DEVNULL, check=True,
        timeout=timeout, text=True, encoding="utf-8",
    ).stdout


def origin_slug(origin):
    for prefix in ORIGIN_PREFIXES:
        if origin.startswith(prefix):
            slug = origin[len(prefix):].removesuffix(".git")
            if SLUG.fullmatch(slug):
                return slug.lower()
    raise ValueError("The checkout origin must be a plain GitHub.com repository remote.")


def capture(repo, repo_id, pr):
    if Path(command(["git", "-C", repo, "rev-parse", "--show-toplevel"], 5).strip()).resolve() != repo:
        raise ValueError("Pass the checkout root for the configured World repo.")
    origin = command(["git", "-C", repo, "config", "--local", "--get", "remote.origin.url"], 5)
    slug = origin_slug(origin.strip())
    result = json.loads(command([
        "gh", "pr", "view", str(pr), "--repo", "github.com/" + slug,
        "--json", "number,headRefOid,url,statusCheckRollup",
    ], 30))
    if (
        result["number"] != pr
        or not HEAD.fullmatch(result["headRefOid"])
        or result["url"].lower() != f"https://github.com/{slug}/pull/{pr}"
    ):
        raise ValueError("GitHub returned a PR outside the requested repository or identity.")
    checks = []
    rows = result["statusCheckRollup"] or []
    if not isinstance(rows, list) or len(rows) > 250:
        raise ValueError("The PR check rollup exceeds the supported receipt bounds.")
    for row in rows:
        if row["__typename"] == "CheckRun":
            name = row["name"]
            status = row["status"]
            conclusion = row.get("conclusion") or ""
        elif row["__typename"] == "StatusContext":
            name, status, conclusion = row["context"], row["state"], ""
        else:
            raise ValueError("GitHub returned an unsupported check type.")
        if not isinstance(name, str) or not name.strip() or len(name) > 512:
            raise ValueError("GitHub returned an invalid check name.")
        if status == "COMPLETED":
            if conclusion and conclusion not in CONCLUSION:
                raise ValueError("GitHub returned an unsupported check conclusion.")
        elif conclusion or status not in STATUS:
            raise ValueError("GitHub returned an unsupported check status.")
        check = {"name": name, "status": status}
        if conclusion:
            check["conclusion"] = conclusion
        checks.append(check)
    return {
        "schema": SCHEMA, "repo": repo_id, "slug": slug, "pr": pr,
        "head": result["headRefOid"],
        "capturedAt": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "checks": checks,
    }


def save(state_dir, repo_id, receipt):
    folder = state_dir
    for name in ("surface-checks", repo_id):
        folder = folder / name
        folder.mkdir(mode=0o700, exist_ok=True)
        if folder.is_symlink() or not folder.is_dir():
            raise ValueError("The receipt directory must be a regular directory.")
    data = (json.dumps(receipt, ensure_ascii=False, separators=(",", ":")) + "\n").encode()
    if len(data) > 128 * 1024:
        raise ValueError("The check receipt exceeds 128 KiB.")
    tmp = None
    try:
        with tempfile.NamedTemporaryFile(dir=folder, prefix=".checks-", delete=False) as out:
            tmp = Path(out.name)
            out.write(data)
        tmp.replace(folder / (receipt["head"] + ".json"))
    finally:
        if tmp is not None:
            tmp.unlink(missing_ok=True)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--state-dir", type=Path, required=True, help="directory containing the daemon's world database")
    parser.add_argument("--repo-id", required=True, help="repository ID in the daemon's follow configuration")
    parser.add_argument("--repo", type=Path, required=True, help="local checkout configured for the repository ID")
    parser.add_argument("--pr", type=int, required=True, help="GitHub pull request number")
    args = parser.parse_args()
    if not REPO_ID.fullmatch(args.repo_id) or args.pr <= 0:
        parser.error("Use a valid configured repository ID and positive PR number.")
    if args.state_dir.is_symlink() or not args.state_dir.is_dir():
        parser.error("The state directory must already exist and be a regular directory.")
    try:
        receipt = capture(args.repo.resolve(strict=True), args.repo_id, args.pr)
        save(args.state_dir, args.repo_id, receipt)
    except (OSError, subprocess.SubprocessError, KeyError, TypeError, ValueError) as exc:
        # Raw subprocess/config errors can contain private checkout or remote data.
        msg = str(exc) if type(exc) is ValueError else "Could not verify the GitHub PR or save its receipt."
        print(msg, file=sys.stderr)
        return 1
    print(f"Imported {len(receipt['checks'])} checks for PR #{receipt['pr']} at {receipt['capturedAt']}.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
