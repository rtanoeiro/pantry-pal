from datetime import datetime

from pantry.categories import Categories
from pantry.pantry_db import PantryDB
from utils.pantry_exceptions import (
    CategoryNotFoundError,
    InvalidItemName,
    InvaliExpiryDate,
    ItemNotFoundError,
)


class Pantry:
    def __init__(self, pantry_db: PantryDB):
        self.available_categories = Categories.available_categories()
        self.pantry_db = pantry_db

    def add_item(self, item_name, category, expiry_date):
        self.validate_all_inputs(item_name, category, expiry_date)

        self.pantry_db.add_item_to_db(
            item_name, category, expiry_date, datetime.today().strftime("%Y-%m-%d")
        )

    def remove_item(self, item_name, category, expiry_date):
        self.validate_all_inputs(item_name, category, expiry_date)

        pantry_items = self.get_pantry_items()
        if not pantry_items:
            raise ItemNotFoundError(item_name)

        for item in pantry_items:
            if item[0] == item_name and item[1] == category and item[2] == expiry_date:
                self.pantry_db.remove_item_from_db(item_name, category, expiry_date)
            else:
                raise ItemNotFoundError(item_name)

    def get_pantry_items(self):
        results = self.pantry_db.check_all_pantry_items()
        return results.fetchall()

    def validate_item_name(self, item_name):
        if not isinstance(item_name, str):
            raise InvalidItemName(item_name)

    def validate_category(self, category):
        if category not in self.available_categories:
            raise CategoryNotFoundError(category, self.available_categories)

    def validate_date(self, date):
        try:
            datetime.strptime(date, "%Y-%m-%d")
        except Exception:
            raise InvaliExpiryDate(date)

    def validate_all_inputs(self, item_name, category, expiry_date):
        self.validate_item_name(item_name)
        self.validate_category(category)
        self.validate_date(expiry_date)
