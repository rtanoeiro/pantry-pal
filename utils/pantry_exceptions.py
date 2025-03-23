class CategoryNotFoundError(Exception):
    """Raised when an item is not found in the pantry."""

    def __init__(self, category, available_categories):
        super().__init__(
            f"Category '{category}' not found in the pantry. Please use one of the available categories {available_categories} "
        )


class ItemNotFoundError(Exception):
    """Raised when an item is expired."""

    def __init__(self, item_name):
        super().__init__(f"Item {item_name} not found, cannot be excluded from pantry.")


class InvalidItemName(Exception):
    def __init__(self, item_name):
        super().__init__(f"Invalid Item {item_name}, please use text")


class InvaliExpiryDate(Exception):
    def __init__(self, expiry_date):
        super().__init__(
            f"Invalid Expiry Date: {expiry_date}, please use the 'YYYY-MM-DD' format"
        )


class QuantityError(Exception):
    def __init__(self, quantity_available, quantity_to_remove):
        super().__init__(
            f"Invalid Quantity Amount: {quantity_to_remove}, please remove less units than available: {quantity_available}"
        )
