// 'use client';

// import React, { createContext, useContext, useState, ReactNode } from 'react';

// interface AuthContextType {
//   isLoggedIn: boolean;
//   login: () => void;
//   logout: () => void;
//   showLogin: boolean;
//   setShowLogin: (show: boolean) => void;
// }

// const AuthContext = createContext<AuthContextType | undefined>(undefined);

// export function AuthProvider({ children }: { children: ReactNode }) {
//   const [isLoggedIn, setIsLoggedIn] = useState(false);
//   const [showLogin, setShowLogin] = useState(false);

//   const login = () => {
//     setIsLoggedIn(true);
//     setShowLogin(false);
//   };

//   const logout = () => {
//     setIsLoggedIn(false);
//   };

//   const value: AuthContextType = {
//     isLoggedIn,
//     login,
//     logout,
//     showLogin,
//     setShowLogin,
//   };

//   return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
// }

// export function useAuth(): AuthContextType {
//   const context = useContext(AuthContext);
//   if (!context) {
//     throw new Error('useAuth must be used within an AuthProvider');
//   }
//   return context;
// }






// // src/app/contexts/AuthContext.tsx
// 'use client';

// import React, { createContext, useContext, ReactNode, useState } from 'react';

// type AuthContextType = {
//   jwtToken: string | null;
//   setJwtToken: (token: string | null) => void;
// };

// const AuthContext = createContext<AuthContextType>({
//   jwtToken: null,
//   setJwtToken: () => {},
// });

// export function useAuth() {
//   return useContext(AuthContext);
// }

// export function AuthProvider({ children }: { children: ReactNode }) {
//   const [jwtToken, setJwtToken] = useState<string | null>(() => {
//     if (typeof window !== 'undefined') {
//       return localStorage.getItem('jwtToken');
//     }
//     return null;
//   });

//   return (
//     <AuthContext.Provider value={{ jwtToken, setJwtToken }}>
//       {children}
//     </AuthContext.Provider>
//   );
// }



// src/app/contexts/AuthContext.tsx
'use client';

import React, { createContext, useContext, useEffect, useState } from 'react';

type AuthContextType = {
  jwtToken: string | null;
  setJwtToken: (token: string | null) => void;
  loaded: boolean;
};

const AuthContext = createContext<AuthContextType>({
  jwtToken: null,
  setJwtToken: () => {},
  loaded: false,
});

export function useAuth() {
  return useContext(AuthContext);
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [jwtToken, setJwtTokenState] = useState<string | null>(null);
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    const storedToken = localStorage.getItem('jwtToken');
    if (storedToken) {
      setJwtTokenState(storedToken);
    }
    setLoaded(true);
  }, []);

  const setJwtToken = (token: string | null) => {
    if (token) {
      localStorage.setItem('jwtToken', token);
    } else {
      localStorage.removeItem('jwtToken');
    }
    setJwtTokenState(token);
  };

  return (
    <AuthContext.Provider value={{ jwtToken, setJwtToken, loaded }}>
      {children}
    </AuthContext.Provider>
  );
}
