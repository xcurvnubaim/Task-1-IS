import { useEffect, useState } from "react";
import Cookies from "js-cookie";
import { useRouter } from "next/router";
import Navbar from "../components/Navbar";
import RequestModal from "../components/RequestModal";

const Request = () => {
    const router = useRouter();
    const [activeTab, setActiveTab] = useState("keluar");
    const [requests, setRequests] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const token = Cookies.get("auth-token");
        if (!token) {
            router.push("/login");
        } else {
            setLoading(false); 
        }

    }, []);

    const handleTabChange = (tab) => {
        setActiveTab(tab);
    };

    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gradient-to-r from-gray-800 to-gray-900">
                <div className="text-white text-xl">Loading...</div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gradient-to-r from-gray-900 via-gray-800 to-gray-700 text-gray-100">
            <Navbar />
            <div className="flex items-center justify-center p-4 md:p-6">
                <div className="bg-gray-800 rounded-xl shadow-xl w-full max-w-6xl p-6 md:p-8 space-y-6">
                    {/* Tab Buttons */}
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

                    {/* Content Based on Active Tab */}
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
    const requests = []; // Placeholder data for requests
    const [isModalOpen, setIsModalOpen] = useState(false); // State to control modal visibility

    const handleAddRequest = () => setIsModalOpen(true); // Open modal on button click

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

            {/* Render Modal if isModalOpen is true */}
            {isModalOpen && (
                <RequestModal
                    isOpen={isModalOpen}
                    onClose={() => setIsModalOpen(false)} // Close modal
                />
            )}

            {/* Table Container */}
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
                                <tr key={index} className="border-b border-gray-700 hover:bg-gray-700 transition">
                                    <td className="px-4 py-2">{index + 1}</td>
                                    <td className="px-4 py-2">{request.username}</td>
                                    <td className="px-4 py-2">{request.status}</td>
                                    <td className="px-4 py-2">
                                        <button className="text-blue-500 hover:text-blue-600 transition">View Details</button>
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
        </div>
    );
};

const RequestMasukContent = () => {
    const requests = []; // Placeholder data for requests

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
                            <tr key={index} className="border-b border-gray-700 hover:bg-gray-700 transition">
                                <td className="px-4 py-2">{index + 1}</td>
                                <td className="px-4 py-2">{request.username}</td>
                                <td className="px-4 py-2">{request.status}</td>
                                <td className="px-4 py-2">
                                    <button className="text-green-500 hover:text-green-600 transition">Approve</button>
                                    <button className="ml-2 text-red-500 hover:text-red-600 transition">Reject</button>
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
