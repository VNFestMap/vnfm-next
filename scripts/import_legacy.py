#!/usr/bin/env python3
"""Import VNFest MySQL dump + clubs/events JSON into vnfm-next Postgres.

Dry-run by default. Writes a report and SQL under --out-dir.
Does not invent NextMoe user ids: memberships/codes/notifications only land
when --user-map maps legacy_id -> oidc_id.
"""

from __future__ import annotations

import argparse
import csv
import json
import os
import re
import subprocess
import sys
from collections import Counter, defaultdict
from datetime import datetime, timezone
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(Path(__file__).resolve().parent))

from legacy_mysql import extract_insert_rows
from legacy_regions import normalize_china_province, normalize_japan_prefecture

REDACTED_INFO = "申请绑定后可见"
MEMBER_STATUSES = {"active", "pending"}
MEMBERSHIP_NOTIF_TYPES = {
    "join_approved",
    "join_rejected",
    "member_kicked",
    "role_changed",
}


def parse_args() -> argparse.Namespace:
    p = argparse.ArgumentParser(description="Import legacy VNFest data (dry-run default)")
    p.add_argument("--dump", type=Path, default=ROOT / "data/www_test_map_vnf_2026-09-09_04-00-01_mysql_data.sql")
    p.add_argument("--clubs-json", type=Path, help="Unredacted data/clubs.json from the old site")
    p.add_argument("--clubs-japan-json", type=Path, help="Unredacted data/clubs_japan.json")
    p.add_argument("--events-json", type=Path)
    p.add_argument("--registrations-json", type=Path)
    p.add_argument("--user-map", type=Path, help="CSV: legacy_id,email,username,oidc_id")
    p.add_argument("--out-dir", type=Path, default=ROOT / "tmp/legacy-import")
    p.add_argument("--owner-id", type=int, default=0, help="OIDC user id used as clubs.created_by fallback")
    p.add_argument("--club-id-start", type=int, default=0, help="0 = max(clubs.id)+1 when applying, else 1000")
    p.add_argument("--event-id-start", type=int, default=0)
    p.add_argument("--allow-redacted-info", action="store_true")
    p.add_argument("--notifications", choices=["none", "membership"], default="membership")
    p.add_argument("--events-mode", choices=["skip", "match-name"], default="match-name")
    p.add_argument("--apply", action="store_true", help="Run the generated SQL with psql")
    p.add_argument("--database-url", default=os.environ.get("DATABASE_URL", ""))
    return p.parse_args()


def load_json_list(path: Path | None, *keys: str) -> list:
    if not path:
        return []
    if not path.exists():
        raise SystemExit(f"missing file: {path}")
    payload = json.loads(path.read_text(encoding="utf-8"))
    if isinstance(payload, list):
        return payload
    if isinstance(payload, dict):
        for key in keys:
            rows = payload.get(key)
            if isinstance(rows, list):
                return rows
    raise SystemExit(f"cannot find list in {path} (tried {keys})")


def load_user_map(path: Path | None) -> dict[int, int]:
    if not path:
        return {}
    mapped: dict[int, int] = {}
    with path.open(newline="", encoding="utf-8") as fh:
        for row in csv.DictReader(fh):
            legacy = int(row["legacy_id"])
            oidc = (row.get("oidc_id") or "").strip()
            if oidc:
                mapped[legacy] = int(oidc)
    return mapped


def sql_str(value: str | None) -> str:
    if value is None:
        return "NULL"
    return "'" + str(value).replace("'", "''") + "'"


def sql_ts(value: str | None) -> str:
    raw = (value or "").strip()
    if not raw:
        return "NOW()"
    raw = raw.replace("T", " ").replace("Z", "")
    if "." in raw:
        raw = raw.split(".", 1)[0]
    try:
        datetime.strptime(raw[:19], "%Y-%m-%d %H:%M:%S")
    except ValueError:
        try:
            datetime.strptime(raw[:10], "%Y-%m-%d")
            raw = raw[:10] + " 00:00:00"
        except ValueError:
            return "NOW()"
    return f"TIMESTAMPTZ '{raw[:19]}+08'"


def sql_bool(value: bool) -> str:
    return "TRUE" if value else "FALSE"


def sql_int(value: int | None) -> str:
    return "NULL" if value is None else str(int(value))


def club_rows_from_json(rows: list, country: str) -> list[dict]:
    out = []
    for item in rows:
        out.append({**item, "country": item.get("country") or country})
    return out


def contact_hidden(item: dict) -> bool:
    if item.get("protected"):
        return True
    if item.get("visible_by_default"):
        return False
    return True


def pick_info(item: dict) -> str:
    info = str(item.get("info") or "")
    if info == REDACTED_INFO:
        return ""
    return info


def next_id(cur: list[int], start: int) -> int:
    cur[0] += 1
    if cur[0] < start:
        cur[0] = start
    return cur[0]


def club_name_keys(name: str) -> list[str]:
    cleaned = re.sub(r"[^\w\u4e00-\u9fff]+", " ", name or "", flags=re.UNICODE)
    parts = [p for p in re.split(r"[_\s]+", cleaned) if len(p) >= 4]
    if name and name not in parts:
        parts.append(name)
    return parts


def match_event_club(title: str, raw_text: str, clubs: list[dict]) -> dict | None:
    blob = f"{title} {raw_text}".lower()
    school_hits = [
        club
        for club in clubs
        if club.get("school") and len(club["school"]) >= 4 and club["school"].lower() in blob
    ]
    if len(school_hits) == 1:
        return school_hits[0]
    scored: list[tuple[int, dict]] = []
    for club in clubs:
        keys = club_name_keys(club["name"])
        if club.get("school"):
            keys.append(club["school"])
        hit = max((len(key) for key in keys if key.lower() in blob), default=0)
        if hit >= 4:
            scored.append((hit, club))
    if not scored:
        return None
    scored.sort(key=lambda item: item[0], reverse=True)
    if len(scored) > 1 and scored[0][0] == scored[1][0] and scored[0][0] < 8:
        return None
    return scored[0][1]


def main() -> None:
    args = parse_args()
    if not args.dump.exists():
        raise SystemExit(f"missing dump: {args.dump}")

    dump_text = args.dump.read_text(encoding="utf-8", errors="replace")
    users = extract_insert_rows(dump_text, "users")
    memberships = extract_insert_rows(dump_text, "club_memberships")
    codes = extract_insert_rows(dump_text, "club_verification_codes")
    notifications = extract_insert_rows(dump_text, "notifications")

    china = club_rows_from_json(load_json_list(args.clubs_json, "data"), "china") if args.clubs_json else []
    japan = club_rows_from_json(load_json_list(args.clubs_japan_json, "data"), "japan") if args.clubs_japan_json else []
    events = load_json_list(args.events_json, "events") if args.events_json else []
    registrations = load_json_list(args.registrations_json, "registrations") if args.registrations_json else []
    user_map = load_user_map(args.user_map)

    clubs = china + japan
    redacted = sum(1 for c in clubs if str(c.get("info") or "") == REDACTED_INFO)
    if clubs and redacted / len(clubs) > 0.2 and not args.allow_redacted_info:
        raise SystemExit(
            f"{redacted}/{len(clubs)} clubs have redacted info ('{REDACTED_INFO}'). "
            "Pass the production JSON files, not the public API dump, or set --allow-redacted-info."
        )

    warnings: list[str] = []
    unknown_regions: list[dict] = []
    club_plan: list[dict] = []
    for item in clubs:
        country = "japan" if item.get("country") == "japan" else "china"
        legacy_id = int(item.get("id") or 0)
        name = str(item.get("display_name") or item.get("name") or "").strip()
        provinces = item.get("provinces") if isinstance(item.get("provinces"), list) else []
        extra_provinces = [str(x).strip() for x in provinces[1:] if str(x).strip()]
        if country == "japan":
            region, ok = normalize_japan_prefecture(str(item.get("prefecture") or item.get("province") or ""))
            province, prefecture = "", region
        else:
            region, ok = normalize_china_province(str(item.get("province") or ""))
            province, prefecture = region, ""
        if region and not ok:
            unknown_regions.append({"country": country, "legacy_id": legacy_id, "value": region})
        remark = str(item.get("remark") or "")
        if extra_provinces:
            extra = "兼: " + "、".join(extra_provinces)
            remark = f"{remark}\n{extra}".strip()
            warnings.append(f"club {country}:{legacy_id} extra provinces stored in remark: {extra_provinces}")
        club_plan.append(
            {
                "legacy_id": legacy_id,
                "legacy_country": country,
                "name": name,
                "school": str(item.get("school") or ""),
                "province": province,
                "prefecture": prefecture,
                "city": str(item.get("city") or ""),
                "type": str(item.get("type") or "school") or "school",
                "info": pick_info(item),
                "remark": remark,
                "logo_url": str(item.get("logo_url") or ""),
                "external_links": str(item.get("external_links") or ""),
                "contact_hidden": contact_hidden(item),
                "created_at": str(item.get("created_at") or ""),
            }
        )

    members_by_club: dict[tuple[str, int], list[dict]] = defaultdict(list)
    skipped_memberships: Counter[str] = Counter()
    for row in memberships:
        rec = {
            "legacy_id": int(row[0]),
            "legacy_user_id": int(row[1]),
            "legacy_club_id": int(row[2]),
            "role": row[3] or "member",
            "status": row[4] or "pending",
            "joined_at": row[5],
            "left_at": row[6],
            "country": row[10] or "china",
            "contact_account": row[12] or row[8] or "",
            "apply_role": row[9] or "member",
            "apply_reason": row[15] or "",
            "reviewed_at": row[17],
            "reviewed_by": int(row[18]) if row[18] else None,
        }
        if rec["status"] not in MEMBER_STATUSES:
            skipped_memberships[rec["status"]] += 1
            continue
        members_by_club[(rec["country"], rec["legacy_club_id"])].append(rec)

    mapped_members = 0
    unmapped_members = 0
    for recs in members_by_club.values():
        for rec in recs:
            if rec["legacy_user_id"] in user_map:
                mapped_members += 1
            else:
                unmapped_members += 1

    owner_id = args.owner_id
    if args.apply and owner_id <= 0:
        raise SystemExit("--owner-id (an existing vnfm/OIDC user id) is required with --apply")

    args.out_dir.mkdir(parents=True, exist_ok=True)
    club_id_start = args.club_id_start or 1000
    event_id_start = args.event_id_start or 1000
    if args.apply:
        if not args.database_url:
            raise SystemExit("DATABASE_URL or --database-url is required with --apply")
        club_id_start = max(club_id_start, query_max_id(args.database_url, "clubs") + 1)
        event_id_start = max(event_id_start, query_max_id(args.database_url, "events") + 1)

    club_ids: dict[tuple[str, int], int] = {}
    counter = [club_id_start - 1]
    sql: list[str] = [
        "BEGIN;",
        "INSERT INTO legacy_user_map (legacy_id, username, email, oidc_id, notes) VALUES",
    ]
    user_values = []
    for row in users:
        legacy_id = int(row[0])
        username = row[5] or ""
        email = row[13] or ""
        oidc = user_map.get(legacy_id)
        notes = row[8] or ""
        user_values.append(
            f"({legacy_id}, {sql_str(username)}, {sql_str(email)}, {sql_int(oidc)}, {sql_str(notes)})"
        )
    sql.append(",\n".join(user_values) + "\nON CONFLICT (legacy_id) DO UPDATE SET email = EXCLUDED.email, oidc_id = COALESCE(EXCLUDED.oidc_id, legacy_user_map.oidc_id);")

    club_values = []
    logos = []
    for club in club_plan:
        new_id = next_id(counter, club_id_start)
        club["new_id"] = new_id
        club_ids[(club["legacy_country"], club["legacy_id"])] = new_id
        created_by = owner_id
        reps = [
            rec
            for rec in members_by_club.get((club["legacy_country"], club["legacy_id"]), [])
            if rec["status"] == "active" and rec["role"] == "representative" and rec["legacy_user_id"] in user_map
        ]
        if reps:
            created_by = user_map[reps[0]["legacy_user_id"]]
        if club["logo_url"]:
            logos.append({"legacy": f"{club['legacy_country']}:{club['legacy_id']}", "url": club["logo_url"]})
        created = club["created_at"] or None
        club_values.append(
            "("
            + ", ".join(
                [
                    str(new_id),
                    sql_str(club["legacy_country"]),
                    sql_str(club["name"]),
                    sql_str(club["school"]),
                    sql_str(club["province"]),
                    sql_str(club["prefecture"]),
                    sql_str(club["city"]),
                    sql_str(club["type"]),
                    sql_str(club["info"]),
                    sql_str(club["remark"]),
                    sql_str(""),
                    sql_str(club["external_links"]),
                    sql_bool(club["contact_hidden"]),
                    str(created_by if created_by > 0 else 0),
                    sql_ts(created),
                    sql_ts(created),
                    str(club["legacy_id"]),
                    sql_str(club["legacy_country"]),
                ]
            )
            + ")"
        )
    if club_values:
        sql.append(
            "INSERT INTO clubs (id, country, name, school, province, prefecture, city, type, info, remark, logo_key, external_links, contact_hidden, created_by, created_at, updated_at, legacy_id, legacy_country) VALUES"
        )
        sql.append(
            ",\n".join(club_values)
            + "\nON CONFLICT (legacy_country, legacy_id) WHERE legacy_id IS NOT NULL DO NOTHING;"
        )
    else:
        warnings.append("no clubs JSON provided — club catalog not imported")

    member_values = []
    for key, recs in members_by_club.items():
        club_id = club_ids.get(key)
        if not club_id:
            warnings.append(f"memberships for missing club {key[0]}:{key[1]} skipped")
            continue
        for rec in recs:
            oidc = user_map.get(rec["legacy_user_id"])
            if not oidc:
                continue
            reviewed_by = user_map.get(rec["reviewed_by"]) if rec["reviewed_by"] else None
            member_values.append(
                "("
                + ", ".join(
                    [
                        str(oidc),
                        str(club_id),
                        sql_str(key[0]),
                        sql_str(rec["role"]),
                        sql_str(rec["status"]),
                        sql_str(rec["contact_account"]),
                        sql_str(rec["apply_reason"]),
                        sql_str(rec["apply_role"]),
                        sql_int(reviewed_by),
                        sql_ts(rec["reviewed_at"]) if rec["reviewed_at"] else "NULL",
                        sql_ts(rec["joined_at"]),
                    ]
                )
                + ")"
            )
    if member_values:
        sql.append(
            "INSERT INTO club_memberships (user_id, club_id, country, role, status, contact_account, apply_reason, apply_role, reviewed_by, reviewed_at, joined_at) VALUES"
        )
        sql.append(",\n".join(member_values) + "\nON CONFLICT (user_id, club_id) DO NOTHING;")

    code_values = []
    skipped_codes = 0
    for row in codes:
        key = (row[9] or "china", int(row[1]))
        club_id = club_ids.get(key)
        created_by = user_map.get(int(row[3])) or owner_id
        if not club_id or created_by <= 0:
            skipped_codes += 1
            continue
        code_values.append(
            "("
            + ", ".join(
                [
                    str(club_id),
                    sql_str(row[2]),
                    sql_int(int(row[4] or 1)),
                    sql_int(int(row[5] or 0)),
                    sql_ts(row[6]) if row[6] else "NULL",
                    sql_bool(row[7] in ("1", 1, True)),
                    str(created_by),
                    sql_ts(row[8]),
                ]
            )
            + ")"
        )
    if code_values:
        sql.append(
            "INSERT INTO club_verification_codes (club_id, code, max_uses, use_count, expires_at, is_active, created_by, created_at) VALUES"
        )
        sql.append(",\n".join(code_values) + "\nON CONFLICT (code) DO NOTHING;")

    event_ids: dict[int, int] = {}
    event_counter = [event_id_start - 1]
    unmatched_events = []
    event_values = []
    if args.events_mode == "match-name" and events:
        sql.append(
            "INSERT INTO events (id, club_id, title, description, location, cover_key, starts_at, ends_at, created_by, created_at, updated_at, legacy_id) VALUES"
        )
        for item in events:
            legacy_event_id = int(item.get("id") or 0)
            title = str(item.get("event") or "").strip()
            matched = match_event_club(title, str(item.get("raw_text") or ""), club_plan)
            if not matched:
                unmatched_events.append({"legacy_id": legacy_event_id, "title": title})
                continue
            new_id = next_id(event_counter, event_id_start)
            event_ids[legacy_event_id] = new_id
            location = ""
            raw = str(item.get("raw_text") or "")
            if " @ " in raw:
                location = raw.split(" @ ", 1)[1]
            created_by = owner_id if owner_id > 0 else 0
            event_values.append(
                "("
                + ", ".join(
                    [
                        str(new_id),
                        str(matched["new_id"]),
                        sql_str(title),
                        sql_str(str(item.get("description") or "")),
                        sql_str(location),
                        sql_str(""),
                        sql_ts(str(item.get("date") or "")),
                        sql_ts(str(item.get("date_end") or "")) if item.get("date_end") else "NULL",
                        str(created_by),
                        sql_ts(str(item.get("created_at") or item.get("date") or "")),
                        sql_ts(str(item.get("updated_at") or item.get("created_at") or "")),
                        str(legacy_event_id) if legacy_event_id else "NULL",
                    ]
                )
                + ")"
            )
        if event_values:
            sql.append(",\n".join(event_values) + "\nON CONFLICT (legacy_id) WHERE legacy_id IS NOT NULL DO NOTHING;")
        else:
            sql.pop()

    reg_values = []
    skipped_regs = 0
    for item in registrations:
        event_id = event_ids.get(int(item.get("event_id") or 0))
        oidc = user_map.get(int(item.get("user_id") or 0))
        if not event_id or not oidc:
            skipped_regs += 1
            continue
        reg_values.append(
            f"({event_id}, {oidc}, {sql_ts(str(item.get('registered_at') or ''))})"
        )
    if reg_values:
        sql.append("INSERT INTO event_registrations (event_id, user_id, created_at) VALUES")
        sql.append(",\n".join(reg_values) + "\nON CONFLICT (event_id, user_id) DO NOTHING;")

    notif_values = []
    skipped_notifs = Counter()
    if args.notifications == "membership":
        for row in notifications:
            ntype = row[2] or ""
            if ntype not in MEMBERSHIP_NOTIF_TYPES:
                skipped_notifs[ntype] += 1
                continue
            oidc = user_map.get(int(row[1]))
            if not oidc:
                skipped_notifs["unmapped_user"] += 1
                continue
            notif_values.append(
                "("
                + ", ".join(
                    [
                        str(oidc),
                        sql_str(ntype),
                        sql_str(row[3] or ""),
                        sql_str(row[4] or ""),
                        sql_str(""),
                        sql_str(row[6] or ""),
                        sql_int(int(row[7] or 0)),
                        sql_bool(row[8] in ("1", 1, True)),
                        sql_ts(row[9]),
                        sql_ts(row[10]) if row[10] else "NULL",
                    ]
                )
                + ")"
            )
        if notif_values:
            sql.append(
                "INSERT INTO notifications (user_id, type, title, message, link, related_type, related_id, is_read, created_at, read_at) VALUES"
            )
            sql.append(",\n".join(notif_values) + ";")
    else:
        skipped_notifs["disabled"] = len(notifications)

    if club_values:
        sql.append("SELECT setval(pg_get_serial_sequence('clubs','id'), (SELECT MAX(id) FROM clubs));")
    if event_values:
        sql.append("SELECT setval(pg_get_serial_sequence('events','id'), (SELECT MAX(id) FROM events));")
    sql.append("COMMIT;")

    if club_values and owner_id <= 0:
        warnings.append("clubs.created_by is 0 in the SQL — set --owner-id before --apply")

    sql_path = args.out_dir / "import.sql"
    sql_path.write_text("\n".join(sql) + "\n", encoding="utf-8")

    map_csv = args.out_dir / "users_extracted.csv"
    with map_csv.open("w", newline="", encoding="utf-8") as fh:
        writer = csv.DictWriter(fh, fieldnames=["legacy_id", "email", "username", "oidc_id", "notes"])
        writer.writeheader()
        for row in users:
            writer.writerow(
                {
                    "legacy_id": row[0],
                    "email": row[13] or "",
                    "username": row[5] or "",
                    "oidc_id": user_map.get(int(row[0]), ""),
                    "notes": row[8] or "",
                }
            )

    report = {
        "generated_at": datetime.now(timezone.utc).isoformat(),
        "sources": {
            "dump": str(args.dump),
            "clubs_json": str(args.clubs_json) if args.clubs_json else None,
            "clubs_japan_json": str(args.clubs_japan_json) if args.clubs_japan_json else None,
            "events_json": str(args.events_json) if args.events_json else None,
            "user_map": str(args.user_map) if args.user_map else None,
        },
        "counts": {
            "users_in_dump": len(users),
            "users_with_email": sum(1 for r in users if r[13]),
            "users_mapped": len(user_map),
            "clubs_json": len(club_plan),
            "memberships_importable": mapped_members,
            "memberships_unmapped_user": unmapped_members,
            "memberships_skipped_status": dict(skipped_memberships),
            "codes_sql": len(code_values),
            "codes_skipped": skipped_codes,
            "events_matched": len(event_values),
            "events_unmatched": len(unmatched_events),
            "registrations_sql": len(reg_values),
            "registrations_skipped": skipped_regs,
            "notifications_sql": len(notif_values),
            "notifications_skipped": dict(skipped_notifs),
            "logos_remote": len(logos),
            "unknown_regions": len(unknown_regions),
        },
        "unknown_regions": unknown_regions,
        "unmatched_events": unmatched_events,
        "logos": logos,
        "warnings": warnings,
        "open_questions": [
            "Old user ids are not NextMoe OIDC ids. Fill oidc_id in users_extracted.csv and pass --user-map.",
            "Public clubs.php hides info as 申请绑定后可见. Need the server-side JSON files.",
            "Events have no club_id; matching is by name in title/raw_text.",
            "Remote logo URLs are not uploaded to R2.",
            "Site-level super_admin is not imported as vnfm authorization.",
        ],
    }
    report_path = args.out_dir / "report.json"
    report_path.write_text(json.dumps(report, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")

    print(json.dumps(report["counts"], ensure_ascii=False, indent=2))
    print(f"wrote {sql_path}")
    print(f"wrote {report_path}")
    print(f"wrote {map_csv} (fill oidc_id, then rerun with --user-map)")
    if args.apply:
        if owner_id <= 0:
            raise SystemExit("refusing to apply with created_by=0")
        run_psql(args.database_url, sql_path)


def query_max_id(dsn: str, table: str) -> int:
    out = subprocess.check_output(
        ["psql", dsn, "-tAc", f"SELECT COALESCE(MAX(id),0) FROM {table}"],
        text=True,
    )
    return int(out.strip() or 0)


def run_psql(dsn: str, sql_path: Path) -> None:
    subprocess.check_call(["psql", dsn, "-v", "ON_ERROR_STOP=1", "-f", str(sql_path)])
    print("applied", sql_path)


if __name__ == "__main__":
    main()
