import unittest

from pantry.pantry import Pantry
from utils.pantry_exceptions import (
    ItemNotFoundError,
    CategoryNotFoundError,
    InvalidItemName,
    InvaliExpiryDate,
)
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

    def test_add_multiple_items(self):
        self.pantry.add_item("rice 1 kg", "grains", "2025-12-01")
        self.pantry.add_item("rice 2 kg", "grains", "2025-12-01")
        self.pantry.add_item("rice 3 kg", "grains", "2025-12-01")

        pantry_items = self.pantry.get_pantry_items()
        self.assertEqual(len(pantry_items), 3)

    def test_add_items_different_categories(self):
        self.pantry.add_item("rice 1 kg", "grains", "2025-12-01")
        self.pantry.add_item("apple", "fruits", "2025-12-01")

        pantry_items = self.pantry.get_pantry_items()
        self.assertEqual(len(pantry_items), 2)

    def test_add_item_invalid_expiry_date(self):
        self.assertRaises(
            InvaliExpiryDate, self.pantry.add_item, "rice 1 kg", "grains", "2025-13-01"
        )

    def test_add_item_invalid_name(self):
        self.assertRaises(
            InvalidItemName, self.pantry.add_item, None, "grains", "2025-13-01"
        )

    def test_add_item_invalid_category_none(self):
        self.assertRaises(
            CategoryNotFoundError, self.pantry.add_item, "rice", None, "2025-13-01"
        )

    def test_add_item_invalid_date_none(self):
        self.assertRaises(
            InvaliExpiryDate, self.pantry.add_item, "rice", "grains", None
        )


if __name__ == "__main__":
    unittest.main()
