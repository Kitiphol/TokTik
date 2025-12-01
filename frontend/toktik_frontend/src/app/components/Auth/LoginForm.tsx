// 'use client';

// import { useState } from 'react';
// import { useRouter } from 'next/navigation';

// export default function LoginForm({
//   onClose,
//   onSuccessLogin,
// }: {
//   onClose: () => void;
//   onSuccessLogin: (token: string) => void;
// }) {
//   const [username, setUsername] = useState('');
//   const [password, setPassword] = useState('');
//   const router = useRouter();

//   const handleLogin = async () => {
//     if (!username || !password) {
//       alert('Please enter both username and password.');
//       return;
//     }

//     try {
//       const response = await fetch('http://localhost:8080/api/login', {
//         method: 'POST',
//         headers: { 'Content-Type': 'application/json' },
//         credentials: 'include',
//         mode: 'cors',
//         body: JSON.stringify({ username, password }),
//       });

//       const data = await response.json();

//       if (response.ok) {
//         const token = data.token;
//         if (!token) throw new Error('No token returned from backend.');

//         localStorage.setItem('jwtToken', token);
//         localStorage.setItem('username', username);

//         alert('Login successful!');
//         onSuccessLogin(token);
//         setTimeout(() => router.push('/'), 1500);
//       } else {
//         alert(data.error || data.message || 'Login failed');
//       }
//     } catch (err) {
//       alert('Login error: ' + (err instanceof Error ? err.message : 'Unknown error'));
//     }
//   };

//   return (
//     <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
//       <div className="bg-white p-8 rounded-xl shadow-xl w-full max-w-md">
//         <h2 className="text-2xl font-bold mb-6 text-center text-black">Log In</h2>
//         <input
//           className="w-full p-2 mb-3 border rounded text-black"
//           placeholder="Username"
//           value={username}
//           onChange={(e) => setUsername(e.target.value)}
//         />
//         <input
//           className="w-full p-2 mb-4 border rounded text-black"
//           type="password"
//           placeholder="Password"
//           value={password}
//           onChange={(e) => setPassword(e.target.value)}
//         />
//         <button
//           className="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600"
//           onClick={handleLogin}
//         >
//           Log In
//         </button>

//         <button
//           className="mt-4 w-full text-sm text-gray-500 hover:underline"
//           onClick={onClose}
//         >
//           Cancel
//         </button>
//       </div>
//     </div>
//   );
// }



// components/Auth/LoginForm.tsx
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';

export default function LoginForm({ onClose, onSuccessLogin }: { onClose: () => void; onSuccessLogin: (token: string) => void }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const router = useRouter();

  const handleLogin = async () => {
    if (!username || !password) {
      alert('Please enter both username and password.');
      return;
    }

    try {


      // const response = await fetch('http://127.0.0.1/api/user/login', {
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        // credentials: 'include',
        mode: 'cors',
        body: JSON.stringify({ username, password }),
      });

      const data = await response.json();

      if (response.ok) {
        const token = data.token;
        if (!token) throw new Error('No token returned from backend.');

        localStorage.setItem('jwtToken', token);
        localStorage.setItem('username', username);

        alert('Login successful!');
        onSuccessLogin(token);
        setTimeout(() => router.push('/'), 1500);
      } else {
        alert(data.error || data.message || 'Login failed');
      }
    } catch (err) {
      alert('Login error: ' + (err instanceof Error ? err.message : 'Unknown error'));
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white p-8 rounded-xl shadow-xl w-full max-w-md">
        <h2 className="text-2xl font-bold mb-6 text-center text-black">Log In</h2>
        <input
          className="w-full p-2 mb-3 border rounded text-black"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <input
          className="w-full p-2 mb-4 border rounded text-black"
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <button
          className="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600"
          onClick={handleLogin}
        >
          Log In
        </button>

        <button
          className="mt-4 w-full text-sm text-gray-500 hover:underline"
          onClick={onClose}
        >
          Cancel
        </button>
      </div>
    </div>
  );
} 
