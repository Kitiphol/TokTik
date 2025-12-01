 // // src/app/components/Home.tsx
 "use client";

import React, { useState, useEffect, useCallback , useMemo  } from 'react';
import LoginForm from './Auth/LoginForm';
import Sidebar from './Util/Sidebar';
import { ImEye, ImArrowLeft2, ImArrowRight2 } from 'react-icons/im';
import { BsFillHandThumbsUpFill } from 'react-icons/bs';
// import Hls from 'hls.js';
import {HlsPlayer} from './HlsPlayer/hls';
import { jwtDecode } from "jwt-decode";



import { useSocket } from '../contexts/SocketContext';  // Import your socket context hook
import NotificationBell from './Util/NotificationBell';


export type VideoThumbInfo = {
  videoID: string;
  thumbnailURL: string;
  title: string;
  totalLikeCount: number;
  totalViewCount: number;
  hasLiked: boolean;
};

type VideoThumbInfoBackend = {
  videoID: string;
  thumbnailURL: string;
  title: string;
  totalLikeCount: number;
  totalViewCount: number;
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

export default function Home() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [showLoginPopup, setShowLoginPopup] = useState(false);
  const [jwtToken, setJwtToken] = useState<string | null>(null);
  const [videos, setVideos] = useState<VideoThumbInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [playlistLoading, setPlaylistLoading] = useState(false);
  const [selectedVideo, setSelectedVideo] = useState<VideoDetailResponse | null>(null);
  const [selectedVideoIndex, setSelectedVideoIndex] = useState<number | null>(null);
  const [modalHasLiked, setModalHasLiked] = useState(false);
  const [newComment, setNewComment] = useState('');
  const [isSendingComment, setIsSendingComment] = useState(false);



  const { socket, isConnected } = useSocket();



  type Notification = {
    id: string;
    message: string;
    read: boolean;
    createdAt: string;
  };

  const [notifications, setNotifications] = useState<Notification[]>([]);



  const onView = useCallback((data: { videoID: string; totalViewCount: number }) => {
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
  }, [selectedVideoIndex, videos]);


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





// const onLike = useCallback(
//   (data: { videoID: string; totalLikeCount: number; username?: string; hasLiked?: boolean }) => {
//     setVideos(prevVideos =>
//       prevVideos.map(v =>
//         v.videoID === data.videoID
//             ? { ...v, totalLikeCount: data.totalLikeCount, hasLiked: data.hasLiked! }
//             : v
//         )
//     );


//     if (selectedVideoIndex !== null && videos[selectedVideoIndex].videoID === data.videoID) {
//         setSelectedVideo(sv => sv ? { ...sv, totalLikeCount: data.totalLikeCount } : sv);
//         setModalHasLiked(data.hasLiked!);
//       }
//     },
//   [selectedVideoIndex, videos]
// );


const onLike = useCallback((data: { videoID: string; totalLikeCount: number; userID?: string }) => {
  setVideos(prev =>
    prev.map(v => {
      if (v.videoID !== data.videoID) return v;

      // If the event is from current user, update hasLiked accordingly
      if (data.userID === jwtTokenUserID) {
        return {
          ...v,
          totalLikeCount: data.totalLikeCount,
          hasLiked: true,
        };
      } else {
        // For other users' like events, just update count, keep hasLiked as is
        return {
          ...v,
          totalLikeCount: data.totalLikeCount,
        };
      }
    })
  );

  if (
    selectedVideoIndex !== null &&
    videos[selectedVideoIndex]?.videoID === data.videoID
  ) {
    setSelectedVideo(sv => {
      if (!sv) return sv;

      // same logic here
      if (data.userID === jwtTokenUserID) {
        setModalHasLiked(true);
      }
      return {
        ...sv,
        totalLikeCount: data.totalLikeCount,
      };
    });
  }
}, [jwtTokenUserID, selectedVideoIndex, videos]);











  // const onComment = useCallback((data: { videoID: string; comment: CommentDTO }) => {
  //   if (
  //     selectedVideoIndex !== null &&
  //     videos[selectedVideoIndex]?.videoID === data.videoID
  //   ) {
  //     setSelectedVideo(sv =>
  //       sv ? { ...sv, comments: [...sv.comments, data.comment] } : sv
  //     );
  //   }
  // }, [selectedVideoIndex, videos]);

  const onComment = useCallback((data: { videoID: string; comment: CommentDTO }) => {
    if (
      selectedVideoIndex !== null &&
      videos[selectedVideoIndex]?.videoID === data.videoID
    ) {
      setSelectedVideo(sv =>
        sv && !sv.comments.some(c => c.id === data.comment.id)
          ? { ...sv, comments: [...sv.comments, data.comment] }
          : sv
      );
    }
  }, [selectedVideoIndex, videos]);


  
  type IncomingNotification = string | {
    id?: string;
    message: string;
  };


  const [isOpen, setIsOpen] = useState(false);


  const unreadCount = notifications.filter(n => !n.read).length;

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

// Run this effect only when closing the bell
// useEffect(() => {
//   if (!isOpen && notifications.length) {
//     const markRead = async () => {
//       const toMark = notifications.filter(n => !n.id.startsWith('temp-'));
//       if (toMark.length === 0) return;

//       await Promise.all(
//         toMark.map(n => {
//           const parts = n.id.split("-");
//           const uuid = parts.slice(1, 6).join("-"); 
//           return fetch(`/api/ws/notifications/${uuid}/read`, {
//             method: 'POST',
//             headers: { Authorization: `Bearer ${jwtToken}` },
//           });
//         })
//       );

//       // Only remove the read notifications (keep temp ones)
//       setNotifications(prev => prev.filter(n => n.id.startsWith('temp-')));
//     };

//     markRead();
//   }
// }, [isOpen, notifications, jwtToken]);


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


  



  const fetchVideos = async (token: string | null) => {
    setLoading(true);
    try {

      // const res = await fetch('http://localhost/api/video/videos');
      const res = await fetch('/api/video/videos');
      const list: VideoThumbInfoBackend[] = res.ok ? await res.json() : [];
      const sorted = [...list].sort((a, b) => a.videoID.localeCompare(b.videoID));
      const thumbs = await Promise.all(
        sorted.map(async (v) => {
          let likes = v.totalLikeCount;
          let hasLiked = false;
          if (token) {


            // const statusRes = await fetch(`http://localhost/api/ws/likes/videos/${v.videoID}`,
            const statusRes = await fetch(`/api/ws/likes/videos/${v.videoID}`,

              
              
              { headers: { Authorization: `Bearer ${token}` } }
            );
            if (statusRes.ok) {
              const data = await statusRes.json();
              likes = data.likes;
              hasLiked = data.hasLiked;
            }
          }
          return {
            videoID: v.videoID,
            thumbnailURL: v.thumbnailURL,
            title: v.title,
            totalLikeCount: likes,
            totalViewCount: v.totalViewCount,
            hasLiked,
          };
        })
      );
      setVideos(thumbs);
    } catch (err) {
      console.error('fetchVideos error:', err);
    } finally {
      setLoading(false);
    }
  };







  useEffect(() => {
    // only in browser
    if (typeof window === 'undefined') return;
    const token = localStorage.getItem('jwtToken');
    setJwtToken(token);
    setIsLoggedIn(!!token);
  }, []);




  

  useEffect(() => {
    if (!jwtToken) return;

    (async () => {
      try {
        const res = await fetch('/api/ws/notifications/unread', {
          headers: { Authorization: `Bearer ${jwtToken}` },
        });
        if (!res.ok) throw new Error('Failed to fetch unread notifications');
        const { notifications: unread }: { notifications: Notification[] } = await res.json();
        setNotifications(unread);
      } catch (err) {
        console.error('Notification fetch error:', err);
      }
    })();
  }, [jwtToken]);



  useEffect(() => {
      if (!socket || !isConnected) return;

      
      socket.onAny((event, payload) => {
        console.log('[WS]', event, payload);
      });
      socket.on('video:view', onView);
      socket.on('video:like', onLike);
      socket.on('video:comment', onComment);
      socket.on('notification', onNotification);

      return () => {
        socket.off('video:view', onView);
        socket.off('video:like', onLike);
        socket.off('video:comment', onComment);
        socket.off('notification', onNotification);
         socket.offAny();
      };
    }, [socket, isConnected, onView, onLike, onComment, onNotification]);





  useEffect(() => {
    const token = localStorage.getItem('jwtToken');
    if (token) {
      setIsLoggedIn(true);
      setJwtToken(token);
      fetchVideos(token);
    } else {
      fetchVideos(null);
    }
  }, []);

  // useEffect(() => {
  //   const iv = setInterval(() => fetchVideos(jwtToken), 10000);
  //   return () => clearInterval(iv);
  // }, [jwtToken]);



  const [likeInProgress, setLikeInProgress] = useState(false);

  const toggleLike = async (videoID: string, idx: number) => {
    if (!jwtToken) {
      setShowLoginPopup(true);
      return;
    }
    if (likeInProgress) return;

    setLikeInProgress(true);
    try {
      const res = await fetch(`/api/ws/likes/videos/${videoID}`, {
        method: videos[idx].hasLiked ? 'DELETE' : 'POST',
        headers: { Authorization: `Bearer ${jwtToken}` },
      });

      if (!res.ok) throw new Error('toggle failed');
      const { likes, hasLiked } = await res.json();

      setVideos(prev =>
        prev.map((v, i) =>
          i === idx ? { ...v, totalLikeCount: likes, hasLiked } : v
        )
      );

      if (selectedVideoIndex === idx && selectedVideo) {
        setSelectedVideo({
          ...selectedVideo,
          totalLikeCount: likes,
        });
        setModalHasLiked(hasLiked);
      }

    } catch (err) {
      console.error('toggleLike error:', err);
      alert('Could not update like');
    } finally {
      setLikeInProgress(false);
    }
  };




  

  // const toggleLike = async (videoID: string, idx: number) => {
  //   if (!jwtToken) {
  //     setShowLoginPopup(true);
  //     return;
  //   }
  //   try {
  //     const currentlyLiked = videos[idx].hasLiked;

  //     // const res = await fetch(`http://localhost/api/ws/likes/videos/${videoID}`, {
  //     const res = await fetch(`/api/ws/likes/videos/${videoID}`, {
  //       method: currentlyLiked ? 'DELETE' : 'POST',
  //       headers: { Authorization: `Bearer ${jwtToken}` },
  //     });
  //     if (!res.ok) throw new Error('toggle failed');
  //     const { likes, hasLiked } = await res.json();


  //     setVideos((prev) =>
  //         prev.map((v, i) =>
  //           i === idx
  //             ? {
  //                 ...v,
  //                 totalLikeCount: currentlyLiked ? v.totalLikeCount - 1 : v.totalLikeCount + 1,
  //                 hasLiked: !currentlyLiked,
  //               }
  //             : v
  //         )
  //       );

  //       if (selectedVideoIndex === idx && selectedVideo) {
  //         setSelectedVideo({
  //           ...selectedVideo,
  //           totalLikeCount: currentlyLiked
  //             ? selectedVideo.totalLikeCount - 1
  //             : selectedVideo.totalLikeCount + 1,
  //         });
  //         setModalHasLiked(!currentlyLiked);
  //       }


  //     if (selectedVideoIndex === idx && selectedVideo) {
  //       setSelectedVideo({ ...selectedVideo, totalLikeCount: likes });
  //       setModalHasLiked(hasLiked);
  //     }
  //   } catch (err) {
  //     console.error('toggleLike error:', err);
  //     alert('Could not update like');
  //   }
  // };

  const handleVideoClick = async (videoID: string, idx: number) => {
    if (!jwtToken) { setShowLoginPopup(true); return; }
    setPlaylistLoading(true);
    setSelectedVideo(null);
    setSelectedVideoIndex(idx);
    try {
      const [vRes, lRes, vwRes, cRes] = await Promise.all([

        fetch(`/api/video/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        fetch(`/api/ws/likes/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        fetch(`/api/ws/views/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        fetch(`/api/ws/comments/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),

        // fetch(`http://localhost:8090/api/video/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        // fetch(`http://localhost:8092/api/ws/likes/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        // fetch(`http://localhost:8092/api/ws/views/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
        // fetch(`http://localhost:8092/api/ws/comments/videos/${videoID}`, { headers: { Authorization: `Bearer ${jwtToken}` } }),
      ]);
      if (!vRes.ok) throw new Error();
      const vd = await vRes.json();
      const ld = lRes.ok ? await lRes.json() : { likes: 0, hasLiked: false };
      const vw = vwRes.ok ? await vwRes.json() : { views: 0 };
      const cm = cRes.ok ? await cRes.json() : { comments: [] };
      const userID = jwtTokenUserID;

      const commentsWithIsUser = Array.isArray(cm.comments) 
        ? cm.comments.map((comment: CommentDTO & { userID?: string }) => ({
            ...comment,
            isUser: comment.username === userID || comment.userID === userID, // or whichever identifies owner
          })) 
        : [];

      setSelectedVideo({
        playlist: vd.playlist,
        title: vd.title,
        totalViewCount: vw.views,
        totalLikeCount: ld.likes,
        comments: commentsWithIsUser,
      });

      setModalHasLiked(ld.hasLiked);
    } catch (err) {
      console.error('handleVideoClick error:', err);
      alert('Failed to load video');
    } finally {
      setPlaylistLoading(false);
    }
  };

  
  const handleSendComment = async () => {
    if (!jwtToken || !newComment.trim() || selectedVideoIndex === null || isSendingComment) return;

    setIsSendingComment(true); // prevent race condition
    try {
      const videoID = videos[selectedVideoIndex].videoID;

      const res = await fetch(`/api/ws/comments/videos/${videoID}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${jwtToken}`,
        },
        body: JSON.stringify({ content: newComment.trim() }),
      });

      if (!res.ok) throw new Error();

      const added = await res.json();

      setSelectedVideo(prev => prev
        ? {
            ...prev,
            comments: prev.comments.some(c => c.id === added.comment.id)
              ? prev.comments
              : [...prev.comments, added.comment]
          }
        : prev
      );
      setNewComment('');
    } catch (err) {
      console.error('handleSendComment error:', err);
      alert('Failed to post comment');
    } finally {
      setIsSendingComment(false);
    }
  };


  


  const handleDeleteComment = async (commentID: string) => {
    if (!jwtToken || selectedVideoIndex === null) return;
    if (!confirm("Are you sure you want to delete this comment?")) return;
    try {
      const videoID = videos[selectedVideoIndex].videoID;


      // const res = await fetch(`http://localhost/api/ws/comments/videos/${videoID}/${commentID}`, {

      const res = await fetch(`/api/ws/comments/videos/${videoID}/${commentID}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${jwtToken}` },
      });
      if (res.ok) {
        setSelectedVideo(prev => prev ? { ...prev, comments: prev.comments.filter(c => c.id !== commentID) } : prev);
      } else {
        alert("Failed to delete comment");
      }
    } catch (err) {
      console.error("Error deleting comment:", err);
      alert("An error occurred while deleting.");
    }
  };

  

  return (



    <div className="min-h-screen bg-gray-100">
      <header className="bg-white shadow p-4 flex items-center justify-between">
        <h1 className="text-3xl font-bold text-black">TokTik</h1>
        <NotificationBell
          notifications={notifications}
          showNotifications={isOpen}
          onToggle={handleToggleBell}
          unreadCount={unreadCount}
        />
      </header>

      <main className="p-6 flex">
        <Sidebar
          isLoggedIn={isLoggedIn}
          onSuccessLogin={(token) => {
            localStorage.setItem('jwtToken', token);
            setJwtToken(token);
            setIsLoggedIn(true);
            fetchVideos(token);
          }}
          onLogout={() => {
            localStorage.removeItem('jwtToken');
            setJwtToken(null);
            setIsLoggedIn(false);
            fetchVideos(null);
          }}
          
        />
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6 pl-[250px] w-full">
          {loading ? (
            <div>Loading...</div>
          ) : (
            videos.map((v, i) => (
              <div
                key={v.videoID}
                onClick={() => handleVideoClick(v.videoID, i)}
                className="bg-gray-300 h-56 rounded-lg shadow-md flex flex-col justify-between p-2 cursor-pointer hover:bg-gray-400"
              >
                <img
                  src={v.thumbnailURL}
                  alt={v.title}
                  className="w-full h-32 object-cover rounded-t"
                />
                <div className="text-lg font-bold text-gray-800 truncate px-1 text-center">
                  {v.title}
                </div>
                <div className="flex gap-4 items-center justify-center text-black">
                  <span
                    onClick={(e: React.MouseEvent) => {
                      e.stopPropagation();
                      toggleLike(v.videoID, i);
                    }}
                    className="cursor-pointer"
                  >
                    <BsFillHandThumbsUpFill
                      size={22}
                      color={v.hasLiked ? '#dc2626' : '#6b7280'}
                    />
                  </span>
                  <span>{v.totalLikeCount}</span>
                  <ImEye size={18} />
                  <span>{v.totalViewCount}</span>
                </div>
              </div>
            ))
          )}
        </div>
      </main>

      {showLoginPopup && <LoginForm onClose={() => setShowLoginPopup(false)} onSuccessLogin={(token) => { localStorage.setItem('jwtToken', token); setJwtToken(token); setIsLoggedIn(true); setShowLoginPopup(false); fetchVideos(token); }} />}

      {playlistLoading && <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center"><div className="bg-white p-6 rounded">Loading...</div></div>}

      {selectedVideo && selectedVideoIndex !== null && (
        <div className="fixed inset-0 bg-black bg-opacity-70 flex items-center justify-center">
          <div className="bg-white rounded shadow-lg w-11/12 max-w-4xl p-4 relative flex flex-col md:flex-row">
            {selectedVideoIndex > 0 && <button className="absolute left-2 top-1/2 -translate-y-1/2 text-3xl text-gray-600 hover:text-black" onClick={() => handleVideoClick(videos[selectedVideoIndex-1].videoID, selectedVideoIndex-1)}><ImArrowLeft2 /></button>}
            {selectedVideoIndex < videos.length-1 && <button className="absolute right-2 top-1/2 -translate-y-1/2 text-3xl text-gray-600 hover:text-black" onClick={() => handleVideoClick(videos[selectedVideoIndex+1].videoID, selectedVideoIndex+1)}><ImArrowRight2 /></button>}
            <button className="absolute top-2 right-2 text-2xl text-gray-600 hover:text-black" onClick={() => { setSelectedVideo(null); setSelectedVideoIndex(null); }}>&times;</button>
            <div className="md:w-2/3 w-full flex items-center justify-center"><HlsPlayer playlistUrl={selectedVideo.playlist} /></div>
            <div className="md:w-1/3 w-full p-4 flex flex-col text-black">
              <h2 className="text-xl font-bold truncate align-middle mb-2">{selectedVideo.title}</h2>
              <div className="flex gap-4 items-center mb-4 text-black">
                <span className="cursor-pointer" onClick={() => toggleLike(videos[selectedVideoIndex].videoID, selectedVideoIndex)}>
                  <BsFillHandThumbsUpFill size={22} color={modalHasLiked ? '#dc2626' : '#6b7280'} />
                </span>
                <span>{selectedVideo.totalLikeCount}</span>
                <ImEye size={18} /> <span>{selectedVideo.totalViewCount}</span>
              </div>
              <div className="flex flex-col h-full">
                <div className="flex-grow overflow-y-auto max-h-48 space-y-2">
                  <h3 className="font-semibold mb-2">Comments</h3>
                  {selectedVideo.comments?.length ? (
                    selectedVideo.comments.map((c) => (
                      <div key={c.id} className="border-b pb-1 flex justify-between items-center">
                        <span className="font-bold">{c.username}:</span> {c.content}

                        {c.isUser && (
                          <button onClick={() => handleDeleteComment(c.id)} className="text-red-600 hover:text-red-800 ml-2" title="Delete comment">🗑️</button>
                        )}
                      </div>
                    ))
                  ) : (
                    <div className="text-gray-400">No comments yet.</div>
                  )}
                </div>
                <div className="mt-2">
                  <textarea value={newComment} onChange={(e) => setNewComment(e.target.value)} placeholder="Add a comment..." className="w-full p-2 border border-gray-300 rounded resize-none" rows={2} />
                  <button onClick={handleSendComment} disabled={isSendingComment || !newComment.trim()} className="mt-2 px-4 py-1 bg-blue-600 text-white rounded disabled:opacity-50">
                    {isSendingComment ? 'Sending...' : 'Send'}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
















