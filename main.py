from pantry.pantry import Pantry
from pantry.pantry_db import PantryDB

if __name__ == "__main__":
    my_pantry = Pantry(PantryDB(":memory:"))
    results = my_pantry.add_item("rice 1 kg", "grains", "2025-12-01")
    print(results)
