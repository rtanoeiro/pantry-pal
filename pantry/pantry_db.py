import sqlite3


class PantryDB:
    def __init__(self, database_file):
        self.pantry_db = database_file
        self.db_conn = None
        self._check_db()

    def __enter__(self):
        if self.db_conn:
            return self.db_conn

        self.db_conn = sqlite3.connect(self.pantry_db)
        return self.db_conn

    def __exit__(self):
        if self.db_conn:
            self.db_conn.close()
            self.db_conn = None
            self._check_db()

    def get_cursor(self):
        return sqlite3.connect(self.pantry_db).cursor()

    def _check_db(self):
        with self.__enter__() as conn:
            cursor = conn.cursor()
            cursor.execute(
                """CREATE TABLE IF NOT EXISTS pantry (item_name TEXT, category TEXT, expiry_date TEXT, added_at TEXT);"""
            )
            conn.commit()

    def add_item_to_db(self, item_name, category, expiry_date, today):
        insert_query = """insert into pantry (item_name, category, expiry_date, added_at) values (?, ?, ?, ?);"""
        to_add_tuple = (
            item_name,
            category,
            expiry_date,
            today,
        )

        with self.__enter__() as conn:
            cursor = conn.cursor()
            cursor.execute(insert_query, to_add_tuple)
            conn.commit()

    def remove_item_from_db(self, item_name, category, expiry_date):
        remove_query = """
            delete from pantry
            where item_name = ?
                and category = ?
                and expiry_date = ?;
        """

        to_remove_tuple = (
            item_name,
            category,
            expiry_date,
        )

        with self.__enter__() as conn:
            cursor = conn.cursor()
            cursor.execute(remove_query, to_remove_tuple)
            conn.commit()

    def check_all_pantry_items(self):
        with self.__enter__() as conn:
            cursor = conn.cursor()
            results = cursor.execute("""SELECT * from pantry;""")
            conn.commit()

        return results
