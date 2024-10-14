export default function Register() {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-r from-gray-800 to-gray-900">
        <div className="bg-gray-700 p-10 rounded-lg shadow-lg w-full max-w-md">
          <h2 className="text-3xl font-bold mb-6 text-center text-white">Create an Account</h2>
          <form>
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
              />
            </div>
            <div className="mb-5">
              <label className="block text-gray-300 text-sm font-semibold mb-2" htmlFor="email">
                Email
              </label>
              <input
                className="shadow-sm appearance-none border rounded-lg w-full py-3 px-4 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-amber-500"
                id="email"
                type="email"
                placeholder="Enter your email"
                required
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
              />
            </div>
            <div className="mb-5">
              <label className="block text-gray-300 text-sm font-semibold mb-2" htmlFor="confirm-password">
                Confirm Password
              </label>
              <input
                className="shadow-sm appearance-none border rounded-lg w-full py-3 px-4 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-amber-500"
                id="confirm-password"
                type="password"
                placeholder="Confirm your password"
                required
              />
            </div>
            <div className="mb-5">
              <label className="inline-flex items-center">
                <input
                  type="checkbox"
                  className="form-checkbox h-5 w-5 text-amber-600"
                  required
                />
                <span className="ml-2 text-gray-300">I accept the <a href="#" className="text-amber-400 hover:text-amber-500">Terms and Conditions</a></span>
              </label>
            </div>
            <button
              className="w-full bg-amber-600 hover:bg-amber-700 text-white font-bold py-3 rounded-lg focus:outline-none focus:ring-2 focus:ring-amber-500"
              type="submit"
            >
              Register
            </button>
          </form>
          <p className="mt-4 text-center text-gray-400">
            Already have an account? <a className="text-amber-400 hover:text-amber-500" href="/login">Login</a>
          </p>
        </div>
      </div>
    );
}