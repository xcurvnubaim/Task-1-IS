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
    // const [loading, setLoading] = useState(true);
    const router = useRouter();

    useEffect(() => {
        const fetchData = async () => {
            const token = Cookies.get("auth-token");
            if (!token) {
                router.push("/login");
                return;
            }

            try {
                // Fetch profile data
                const profileRes = await fetch("api", {
                    headers: { Authorization: `Bearer ${token}` }
                });
                const profileData = await profileRes.json();
                setProfile(profileData.data);

                // Fetch file data
                const filesRes = await fetch("api", {
                    headers: { Authorization: `Bearer ${token}` }
                });
                const filesData = await filesRes.json();
                setFiles(filesData.data.files);

                // setLoading(false);
            } catch (error) {
                console.error("Error fetching data:", error);
            }
        };
        fetchData();
    }, [router]);

    // if (loading) {
    //     return (
    //         <div className="min-h-screen flex items-center justify-center bg-gradient-to-r from-gray-800 to-gray-900">
    //             <div className="text-white text-xl">Loading...</div>
    //         </div>
    //     );
    // }

    return (
        <div className="min-h-screen bg-gradient-to-r from-gray-800 to-gray-900">
            <Navbar />
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
                                {files?.map((file, index) => (
                                    <tr key={index} className="border-b border-gray-600">
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
    // Placeholder function for handling file downloads
    alert(`Downloading file: ${filename}`);
}
