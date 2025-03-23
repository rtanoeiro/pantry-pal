import unittest

from pantry import pantry


class TestPantry(unittest.TestCase):
    def setUp(self):
        self.pantry = pantry()
        self.today = "2025-01-01"
        return self.pantry

    def test_add_item(self):
        self.pantry.add_item("rice 1 kg", "2025-12-01")
