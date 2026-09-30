from datetime import datetime
from pathlib import Path
from tempfile import TemporaryDirectory
import unittest

from set_dates import set_dates


class SetDatesTest(unittest.TestCase):
    def test_published_posts_only_and_repeat_builds(self):
        with TemporaryDirectory() as directory:
            root = Path(directory)
            originals = {
                "existing.md": b"+++\ndate = '2020-01-02T03:04:05+09:00'\ndraft = false\n+++\nBody\n",
                "draft.md": b"+++\ndraft = true\n+++\nBody\n",
                "_index.md": b"+++\ntitle = 'Posts'\n+++\n",
                "new/index.md": b"+++\ndraft = false\ntitle = 'New'\n+++\n\nBody +++\n",
                "implicit.md": b"+++\ntitle = 'Implicit'\n+++\n",
                "windows.md": b"+++\r\ndraft = false\r\n+++\r\nBody\r\n",
                "nested.md": b"+++\n[params]\ndate = 'custom'\n+++\n",
            }
            for name, content in originals.items():
                path = root / name
                path.parent.mkdir(parents=True, exist_ok=True)
                path.write_bytes(content)
            now = datetime.fromisoformat("2026-09-30T12:34:56+09:00")
            changed = set_dates(root, now)
            self.assertEqual({p.relative_to(root).as_posix() for p in changed},
                             {"new/index.md", "implicit.md", "windows.md", "nested.md"})
            for name, original in originals.items():
                newline = b"\r\n" if b"\r\n" in original else b"\n"
                expected = original
                if root / name in changed:
                    expected = original.replace(b"+++" + newline,
                        b"+++" + newline + b"date = '2026-09-30T12:34:56+09:00'" + newline, 1)
                self.assertEqual((root / name).read_bytes(), expected)
            self.assertEqual(set_dates(root, datetime.fromisoformat("2026-10-01T12:00:00+09:00")), [])

    def test_invalid_front_matter_does_not_partially_write(self):
        with TemporaryDirectory() as directory:
            root = Path(directory)
            original = b"+++\ndraft = false\n+++\n"
            (root / "a.md").write_bytes(original)
            for invalid in (b"---\ntitle: bad\n---\n", b"+++\ntitle = 'bad'\n", b"+++\ndraft = !!!\n+++\n"):
                with self.subTest(invalid=invalid):
                    (root / "z.md").write_bytes(invalid)
                    with self.assertRaises(ValueError):
                        set_dates(root, datetime.fromisoformat("2026-09-30T12:00:00+09:00"))
                    self.assertEqual((root / "a.md").read_bytes(), original)


if __name__ == "__main__":
    unittest.main()
