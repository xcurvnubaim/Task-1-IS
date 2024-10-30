// pages/detail.js
import { useState, useEffect } from "react";
import Navbar from "../components/Navbar";
import Cookies from "js-cookie";
import { useRouter } from "next/router";

export default function Detail() {
    const [profile, setProfile] = useState({
        fullname: '',
        email: '',
        phone: '',
        address: '',
        nik: ''
    });
    const [files, setFiles] = useState([]);
    const [aesKey, setAesKey] = useState(''); // State for AES key
    const [error, setError] = useState(null);
    const router = useRouter();
    const { id } = router.query; // Retrieve id from query parameters

    const handleAesKeyChange = (e) => setAesKey(e.target.value);

    const fetchData = async () => {
        const token = Cookies.get("auth-token");
        if (!token) {
            router.push("/login");
            return;
        }

        try {
            // Fetch profile and file data with the provided id and AES key
            const response = await fetch(`http://localhost:3000/api/v1/share-request/${id}?aes_key=${encodeURIComponent(aesKey)}`, {
                method: "GET",
                headers: { 
                    Authorization: `Bearer ${token}`,
                    "Content-Type": "application/json"
                }
            });

            const result = await response.json();

            if (result.status) {
                const userProfile = JSON.parse(result.data.user_profile_json);
                setProfile({
                    fullname: userProfile.fullname,
                    email: userProfile.email,
                    phone: userProfile.phone,
                    address: userProfile.address,
                    nik: userProfile.nik
                });

                setFiles(result.data.files.map(file => ({
                    id: file.FileId,
                    name: file.FileName,
                    encryption: file.EncryptionType
                })));
            } else {
                setError(result.error || "Failed to retrieve share request.");
            }
        } catch (error) {
            setError("Error fetching data: " + error.message);
        }
    };

    useEffect(() => {
        if (id && aesKey) fetchData();
    }, [id, aesKey]);

    return (
        <div className="min-h-screen bg-gradient-to-r from-gray-800 to-gray-900">
            <Navbar />
            <div className="p-6">
                <h1 className="text-2xl font-bold text-white mb-4">Enter AES Key</h1>
                <input
                    type="text"
                    value={aesKey}
                    onChange={handleAesKeyChange}
                    placeholder="Enter AES Key"
                    className="mb-6 p-2 rounded border border-gray-500"
                />
                <button
                    onClick={fetchData}
                    className="bg-blue-500 text-white px-4 py-2 rounded"
                    disabled={!aesKey}
                >
                    Submit
                </button>
                
                {error && <p className="text-red-500 mt-4">{error}</p>}
            </div>

            <div className="flex flex-col md:flex-row items-start p-6 gap-6">
                {/* Profile Section */}
                <div className="bg-gray-700 p-6 rounded-lg shadow-lg w-full md:w-1/3">
                    <h2 className="text-2xl font-bold text-center text-white mb-6">Profile Details</h2>
                    <div className="text-gray-300 mb-4">
                        <span className="block font-semibold">Full Name:</span>
                        <p>{profile.fullname || 'N/A'}</p>
                    </div>
                    <div className="text-gray-300 mb-4">
                        <span className="block font-semibold">NIK:</span>
                        <p>{profile.nik || 'N/A'}</p>
                    </div>
                    <div className="text-gray-300 mb-4">
                        <span className="block font-semibold">Email:</span>
                        <p>{profile.email || 'N/A'}</p>
                    </div>
                    <div className="text-gray-300 mb-4">
                        <span className="block font-semibold">Phone:</span>
                        <p>{profile.phone || 'N/A'}</p>
                    </div>
                    <div className="text-gray-300">
                        <span className="block font-semibold">Address:</span>
                        <p>{profile.address || 'N/A'}</p>
                    </div>
                </div>

                {/* Files Section */}
                <div className="bg-gray-700 p-6 rounded-lg shadow-lg w-full md:w-2/3">
                    <h2 className="text-2xl font-bold text-center text-white mb-6">Files</h2>
                    <div className="bg-gray-800 rounded-lg shadow-md p-4 overflow-x-auto">
                        <table className="min-w-full table-auto text-white text-sm">
                            <thead>
                                <tr>
                                    <th className="px-4 py-2 text-left">File Name</th>
                                    <th className="px-4 py-2 text-left">Encryption</th>
                                    <th className="px-4 py-2 text-left">Action</th>
                                </tr>
                            </thead>
                            <tbody>
                                {files?.map((file) => (
                                    <tr key={file.id} className="border-b border-gray-600">
                                        <td className="px-4 py-2">{file.name}</td>
                                        <td className="px-4 py-2">{file.encryption || "None"}</td>
                                        <td className="px-4 py-2">
                                            <button
                                                className="text-green-500 hover:text-green-600"
                                                onClick={() => handleDownload(file.id, file.name)}
                                            >
                                                ⬇️ Download
                                            </button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    );
}

function handleDownload(fileId, filename) {
    alert(`Downloading file: ${filename}`);
}
