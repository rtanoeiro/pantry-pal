class CategoryNotFoundError(Exception):
    """Raised when an item is not found in the pantry."""

    def __init__(self, category):
        self.category = category
        super().__init__(f"Category '{self.category}' not found in the pantry.")


class ItemDoesNotExistError(Exception):
    """Raised when an item is expired."""

    def __init__(self, item_name):
        self.item_name = item_name
        super().__init__(f"Item {item_name} not found, cannot be excluded from pantry.")
