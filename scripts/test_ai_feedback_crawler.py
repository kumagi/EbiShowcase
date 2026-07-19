import unittest
import sys
import argparse
import json
import tempfile
from unittest.mock import patch
from urllib.parse import parse_qs
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent))
from ai_feedback_crawler import (
    LMStudio,
    Page,
    PageParser,
    StateStore,
    classify_page,
    format_lens_instruction,
    load_gate_review,
    normalize_suggestion,
    normalize_url,
    run_batch,
    submit_feedback,
    validate_gate_review_response,
    validate_args,
)


class FeedbackCrawlerTests(unittest.TestCase):
    def test_normalize_url_drops_fragment_and_adds_directory_slash(self):
        self.assertEqual(
            normalize_url("https://example.test/EbiShowcase/ja/#play"),
            "https://example.test/EbiShowcase/ja/",
        )

    def test_classify_page_distinguishes_catalog_page_kinds(self):
        self.assertEqual(classify_page("https://example.test/EbiShowcase/en/games/flappy/"), "core")
        self.assertEqual(classify_page("https://example.test/EbiShowcase/en/tracks/racing/lap-gates/"), "track")
        self.assertEqual(classify_page("https://example.test/EbiShowcase/en/tracks/visual-effects/vfx-walk/"), "vfx")
        self.assertEqual(classify_page("https://example.test/EbiShowcase/en/guides/setup/"), "guide")

    def test_inapplicable_gate_skips_model_request(self):
        page = Page("https://example.test/EbiShowcase/en/tracks/racing/lap-gates/", "Lesson", [], "body", [], None, None, "en")
        model = LMStudio(
            "http://127.0.0.1:1234/v1",
            "demo",
            1,
            gate_review=True,
            gate_ids=["vfx.play-vs-fx-split"],
            gate_applies_to=["vfx"],
            gate_languages=["ja", "en"],
        )
        with patch.object(model, "_json_request") as request:
            result = model.suggest(page)
        self.assertTrue(result.startswith("[pass]"))
        request.assert_not_called()

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

    def test_catalog_supplies_low_confidence_policy_and_structured_pedagogy_gate(self):
        lenses, policy = load_gate_review("pedagogy")
        code_gate = next(gate for gate in lenses if gate["id"] == "pedagogy.code-matches-impl")
        self.assertTrue(any("Never guess" in rule for rule in policy))
        self.assertIn("Two visible, path-labelled artifacts", code_gate["fail_when"])
        self.assertIn("only one side", code_gate["do_not_flag"])
        self.assertIn("both conflicting code fragments", code_gate["evidence_required"])
        exact, _ = load_gate_review("pedagogy.code-matches-impl")
        self.assertEqual([gate["id"] for gate in exact], ["pedagogy.code-matches-impl"])

    def test_gate_instruction_has_one_unambiguous_json_contract(self):
        lens = {
            "id": "pedagogy.code-matches-impl",
            "severity": "fail",
            "fail_when": "Two shown files conflict.",
            "do_not_flag": "Only one file is shown.",
            "evidence_required": "Quote both files.",
        }
        instruction = format_lens_instruction([lens], review_policy=["Never guess unseen files."])
        self.assertIn("OUTPUT: Return exactly one JSON object", instruction)
        self.assertIn("FAIL ONLY WHEN", instruction)
        self.assertIn("DO NOT FLAG", instruction)
        self.assertIn("EVIDENCE REQUIRED", instruction)
        self.assertNotIn('one key, suggestion', instruction)

    def test_gate_review_rejects_inference_without_exact_page_quote(self):
        page = Page("https://example.test/lesson/", "Lesson", ["REAL GO"], "main.go calls Run(Config).", [], None, None, "en")
        raw = json.dumps(
            {
                "gate_id": "pedagogy.code-matches-impl",
                "verdict": "fail",
                "evidence": "The hidden implementation probably uses a different loop.",
                "fix": "Rewrite the snippet.",
            }
        )
        checked = validate_gate_review_response(raw, page, {"pedagogy.code-matches-impl"})
        self.assertEqual(json.loads(checked)["verdict"], "pass")

    def test_gate_review_keeps_failure_with_exact_quote_and_direct_fix(self):
        page = Page("https://example.test/lesson/", "Lesson", [], "Draw says: score++", [], None, None, "en")
        raw = json.dumps(
            {
                "gate_id": "loop.draw-is-projection",
                "verdict": "fail",
                "evidence": "score++",
                "fix": "Move score++ to Update.",
            }
        )
        checked = validate_gate_review_response(raw, page, {"loop.draw-is-projection"})
        self.assertEqual(json.loads(checked)["verdict"], "fail")

    def test_gate_review_accepts_multiple_exact_quotes_separated_by_pipe(self):
        page = Page("https://example.test/lesson/", "Lesson", [], "ENTRY: Run(Config) EXCERPT: score++", [], None, None, "en")
        raw = json.dumps(
            {
                "gate_id": "pedagogy.code-matches-impl",
                "verdict": "fail",
                "evidence": "ENTRY: Run(Config) | EXCERPT: score++",
                "fix": "Label the excerpt as conceptual or show the matching implementation.",
            }
        )
        checked = validate_gate_review_response(raw, page, {"pedagogy.code-matches-impl"})
        self.assertEqual(json.loads(checked)["verdict"], "fail")

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

    def test_same_review_is_not_submitted_twice_even_with_force(self):
        page = Page("https://example.test/lesson/", "Lesson", [], "body", [], None, None, "en")

        class Crawler:
            def __init__(self, *args, **kwargs):
                pass

            def crawl(self, seeds, store):
                return [page]

        class Model:
            calls = 0

            def suggest(self, page):
                self.calls += 1
                return "Add a clearer goal."

        args = argparse.Namespace(
            base_url="https://example.test/", timeout=1, delay=0, max_pages=1,
            seed=[], lens_signature="goal", force=True, submit=True,
        )
        model = Model()
        with tempfile.TemporaryDirectory() as directory:
            store = StateStore(Path(directory) / "state.sqlite3")
            with patch("ai_feedback_crawler.PageCrawler", Crawler), patch("ai_feedback_crawler.submit_feedback") as submit:
                run_batch(args, store, model)
                run_batch(args, store, model)
            self.assertEqual(model.calls, 1)
            submit.assert_called_once()

    def test_changed_review_clears_old_submission_marker(self):
        page = Page("https://example.test/lesson/", "Lesson", [], "body", [], None, None, "en")
        with tempfile.TemporaryDirectory() as directory:
            store = StateStore(Path(directory) / "state.sqlite3")
            store.save(page, "First review", content_hash="review-hash-1")
            store.mark_submitted(page.url)
            self.assertIsNotNone(store.get(page.url)[2])

            store.save(page, "Changed review", content_hash="review-hash-2")
            self.assertIsNone(store.get(page.url)[2])

    def test_validation_rejects_unbounded_or_large_submissions(self):
        invalid = (
            argparse.Namespace(max_pages=1, submit=True, once=False, force=False),
            argparse.Namespace(max_pages=1, submit=True, once=True, force=True),
            argparse.Namespace(max_pages=26, submit=True, once=True, force=False),
        )
        for args in invalid:
            with self.subTest(args=args), self.assertRaises(SystemExit):
                validate_args(args)


if __name__ == "__main__":
    unittest.main()
