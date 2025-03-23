from pantry.categories import Categories
from utils.pantry_exceptions import CategoryNotFoundError
from datetime import datetime
from pantry.pantry_db import PantryDB


class Pantry:
    def __init__(self, pantry_db: PantryDB):
        self.available_categories = Categories.available_categories()
        self.pantry_db = pantry_db

    def add_item(self, item_name, category, expiry_date):
        if category not in self.available_categories:
            CategoryNotFoundError(category)

        self.pantry_db.add_item_to_db(
            item_name, category, expiry_date, datetime.today().strftime("%Y-%m-%d")
        )

    def remove_item(self, item_name, category, expiry_date):
        if category not in self.available_categories:
            CategoryNotFoundError(category)

        self.pantry_db.remove_item_from_db(item_name, category, expiry_date)

    def check_pantry_items(self):
        results = self.pantry_db.check_all_pantry_items()
        return results.fetchall()
