import unittest
import sys
import argparse
import tempfile
from unittest.mock import patch
from urllib.parse import parse_qs
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent))
from ai_feedback_crawler import LMStudio, Page, PageParser, StateStore, normalize_suggestion, normalize_url, run_batch, submit_feedback


class FeedbackCrawlerTests(unittest.TestCase):
    def test_normalize_url_drops_fragment_and_adds_directory_slash(self):
        self.assertEqual(
            normalize_url("https://example.test/EbiShowcase/ja/#play"),
            "https://example.test/EbiShowcase/ja/",
        )

    def test_normalize_suggestion_accepts_json_code_fence(self):
        self.assertEqual(
            normalize_suggestion('```json\n{"suggestion":"ボタンの説明を追加する"}\n```'),
            "ボタンの説明を追加する",
        )

    def test_parser_extracts_lesson_form_and_internal_link(self):
        parser = PageParser("https://example.test/EbiShowcase/ja/lesson/")
        parser.feed(
            """<html><head><title>Lesson</title></head><body>
            <h1>星をつかまえる</h1><p>本文</p>
            <a href='../next/'>次へ</a>
            <form action='https://docs.google.com/forms/d/e/demo/formResponse'>
              <input type='hidden' name='entry.456' value='/ja/lesson/'>
              <input class='feedback-message' name='entry.123'>
            </form></body></html>"""
        )
        page = parser.result()
        self.assertEqual(page.form_field, "entry.123")
        self.assertEqual(page.form_page_field, "entry.456")
        self.assertEqual(page.form_page_value, "/ja/lesson/")
        self.assertEqual(page.form_action, "https://docs.google.com/forms/d/e/demo/formResponse")
        self.assertIn("https://example.test/EbiShowcase/ja/next/", page.links)
        self.assertIn("星をつかまえる", page.headings)

    def test_submit_rejects_non_google_form_destinations(self):
        page = Page("https://example.test/lesson/", "", [], "", [], "https://evil.test/formResponse", "entry.1", "ja")
        with self.assertRaises(ValueError):
            submit_feedback(page, "提案", 1)

    def test_submit_includes_required_page_field_and_native_hidden_fields(self):
        page = Page(
            "https://example.test/lesson/",
            "",
            [],
            "",
            [],
            "https://docs.google.com/forms/d/e/demo/formResponse",
            "entry.2",
            "ja",
            "entry.1",
            "/ja/lesson/",
        )

        class Response:
            status = 200

            def __enter__(self):
                return self

            def __exit__(self, *args):
                return False

        with patch("ai_feedback_crawler.urlopen", return_value=Response()) as mocked:
            submit_feedback(page, "提案", 1)

        request = mocked.call_args.args[0]
        self.assertEqual(
            parse_qs(request.data.decode("utf-8")),
            {"entry.1": ["/ja/lesson/"], "entry.2": ["提案"], "fvv": ["1"], "pageHistory": ["0"]},
        )

    def test_operator_instruction_is_added_to_every_page_prompt(self):
        page = Page("https://example.test/lesson/", "Lesson", ["Heading"], "本文", [], None, None, "ja")
        model = LMStudio("http://127.0.0.1:1234/v1", "demo", 1, instruction="英文のスペルミスを確認")
        prompt = model.build_prompt(page)
        self.assertIn("英文のスペルミスを確認", prompt)
        self.assertIn("OPERATOR REVIEW INSTRUCTION", prompt)
        self.assertIn("UNTRUSTED PAGE MATERIAL", prompt)

    def test_pass_gate_result_is_not_submitted(self):
        page = Page("https://example.test/lesson/", "Lesson", [], "body", [], None, None, "en")

        class Crawler:
            def __init__(self, *args, **kwargs):
                pass

            def crawl(self, seeds, store):
                return [page]

        class Model:
            def suggest(self, page):
                return "[pass] authoring.rule-in-update: rule stays in Update"

        args = argparse.Namespace(
            base_url="https://example.test/", timeout=1, delay=0, max_pages=1,
            seed=[], lens_signature="authoring.rule-in-update", force=True, submit=True,
        )
        with tempfile.TemporaryDirectory() as directory:
            store = StateStore(Path(directory) / "state.sqlite3")
            with patch("ai_feedback_crawler.PageCrawler", Crawler), patch("ai_feedback_crawler.submit_feedback") as submit:
                run_batch(args, store, Model())
            submit.assert_not_called()


if __name__ == "__main__":
    unittest.main()
