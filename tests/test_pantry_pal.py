import unittest

from pantry.pantry import Pantry
from utils.pantry_exceptions import ItemDoesNotExistError, CategoryNotFoundError
from pantry.pantry_db import PantryDB


class TestPantry(unittest.TestCase):
    def setUp(self):
        self.pantry = Pantry(PantryDB(":memory:"))
        self.today = "2025-01-01"
        return self.pantry

    def test_add_item(self):
        results = self.pantry.add_item("rice 1 kg", "grains", "2025-12-01")

        self.assertEqual(results[0], "rice 1 kg")
        self.assertEqual(results[1], "grains")
        self.assertEqual(results[2], "2025-12-01")

    def test_remove_item(self):
        self.pantry.add_item("rice 1 kg", "grains", "2025-12-01")
        self.pantry.remove_item("rice 1 kg", "grains", "2025-12-01")

        self.assertEqual(len(self.pantry["grains"]), 0)

    def test_remove_item_doesnt_exist(self):
        self.pantry.add_item("rice 1 kg", "grains", "2025-12-01")
        print(self.pantry)
        self.pantry.remove_item("rice 3 kg", "grains", "2025-12-01")

        self.assertRaises(ItemDoesNotExistError)("Item does not exist in pantry")
        self.assertEqual(len(self.pantry["grains"]), 1)
