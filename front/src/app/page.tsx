"use client";

import { useState, useEffect, useMemo } from "react";

function CMPCostedDate(a: any, b: any) {
	if (a.PostedAt < b.PostedAt) return 1;
	if (a.PostedAt > b.PostedAt) return -1;
	if (a < b) return 1;
	return 0;
}

export default function Home() {
	const [jobs, setJobs] = useState<any[]>([]);
	const [cronScrapingJob, setCronScrapingJob] = useState<any>([]);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);
	const [filterStatus, setFilterStatus] = useState<string>("null");

	const API_BASE_URL = "http://localhost:8077/api/v1";

	const jobsFiltered = useMemo(() => {
		if (filterStatus === "all") return jobs;
		return jobs.filter((job) => job.Status.toLowerCase() === filterStatus);
	}, [filterStatus,jobs]);

	const statusListChange = (s: string) => {
		if (s === filterStatus) {
			setFilterStatus("all");
		} else {
			setFilterStatus(s);
		}
	};

	const fetchCronJobsStatus = async () => {
		const response = await fetch(`${API_BASE_URL}/cron/scraping/status`);
		if (!response.ok) throw new Error("Failed to fetch jobs");
		const data = await response.json();

		// update coming
		if (cronScrapingJob.processing === true && data.processing === false) {
			fetchJobs();
		}
		setCronScrapingJob(data);
	};

	const activeCronScrapingJob = async () => {
		const response = await fetch(`${API_BASE_URL}/cron/scraping/active`);
		if (!response.ok) throw new Error("Failed to fetch jobs");

		setCronScrapingJob({ ...cronScrapingJob, processing: true });
	};

	const fetchJobs = async () => {
		try {
			const response = await fetch(`${API_BASE_URL}/job`);
			if (!response.ok) throw new Error("Failed to fetch jobs");
			const data = await response.json();
			setJobs(Array.isArray(data) ? data : []);
		} catch (err) {
			setError(err instanceof Error ? err.message : "An error occurred");
		} finally {
			setLoading(false);
		}
	};

	const updateStatus = async (id: string, newStatus: string) => {
		try {
			const response = await fetch(`${API_BASE_URL}/job/${id}/status`, {
				method: "PUT",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify({ status: newStatus }),
			});

			if (!response.ok) throw new Error("Failed to update status");

			setJobs((prevJobs) =>
				prevJobs.map((job) => {
					const jobId = job.ID || job.Id;
					if (jobId === id) {
						return { ...job, status: newStatus, Status: newStatus };
					}
					return job;
				}),
			);
		} catch (err) {
			alert(err instanceof Error ? err.message : "Failed to update status");
		}
	};

	const getStatusClasses = (status: string) => {
		switch (status?.toLowerCase()) {
			case "new":
				return "bg-blue-50 text-blue-700 border-blue-200";
			case "viewed":
				return "bg-slate-50 text-slate-600 border-slate-200";
			case "favorite":
				return "bg-yellow-50 text-yellow-700 border-yellow-200";
			case "registered":
				return "bg-green-50 text-green-700 border-green-200";
			case "interview":
				return "bg-amber-50 text-amber-700 border-amber-200";
			case "rejected":
				return "bg-red-50 text-red-700 border-red-200";
			case "offered":
				return "bg-emerald-50 text-emerald-700 border-emerald-200";
			case "optional":
				return "bg-indigo-50 text-indigo-700 border-indigo-200";
			default:
				return "bg-gray-50 text-gray-600 border-gray-300";
		}
	};

	useEffect(() => {
		fetchJobs();
		fetchCronJobsStatus();

		const interval = setInterval(fetchCronJobsStatus, 10000); // 10 sec
		return () => clearInterval(interval);
	}, []);

	if (loading)
		return (
			<div className="flex min-h-screen items-center justify-center font-sans">
				Loading jobs...
			</div>
		);
	if (error)
		return (
			<div className="flex min-h-screen items-center justify-center font-sans text-red-500">
				Error: {error}
			</div>
		);

	const statusCounts = jobs.reduce(
		(acc, job) => {
			const s = (job.status || job.Status || "new").toLowerCase();
			acc[s] = (acc[s] || 0) + 1;
			return acc;
		},
		{} as Record<string, number>,
	);

	const statusList = ["new", "viewed", "favorite", "registered", "interview", "rejected", "offered", "optional"];

	return (
		<div className="min-h-screen bg-gray-100 p-4 md:p-8 font-sans">
			<div className="mx-auto">
				<header className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold text-gray-800">Job Interview Dashboard</h1>
					<div className="flex gap-1">
						<button
              disabled={cronScrapingJob.processing === true || cronScrapingJob.processing === undefined}
							onClick={() => activeCronScrapingJob()}
							className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors shadow-sm"
						>
							{cronScrapingJob.processing ? (
								<span className="flex items-center gap-2">
                  <svg className="animate-spin h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
									</svg>
									server กำลังดึงข้อมูลใหม่...
								</span>
							) : (
                `ดึงข้อมูลล่าสุด ${cronScrapingJob.time === null ? "N/A" : new Date(cronScrapingJob.time).toLocaleTimeString("th-TH", { hour: 'numeric', minute: 'numeric', second: 'numeric' })}`
							)}
						</button>
					</div>
				</header>

				<div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
					{statusList.map((s) => (
						<div
							key={s}
							onClick={() => statusListChange(s)}
							className={`${s == filterStatus ? "bg-gray-100" : "bg-white"} p-4 rounded-xl shadow-sm border border-gray-200 flex flex-col items-center transition-transform hover:scale-105`}
						>
							<span className="text-xs font-bold uppercase text-gray-500 mb-1">
								{s}
							</span>
							<span
								className={`text-3xl font-extrabold ${
									s === "new"
										? "text-blue-600"
										: s === "viewed"
											? "text-slate-600"
											: s === "favorite"
											? "text-yellow-600"
											: s === "registered"
											? "text-green-600"
											: s === "interview"
												? "text-amber-600"
												: s === "rejected"
													? "text-red-600"
													: s === "offered"
													? "text-emerald-600"
													: "text-indigo-600"
								}`}
							>
								{statusCounts[s] || 0}
							</span>
						</div>
					))}
				</div>

				<div className="bg-white shadow-xl rounded-xl overflow-hidden border border-gray-200">
					<div className="overflow-x-auto">
						<table className="min-w-full divide-y divide-gray-200">
							<thead className="bg-gray-50">
								<tr>
									<th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
										Job & Company
									</th>
									<th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
										Source
									</th>
									<th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
										Posted
									</th>
									<th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
										Location
									</th>
									<th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
										Skills
									</th>
									<th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
										Salary
									</th>
									<th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
										Status
									</th>
									<th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">
										Action
									</th>
								</tr>
							</thead>
							<tbody className="bg-white divide-y divide-gray-200">
								{(filterStatus === "null" ? jobs : jobsFiltered)
									.sort(CMPCostedDate)
									.map((job) => {
										const id = job.ID || job.Id;
										const title = job.title || job.Title;
										const company = job.company_name || job.CompanyName;
										const source = job.source || job.Source;
										const row_posted = job.posted_at || job.PostedAt;
										const posted = new Date(row_posted);
										const rawSalary = job.salary || job.Salary || "";
										const salary =
											rawSalary.length > 30
												? rawSalary.substring(0, 30) + "..."
												: rawSalary;
										const status = job.status || job.Status;
										const url = job.url || job.URL;
										const rawLocation = job.location || job.Location || "";
										const location =
											rawLocation.length > 50
												? rawLocation.substring(0, 50) + "..."
												: rawLocation;
										const skills = job.skills || job.Skills;
										const languages =
											skills?.languages || skills?.Languages || [];
										const frameworks =
											skills?.frameworks || skills?.Frameworks || [];
										const databases =
											skills?.databases || skills?.Databases || [];

										const allSkills = [
											...languages.map((s: string) => ({
												name: s,
												type: "lang",
											})),
											...frameworks.map((s: string) => ({
												name: s,
												type: "fw",
											})),
											...databases.map((s: string) => ({
												name: s,
												type: "db",
											})),
										].slice(0, 5);

										return (
											<tr
												key={id}
												className="hover:bg-gray-50 transition-colors"
											>
												<td className="px-6 py-4">
													<div className="text-sm font-bold text-gray-900">
														{title}
													</div>
													<div className="text-sm text-gray-500">{company}</div>
												</td>
												<td className="px-6 py-4 whitespace-nowrap">
													<span
														className={`px-3 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${
															source === "jobsdb"
																? "bg-blue-100 text-blue-800"
																: "bg-orange-100 text-orange-800"
														}`}
													>
														{source}
													</span>
												</td>
												<td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
													{posted.toLocaleDateString("th-TH", {
														year: "numeric",
														month: "long",
														day: "numeric",
													}) || "N/A"}
												</td>
												<td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
													{location || "N/A"}
												</td>
												<td className="px-6 py-4 text-sm text-gray-600">
													<div className="flex flex-wrap gap-1 max-w-xs">
														{allSkills.map((skill, idx) => (
															<span
																key={idx}
																className={`px-2 py-0.5 rounded text-[10px] font-bold uppercase border ${
																	skill.type === "lang"
																		? "bg-blue-50 text-blue-700 border-blue-100"
																		: skill.type === "fw"
																			? "bg-green-50 text-green-700 border-green-100"
																			: "bg-purple-50 text-purple-700 border-purple-100"
																}`}
															>
																{skill.name}
															</span>
														))}
														{allSkills.length === 0 && (
															<span className="text-gray-400 italic text-xs">
																No skills analyzed
															</span>
														)}
													</div>
												</td>
												<td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
													{salary || "Not specified"}
												</td>
												<td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
													<select
														value={status}
														onChange={(e) => updateStatus(id, e.target.value)}
														className={`block w-full text-xs font-bold uppercase border focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md px-2 py-1 transition-colors cursor-pointer ${getStatusClasses(
															status,
														)}`}
													>
														<option value="new">New</option>
														<option value="viewed">Viewed</option>
														<option value="favorite">Favorite</option>
														<option value="registered">Registered</option>
														<option value="interview">Interview</option>
														<option value="rejected">Rejected</option>
														<option value="offered">Offered</option>
														<option value="optional">Optional</option>
													</select>
												</td>
												<td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
													<a
														href={url}
														target="_blank"
														rel="noopener noreferrer"
														className="inline-block bg-blue-600 text-white px-3 py-1 rounded hover:bg-blue-700 transition-colors shadow-sm text-xs font-bold uppercase"
														onClick={() => {
															if (status === "new") {
																updateStatus(id, "viewed");
															}
														}}
													>
														View Job
													</a>
												</td>
											</tr>
										);
									})}
							</tbody>
						</table>
						{jobs.length === 0 && (
							<div className="p-12 text-center text-gray-500 italic">
								No jobs found. Start scraping to see results!
							</div>
						)}
					</div>
				</div>
			</div>
		</div>
	);
}
