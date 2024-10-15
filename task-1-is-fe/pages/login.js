import { useState } from 'react';
import { useRouter } from 'next/router';
import Cookies from 'js-cookie'; // Import js-cookie

export default function Login() {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState(null);
    const router = useRouter();

    const handleLogin = async (event) => {
        event.preventDefault(); // Prevent default form submission behavior

        const res = await fetch("http://localhost:3000/api/v1/auth/login", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password }),
        });

        if (res.ok) {
            const data = await res.json(); // Assuming your API returns a token
            Cookies.set('auth-token', data.token, { expires: 7 }); // Set token in cookies with a 7-day expiration
            router.push('/dashboard'); // Redirect to dashboard
        } else {
            const data = await res.json();
            setError(data.message || 'Login failed. Please try again.');
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-gradient-to-r from-gray-800 to-gray-900">
            <div className="bg-gray-700 p-10 rounded-lg shadow-lg w-full max-w-md">
                <h2 className="text-3xl font-bold mb-6 text-center text-white">Sign in to your account</h2>
                {error && <p className="text-red-500 text-center">{error}</p>}
                <form onSubmit={handleLogin}>
                    <div className="mb-5">
                        <label className="block text-gray-300 text-sm font-semibold mb-2" htmlFor="username">
                            Username
                        </label>
                        <input
                            className="shadow-sm appearance-none border rounded-lg w-full py-3 px-4 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-amber-500"
                            id="username"
                            type="text"
                            placeholder="Enter your username"
                            required
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                        />
                    </div>
                    <div className="mb-5">
                        <label className="block text-gray-300 text-sm font-semibold mb-2" htmlFor="password">
                            Password
                        </label>
                        <input
                            className="shadow-sm appearance-none border rounded-lg w-full py-3 px-4 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-amber-500"
                            id="password"
                            type="password"
                            placeholder="Enter your password"
                            required
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                        />
                    </div>
                    <div className="flex items-center justify-between mb-5">
                        <a className="text-sm text-amber-400 hover:text-amber-500" href="#">
                            Forgot Password?
                        </a>
                    </div>
                    <button
                        className="w-full bg-amber-600 hover:bg-amber-700 text-white font-bold py-3 rounded-lg focus:outline-none focus:ring-2 focus:ring-amber-500"
                        type="submit"
                    >
                        Login
                    </button>
                </form>
                <p className="mt-4 text-center text-gray-400">
                    Don’t have an account? <a className="text-amber-400 hover:text-amber-500" href="/register">Sign Up</a>
                </p>
            </div>
        </div>
    );
}
