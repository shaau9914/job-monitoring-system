const API_BASE_URL = "http://localhost:8080";

function getJobIdFromUrl() {
    const params = new URLSearchParams(window.location.search);
    return params.get("id");
}

async function fetchJobDetail() {
    const jobId = getJobIdFromUrl();

    if (!jobId) {
        alert("ジョブIDが指定されていません");
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/jobs/${jobId}`);

        if (!response.ok) {
            throw new Error("ジョブ詳細の取得に失敗しました");
        }

        const job = await response.json();
        displayJobDetail(job);

    } catch (error) {
        console.error(error);
        alert("ジョブ詳細を取得できませんでした。Go APIが起動しているか確認してください。");
    }
}

function displayJobDetail(job) {
    document.getElementById("jobId").textContent = job.id;
    document.getElementById("jobName").textContent = job.job_name;
    document.getElementById("executedAt").textContent = job.executed_at;

    const statusElement = document.getElementById("status");
    statusElement.textContent = job.status;
    statusElement.className = job.status === "failed" ? "failed" : "success";

    document.getElementById("errorCode").textContent = job.error_code ?? "-";
    document.getElementById("errorMessage").textContent = job.error_message ?? "-";

    const runbookSection = document.getElementById("runbookSection");

    if (job.runbook) {
        runbookSection.style.display = "block";
        document.getElementById("runbookTitle").textContent = job.runbook.title;
        document.getElementById("runbookCause").textContent = job.runbook.cause;
        document.getElementById("runbookAction").textContent = job.runbook.action;
    } else {
        runbookSection.style.display = "none";
    }
}

fetchJobDetail();