'use client';

import React from 'react';
import { FaBell } from 'react-icons/fa';

type Notification = {
  id: string;
  message: string;
  read: boolean;
  createdAt: string;
};

type NotificationBellProps = {
  notifications: Notification[];
  showNotifications: boolean;
  onToggle: () => void;
  unreadCount: number;  // NEW prop
};

const NotificationBell: React.FC<NotificationBellProps> = ({
  notifications,
  showNotifications,
  onToggle,
  unreadCount,
}) => {
  return (
    <div className="relative">
      <button
        className="text-gray-600 hover:text-black relative"
        onClick={onToggle}
        aria-label="Toggle notifications"
      >
        <FaBell size={24} />
        {unreadCount > 0 && (
          <span className="absolute -top-1 -right-1 bg-red-600 text-white text-xs w-4 h-4 flex items-center justify-center rounded-full">
            {unreadCount > 9 ? '9+' : unreadCount}
          </span>
        )}
      </button>
      {showNotifications && (
        <div className="absolute right-0 mt-2 w-64 bg-white border shadow-lg rounded z-10 max-h-60 overflow-y-auto">
          {notifications.length > 0 ? (
            notifications.map((notif) => (
              <div
                key={notif.id}
                className={`p-2 border-b text-sm text-black ${
                  !notif.read ? 'bg-gray-100 font-semibold' : ''
                }`}
              >
                <div>{notif.message}</div>
                <div className="text-xs text-gray-400">
                  {new Date(notif.createdAt).toLocaleString()}
                </div>
              </div>
            ))
          ) : (
            <div className="p-2 text-sm text-gray-500">No notifications</div>
          )}
        </div>
      )}
    </div>
  );
};

export default NotificationBell;
