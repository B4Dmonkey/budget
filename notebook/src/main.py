import sqlite3

# * The purpose of this file is to serve as a helper for notebooks

def ping_db(conn) -> bool:
    try:
        cursor = conn.cursor()
        cursor.execute("SELECT 1")
        cursor.fetchone()
        print("Database connection is active.")
        return True
    except sqlite3.Error as e:
        print(f"Database connection failed: {e}")
        return False


def main():
    conn = sqlite3.connect("notebook.db")
    ping_db(conn)
    conn.close()


if __name__ == "__main__":
    main()
