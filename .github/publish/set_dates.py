"""Set publication dates once, preserving existing TOML front matter and body."""

from datetime import datetime
from pathlib import Path
import tomllib
from zoneinfo import ZoneInfo


def set_dates(root: Path, now: datetime) -> list[Path]:
    updates = []
    for path in sorted(root.rglob("*.md")):
        if path.name == "_index.md":
            continue
        content = path.read_bytes().decode("utf-8")
        lines = content.splitlines(keepends=True)
        if not lines or lines[0].strip() != "+++":
            raise ValueError(f"{path}: expected TOML front matter")
        end = next((i for i in range(1, len(lines)) if lines[i].strip() == "+++"), None)
        if end is None:
            raise ValueError(f"{path}: unclosed front matter")
        metadata = tomllib.loads("".join(lines[1:end]))
        if "date" in metadata or metadata.get("draft", False):
            continue
        newline = "\r\n" if lines[0].endswith("\r\n") else "\n"
        date = now.isoformat(timespec="seconds")
        updates.append((path, lines[0] + f"date = '{date}'{newline}" + "".join(lines[1:])))
    # Validate all front matter before changing any file.
    for path, content in updates:
        path.write_bytes(content.encode("utf-8"))
    return [path for path, _ in updates]


if __name__ == "__main__":
    for path in set_dates(Path("content/post"), datetime.now(ZoneInfo("Asia/Tokyo"))):
        print(path)
