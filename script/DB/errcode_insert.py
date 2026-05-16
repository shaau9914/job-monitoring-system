import json
import sqlite3
from pathlib import Path
from datetime import datetime

JOB_NAME = "import_status_data"

INPUT_FILE = Path(r"D:\portfolio\input\status_data.json")
DB_PATH = Path(r"D:\portfolio\DB\20260429_portfolioDB.db")

REQUIRED_FIELDS = ["id", "system_name", "check_target", "status", "checked_at"]


def insert_job_result(status, error_code=None, error_message=None):
    executed_at = datetime.now().isoformat(timespec="seconds")

    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()

    cursor.execute(
        """
        INSERT INTO job_result (
            job_name,
            status,
            executed_at,
            error_code,
            error_message
        ) VALUES (?, ?, ?, ?, ?)
        """,
        (
            JOB_NAME,
            status,
            executed_at,
            error_code,
            error_message,
        ),
    )

    conn.commit()
    conn.close()


def validate_record(record):
    if not isinstance(record, dict):
        return "failed", "INVALID_FORMAT", "入力形式不正"

    missing_fields = [field for field in REQUIRED_FIELDS if field not in record]

    if missing_fields:
        missing_text = ", ".join(missing_fields)
        return (
            "failed",
            "INVALID_FORMAT",
            f"入力形式不正。不足項目: {missing_text}",
        )

    if record.get("status") == "failed":
        return "failed", "JOB_FAILED", "ジョブ結果失敗"

    return "success", None, None


def main():
    if not INPUT_FILE.exists():
        insert_job_result(
            status="failed",
            error_code="FILE_NOT_FOUND",
            error_message="入力ファイル未配置",
        )
        print("failed: FILE_NOT_FOUND")
        return

    try:
        with INPUT_FILE.open("r", encoding="utf-8") as f:
            data = json.load(f)
    except json.JSONDecodeError:
        insert_job_result(
            status="failed",
            error_code="INVALID_FORMAT",
            error_message="入力形式不正",
        )
        print("failed: INVALID_FORMAT")
        return

    if not data:
        insert_job_result(
            status="failed",
            error_code="NO_DATA",
            error_message="入力データなし",
        )
        print("failed: NO_DATA")
        return

    if not isinstance(data, list):
        insert_job_result(
            status="failed",
            error_code="INVALID_FORMAT",
            error_message="入力形式不正",
        )
        print("failed: INVALID_FORMAT")
        return

    for record in data:
        status, error_code, error_message = validate_record(record)

        insert_job_result(
            status=status,
            error_code=error_code,
            error_message=error_message,
        )

        print(
            f"{status}: "
            f"record_id={record.get('id') if isinstance(record, dict) else '-'} "
            f"error_code={error_code}"
        )


if __name__ == "__main__":
    main()