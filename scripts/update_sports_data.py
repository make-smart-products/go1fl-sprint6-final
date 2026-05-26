#!/usr/bin/env python3
"""Refresh the sports statistics snapshot used by the static dashboard.

The project intentionally keeps the public site dependency-free. This script is
safe to run from cron or GitHub Actions every hour: it preserves the curated
statistics in data/sports-stats.json, refreshes timestamps, and checks the
official league pages that are used as primary sources.
"""

from __future__ import annotations

import json
import re
import sys
from datetime import datetime, timedelta, timezone
from pathlib import Path
from typing import Any
from urllib.error import HTTPError, URLError
from urllib.request import Request, urlopen


ROOT = Path(__file__).resolve().parents[1]
DATA_FILE = ROOT / "data" / "sports-stats.json"
REQUEST_TIMEOUT_SECONDS = 20
USER_AGENT = (
    "SportStatRF/1.0 (+https://github.com/Yandex-Practicum/go1fl-sprint6-final; "
    "hourly source health check)"
)


def utc_now() -> datetime:
    return datetime.now(timezone.utc).replace(microsecond=0)


def isoformat(value: datetime) -> str:
    return value.isoformat().replace("+00:00", "Z")


def extract_title(html: str) -> str:
    match = re.search(r"<title[^>]*>(.*?)</title>", html, flags=re.IGNORECASE | re.DOTALL)
    if not match:
        return ""

    title = re.sub(r"\s+", " ", match.group(1)).strip()
    return title[:180]


def check_source(source: dict[str, Any], checked_at: datetime) -> dict[str, Any]:
    request = Request(
        source["url"],
        headers={
            "User-Agent": USER_AGENT,
            "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
        },
    )

    try:
        with urlopen(request, timeout=REQUEST_TIMEOUT_SECONDS) as response:
            body = response.read(32768).decode("utf-8", errors="ignore")
            status = response.status
            return {
                "id": source["id"],
                "name": source["name"],
                "ok": 200 <= status < 400,
                "status": status,
                "checkedAt": isoformat(checked_at),
                "title": extract_title(body),
            }
    except HTTPError as error:
        return {
            "id": source["id"],
            "name": source["name"],
            "ok": False,
            "status": error.code,
            "checkedAt": isoformat(checked_at),
            "error": str(error.reason),
        }
    except (OSError, URLError) as error:
        return {
            "id": source["id"],
            "name": source["name"],
            "ok": False,
            "status": 0,
            "checkedAt": isoformat(checked_at),
            "error": str(error),
        }


def refresh_snapshot() -> dict[str, Any]:
    with DATA_FILE.open("r", encoding="utf-8") as file:
        data = json.load(file)

    now = utc_now()
    interval = int(data.get("refreshIntervalMinutes", 60))

    data["generatedAt"] = isoformat(now)
    data["nextUpdateAt"] = isoformat(now + timedelta(minutes=interval))
    data["sourceHealth"] = [check_source(source, now) for source in data["sources"]]

    return data


def main() -> int:
    data = refresh_snapshot()
    with DATA_FILE.open("w", encoding="utf-8") as file:
        json.dump(data, file, ensure_ascii=False, indent=2)
        file.write("\n")

    healthy = sum(1 for source in data["sourceHealth"] if source["ok"])
    print(f"Updated {DATA_FILE.relative_to(ROOT)}; sources healthy: {healthy}/{len(data['sourceHealth'])}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
