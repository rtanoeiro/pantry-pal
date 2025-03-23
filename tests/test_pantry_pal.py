import unittest

from pantry.pantry import Pantry
from utils.pantry_exceptions import ItemNotFoundError, CategoryNotFoundError
from pantry.pantry_db import PantryDB


class TestPantry(unittest.TestCase):
    def setUp(self):
        self.pantry = Pantry(PantryDB(":memory:"))
        self.today = "2025-01-01"
        return self.pantry

    def test_add_item(self):
        self.pantry.add_item("rice 1 kg", "grains", "2025-12-01")
        pantry_items = self.pantry.get_pantry_items()

        self.assertEqual(pantry_items[0][0], "rice 1 kg")
        self.assertEqual(pantry_items[0][1], "grains")
        self.assertEqual(pantry_items[0][2], "2025-12-01")

    def test_remove_item(self):
        self.pantry.add_item("rice 1 kg", "grains", "2025-12-01")
        self.pantry.remove_item("rice 1 kg", "grains", "2025-12-01")

        pantry_items = self.pantry.get_pantry_items()
        self.assertEqual(len(pantry_items), 0)

    def test_remove_item_doesnt_exist(self):
        self.assertRaises(
            ItemNotFoundError,
            self.pantry.remove_item,
            "rice 3 kg",
            "grains",
            "2025-12-01",
        )

    def test_add_item_category_doesnt_found(self):
        self.assertRaises(
            CategoryNotFoundError,
            self.pantry.add_item,
            "rice 1 kg",
            "grain",
            "2025-12-01",
        )
