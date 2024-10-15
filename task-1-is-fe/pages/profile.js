import { useState, useEffect } from 'react';
import Navbar from '../components/Navbar'; // Import the Navbar component

export default function Profile() {
  // States for profile fields
  const [fullname, setFullname] = useState('');
  const [email, setEmail] = useState('');
  const [phone, setPhone] = useState('');
  const [address, setAddress] = useState('');
  const [nik, setNik] = useState('');
  const [isEditing, setIsEditing] = useState(false);

  // Fetch user data from the API when the component mounts
  useEffect(() => {
    const fetchData = async () => {
      const res = await fetch('/api/account'); // Adjust this based on your API endpoint
      const data = await res.json();
      
      setFullname(data.fullname || '');
      setEmail(data.email || '');
      setPhone(data.phone || '');
      setAddress(data.address || '');
      setNik(data.nik || '');
    };

    fetchData();
  }, []);

  // Function to handle profile updates
  const handleSubmit = async (e) => {
    e.preventDefault();

    const res = await fetch('/api/account', {
      method: 'PUT', // Using PUT for updating profile
      headers: {
        'Content-Type': 'application/json',
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
      setIsEditing(false); // Exit editing mode after successful update
    } else {
      alert('Error updating profile.');
    }
  };

  return (
    <div className="min-h-screen flex flex-col bg-gradient-to-r from-gray-800 to-gray-900">
      <Navbar /> {/* Add the Navbar component here */}

      <div className="flex items-center justify-center flex-grow">
        <div className="bg-gray-700 p-10 rounded-lg shadow-lg w-full max-w-md">
          <h2 className="text-3xl font-bold mb-6 text-center text-white">Your Profile</h2>
          {!isEditing ? (
            <div>
              <div className="mb-5">
                <span className="block text-gray-300 text-sm font-semibold">Full Name:</span>
                <p className="text-white">{fullname || 'N/A'}</p>
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
              <div className="mb-5">
                <span className="block text-gray-300 text-sm font-semibold">NIK:</span>
                <p className="text-white">{nik || 'N/A'}</p>
              </div>
              <button
                className="w-full bg-amber-600 text-white p-3 rounded-lg hover:bg-amber-700"
                onClick={() => setIsEditing(true)}
              >
                Edit Profile
              </button>
            </div>
          ) : (
            <form onSubmit={handleSubmit}>
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
                <label className="block text-gray-300 text-sm font-semibold mb-2">Email</label>
                <input
                  className="w-full p-2 rounded-lg"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  disabled
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
              <div className="mb-5">
                <label className="block text-gray-300 text-sm font-semibold mb-2">NIK</label>
                <input
                  className="w-full p-2 rounded-lg"
                  type="text"
                  value={nik}
                  onChange={(e) => setNik(e.target.value)}
                />
              </div>
              <div className="flex justify-between">
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
