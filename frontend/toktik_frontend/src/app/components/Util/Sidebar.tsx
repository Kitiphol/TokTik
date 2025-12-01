'use client';

import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import { FaHome, FaUpload, FaVideo, FaUser } from 'react-icons/fa'; // imported icons
import Button from './Button';
import Button2 from './Button2';
import LoginForm from '../Auth/LoginForm';
import UploadForm from './UploadForm';
import UpdateUserForm from './UpdateUser';

interface SidebarProps {
  isLoggedIn: boolean;
  onSuccessLogin: (token: string) => void;
  onLogout: () => void;
}

export default function Sidebar({
  isLoggedIn,
  onSuccessLogin,
  onLogout
}: SidebarProps) {
  const [jwtToken, setJwtToken] = useState<string | null>(null);
  const [showLogin, setShowLogin] = useState(false);
  const [showUpload, setShowUpload] = useState(false);
  const [showUpdateForm, setShowUpdateForm] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem('jwtToken');
    if (token) {
      setJwtToken(token);
    }
  }, []);

  const handleLoginSuccess = (token: string) => {
    localStorage.setItem('jwtToken', token);
    setJwtToken(token);
    setShowLogin(false);
    onSuccessLogin(token);
  };

  return (
    <>
      <aside className="fixed top-0 left-0 h-screen w-64 bg-black text-white p-4 flex flex-col gap-4">
        <h2 className="text-xl font-bold">Sidebar</h2>

        {/* Home Button */}
        <Link
          href="/"
          className="flex items-center gap-2 text-white py-2 px-4 rounded text-center bg-transparent hover:bg-white/10 transition"
        >
          <FaHome className="text-xl" />
          <span>Home</span>
        </Link>

        {/* Upload Button (always visible) */}
        <button
          onClick={() => {
            if (jwtToken) {
              setShowUpload(true);
            } else {
              setShowLogin(true);
            }
          }}
          className="flex items-center gap-2 text-white py-2 px-4 rounded text-center bg-transparent hover:bg-white/10 transition"
        >
          <FaUpload className="text-xl" />
          <span>Upload</span>
        </button>

        {/* Show only if not logged in */}
        {!isLoggedIn && (
          <>
            <Button
              label="Sign Up"
              href="/signup"
              className="bg-green-600 hover:bg-green-700 text-white py-2 px-4 rounded text-center"
            />
            <button
              onClick={() => setShowLogin(true)}
              className="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded text-center"
            >
              Log In
            </button>
          </>
        )}

        {/* Show if logged in */}
        {isLoggedIn && (
          <>
            {/* My Videos Link */}
            <Link
              href="/myvideo"
              className="flex items-center gap-2 text-white py-2 px-4 rounded text-center bg-transparent hover:bg-white/10 transition"
            >
              <FaVideo className="text-xl" />
              <span>My Videos</span>
            </Link>

            <button
              onClick={() => setShowUpdateForm(true)}
              className="flex items-center gap-2 text-white py-2 px-4 rounded text-center bg-transparent hover:bg-white/10 transition"
            >
              <FaUser className="text-xl" />
              <span>Update Profile</span>
            </button>

            <button
              onClick={() => {
                localStorage.removeItem('jwtToken');
                localStorage.removeItem('username');
                setJwtToken(null);
                onLogout();
                window.location.href = '/';
              }}
              className="bg-red-600 hover:bg-red-700 text-white py-2 px-4 rounded text-center"
            >
              Sign Out
            </button>
          </>
        )}

        {/* Ping Backend Button */}
        <Button2
          label="Ping Backend"
          href="#"
          className="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600 transition mx-auto mt-6"
          onClick={(e) => {
            e.preventDefault();
            fetch("/api/auth/ping", {
              mode: "cors",
              headers: {
                Authorization: `Bearer ${jwtToken ?? ''}`
              }
            })
              .then((res) => {
                if (!res.ok) throw new Error(`HTTP error! Status: ${res.status}`);
                return res.text();
              })
              .then((data) => alert(`Response from backend: ${data}`))
              .catch((err) =>
                alert(`Error: ${err instanceof Error ? err.message : 'Unknown error'}`)
              );
          }}
        />
      </aside>

      {/* Popups */}
      {showLogin && (
        <LoginForm
          onClose={() => setShowLogin(false)}
          onSuccessLogin={handleLoginSuccess}
        />
      )}
      {showUpload && <UploadForm onClose={() => setShowUpload(false)} />}
      {showUpdateForm && (
        <UpdateUserForm onClose={() => setShowUpdateForm(false)} jwtToken={jwtToken} />
      )}
    </>
  );
}
