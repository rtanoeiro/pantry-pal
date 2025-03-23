#!/bin/bash

# Define color codes
RESET='\033[0m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'

# Function to display the menu
show_menu() {
    printf "${GREEN}Welcome to PantryPal! What would you like to do today?${RESET}\n"
    printf "${YELLOW}1. Add an item${RESET}\n"
    printf "${YELLOW}2. Remove an item${RESET}\n"
    printf "${YELLOW}3. Delete Items from table${RESET}\n"
    printf "${YELLOW}4. Check Available Items${RESET}\n"
    printf "${RED}5. Exit${RESET}\n"
    printf "${BLUE}Enter your choice (1/2/3/4/5): ${RESET}"
}

# Function to add an item
add_item() {
    printf "${GREEN}Enter the item name: ${RESET}"
    read item_name
    printf "${GREEN}Enter the item category: ${RESET}"
    read category
    printf "${GREEN}Enter the item quantity: ${RESET}"
    read quantity
    printf "${GREEN}Enter the expiry date (YYYY-MM-DD): ${RESET}"
    read expiry_date

    # Call the Python script for adding an item
    python3 main.py add "$item_name" "$category" "$quantity" "$expiry_date"
}

# Function to remove an item
remove_item() {
    printf "${GREEN}Enter the item name to remove: ${RESET}"
    read item_name
    printf "${GREEN}Enter the item category: ${RESET}"
    read category
    printf "${GREEN}Enter the item quantity: ${RESET}"
    read quantity
    printf "${GREEN}Enter the expiry date (YYYY-MM-DD): ${RESET}"
    read expiry_date

    # Call the Python script for removing an item
    python3 main.py remove "$item_name" "$category" "$quantity" "$expiry_date"
}

delete_item() {
    printf "${GREEN}Enter the item name to delete: ${RESET}"
    read item_name
    printf "${GREEN}Enter the item category: ${RESET}"
    read category
    printf "${GREEN}Enter the expiry date (YYYY-MM-DD): ${RESET}"
    read expiry_date

    # Call the Python script for removing an item
    python3 main.py delete "$item_name" "$category" "0" "$expiry_date"
}

# Function to check items
check_items() {
    # Call the Python script for checking items
    python3 main.py check "" "" "0" ""
}

while true; do
    show_menu
    read choice

    if [[ $choice == 1 ]]; then
        add_item
    elif [[ $choice == 2 ]]; then
        remove_item
    elif [[ $choice == 3 ]]; then
        delete_item
    elif [[ $choice == 4 ]]; then
        check_items
    elif [[ $choice == 5 ]]; then
        printf "${GREEN}Exiting PantryPal. Goodbye!${RESET}\n"
        break
    else
        printf "${RED}Invalid choice. Please try again.${RESET}\n"
    fi

    # Add a small delay for better user experience
    sleep 1
done
