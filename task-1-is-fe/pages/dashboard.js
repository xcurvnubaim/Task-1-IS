import { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import Cookies from 'js-cookie';
import Navbar from '../components/Navbar';

export default function Dashboard() {
  const [files, setFiles] = useState([
    { name: "My Dream", type: "folder", size: "41Gb", modified: "Nov 12, 2022", encryption: "AES-256" },
    { name: "My Projects", type: "folder", size: "32Gb", modified: "Nov 12, 2022", encryption: "AES-256" },
    { name: "Home Design3.mp4", type: "video", size: "421Mb", modified: "Nov 12, 2022", encryption: "None" },
    // Add more file objects as needed
  ]);

  const [loading, setLoading] = useState(true); // Loading state
  const router = useRouter();

  useEffect(() => {
    const token = Cookies.get('auth-token'); // Get the token from cookies

    // Check if the token is present
    if (!token) {
      // Redirect to login if token is not found
      router.push('/login');
    } else {
      setLoading(false); // Token found, stop loading
    }
  }, [router]);

  if (loading) {
    // Show a loading message or spinner while checking for token
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

          {/* Upload Button */}
          <div className="flex justify-end mb-4">
            <button className="bg-amber-600 text-white px-4 py-2 md:px-6 md:py-3 rounded-xl shadow-lg hover:bg-amber-700 focus:outline-none focus:ring-2 focus:ring-amber-500">
              Upload
            </button>
          </div>

          {/* File List */}
          <div className="bg-gray-800 rounded-lg shadow-md p-4 overflow-x-auto">
            <table className="min-w-full table-auto text-white text-sm md:text-base">
              <thead>
                <tr>
                  <th className="px-2 md:px-4 py-2 text-left">File Name</th>
                  <th className="px-2 md:px-4 py-2 text-left">Encryption</th>
                  <th className="px-2 md:px-4 py-2 text-left">Action</th>
                </tr>
              </thead>
              <tbody>
                {files.map((file, index) => (
                  <tr key={index} className="border-b border-gray-600">
                    <td className="px-2 md:px-4 py-2">{file.type === 'folder' ? '📁' : '📄'} {file.name}</td>
                    <td className="px-2 md:px-4 py-2">{file.encryption}</td>
                    <td className="px-2 md:px-4 py-2">
                      <button className="text-green-500 hover:text-green-600">⬇️ Download</button>
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
