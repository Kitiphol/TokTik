'use client';

import React, { useState } from 'react';

type UpdateUserFormProps = {
  onClose: () => void;
  jwtToken: string | null;
};

export default function UpdateUserForm({ onClose, jwtToken }: UpdateUserFormProps) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [message, setMessage] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!jwtToken) {
      alert('You must be logged in to update your user.');
      return;
    }

    if (password !== confirmPassword) {
      setMessage('Passwords do not match.');
      return;
    }

    try {

      // const response = await fetch('http://localhost/api/user/updateUser', {
      const response = await fetch('/api/user/updateUser', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${jwtToken}`,
        },
        body: JSON.stringify({ username, password, confirm_password: confirmPassword }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText);
      }

    
      const result = await response.text();
      alert(`Update successful: ${result}`);
      onClose();
    } catch (err: unknown) {
      if (err instanceof Error) {
        setMessage(`Error: ${err.message}`);
      } else {
        setMessage('An unknown error occurred');
      }
    }
  };

  return (
    <div className="fixed top-0 left-0 w-screen h-screen bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white text-black p-6 rounded-lg w-full max-w-md">
        <h2 className="text-xl font-bold mb-4">Update Your Profile</h2>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <label>
            New Username:
            <input
              type="text"
              className="w-full p-2 border rounded mt-1"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
          </label>

          <label>
            New Password:
            <input
              type="password"
              className="w-full p-2 border rounded mt-1"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </label>

          <label>
            Confirm Password:
            <input
              type="password"
              className="w-full p-2 border rounded mt-1"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
            />
          </label>

          <div className="flex justify-end gap-2">
            <button
              type="button"
              onClick={onClose}
              className="bg-gray-400 hover:bg-gray-500 text-white px-4 py-2 rounded"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded"
            >
              Update
            </button>
          </div>
        </form>

        {message && <p className="mt-4 text-red-500">{message}</p>}
      </div>
    </div>
  );
}
