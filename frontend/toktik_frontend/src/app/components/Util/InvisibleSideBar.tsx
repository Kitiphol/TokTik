import React from 'react';

function Invisiblesidebar() {
  return (
    <aside className="h-screen w-64 opacity-0">
      {/* Transparent but reserves 256px width space */}
    </aside>
  );
}

export default Invisiblesidebar;
