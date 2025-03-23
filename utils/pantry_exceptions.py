class CategoryNotFoundError(Exception):
    """Raised when an item is not found in the pantry."""

    def __init__(self, category, available_categories):
        self.category = category
        super().__init__(
            f"Category '{self.category}' not found in the pantry. Please use one of the available categories {available_categories} "
        )


class ItemNotFoundError(Exception):
    """Raised when an item is expired."""

    def __init__(self, item_name):
        self.item_name = item_name
        super().__init__(f"Item {item_name} not found, cannot be excluded from pantry.")


class InvalidItemName(Exception):
    def __init__(self, item_name):
        self.item_name = item_name
        super().__init__(f"Invalid Item {item_name}, please use text")


class InvaliExpiryDate(Exception):
    def __init__(self, expiry_date):
        self.expiry_date = expiry_date
        super().__init__(
            f"Invalid Expiry Date: {expiry_date}, please use the 'YYYY-MM-DD' format"
        )
