import { useState } from 'react';
import Link from 'next/link';
import { HiOutlineUserCircle } from 'react-icons/hi';
import { useRouter } from 'next/router';
import Cookies from 'js-cookie';

const Navbar = () => {
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const router = useRouter();

  const toggleDropdown = () => {
    setDropdownOpen(!dropdownOpen);
  };

  const handleLogout = () => {
    // Clear the auth token from cookies
    Cookies.remove('auth-token'); // Remove the token

    // Redirect to the login page after logout
    router.push('/login');
  };

  return (
    <nav className="bg-amber-500 dark:bg-gray-800 p-4 flex justify-between items-center shadow-md">
      {/* Dashboard Link */}
      <Link href="/dashboard" className="text-white text-xl font-semibold dark:text-white hover:text-gray-300 transition-colors duration-200 ease-in-out">
        Dashboard
      </Link>

      <div className="relative">
        {/* Profile Icon using react-icons */}
        <button
          onClick={toggleDropdown}
          className="text-white text-3xl focus:outline-none hover:text-gray-300 transition-colors duration-200 ease-in-out"
        >
          <HiOutlineUserCircle className="w-8 h-8" />
        </button>

        {/* Dropdown Menu */}
        {dropdownOpen && (
          <div className="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-700 rounded-lg shadow-lg z-50">
            <ul className="py-2">
              <li>
                <Link href="/profile" className="block px-4 py-2 text-gray-700 dark:text-white hover:bg-amber-500 dark:hover:bg-gray-600 rounded-t-lg transition-colors duration-200">
                  Profile
                </Link>
              </li>
              <li>
                <button
                  onClick={handleLogout}
                  className="w-full text-left block px-4 py-2 text-gray-700 dark:text-white hover:bg-amber-500 dark:hover:bg-gray-600 rounded-b-lg transition-colors duration-200"
                >
                  Log Out
                </button>
              </li>
            </ul>
          </div>
        )}
      </div>
    </nav>
  );
};

export default Navbar;
