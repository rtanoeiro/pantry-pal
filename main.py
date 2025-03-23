from pantry.pantry import Pantry
from pantry.pantry_db import PantryDB

if __name__ == "__main__":
    my_pantry = Pantry(PantryDB("pantry.db"))
    my_pantry.add_item("rice 1 kg", "grains", 1, "2025-12-01")
    my_pantry.add_item("chips", "snacks", 5, "2025-12-01")
    my_pantry.remove_item("rice 1 kg", "grains", 1, "2025-12-01")
    my_pantry.add_item("chips - vinegar and salt", "snacks", 2, "2025-12-01")
