import { useState, useEffect } from "react";
import { Doughnut } from "react-chartjs-2";
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from "chart.js";
import type { ChartData } from "chart.js"

ChartJS.register(ArcElement, Tooltip, Legend);

type MetricsData = {
  active_users: number;
  events_per_min: number;
  avg_duration: number;
  top_elements: { element: string; count: number }[];
};

function App() {
  const [activeUsers, setActiveUsers] = useState(0);
  const [eventsMin, setEventsMin] = useState(0);
  const [avgEngagement, setAvgEngagement] = useState(0);
  const [topElements, setTopElements] = useState<
    { element: string; count: number }[]
  >([]);
  const [chartData, setChartData] = useState<ChartData<"doughnut">>({
    labels: [],
    datasets: [
      {
        data: [],
        backgroundColor: ["#10B981", "#3B82F6", "#EF4444", "#F59E0B", "#8B5CF6"],
      },
    ],
  });

  useEffect(() => {
    const fetchMetrics = () => {
      fetch("http://localhost:8081/metrics")
        .then((res) => res.json())
        .then((data: MetricsData) => {
          console.log(data)
          setActiveUsers(data.active_users);
          setEventsMin(data.events_per_min);
          setAvgEngagement(data.avg_duration);
          setTopElements(data.top_elements);

          setChartData({
            labels: data.top_elements.map((e) => e.element),
            datasets: [
              {
                data: data.top_elements.map((e) => e.count),
                backgroundColor: [
                  "#10B981",
                  "#3B82F6",
                  "#EF4444",
                  "#F59E0B",
                  "#8B5CF6",
                ],
              },
            ],
          });
        })
        .catch((err) => console.error(err));
    };

    fetchMetrics();
    const interval = setInterval(fetchMetrics, 3000);
    return () => clearInterval(interval);
  }, []);

  const sendEvent = (action: string, element: string) => {
    const event = {
      user_id: "user_" + Math.floor(Math.random() * 1000),
      action,
      element,
      duration: action === "play" ? Math.random() * 10 : 0,
    };

    fetch("http://localhost:8081/event", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(event),
    });
  };

  return (
    <section className="bg-gray-900 text-gray-100 min-h-screen">
      <div className="container mx-auto p-6">
        <h1 className="text-3xl font-bold mb-6 text-emerald-400">
          User Engagement Analytics
        </h1>

        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <div className="bg-gray-800 p-4 rounded-2xl shadow-lg border-l-4 border-emerald-500">
            <h3 className="font-medium text-gray-400">Active Users (5m)</h3>
            <p className="text-4xl font-bold text-emerald-400">{activeUsers}</p>
          </div>
          <div className="bg-gray-800 p-4 rounded-2xl shadow-lg border-l-4 border-emerald-500">
            <h3 className="font-medium text-gray-400">Events/Min</h3>
            <p className="text-4xl font-bold text-emerald-400">
              {eventsMin.toFixed(1)}
            </p>
          </div>
          <div className="bg-gray-800 p-4 rounded-2xl shadow-lg border-l-4 border-emerald-500">
            <h3 className="font-medium text-gray-400">Avg Engagement</h3>
            <p className="text-4xl font-bold text-emerald-400">
              {avgEngagement.toFixed(1)}s
            </p>
          </div>
        </div>

        {/* Charts + List */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-gray-800 p-4 rounded-2xl shadow">
            <h3 className="font-semibold mb-4 text-gray-200">
              Top Interactive Elements
            </h3>
            <ul className="divide-y divide-gray-700">
              {topElements.map((item) => (
                <li
                  key={item.element}
                  className="py-2 flex justify-between text-gray-300"
                >
                  <span className="font-medium">{item.element}</span>
                  <span className="text-gray-500">
                    {item.count} interactions
                  </span>
                </li>
              ))}
            </ul>
          </div>

          <div className="bg-gray-800 p-4 rounded-2xl shadow">
            <h3 className="font-semibold mb-4 text-gray-200">
              Engagement Distribution
            </h3>
            <Doughnut data={chartData} />
          </div>
        </div>

        {/* Buttons */}
        <div className="mt-8 bg-gray-800 rounded-2xl shadow p-4">
          <h3 className="font-semibold mb-4 text-gray-200">Send Test Event</h3>
          <div className="flex flex-wrap gap-3">
            <button
              onClick={() => sendEvent("play", "video_player")}
              className="bg-emerald-600 hover:bg-emerald-500 text-white px-4 py-2 rounded-lg shadow"
            >
              Play Video
            </button>
            <button
              onClick={() => sendEvent("pause", "video_player")}
              className="bg-yellow-600 hover:bg-yellow-500 text-white px-4 py-2 rounded-lg shadow"
            >
              Pause Video
            </button>
            <button
              onClick={() => sendEvent("click", "subscribe_button")}
              className="bg-emerald-500 hover:bg-emerald-400 text-white px-4 py-2 rounded-lg shadow"
            >
              Subscribe
            </button>
            <button
              onClick={() => sendEvent("click", "like_button")}
              className="bg-red-600 hover:bg-red-500 text-white px-4 py-2 rounded-lg shadow"
            >
              Like
            </button>
          </div>
        </div>
      </div>
    </section>
  );
}

export default App;
