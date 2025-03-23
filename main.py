from pantry.pantry import Pantry
from pantry.pantry_db import PantryDB

if __name__ == "__main__":
    my_pantry = Pantry(PantryDB(":memory:"))
    my_pantry.remove_item("rice 1 kg", "grains", "2025-12-01")
    items = my_pantry.get_pantry_items()
    print(items)
