import { useState } from "react";
import { useRouter } from "next/router";
import Cookies from "js-cookie";

const RequestModal = ({ isOpen, onClose }) => {
    const [username, setUsername] = useState("");
    const [userId, setUserId] = useState(null);
    const [isUsernameValid, setIsUsernameValid] = useState(null);
    const [errorMessage, setErrorMessage] = useState("");
    const router = useRouter();

    const handleUsernameChange = (e) => {
        setUsername(e.target.value);
        setIsUsernameValid(null);
        setErrorMessage("");
    };

    const handleCheckUsername = async () => {
        if (!username) {
            setErrorMessage("Please enter a username.");
            return;
        }

        try {
            const token = Cookies.get("auth-token");
            if (!token) {
                router.push("/login");
                return;
            }

            const response = await fetch(`http://localhost:3000/api/v1/auth/username/${username}`, {
                method: "GET",
                headers: {
                    Authorization: `Bearer ${token}`,
                    "Content-Type": "application/json",
                },
            });
            const result = await response.json();

            if (result.status) {
                setIsUsernameValid(true);
                setUserId(result.data.id);
                setErrorMessage("");
            } else {
                setIsUsernameValid(false);
                setErrorMessage(result.error || "User not found.");
                setUserId(null);
            }
        } catch (error) {
            console.error("Error checking username:", error);
            setErrorMessage("Error validating username.");
        }
    };

    const handleRequestSubmit = async () => {
        if (isUsernameValid && userId) {
            try {
                const token = Cookies.get("auth-token");
                const response = await fetch("http://localhost:3000/api/v1/share-request/", {
                    method: "POST",
                    headers: {
                        Authorization: `Bearer ${token}`,
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify({
                        request_to: userId,
                    }),
                });

                if (!response.ok) {
                    throw new Error("Failed to submit request.");
                }

                const data = await response.json();
                console.log("Request submitted:", data);
            } catch (error) {
                console.error("Error submitting request:", error);
                setErrorMessage("Error submitting request.");
            } finally {
                onClose();
            }
        } else {
            setErrorMessage("Please enter a valid username.");
        }
    };

    return (
        isOpen && (
            <div className="fixed inset-0 bg-gray-900 bg-opacity-75 flex items-center justify-center z-50">
                <div className="bg-gray-800 rounded-lg p-6 w-full max-w-md mx-4 space-y-4">
                    <h2 className="text-xl font-semibold text-gray-100">Request To</h2>
                    
                    <div className="flex items-center space-x-2">
                        <input
                            type="text"
                            value={username}
                            onChange={handleUsernameChange}
                            className="w-full px-4 py-2 rounded-md text-gray-900 focus:outline-none"
                            placeholder="Enter username"
                        />
                        <button
                            onClick={handleCheckUsername}
                            className="bg-blue-600 text-white px-3 py-2 rounded-lg shadow-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 transition"
                        >
                            Check
                        </button>
                    </div>
                    
                    {isUsernameValid === false && <p className="text-red-500">{errorMessage}</p>}
                    {isUsernameValid && <p className="text-green-500">Username is valid.</p>}
                    
                    <div className="flex justify-end space-x-2 mt-4">
                        <button
                            onClick={onClose}
                            className="bg-gray-500 text-white px-4 py-2 rounded-lg shadow-md hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-gray-400 transition"
                        >
                            Cancel
                        </button>
                        <button
                            onClick={handleRequestSubmit}
                            className="bg-amber-600 text-white px-4 py-2 rounded-lg shadow-md hover:bg-amber-700 focus:outline-none focus:ring-2 focus:ring-amber-500 transition"
                            disabled={!isUsernameValid}
                        >
                            Submit Request
                        </button>
                    </div>
                </div>
            </div>
        )
    );
};

export default RequestModal;
