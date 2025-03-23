from enum import Enum


class Categories(Enum):
    GRAINS = "grains"
    PASTA = "pasta"
    PROTEINS = "proteins"
    DAIRY = "dairy"
    FRUITS = "fruits"
    VEGETABLES = "vegetables"
    CANNED_GOODS = "canned goods"
    BAKING = "baking"
    SPICES = "spices"
    OILS = "oils"
    VINEGARS = "vinegars"
    CONDIMENTS = "condiments"
    SNACKS = "snacks"
    SWEETS = "sweets"
    BEVERAGES = "beverages"
    DRIED = "dried"
    MISCELLANEOUS = "miscellaneous"

    @classmethod
    def available_categories(cls):
        return [category.value for category in cls]
