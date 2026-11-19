import os
import unittest


class TestBenchmarkEnso(unittest.TestCase):
    def test_benchmark_enso_module_exists_and_valid_syntax(self):
        path = os.path.join(os.path.dirname(__file__), "benchmark_enso.py")
        self.assertTrue(os.path.exists(path))
        with open(path, "r", encoding="utf-8") as f:
            compile(f.read(), path, "exec")


if __name__ == "__main__":
    unittest.main()
