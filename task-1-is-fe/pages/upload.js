import { useState } from 'react';
import Navbar from '../components/Navbar';
import Cookies from 'js-cookie'; // Import js-cookie
import { useRouter } from 'next/router';

export default function UploadFile() {
  const [file, setFile] = useState(null);
  const [encryptionType, setEncryptionType] = useState('none');
  const [uploading, setUploading] = useState(false);
  const [message, setMessage] = useState('');
  const router = useRouter();

  const handleFileChange = (event) => {
    setFile(event.target.files[0]);
  };

  const handleEncryptionChange = (event) => {
    setEncryptionType(event.target.value);
  };

  const handleUpload = async (event) => {
    event.preventDefault();
    if (!file) {
      setMessage('Please select a file to upload.');
      return;
    }
    
    setUploading(true);
    const formData = new FormData();
    formData.append('file', file);
    formData.append('encryption_type', encryptionType);

    try {
      const token = Cookies.get('auth-token');
      if (!token) {
        setMessage('Unauthorized. Please login first.');
        router.push('/login');
        return; // Ensure no further code is executed
      }
      
      const response = await fetch('http://localhost:3000/api/v1/file/upload', {
        method: 'POST',
        body: formData,
        headers: {
          Authorization: `Bearer ${token}`, // Adjust according to your auth mechanism
        },
      });

      if (response.ok) {
        const data = await response.json();
        setMessage('File uploaded successfully!');
        // Redirect to the dashboard after successful upload
        router.push('/dashboard'); // Update the path to your actual dashboard route
      } else {
        const errorData = await response.json();
        setMessage(`Error: ${errorData.message || 'File upload failed.'}`);
      }
    } catch (error) {
      setMessage(`Error: ${error.message}`);
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-r from-gray-800 to-gray-900">
      <Navbar />
      
      <div className="flex items-center justify-center p-4 md:p-6">
        <div className="bg-gray-700 rounded-lg shadow-lg w-full max-w-6xl p-4 md:p-6">
          <h1 className="text-white text-lg md:text-xl font-bold mb-4">Upload File</h1>
          <form onSubmit={handleUpload}>
            <div className="mb-4">
              <label className="text-white" htmlFor="file">Select File:</label>
              <input
                type="file"
                id="file"
                onChange={handleFileChange}
                className="block w-full text-gray-900"
              />
            </div>
            
            <div className="mb-4">
              <label className="text-white" htmlFor="encryption">Encryption Type:</label>
              <select
                id="encryption"
                value={encryptionType}
                onChange={handleEncryptionChange}
                className="block w-full text-gray-900"
              >
                <option value="none">None</option>
                <option value="des">DES</option>
                <option value="aes">AES</option>
                <option value="rc4">RC4</option>
              </select>
            </div>

            <div className="flex justify-end mb-4">
              <button
                type="submit"
                className={`bg-amber-600 text-white px-4 py-2 md:px-6 md:py-3 rounded-xl shadow-lg hover:bg-amber-700 focus:outline-none focus:ring-2 focus:ring-amber-500 ${uploading ? 'opacity-50 cursor-not-allowed' : ''}`}
                disabled={uploading}
              >
                {uploading ? 'Uploading...' : 'Upload'}
              </button>
            </div>
          </form>

          {message && <p className="text-red-500">{message}</p>}
        </div>
      </div>
    </div>
  );
}
