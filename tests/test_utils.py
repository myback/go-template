import unittest
from pathlib import Path

import go_template


class TestRender(unittest.TestCase):
    def test_add(self):
        values = {"Count": 12, "Material": "Wool"}

        render_data = go_template.render(Path('tests/sample.tmpl'), values).decode()

        with open('tests/test.txt') as f:
            test_data = f.read()

        self.assertEqual(render_data, test_data)


class TestRenderFromValuesFile(unittest.TestCase):
    def test_add(self):
        render_data = go_template.render_from_values_file(
            Path('tests/sample.tmpl'),
            Path('tests/values.yml')
        ).decode()

        with open('tests/test.txt') as f:
            test_data = f.read()

        self.assertEqual(render_data, test_data)


if __name__ == '__main__':
    unittest.main()
