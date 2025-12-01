'use client';

import {
  useState,
  useCallback,
  ChangeEvent,
  DragEvent,
} from 'react';


import axios, { AxiosRequestHeaders } from 'axios';

export default function UploadBox() {

  const [videoFile, setVideoFile] = useState<File | null>(null);
  const [videoURL, setVideoURL] = useState<string | null>(null);
  const [title, setTitle] = useState('');
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  const axiosInstance = axios.create();
  axiosInstance.interceptors.request.use((config) => {
    console.log('Intercepted request headers:', config.headers as AxiosRequestHeaders);
    return config;
  });

  const handleFile = useCallback((file: File) => {
    setError(null);
    setSuccessMsg(null);
    if (file.type !== 'video/mp4') {
      setError('Only MP4 files are allowed.');
      return;
    }
    const url = URL.createObjectURL(file);
    const tempVid = document.createElement('video');
    tempVid.preload = 'metadata';
    tempVid.src = url;
    tempVid.onloadedmetadata = () => {
      if (tempVid.duration > 60) {
        setError('Video must be 60 seconds or shorter.');
        URL.revokeObjectURL(url);
      } else {
        setVideoFile(file);
        setVideoURL(url);
      }
    };
    tempVid.onerror = () => {
      setError('Could not load the selected video.');
      URL.revokeObjectURL(url);
    };
  }, []);

  const onFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const f = e.target.files?.[0];
    if (f) handleFile(f);
  };

  const onDrop = (e: DragEvent<HTMLLabelElement>) => {
    e.preventDefault();
    const f = e.dataTransfer.files?.[0];
    if (f) handleFile(f);
  };

  const handleUpload = async () => {
    if (!videoFile) return setError('Please select a video first.');
    if (!title.trim()) return setError('Please enter a title.');
    const token = localStorage.getItem('jwtToken');
    if (!token) return setError('You must be logged in to upload.');

    setUploading(true);
    setProgress(0);
    setError(null);
    setSuccessMsg(null);

    let videoID: string | undefined;

    try {
      const res = await fetch('/api/video/videos', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ filename: videoFile.name, title: title.trim() }),
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(`Backend Error: ${text}`);
      }

      const { videoID: vid, presignedURL } = await res.json();
      videoID = vid;

      if (!presignedURL) throw new Error('No upload URL returned');

      await axios.put(presignedURL, videoFile, {
        onUploadProgress: (event) => {
          if (event.total) {
            setProgress(Math.round((event.loaded / event.total) * 100));
          }
        },
      });

      const checkRes = await fetch(`/api/video/videos/check?videoID=${videoID}`, {
        method: 'GET',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      const checkData = await checkRes.json();

      if (checkData.uploadStatus === true) {
        setSuccessMsg('Upload successful!');
      } else {
        alert('Upload not found or not complete.');
      }

      setVideoFile(null);
      setVideoURL(null);
      setTitle('');
    } catch (err: unknown) {
      if (videoID) {
        try {
          await fetch(`/api/video/videos/${videoID}`, {
            method: 'DELETE',
            headers: {
              Authorization: `Bearer ${token}`,
            },
          });
        } catch (deleteErr) {
          console.error('Failed to delete video record after upload failure', deleteErr);
        }
      }

      if (err instanceof Error) setError(err.message);
      else setError('An unknown error occurred');
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="bg-white rounded-xl p-6 max-w-2xl w-full mx-auto my-6 shadow">
      <h2 className="text-2xl font-semibold mb-4 text-center">Upload Your Video</h2>

      <label
        htmlFor="video-file"
        onDrop={onDrop}
        onDragOver={(e) => e.preventDefault()}
        className="block mb-4 cursor-pointer border-2 border-dashed border-gray-300 rounded-md p-4 text-center hover:border-blue-500"
      >
        {videoURL ? (
          <video src={videoURL} controls className="max-w-full max-h-48 rounded mb-2" />
        ) : (
          <span className="text-gray-500">Drag & drop or click to select a video</span>
        )}
        <input
          id="video-file"
          type="file"
          accept="video/mp4"
          onChange={onFileChange}
          className="hidden"
          disabled={uploading}
        />
      </label>

      <div className="flex justify-around text-gray-600 text-sm mb-4">
        <div className="flex flex-col items-center">
          <span className="text-2xl">🎥</span>
          <span>MP4 only</span>
        </div>
        <div className="flex flex-col items-center">
          <span className="text-2xl">⏱️</span>
          <span>≤ 60 sec</span>
        </div>
      </div>

      <div className="mb-4">
        <label htmlFor="video-title" className="block mb-1 text-gray-700">
          Video Title
        </label>
        <input
          id="video-title"
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          className="w-full rounded border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
          disabled={uploading}
        />
      </div>

      {uploading && (
        <div className="w-full bg-gray-200 rounded h-4 mb-4">
          <div className="h-full bg-blue-600 rounded" style={{ width: `${progress}%` }} />
        </div>
      )}

      {error && <p className="text-red-600 text-center mb-2">{error}</p>}
      {successMsg && <p className="text-green-600 text-center mb-2">{successMsg}</p>}

      <div className="flex justify-end gap-3">
        <button
          onClick={() => {
            setVideoFile(null);
            setVideoURL(null);
            setTitle('');
            setError(null);
            setSuccessMsg(null);
          }}
          disabled={uploading}
          className="px-4 py-2 border border-gray-400 rounded hover:bg-gray-100 disabled:opacity-50"
        >
          Reset
        </button>
        <button
          onClick={handleUpload}
          disabled={uploading}
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
        >
          {uploading ? `Uploading ${progress}%` : 'Upload'}
        </button>
      </div>
    </div>
  );
}
