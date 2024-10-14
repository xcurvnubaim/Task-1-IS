export default function Login() {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-r from-gray-800 to-gray-900">
        <div className="bg-gray-700 p-10 rounded-lg shadow-lg w-full max-w-md">
          <h2 className="text-3xl font-bold mb-6 text-center text-white">Sign in to your account</h2>
          <form>
            <div className="mb-5">
              <label className="block text-gray-300 text-sm font-semibold mb-2" htmlFor="username">
                Username
              </label>
              <input
                className="shadow-sm appearance-none border rounded w-full py-3 px-4 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-amber-500"
                id="username"
                type="text"
                placeholder="Enter your username"
                required
              />
            </div>
            <div className="mb-5">
              <label className="block text-gray-300 text-sm font-semibold mb-2" htmlFor="password">
                Password
              </label>
              <input
                className="shadow-sm appearance-none border rounded w-full py-3 px-4 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-amber-500"
                id="password"
                type="password"
                placeholder="Enter your password"
                required
              />
            </div>
            <div className="flex items-center justify-between mb-5">
              <a
                className="text-sm text-amber-400 hover:text-amber-500"
                href="#"
              >
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
            Don’t have an account? <a className="text-amber-400 hover:text-amber-500" href="#">Sign Up</a>
          </p>
        </div>
      </div>
    );
}