// 'use client';

// import { useState } from 'react';
// import { useRouter } from 'next/navigation';

// export default function SignupPage() {
//   const router = useRouter();
//   const [username, setUsername] = useState('');
//   const [email, setEmail] = useState('');
//   const [password, setPassword] = useState('');
//   const [confirmPassword, setConfirmPassword] = useState('');

//   const handleSignup = async (e: React.FormEvent) => {
//     e.preventDefault();

//     try {
//       const res = await fetch('http://localhost:8080/api/register', {
//         method: 'POST',
//         headers: {
//           'Content-Type': 'application/json',
//         },
//         mode: 'cors',
//         body: JSON.stringify({
//           username,
//           email,
//           password,
//           confirm_password: confirmPassword,
//         }),
//       });

//       const data = await res.json(); // parse JSON response

//       if (!res.ok) {
//         alert(data.error || data.message || 'Sign up failed');
//         return;
//       }

//       const token = data.token;
//       const returnedUsername = data.username; // username from response

//       if (!token || !returnedUsername) {
//         alert('Invalid response from server');
//         return;
//       }

//       // Store token and username from server response
//       localStorage.setItem('jwtToken', token);
//       localStorage.setItem('username', returnedUsername);

//       alert('Signup successful! Redirecting...');
//       router.push('/'); // Redirect home or wherever
//     } catch (err: unknown) {
//       if (err instanceof Error) {
//         alert('Signup failed: ' + err.message);
//       } else {
//         alert('Signup failed: Unknown error');
//       }
//     }
//   };

//   return (
//     <div className="min-h-screen flex flex-col items-center justify-center bg-gray-100 px-4">
//       <div className="w-full max-w-md bg-white p-8 rounded shadow">
//         <h2 className="text-2xl font-bold mb-6 text-center">Sign Up</h2>
//         <form onSubmit={handleSignup}>
//           <input
//             type="text"
//             placeholder="Username"
//             className="w-full p-2 mb-4 border rounded text-black"
//             value={username}
//             onChange={(e) => setUsername(e.target.value)}
//             required
//           />
//           <input
//             type="email"
//             placeholder="Email"
//             className="w-full p-2 mb-4 border rounded text-black"
//             value={email}
//             onChange={(e) => setEmail(e.target.value)}
//             required
//           />
//           <input
//             type="password"
//             placeholder="Password"
//             className="w-full p-2 mb-4 border rounded text-black"
//             value={password}
//             onChange={(e) => setPassword(e.target.value)}
//             required
//           />
//           <input
//             type="password"
//             placeholder="Confirm Password"
//             className="w-full p-2 mb-4 border rounded text-black"
//             value={confirmPassword}
//             onChange={(e) => setConfirmPassword(e.target.value)}
//             required
//           />
//           <button
//             type="submit"
//             className="w-full bg-green-600 hover:bg-green-700 text-white py-2 rounded"
//           >
//             Sign Up
//           </button>
//         </form>
//       </div>
//     </div>
//   );
// }


// src/app/signup/page.tsx
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';

export default function SignupPage() {
  const router = useRouter();
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');

  const handleSignup = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      const res = await fetch('/api/auth/register', {
      // const res = await fetch('http://127.0.0.1/api/auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        mode: 'cors',
        body: JSON.stringify({
          username,
          email,
          password,
          confirm_password: confirmPassword,
        }),
      });

      const data = await res.json();

      if (!res.ok) {
        alert(data.error || data.message || 'Sign up failed');
        return;
      }

      const token = data.token;
      const returnedUsername = data.username;

      if (!token || !returnedUsername) {
        alert('Invalid response from server');
        return;
      }

      localStorage.setItem('jwtToken', token);
      localStorage.setItem('username', returnedUsername);

      alert('Signup successful! Redirecting...');
      router.push('/');
    } catch (err: unknown) {
      if (err instanceof Error) {
        alert('Signup failed: ' + err.message);
      } else {
        alert('Signup failed: Unknown error');
      }
    }
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-100 px-4">
      <div className="w-full max-w-md bg-white p-8 rounded shadow">
        <h2 className="text-2xl font-bold mb-6 text-center">Sign Up</h2>
        <form onSubmit={handleSignup}>
          <input
            type="text"
            placeholder="Username"
            className="w-full p-2 mb-4 border rounded text-black"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
          <input
            type="email"
            placeholder="Email"
            className="w-full p-2 mb-4 border rounded text-black"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          <input
            type="password"
            placeholder="Password"
            className="w-full p-2 mb-4 border rounded text-black"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          <input
            type="password"
            placeholder="Confirm Password"
            className="w-full p-2 mb-4 border rounded text-black"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            required
          />
          <button
            type="submit"
            className="w-full bg-green-600 hover:bg-green-700 text-white py-2 rounded"
          >
            Sign Up
          </button>
        </form>
      </div>
    </div>
  );
}



