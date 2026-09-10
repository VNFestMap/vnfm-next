from __future__ import annotations

from pathlib import Path


def _decode_mysql_string(raw: str) -> str:
    out: list[str] = []
    i = 0
    while i < len(raw):
        ch = raw[i]
        if ch != "\\":
            out.append(ch)
            i += 1
            continue
        if i + 1 >= len(raw):
            out.append(ch)
            break
        nxt = raw[i + 1]
        mapping = {
            "0": "\0",
            "b": "\b",
            "n": "\n",
            "r": "\r",
            "t": "\t",
            "Z": "\x1a",
            "\\": "\\",
            "'": "'",
            '"': '"',
        }
        out.append(mapping.get(nxt, nxt))
        i += 2
    return "".join(out)


def parse_mysql_tuples(body: str) -> list[list[str | None]]:
    rows: list[list[str | None]] = []
    current: list[str | None] = []
    i = 0
    n = len(body)
    while i < n:
        ch = body[i]
        if ch in " \t\r\n,":
            i += 1
            continue
        if ch == ";":
            break
        if ch != "(":
            i += 1
            continue
        i += 1
        current = []
        while i < n:
            while i < n and body[i] in " \t\r\n":
                i += 1
            if i >= n:
                break
            if body[i] == ")":
                rows.append(current)
                i += 1
                break
            if body[i] == ",":
                i += 1
                continue
            if body.startswith("NULL", i) and (
                i + 4 == n or body[i + 4] in ",) \t\r\n"
            ):
                current.append(None)
                i += 4
                continue
            if body[i] == "'":
                i += 1
                buf: list[str] = []
                while i < n:
                    if body[i] == "\\" and i + 1 < n:
                        buf.append(body[i : i + 2])
                        i += 2
                        continue
                    if body[i] == "'" and i + 1 < n and body[i + 1] == "'":
                        buf.append("''")
                        i += 2
                        continue
                    if body[i] == "'":
                        i += 1
                        break
                    buf.append(body[i])
                    i += 1
                value = "".join(buf).replace("''", "'")
                current.append(_decode_mysql_string(value))
                continue
            start = i
            while i < n and body[i] not in ",)":
                i += 1
            current.append(body[start:i].strip())
    return rows


def extract_insert_rows(dump_text: str, table: str) -> list[list[str | None]]:
    needle = f"INSERT INTO `{table}` VALUES"
    rows: list[list[str | None]] = []
    start = 0
    while True:
        idx = dump_text.find(needle, start)
        if idx < 0:
            break
        body_start = idx + len(needle)
        i = body_start
        n = len(dump_text)
        depth = 0
        in_str = False
        while i < n:
            ch = dump_text[i]
            if in_str:
                if ch == "\\" and i + 1 < n:
                    i += 2
                    continue
                if ch == "'" and i + 1 < n and dump_text[i + 1] == "'":
                    i += 2
                    continue
                if ch == "'":
                    in_str = False
                i += 1
                continue
            if ch == "'":
                in_str = True
                i += 1
                continue
            if ch == "(":
                depth += 1
            elif ch == ")":
                depth -= 1
            elif ch == ";" and depth == 0:
                rows.extend(parse_mysql_tuples(dump_text[body_start:i]))
                i += 1
                break
            i += 1
        start = i
    return rows


def load_dump_tables(path: Path, tables: list[str]) -> dict[str, list[list[str | None]]]:
    text = path.read_text(encoding="utf-8", errors="replace")
    return {name: extract_insert_rows(text, name) for name in tables}
