import sqlite3
from pathlib import Path
from datetime import datetime

# DB作成先
DB_DIR = Path(r"D:\portfolio\DB")
DB_DIR.mkdir(parents=True, exist_ok=True)

# DBファイル名：YYYYMMDD_portfolioDB.db
today = datetime.now().strftime("%Y%m%d")
DB_PATH = DB_DIR / f"{today}_portfolioDB.db"


def create_tables():
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()

    # job_result テーブル作成
    cursor.execute("""
    CREATE TABLE IF NOT EXISTS job_result (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        job_name TEXT NOT NULL,
        status TEXT NOT NULL CHECK (status IN ('success', 'failed')),
        executed_at TEXT NOT NULL,
        error_code TEXT,
        error_message TEXT
    );
    """)

    # run_books テーブル作成
    cursor.execute("""
    CREATE TABLE IF NOT EXISTS run_books (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        error_code TEXT NOT NULL UNIQUE,
        title TEXT NOT NULL,
        cause TEXT NOT NULL,
        action TEXT NOT NULL
    );
    """)

    conn.commit()
    conn.close()

    print(f"DB作成完了: {DB_PATH}")


if __name__ == "__main__":
    create_tables()