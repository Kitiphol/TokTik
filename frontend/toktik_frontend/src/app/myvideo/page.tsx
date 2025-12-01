// src/app/myvideo/page.tsx
'use client';

import React, { useEffect, useState, useMemo } from 'react';
import { ImEye, ImArrowLeft2, ImArrowRight2 } from 'react-icons/im';
import { BsFillHandThumbsUpFill } from 'react-icons/bs';
import { useRouter } from 'next/navigation';
import { HlsPlayer } from '../components/HlsPlayer/hls'; // Your existing HlsPlayer component
import UploadBox from '../components/Util/UploadBox';
import NotificationBell from '../components/Util/NotificationBell';
import { FaTrash } from 'react-icons/fa';
import {jwtDecode} from "jwt-decode";

import { useCallback } from 'react';

import { useSocket } from '../contexts/SocketContext'; 

import Sidebar from '../components/Util/Sidebar';


type VideoThumbInfo = {
  videoID: string;
  thumbnailURL: string;
  title: string;
  totalLikeCount: number;
  totalViewCount: number;
  hasLiked: boolean;
};

type Notification = {
  id: string;
  message: string;
  read: boolean;
  createdAt: string;
};

type CommentDTO = {
  id: string;
  content: string;
  username: string;
  isUser: boolean;
};

type VideoDetailResponse = {
  playlist: string;
  title: string;
  totalViewCount: number;
  totalLikeCount: number;
  comments: CommentDTO[];
};



// type NotificationPayload = string | { message: string; [key: string]: unknown };
type IncomingNotification = string | {
    id?: string;
    message: string;
  };





export default function MyVideoPage() {
  const router = useRouter();


  const { socket, isConnected } = useSocket();   

  // const [notifications, setNotifications] = useState<string[]>([]);
  const [notifications, setNotifications] = useState<Notification[]>([]);



  const [videos, setVideos] = useState<VideoThumbInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Modal and video details states
  const [selectedVideo, setSelectedVideo] = useState<VideoDetailResponse | null>(null);
  const [selectedVideoIndex, setSelectedVideoIndex] = useState<number | null>(null);
  const [modalHasLiked, setModalHasLiked] = useState(false);
  const [newComment, setNewComment] = useState('');
  const [isSendingComment, setIsSendingComment] = useState(false);
  const [playlistLoading, setPlaylistLoading] = useState(false);


  const [jwtToken, setJwtToken] = useState<string | null>(null);
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem('jwtToken');
    setIsLoggedIn(!!token);
  }, []);

  const handleSuccessLogin = (token: string) => {
    localStorage.setItem('jwtToken', token);
    setIsLoggedIn(true);
  };

  const handleLogout = () => {
    setIsLoggedIn(false);
  };







  const getUserIDFromJWT = (token: string | null): string | null => {
    try {
      if (!token) return null;
      const decoded: { sub?: string } = jwtDecode(token);
      return decoded?.sub ?? null;
    } catch {
      return null;
    }
  };

  const jwtTokenUserID = useMemo(() => getUserIDFromJWT(jwtToken), [jwtToken]);


  

  // Fetch user's videos on mount
  useEffect(() => {
    const fetchMyVideos = async () => {
      setLoading(true);
      setError(null);

      if (!jwtToken) {
        setError('You must be logged in to view your videos.');
        setLoading(false);
        return;
      }

      try {
        // const res = await fetch('http://localhost/api/video/videos/user', {
        const res = await fetch('/api/video/videos/user', {
          headers: { Authorization: `Bearer ${jwtToken}` },
        });
        if (!res.ok) throw new Error(`Error fetching videos: ${res.statusText}`);
        const data: VideoThumbInfo[] = await res.json();
        // setVideos(data);
        setVideos(Array.isArray(data) ? data : []);

      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unknown error');
      } finally {
        setLoading(false);
      }
    };

    fetchMyVideos();
  }, [jwtToken]);


  useEffect(() => {
    // only in browser
    if (typeof window === 'undefined') return;
    const token = localStorage.getItem('jwtToken');
    setJwtToken(token);
    setIsLoggedIn(!!token);
  }, []);






  const onNotification = useCallback((payload: IncomingNotification) => {
  if (!jwtToken) return;

  const message = typeof payload === 'string' ? payload : payload.message;
  const id = typeof payload === 'object' && payload.id ? payload.id : `temp-${Date.now()}`;

  const newNotif: Notification = {
    id,
    message,
    read: false,
    createdAt: new Date().toISOString(),
  };

  setNotifications(prev => [newNotif, ...prev]);

  // ✅ Immediately mark real notifications as read (not temp- ones)
  if (!id.startsWith('temp-')) {
    fetch(`/api/ws/notifications/${id}/read`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${jwtToken}` },
    }).catch(err => {
      console.error('Failed to mark notification as read:', err);
    });
  }

}, [jwtToken]);





  


useEffect(() => {
    if (!jwtToken) return;
    (async () => {
      try {
        const res = await fetch('/api/ws/notifications/unread', {
          headers: { Authorization: `Bearer ${jwtToken}` }
        });
        if (!res.ok) throw new Error();
        const { notifications: unread }: { notifications: Notification[] } = await res.json();
        setNotifications(unread);
      } catch (err) {
        console.error('Notify fetch error', err);
      }
    })();
  }, [jwtToken]);









  const onLike = useCallback((data: { videoID: string; totalLikeCount: number; userID?: string }) => {
  setVideos(prev =>
    prev.map(v =>
      v.videoID === data.videoID
        ? {
            ...v,
            totalLikeCount: data.totalLikeCount,
            hasLiked: data.userID === jwtTokenUserID,
          }
        : v
    )
  );

  if (
    selectedVideoIndex !== null &&
    videos[selectedVideoIndex]?.videoID === data.videoID
  ) {
    setSelectedVideo(sv =>
      sv
        ? {
            ...sv,
            totalLikeCount: data.totalLikeCount,
          }
        : sv
    );
    if (data.userID === jwtTokenUserID) {
      setModalHasLiked(true);
    }
  }
}, [selectedVideoIndex, videos, jwtTokenUserID]);





  


useEffect(() => {
    if (!socket || !isConnected) return;

    // Views
    const onView = (data: { videoID: string; totalViewCount: number }) => {
      setVideos(prev =>
        prev.map(v =>
          v.videoID === data.videoID ? { ...v, totalViewCount: data.totalViewCount } : v
        )
      );
      if (
        selectedVideoIndex !== null &&
        videos[selectedVideoIndex]?.videoID === data.videoID
      ) {
        setSelectedVideo(sv => (sv ? { ...sv, totalViewCount: data.totalViewCount } : sv));
      }
    };

    


    // Comments
    const onComment = (data: { videoID: string; comment: CommentDTO }) => {
      if (
        selectedVideoIndex !== null &&
        videos[selectedVideoIndex]?.videoID === data.videoID
      ) {
        setSelectedVideo(sv =>
          sv ? { ...sv, comments: [...sv.comments, data.comment] } : sv
        );
      }
    };




    
    

    socket.on('video:view', onView);
    socket.on('video:like', onLike);
    socket.on('video:comment', onComment);
    socket.on('notification', onNotification);

    return () => {
      socket.off('video:view', onView);
      socket.off('video:like', onLike);
      socket.off('video:comment', onComment);
      socket.off('notification', onNotification);
    };
  }, [socket, isConnected, videos, selectedVideoIndex, onLike]);







  // Delete video handler
  const handleDeleteVideo = async (videoID: string) => {
    if (!jwtToken) return alert('You must be logged in.');

    if (!confirm('Are you sure you want to delete this video?')) return;

    try {
    
        // const res = await fetch(`http://localhost/api/video/videos/${videoID}`, {

      const res = await fetch(`/api/video/videos/${videoID}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${jwtToken}` },
      });

      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || 'Failed to delete video');
      }

      // Remove from UI
      setVideos((prev) => prev.filter((v) => v.videoID !== videoID));

      // Close modal if the deleted video is open
      if (selectedVideoIndex !== null && videos[selectedVideoIndex].videoID === videoID) {
        setSelectedVideo(null);
        setSelectedVideoIndex(null);
      }
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Unknown error');
    }
  };

  // Open video modal & fetch video details + comments + likes + views
  const handleVideoClick = async (videoID: string, idx: number) => {
    if (!jwtToken) {
      alert('You must be logged in to watch videos.');
      return;
    }

    setPlaylistLoading(true);
    setSelectedVideo(null);
    setSelectedVideoIndex(idx);

    try {
      const [vRes, lRes, vwRes, cRes] = await Promise.all([
        // fetch(`http://localhost/api/video/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        fetch(`/api/video/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),

        // fetch(`http://localhost/api/ws/likes/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        fetch(`/api/ws/likes/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),


        // fetch(`http://localhost/api/ws/views/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        fetch(`/api/ws/views/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),


        // fetch(`http://localhost/api/ws/comments/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        fetch(`/api/ws/comments/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
      ]);

      if (!vRes.ok) throw new Error('Failed to load video details');

      const vd = await vRes.json();
      const ld = lRes.ok ? await lRes.json() : { likes: 0, hasLiked: false };
      const vw = vwRes.ok ? await vwRes.json() : { views: 0 };
      const cm = cRes.ok ? await cRes.json() : { comments: [] };

      setSelectedVideo({
        playlist: vd.playlist,
        title: vd.title,
        totalViewCount: vw.views,
        totalLikeCount: ld.likes,
        comments: Array.isArray(cm.comments) ? cm.comments : [],
      });

      setModalHasLiked(ld.hasLiked);
    } catch (err: unknown) {
    console.error('Error loading video details:', err);
      alert('Failed to load video');
      setSelectedVideo(null);
      setSelectedVideoIndex(null);
    } finally {
      setPlaylistLoading(false);
    }
  };

  // Like/unlike toggle
  const toggleLike = async (videoID: string, idx: number) => {
    if (!jwtToken) {
      alert('You must be logged in.');
      return;
    }

    try {
      const currentlyLiked = videos[idx].hasLiked;
      // const res = await fetch(`http://localhost/api/ws/likes/videos/${videoID}`, {

      const res = await fetch(`/api/ws/likes/videos/${videoID}`, {
        method: currentlyLiked ? 'DELETE' : 'POST',
        headers: { Authorization: `Bearer ${jwtToken}` },
      });

      if (!res.ok) throw new Error('Toggle like failed');

      const { likes, hasLiked } = await res.json();

      // Update videos list likes
      // setVideos((prev) =>
      //   prev.map((v, i) => (i === idx ? { ...v, totalLikeCount: likes, hasLiked } : v))
      // );


      // Optimistic update
      setVideos((prev) =>
        prev.map((v, i) =>
          i === idx
            ? {
                ...v,
                totalLikeCount: currentlyLiked ? v.totalLikeCount - 1 : v.totalLikeCount + 1,
                hasLiked: !currentlyLiked,
              }
            : v
        )
      );

      if (selectedVideoIndex === idx && selectedVideo) {
        setSelectedVideo({
          ...selectedVideo,
          totalLikeCount: currentlyLiked
            ? selectedVideo.totalLikeCount - 1
            : selectedVideo.totalLikeCount + 1,
        });
        setModalHasLiked(!currentlyLiked);
      }





      // Update modal likes if open for this video
      if (selectedVideoIndex === idx && selectedVideo) {
        setSelectedVideo({ ...selectedVideo, totalLikeCount: likes });
        setModalHasLiked(hasLiked);
      }
    } catch {
      alert('Could not update like');
    }
  };

  // Send a comment
  const handleSendComment = async () => {
    if (!jwtToken || !newComment.trim() || selectedVideoIndex === null) return;

    try {
      setIsSendingComment(true);

      const videoID = videos[selectedVideoIndex].videoID;

      // const res = await fetch(`http://localhost/api/ws/comments/videos/${videoID}`, {

      const res = await fetch(`/api/ws/comments/videos/${videoID}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${jwtToken}`,
        },
        body: JSON.stringify({ content: newComment }),
      });

      if (!res.ok) throw new Error('Failed to post comment');

      const added = await res.json();

      // setSelectedVideo((prev) =>
      //   prev ? { ...prev, comments: [...prev.comments, added] } : prev
      // );
      setSelectedVideo(prev => prev ? { ...prev, comments: [...prev.comments, added.comment] } : prev);


      setNewComment('');
    } catch {
      alert('Failed to post comment');
    } finally {
      setIsSendingComment(false);
    }
  };


  const unreadCount = notifications.filter(n => !n.read).length;
  const [isOpen, setIsOpen] = useState(false);
  

  const handleToggleBell = async () => {
  const newIsOpen = !isOpen;
  setIsOpen(newIsOpen);


  if (newIsOpen && unreadCount > 0 && jwtToken) {
    try {
      await Promise.all(
        notifications
          .filter(n => !n.read && !n.id.startsWith('temp-'))
          .map(n =>
            fetch(`/api/ws/notifications/${n.id}/read`, {
              method: 'POST',
              headers: { Authorization: `Bearer ${jwtToken}` },
            })
          )
      );

      // Mark all as read in local state
      setNotifications(prev =>
        prev.map(n => ({ ...n, read: true }))
      );
    } catch (err) {
      console.error('Failed to mark notifications as read', err);
    }
  }
};



// useEffect(() => {
//     if (!isOpen && notifications.length) {
//       const markRead = async () => {
//         const toMark = notifications.filter(n => !n.id.startsWith('temp-'));
//         if (toMark.length === 0) return;
  
//         await Promise.all(
//           toMark.map(n => {
//             const parts = n.id.split("-");
//             const uuid = parts.slice(1, 6).join("-"); 
//             return fetch(`/api/ws/notifications/${uuid}/read`, {
//               method: 'POST',
//               headers: { Authorization: `Bearer ${jwtToken}` },
//             });
//           })
//         );
  
//         // Only remove the read notifications (keep temp ones)
//         setNotifications(prev => prev.filter(n => n.id.startsWith('temp-')));
//       };
  
//       markRead();
//     }
//   }, [isOpen, notifications, jwtToken]);




  // Delete a comment
  const handleDeleteComment = async (commentID: string) => {
    if (!jwtToken || selectedVideoIndex === null) return;
    if (!confirm('Are you sure you want to delete this comment?')) return;

    try {
      const videoID = videos[selectedVideoIndex].videoID;

      // const res = await fetch(`http://localhost/api/ws/comments/videos/${videoID}/${commentID}`, {
      const res = await fetch(`api/ws/comments/videos/${videoID}/${commentID}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${jwtToken}` },
      });

      if (res.ok) {
        setSelectedVideo((prev) =>
          prev ? { ...prev, comments: prev.comments.filter((c) => c.id !== commentID) } : prev
        );
      } else {
        alert('Failed to delete comment');
      }
    } catch {
      alert('An error occurred while deleting the comment.');
    }
  };

  // if (loading) return <div className="p-4">Loading your videos...</div>;
  // if (error) return <div className="p-4 text-red-600">Error: {error}</div>;
  // if (!videos.length) return <div className="p-4">You have no uploaded videos.</div>;


  

  return (

    <div className="flex">
    <Sidebar
      isLoggedIn={isLoggedIn}
      onSuccessLogin={handleSuccessLogin}
      onLogout={handleLogout}
    />

      {/* <div className="p-4 min-h-screen bg-gray-100"> */}
      {/* <div className="p-4 min-h-screen bg-gray-100 ml-64"> */}
      <div className="ml-64 p-4 bg-gray-100 min-h-screen w-full">

        <div className="flex justify-end mb-4">
          <button
            onClick={() => router.push('/')}
            className="text-blue-600 hover:text-blue-800 font-semibold underline"
          >
            Home
          </button>


          <NotificationBell
            notifications={notifications}
            showNotifications={isOpen}
            onToggle={handleToggleBell}
            unreadCount={unreadCount}
          />

        </div>

        <UploadBox />

        {loading ? (
          <div className="text-gray-600">Loading your videos...</div>
        ) : error ? (
          <div className="text-red-600">Error: {error}</div>
        ) : (
          <>
            <h1 className="text-2xl font-bold mb-4 text-black">My Videos</h1>
            {videos.length === 0 ? (
              <div className="p-4 text-gray-600">You have no uploaded videos.</div>
            ) : (





              <div className="flex flex-wrap gap-4">



                {videos.map((video, idx) => (
                  <div
                    key={video.videoID}
                    className="relative flex items-center justify-between border rounded p-2 shadow hover:shadow-lg bg-white w-full"
                    onClick={() => handleVideoClick(video.videoID, idx)}
                  >
                    <div className="flex items-center gap-4 w-full">
                      <img
                        src={video.thumbnailURL}
                        alt={video.title}
                        className="w-24 h-16 object-cover rounded"
                      />

                      <div className="flex flex-col flex-grow">
                        <h2 className="font-semibold text-black truncate">{video.title}</h2>
                        <div className="flex items-center gap-4 text-sm text-gray-600 mt-1">
                          <span
                            onClick={(e) => {
                              e.stopPropagation();
                              toggleLike(video.videoID, idx);
                            }}
                            className="flex items-center gap-1 cursor-pointer select-none"
                          >
                            <BsFillHandThumbsUpFill
                              color={video.hasLiked ? '#dc2626' : '#6b7280'}
                            />
                            {video.totalLikeCount}
                          </span>
                          <span className="flex items-center gap-1">
                            <ImEye />
                            {video.totalViewCount}
                          </span>
                        </div>
                      </div>
                    </div>

                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        handleDeleteVideo(video.videoID);
                      }}
                      className="text-red-500 hover:text-red-700 ml-2 self-start"
                      title="Delete video"
                    >
                      <FaTrash />
                    </button>
                  </div>
                ))}





              </div>






            )}
          </>
        )}

        {/* Modal Video Player */}
        {playlistLoading && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white p-6 rounded">Loading video...</div>
          </div>
        )}

        {selectedVideo && selectedVideoIndex !== null && (
          <div className="fixed inset-0 bg-black bg-opacity-70 flex items-center justify-center z-50 p-4">
            <div className="bg-white rounded shadow-lg max-w-5xl w-full flex flex-col md:flex-row relative">
              {/* Navigation Buttons */}
              {selectedVideoIndex > 0 && (
                <button
                  className="absolute left-2 top-1/2 -translate-y-1/2 text-3xl text-gray-600 hover:text-black"
                  onClick={() =>
                    handleVideoClick(
                      videos[selectedVideoIndex - 1].videoID,
                      selectedVideoIndex - 1
                    )
                  }
                  aria-label="Previous video"
                >
                  <ImArrowLeft2 />
                </button>
              )}
              {selectedVideoIndex < videos.length - 1 && (
                <button
                  className="absolute right-2 top-1/2 -translate-y-1/2 text-3xl text-gray-600 hover:text-black"
                  onClick={() =>
                    handleVideoClick(
                      videos[selectedVideoIndex + 1].videoID,
                      selectedVideoIndex + 1
                    )
                  }
                  aria-label="Next video"
                >
                  <ImArrowRight2 />
                </button>
              )}

              {/* Close modal */}
              <button
                className="absolute top-2 right-2 text-2xl text-gray-600 hover:text-black"
                onClick={() => {
                  setSelectedVideo(null);
                  setSelectedVideoIndex(null);
                  setNewComment('');
                }}
                aria-label="Close video modal"
              >
                &times;
              </button>

              {/* Video player */}
              <div className="md:w-2/3 w-full flex items-center justify-center bg-black rounded">
                <HlsPlayer playlistUrl={selectedVideo.playlist} />
              </div>

              {/* Video info & comments */} 
              <div className="md:w-1/3 w-full p-4 flex flex-col">
                <h2 className="text-xl font-bold truncate mb-2 text-black" >
                  {selectedVideo.title}
                </h2>

                {/* Like and views */}
                <div className="flex gap-4 items-center mb-4 text-black">
                  <span
                    onClick={() =>
                      toggleLike(
                        videos[selectedVideoIndex].videoID,
                        selectedVideoIndex
                      )
                    }
                    className="cursor-pointer"
                    aria-label={modalHasLiked ? 'Unlike video' : 'Like video'}
                  >
                    <BsFillHandThumbsUpFill
                      size={22}
                      color={modalHasLiked ? '#dc2626' : '#6b7280'}
                    />
                  </span>
                  <span>{selectedVideo.totalLikeCount}</span>
                  <ImEye size={18} />
                  <span>{selectedVideo.totalViewCount}</span>
                </div>

                {/* Comments section */}
                <div className="flex flex-col flex-grow text-black">
                  <h3 className="font-semibold mb-2 text-black">Comments</h3>
                  <div className="flex-grow overflow-y-auto max-h-48 space-y-2 border border-gray-300 p-2 rounded text-black">
                    {selectedVideo.comments.length ? (
                      selectedVideo.comments.map((c) => (
                        <div
                          key={c.id}
                          className="flex justify-between items-center border-b pb-1"
                        >
                          <span>
                            <strong>{c.username}:</strong> {c.content}
                          </span>
                          {c.isUser && (
                            <button
                              onClick={() => handleDeleteComment(c.id)}
                              className="text-red-600 hover:text-red-800 ml-2"
                              title="Delete comment"
                              aria-label="Delete comment"
                            >
                              🗑️
                            </button>
                          )}
                        </div>
                      ))
                    ) : (
                      <div className="text-gray-400">No comments yet.</div>
                    )}
                  </div>

                  <textarea
                    value={newComment}
                    onChange={(e) => setNewComment(e.target.value)}
                    placeholder="Add a comment..."
                    className="mt-2 p-2 border border-gray-300 rounded resize-none"
                    rows={2}
                    aria-label="Add a comment"
                  />

                  <button
                    onClick={handleSendComment}
                    disabled={isSendingComment || !newComment.trim()}
                    className="mt-2 px-4 py-1 bg-blue-600 text-white rounded disabled:opacity-50"
                  >
                    {isSendingComment ? 'Sending...' : 'Send'}
                  </button>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}