from datetime import datetime

from pantry.categories import Categories
from pantry.pantry_db import PantryDB
from utils.pantry_exceptions import (
    CategoryNotFoundError,
    InvalidItemName,
    InvaliExpiryDate,
    ItemNotFoundError,
    QuantityError,
)


class Pantry:
    def __init__(self, pantry_db: PantryDB):
        self.available_categories = Categories.available_categories()
        self.pantry_db = pantry_db

    def add_item(self, item_name: str, category: str, quantity: int, expiry_date: str):
        self.validate_all_inputs(item_name, category, quantity, expiry_date)

        pantry_item = self.get_specific_pantry_item(item_name, category, expiry_date)
        if pantry_item:
            self.pantry_db.update_item_in_db(
                item_name, category, pantry_item[2] + quantity, expiry_date
            )
            return

        self.pantry_db.add_item_to_db(
            item_name,
            category,
            quantity,
            expiry_date,
            datetime.today().strftime("%Y-%m-%d"),
        )

    def remove_item(
        self, item_name: str, category: str, quantity: int, expiry_date: str
    ):
        self.validate_all_inputs(item_name, category, quantity, expiry_date)

        pantry_item = self.get_specific_pantry_item(item_name, category, expiry_date)
        if not pantry_item:
            raise ItemNotFoundError(item_name)

        if pantry_item and pantry_item[2] >= quantity:
            self.pantry_db.update_item_in_db(
                item_name, category, pantry_item[2] - quantity, expiry_date
            )
            return
        elif pantry_item and pantry_item[2] < quantity:
            raise QuantityError(pantry_item[2], quantity)
        else:
            raise ItemNotFoundError(item_name)

    def get_all_pantry_items(self):
        results = self.pantry_db.check_all_pantry_items()
        header = f"| {'Item Name':<41}| {'Category':<21}| {'Quantity':<9}| {'Expiry Date':<12}|"
        separator = f"+{'-' * 42}+{'-' * 22}+{'-' * 10}+{'-' * 13}+"
        print(separator)
        print(header)
        print(separator)
        for row in results.fetchall():
            # Print in a nice format:
            print(f"| {row[0]:<41}| {row[1]:<21}| {row[2]:<9}| {row[3]:<12}|")
            print(separator)
        return results.fetchall()

    def get_specific_pantry_item(self, item_name, category, expiry_date):
        results = self.pantry_db.check_specific_pantry_item(
            item_name, category, expiry_date
        )
        return results.fetchone()

    def validate_item_name(self, item_name):
        if not isinstance(item_name, str):
            raise InvalidItemName(item_name)

    def validate_category(self, category):
        if category not in self.available_categories:
            raise CategoryNotFoundError(category, self.available_categories)

    def validate_quantity(self, quantity):
        if not isinstance(quantity, int):
            raise ValueError("Quantity must be a positive number")

        if quantity < 1:
            raise ValueError("Quantity must be a positive number")

    def validate_date(self, date):
        try:
            datetime.strptime(date, "%Y-%m-%d")
        except Exception:
            raise InvaliExpiryDate(date)

    def validate_all_inputs(self, item_name, category, quantity, expiry_date):
        self.validate_item_name(item_name)
        self.validate_category(category)
        self.validate_quantity(quantity)
        self.validate_date(expiry_date)
