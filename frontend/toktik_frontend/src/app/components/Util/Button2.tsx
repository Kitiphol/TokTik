import React from "react";
import Link from "next/link";

interface ButtonProps {
  label: string;
  href?: string;
  className?: string;
  onClick?: (e: React.MouseEvent) => void; // <-- add this line
}

function Button2({
  label,
  href = "",
  className = "bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition place-content-center",
  onClick, // <-- accept onClick
}: ButtonProps) {
  const handleClick = async (e: React.MouseEvent) => {
    e.preventDefault();

    if (onClick) {
      onClick(e); // if custom onClick passed, call it
    } else {
      try {
        // const response = await fetch("http://localhost/api/auth/ping", {
        const response = await fetch("/api/auth/ping", {
          mode: "cors", // good to add this explicitly
        });
        if (!response.ok) throw new Error("Network response was not ok");
        const data = await response.text();
        alert("Message from API: " + data);
      } catch (error) {
        alert("Error fetching data: " + error);
      }
    }
  };

  return (
    <Link href={href} className={className} onClick={handleClick}>
      {label}
    </Link>
  );
}

export default Button2;
