import { useState } from 'react';
import Link from 'next/link';
import { HiOutlineUserCircle } from 'react-icons/hi';
import { useRouter } from 'next/router';

const Navbar = () => {
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const router = useRouter();

  const toggleDropdown = () => {
    setDropdownOpen(!dropdownOpen);
  };

  const handleLogout = async () => {
    const res = await fetch('/api/logout', {
      method: 'POST',
    });

    if (res.ok) {
      // Redirect to the login page after logout
      router.push('/login');
    } else {
      console.error('Logout failed');
    }
  };

  return (
    <nav className="bg-amber-500 dark:bg-gray-800 p-4 flex justify-end">
      <div className="relative">
        {/* Account Icon using react-icons */}
        <button
          onClick={toggleDropdown}
          className="text-white text-3xl focus:outline-none"
        >
          <HiOutlineUserCircle className="w-8 h-8" />
        </button>

        {/* Dropdown Menu */}
        {dropdownOpen && (
          <div className="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-700 rounded-md shadow-lg z-50">
            <ul className="py-1">
              <li>
                <Link href="/account" className="block px-4 py-2 text-gray-700 dark:text-white hover:bg-gray-100 dark:hover:bg-gray-600">
                  Account
                </Link>
              </li>
              <li>
                <button
                  onClick={handleLogout}
                  className="w-full text-left block px-4 py-2 text-gray-700 dark:text-white hover:bg-gray-100 dark:hover:bg-gray-600"
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
