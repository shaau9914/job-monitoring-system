const API_BASE_URL = "http://localhost:8080";

async function fetchJobs() {
    try {
        const response = await fetch(`${API_BASE_URL}/jobs`);

        if (!response.ok) {
            throw new Error("ジョブ一覧の取得に失敗しました");
        }

        const jobs = await response.json();
        displayJobs(jobs);

    } catch (error) {
        console.error(error);
        alert("ジョブ一覧を取得できませんでした。Go APIが起動しているか確認してください。");
    }
}

function displayJobs(jobs) {
    const tableBody = document.getElementById("jobTableBody");
    tableBody.innerHTML = "";

    if (jobs.length === 0) {
        tableBody.innerHTML = `
            <tr>
                <td colspan="6">データがありません</td>
            </tr>
        `;
        return;
    }

    jobs.forEach(job => {
        const row = document.createElement("tr");

        const statusClass = job.status === "failed" ? "failed" : "success";
        const errorCode = job.error_code === null ? "-" : job.error_code;

        const formattedDate = job.executed_at.replace("T", " ");

        row.innerHTML = `
            <td>${job.id}</td>
            <td>${job.job_name}</td>
            <td class="${statusClass}">${job.status}</td>
            <td>${formattedDate}</td>
            <td>${errorCode}</td>
            <td><a href="detail.html?id=${job.id}">詳細</a></td>
        `;

        tableBody.appendChild(row);
    });
}

fetchJobs();