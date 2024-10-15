import { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import Cookies from 'js-cookie';
import Navbar from '../components/Navbar';

export default function Profile() {
    const [profileData, setProfileData] = useState(null);
    const [loading, setLoading] = useState(true);
    const router = useRouter();

    useEffect(() => {
        const token = Cookies.get('auth-token'); // Get the token from cookies
        if (!token) {
            router.push('/login'); // Redirect to login if token is not present
        } else {
            fetchProfileData(token); // Fetch profile data using the token
        }
    }, [router]);

    const fetchProfileData = async (token) => {
        try {
            const response = await fetch("http://localhost:3000/api/v1/profile/", {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`, // Include the token in the Authorization header
                },
            });
            if (!response.ok) throw new Error('Failed to fetch profile data');
            const data = await response.json();
            setProfileData(data); // Update state with the profile data
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false); // Set loading to false after fetching
        }
    };

    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gradient-to-r from-gray-800 to-gray-900">
                <div className="text-white text-xl">Loading...</div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gradient-to-r from-gray-800 to-gray-900">
            <Navbar />

            <div className="flex items-center justify-center p-4 md:p-6">
                <div className="bg-gray-700 rounded-lg shadow-lg w-full max-w-6xl p-4 md:p-6">
                    <h1 className="text-white text-2xl mb-4">Profile</h1>
                    {profileData ? (
                        <div>
                            <h2 className="text-white text-xl">Username: {profileData.username}</h2>
                            <p className="text-white">Email: {profileData.email}</p>
                            <p className="text-white">Phone: {profileData.phone}</p>
                            <p className="text-white">Address: {profileData.address}</p>
                            <p className="text-white">NIK: {profileData.nik}</p>
                            {/* Add more profile fields as needed */}
                        </div>
                    ) : (
                        <div className="text-white">Profile data not available.</div>
                    )}
                </div>
            </div>
        </div>
    );
}
