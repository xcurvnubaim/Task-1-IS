import { useEffect, useState } from "react";
import Cookies from "js-cookie";
import { useRouter } from "next/router";
import Navbar from "../components/Navbar";
import RequestModal from "../components/RequestModal";

const Request = () => {
    const router = useRouter();
    const [activeTab, setActiveTab] = useState("keluar");
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const token = Cookies.get("auth-token");
        if (!token) {
            router.push("/login");
        } else {
            setLoading(false);
        }
    }, []);

    const handleTabChange = (tab) => setActiveTab(tab);

    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gradient-to-r from-gray-800 to-gray-900">
                <div className="text-white text-xl animate-pulse">Loading...</div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gradient-to-r from-gray-900 via-gray-800 to-gray-700 text-gray-100">
            <Navbar />
            <div className="flex items-center justify-center p-4 md:p-6">
                <div className="bg-gray-800 rounded-xl shadow-xl w-full max-w-6xl p-6 md:p-8 space-y-6">
                    <div className="flex space-x-6">
                        {["keluar", "masuk"].map((tab) => (
                            <button
                                key={tab}
                                onClick={() => handleTabChange(tab)}
                                className={`flex-1 py-3 rounded-lg transition-colors duration-300
                                    ${activeTab === tab ? "bg-amber-500 text-white" : "bg-gray-700 text-gray-400 hover:bg-gray-600"}`}
                            >
                                {tab === "keluar" ? "Request Keluar" : "Request Masuk"}
                            </button>
                        ))}
                    </div>

                    <div className="bg-gray-900 rounded-lg shadow-lg p-4">
                        {activeTab === "keluar" ? (
                            <RequestKeluarContent />
                        ) : (
                            <RequestMasukContent />
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};

const RequestKeluarContent = () => {
    const [requests, setRequests] = useState([]);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [loading, setLoading] = useState(true);
    const router = useRouter(); // Initialize router

    useEffect(() => {
        const fetchOutgoingRequests = async () => {
            const token = Cookies.get("auth-token");

            try {
                const response = await fetch("http://localhost:3000/api/v1/share-request/by-me", {
                    headers: {
                        Authorization: `Bearer ${token}`,
                        "Content-Type": "application/json",
                    },
                });
                const result = await response.json();
                setRequests(result.data.request || []);
            } catch (error) {
                console.error("Error fetching outgoing requests:", error);
                setRequests([]);
            } finally {
                setLoading(false);
            }
        };

        fetchOutgoingRequests();
    }, []);

    const handleAddRequest = () => setIsModalOpen(true);

    const handleViewDetails = (id) => {
        // Navigate to detail page with id as a query parameter
        router.push(`/detail?id=${id}`);
    };

    return (
        <div>
            <div className="flex justify-end mb-6">
                <button
                    onClick={handleAddRequest}
                    className="bg-amber-600 text-white px-5 py-2 rounded-lg shadow-md hover:bg-amber-700 focus:outline-none focus:ring-2 focus:ring-amber-500 transition-all duration-200"
                >
                    Add Request
                </button>
            </div>

            {isModalOpen && (
                <RequestModal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} />
            )}

            {loading ? (
                <div className="text-center text-gray-400">Loading...</div>
            ) : (
                <div className="bg-gray-800 rounded-lg shadow-md p-4 overflow-x-auto">
                    <table className="min-w-full table-auto text-sm md:text-base">
                        <thead className="bg-gray-700 text-gray-200">
                            <tr>
                                <th className="px-4 py-2 text-left">No</th>
                                <th className="px-4 py-2 text-left">Username</th>
                                <th className="px-4 py-2 text-left">Status</th>
                                <th className="px-4 py-2 text-left">Detail</th>
                            </tr>
                        </thead>
                        <tbody>
                            {requests.length > 0 ? (
                                requests.map((request, index) => (
                                    <tr key={request.id} className="border-b border-gray-700 hover:bg-gray-700 transition">
                                        <td className="px-4 py-2">{index + 1}</td>
                                        <td className="px-4 py-2">{request.request_to_name}</td>
                                        <td className="px-4 py-2">{request.status}</td>
                                        <td className="px-4 py-2">
                                            <button 
                                                className="text-blue-500 hover:text-blue-600 transition" 
                                                onClick={() => handleViewDetails(request.id)} // Navigate to detail page
                                            >
                                                View Details
                                            </button>
                                        </td>
                                    </tr>
                                ))
                            ) : (
                                <tr>
                                    <td colSpan="4" className="text-center py-4 text-gray-400">No outgoing requests.</td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                </div>
            )}
        </div>
    );
};

const RequestMasukContent = () => {
    const [requests, setRequests] = useState([]);

    useEffect(() => {
        const fetchRequests = async () => {
            try {
                const token = Cookies.get("auth-token");
                const response = await fetch("http://localhost:3000/api/v1/share-request/to-me", {
                    headers: {
                        Authorization: `Bearer ${token}`,
                        "Content-Type": "application/json",
                    },
                });

                const result = await response.json();
                setRequests(result.data.request || []);
            } catch (error) {
                console.error("Error fetching requests:", error);
                setRequests([]);
            }
        };

        fetchRequests();
    }, []);

    const handleApprove = async (id) => {
        try {
            const token = Cookies.get("auth-token");
            const response = await fetch("http://localhost:3000/api/v1/share-request/", {
                method: "PUT",
                headers: {
                    Authorization: `Bearer ${token}`,
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ id, status: "accepted" }),
            });

            const result = await response.json();
            if (result.status) {
                setRequests((prevRequests) => prevRequests.filter((request) => request.id !== id));
            } else {
                console.error("Error approving request:", result.error);
            }
        } catch (error) {
            console.error("Error approving request:", error);
        }
    };

    const handleReject = async (id) => {
        try {
            const token = Cookies.get("auth-token");
            const response = await fetch("http://localhost:3000/api/v1/share-request/", {
                method: "PUT",
                headers: {
                    Authorization: `Bearer ${token}`,
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ id, status: "rejected" }),
            });

            const result = await response.json();
            if (result.status) {
                setRequests((prevRequests) => prevRequests.filter((request) => request.id !== id));
            } else {
                console.error("Error rejecting request:", result.error);
            }
        } catch (error) {
            console.error("Error rejecting request:", error);
        }
    };

    return (
        <div>
            <table className="min-w-full table-auto text-sm md:text-base">
                <thead className="bg-gray-700 text-gray-200">
                    <tr>
                        <th className="px-4 py-2 text-left">No</th>
                        <th className="px-4 py-2 text-left">Username</th>
                        <th className="px-4 py-2 text-left">Status</th>
                        <th className="px-4 py-2 text-left">Action</th>
                    </tr>
                </thead>
                <tbody>
                    {requests.length > 0 ? (
                        requests.map((request, index) => (
                            <tr key={request.id} className="border-b border-gray-700 hover:bg-gray-700 transition">
                                <td className="px-4 py-2">{index + 1}</td>
                                <td className="px-4 py-2">{request.request_by_name}</td>
                                <td className="px-4 py-2">{request.status}</td>
                                <td className="px-4 py-2">
                                    <button
                                        onClick={() => handleApprove(request.id)}
                                        disabled={request.status !== "pending"}
                                        className={`text-green-500 font-semibold transition px-3 py-1 rounded-lg ${request.status === "pending" ? "hover:bg-green-600" : "text-gray-500 cursor-not-allowed"
                                            }`}
                                    >
                                        Approve
                                    </button>
                                    <button
                                        onClick={() => handleReject(request.id)}
                                        disabled={request.status !== "pending"}
                                        className={`ml-2 text-red-500 font-semibold transition px-3 py-1 rounded-lg ${request.status === "pending" ? "hover:bg-red-600" : "text-gray-500 cursor-not-allowed"
                                            }`}
                                    >
                                        Reject
                                    </button>
                                </td>
                            </tr>
                        ))
                    ) : (
                        <tr>
                            <td colSpan="4" className="text-center py-4 text-gray-400">No incoming requests.</td>
                        </tr>
                    )}
                </tbody>
            </table>
        </div>
    );
};


export default Request;
