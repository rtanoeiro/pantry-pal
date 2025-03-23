from pantry.pantry import Pantry
from pantry.pantry_db import PantryDB
import sys


def main():
    action = sys.argv[1]
    item_name = sys.argv[2]
    category = sys.argv[3]
    quantity = int(sys.argv[4])
    expiry_date = sys.argv[5]

    my_pantry = Pantry(PantryDB("pantry.db"))
    if action == "add":
        print(f"Adding item: {item_name}, Category: {category}, Expiry: {expiry_date}")
        my_pantry.add_item(item_name, category, quantity, expiry_date)
    elif action == "remove":
        print(
            f" Trying to remove {quantity} item: {item_name}, Category: {category}, Expiry: {expiry_date}"
        )
        my_pantry.add_item(item_name, category, quantity, expiry_date)
    elif action == "check":
        my_pantry.get_all_pantry_items()
    else:
        print("Invalid option, try again")


if __name__ == "__main__":
    main()
