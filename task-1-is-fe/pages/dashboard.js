import { useEffect, useState } from "react";
import Navbar from "../components/Navbar";
import { useRouter } from "next/router";
import Cookies from "js-cookie"; // Import js-cookie


export default function Dashboard() {
  const [files, setFiles] = useState([]);
  const [loading, setLoading] = useState(true); // Loading state
  const router = useRouter();

  // Fetch files from the API
  const fetchFiles = async () => {
    try {
      // Token for authentication
      const token = Cookies.get("auth-token");
      if (!token) {
        // Redirect to login if token is not set
        router.push("/login");
        return null;
      }
      const response = await fetch("http://localhost:3000/api/v1/file/", {
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        throw new Error("Failed to fetch files");
      }

      const data = await response.json();
      setFiles(
        data.data.files.map((file) => ({
          id: file.file_id,
          name: file.file_name,
          encryption: file.encryption_type,
        }))
      );

      setLoading(false); // Set loading to false after fetching files
    } catch (error) {
      console.error("Error fetching files:", error);
    }
  };

  // Use effect to fetch files when component mounts
  useEffect(() => {
    fetchFiles();
  }, []);

  // Function to handle file download
  const handleDownload = async (fileId, filename) => {
    try {
      const token = Cookies.get("auth-token");
      const response = await fetch(
        `http://localhost:3000/api/v1/file/download/${fileId}`,
        {
          method: "GET",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error("Failed to download file");
      }

      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);

      // Get the filename from the response headers
      const contentDisposition = response.headers.get("Content-Disposition");
      if (
        contentDisposition &&
        contentDisposition.indexOf("attachment") !== -1
      ) {
        const matches = /filename="?(.+)"?/.exec(contentDisposition);
        if (matches != null && matches[1]) {
          filename = matches[1];
        }
      }

      const a = document.createElement("a");
      a.href = url;
      a.download = filename; // Set the filename from the response
      document.body.appendChild(a);
      a.click();
      a.remove();
      window.URL.revokeObjectURL(url);
    } catch (error) {
      console.error("Error downloading file:", error);
    }
  };



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
            <a href="/upload">
              <button className="bg-amber-600 text-white px-4 py-2 md:px-6 md:py-3 rounded-xl shadow-lg hover:bg-amber-700 focus:outline-none focus:ring-2 focus:ring-amber-500">
                Upload
              </button>
            </a>
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
                {files?.map((file, index) => (
                  <tr key={index} className="border-b border-gray-600">
                    <td className="px-2 md:px-4 py-2">
                      {file.type === "folder" ? "📁" : "📄"} {file.name}
                    </td>
                    <td className="px-2 md:px-4 py-2">
                      {file.encryption || "None"}
                    </td>
                    <td className="px-2 md:px-4 py-2">
                      <button
                        className="text-green-500 hover:text-green-600"
                        onClick={() => handleDownload(file.id, file.name)} // Call the download function on click
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
