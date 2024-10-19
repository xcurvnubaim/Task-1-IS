import { useState, useEffect } from "react";
import { useRouter } from "next/router";
import Cookies from "js-cookie"; // Import js-cookie
import Navbar from "../components/Navbar";

export default function Profile() {
    const [fullname, setFullname] = useState('');
    const [email, setEmail] = useState('');
    const [phone, setPhone] = useState('');
    const [address, setAddress] = useState('');
    const [nik, setNik] = useState('');
    const [isEditing, setIsEditing] = useState(false);
    const [loading, setLoading] = useState(true); // Loading state
    const router = useRouter();

    useEffect(() => {
        const fetchData = async () => {
            // Get the token from cookies
            const token = Cookies.get("auth-token");

            // If no token, redirect to login
            if (!token) {
                router.push("/login");
                return;
            }

            try {
                const res = await fetch("http://localhost:3000/api/v1/profile/", {
                    headers: {
                        Authorization: `Bearer ${token}`, // Include token in Authorization header
                    },
                });

                if (!res.ok) {
                    throw new Error("Failed to fetch profile");
                }

                const data = await res.json();
                console.log(data);
                setFullname(data.data.fullname || '');
                setEmail(data.data.email || '');
                setPhone(data.data.phone || '');
                setAddress(data.data.address || '');
                setNik(data.data.nik || '');

                setLoading(false); // Stop loading after fetching the data
            } catch (error) {
                console.error("Error fetching profile:", error);
                // Handle error case
            }
        };

        fetchData();
    }, [router, isEditing]);

    const handleSubmit = async (e) => {
        e.preventDefault();

        const token = Cookies.get("auth-token"); // Retrieve token from cookies for update request
        try {
            const res = await fetch('http://localhost:3000/api/v1/profile/', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${token}`, // Include the token in the Authorization header
                },
                body: JSON.stringify({
                    fullname,
                    email,
                    phone,
                    address,
                    nik,
                }),
            });

            if (res.ok) {
                alert('Profile updated successfully!');
                setIsEditing(false); // Exit editing mode
            } else {
                alert('Error updating profile.');
            }
        } catch (error) {
            console.error("Error updating profile:", error);
        }
    };

    if (loading) {
        // Show a loading message or spinner while fetching the data
        return (
            <div className="min-h-screen flex items-center justify-center bg-gradient-to-r from-gray-800 to-gray-900">
                <div className="text-white text-xl">Loading...</div>
            </div>
        );
    }

    return (
        <div className="min-h-screen flex flex-col bg-gradient-to-r from-gray-800 to-gray-900">
            <Navbar />

            <div className="flex items-center justify-center flex-grow">
                <div className="bg-gray-700 p-10 rounded-lg shadow-lg w-full max-w-3xl">
                    <h2 className="text-3xl font-bold mb-6 text-center text-white">Your Profile</h2>
                    {!isEditing ? (
                        <div className="grid grid-cols-2 gap-6">
                            <div className="mb-5">
                                <span className="block text-gray-300 text-sm font-semibold">Full Name:</span>
                                <p className="text-white">{fullname || 'N/A'}</p>
                            </div>
                            <div className="mb-5">
                                <span className="block text-gray-300 text-sm font-semibold">NIK:</span>
                                <p className="text-white">{nik || 'N/A'}</p>
                            </div>
                            <div className="mb-5">
                                <span className="block text-gray-300 text-sm font-semibold">Email:</span>
                                <p className="text-white">{email || 'N/A'}</p>
                            </div>
                            <div className="mb-5">
                                <span className="block text-gray-300 text-sm font-semibold">Phone:</span>
                                <p className="text-white">{phone || 'N/A'}</p>
                            </div>
                            <div className="mb-5">
                                <span className="block text-gray-300 text-sm font-semibold">Address:</span>
                                <p className="text-white">{address || 'N/A'}</p>
                            </div>
                            <div className="col-span-2">
                                <button
                                    className="w-full bg-amber-600 text-white p-3 rounded-lg hover:bg-amber-700"
                                    onClick={() => setIsEditing(true)}
                                >
                                    Edit Profile
                                </button>
                            </div>
                        </div>
                    ) : (
                        <form onSubmit={handleSubmit} className="grid grid-cols-2 gap-6">
                            <div className="mb-5">
                                <label className="block text-gray-300 text-sm font-semibold mb-2">Full Name</label>
                                <input
                                    className="w-full p-2 rounded-lg"
                                    type="text"
                                    value={fullname}
                                    onChange={(e) => setFullname(e.target.value)}
                                />
                            </div>
                            <div className="mb-5">
                                <label className="block text-gray-300 text-sm font-semibold mb-2">NIK</label>
                                <input
                                    className="w-full p-2 rounded-lg"
                                    type="text"
                                    value={nik}
                                    onChange={(e) => setNik(e.target.value)}
                                />
                            </div>
                            <div className="mb-5">
                                <label className="block text-gray-300 text-sm font-semibold mb-2">Email</label>
                                <input
                                    className="w-full p-2 rounded-lg"
                                    type="email"
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                />
                            </div>
                            <div className="mb-5">
                                <label className="block text-gray-300 text-sm font-semibold mb-2">Phone</label>
                                <input
                                    className="w-full p-2 rounded-lg"
                                    type="text"
                                    value={phone}
                                    onChange={(e) => setPhone(e.target.value)}
                                />
                            </div>
                            <div className="mb-5">
                                <label className="block text-gray-300 text-sm font-semibold mb-2">Address</label>
                                <input
                                    className="w-full p-2 rounded-lg"
                                    type="text"
                                    value={address}
                                    onChange={(e) => setAddress(e.target.value)}
                                />
                            </div>
                            <div className="col-span-2 flex justify-between">
                                <button
                                    className="bg-amber-600 text-white p-3 rounded-lg hover:bg-amber-700"
                                    type="submit"
                                >
                                    Save Changes
                                </button>
                                <button
                                    className="bg-gray-500 text-white p-3 rounded-lg hover:bg-gray-600"
                                    type="button"
                                    onClick={() => setIsEditing(false)}
                                >
                                    Cancel
                                </button>
                            </div>
                        </form>
                    )}
                </div>
            </div>
        </div>
    );
}
