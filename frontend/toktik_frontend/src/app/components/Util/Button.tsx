"use client";

import React from 'react'
import Link from 'next/link'

interface ButtonProps {
    label: string;
    href: string;
    className?: string;
}



function Button({ 
    label, 
    href, 
    className = 'bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition place-content-center'
    
}: ButtonProps)
{
  return (
    <Link href={href} className={className}> 
        {label}
    </Link>
  )
}

export default Button