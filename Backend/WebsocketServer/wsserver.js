console.log('App starting...');

require('dotenv').config();

console.log('dotenv loaded');
console.log('PORT:', process.env.PORT);
console.log('JWT_SECRET:', process.env.JWT_SECRET);

const { createServer } = require('http');
const express = require('express');
const { Server } = require('socket.io');
const Redis = require('ioredis');
const jwt = require('jsonwebtoken');

const PORT = process.env.PORT || 8080;
const JWT_SECRET = process.env.JWT_SECRET;

console.log(`[Startup] PORT=${PORT}`);
console.log(`[Startup] JWT_SECRET=${JWT_SECRET ? '***set***' : '***NOT SET***'}`);

// -------------------- Redis Setup --------------------

const redis = new Redis.Cluster(
  [
    { host: process.env.REDIS_HOST_0, port: parseInt(process.env.REDIS_PORT, 10) },
    { host: process.env.REDIS_HOST_1, port: parseInt(process.env.REDIS_PORT, 10) },
    { host: process.env.REDIS_HOST_2, port: parseInt(process.env.REDIS_PORT, 10) },
    { host: process.env.REDIS_HOST_3, port: parseInt(process.env.REDIS_PORT, 10) },
    { host: process.env.REDIS_HOST_4, port: parseInt(process.env.REDIS_PORT, 10) },
    { host: process.env.REDIS_HOST_5, port: parseInt(process.env.REDIS_PORT, 10) },
  ],
  {
    scaleReads: 'all',
    redisOptions: {
      password: process.env.REDIS_PASSWORD,
    },
  }
);

redis.on('connect', () => console.log('[Redis] Connected to Redis cluster'));
redis.on('error', (err) => console.error('[Redis] Error:', err));

// Standalone Redis for Pub/Sub
const redisSub = new Redis({
  host: process.env.REDIS_PUBSUB_HOST || process.env.REDIS_HOST_0,
  port: parseInt(process.env.REDIS_PORT, 10),
  password: process.env.REDIS_PASSWORD,
});

redisSub.on('connect', () => console.log('[RedisSub] Connected to Redis pub/sub node'));
redisSub.on('error', (err) => console.error('[RedisSub] Error:', err));

// -------------------- Express App --------------------

const app = express();
app.get('/health', (_, res) => res.status(200).send('OK'));

const httpServer = createServer(app);

// -------------------- Socket.IO Setup --------------------

const io = new Server(httpServer, {
  cors: {
    origin: '*',
    credentials: true,
  },
  transports: ['websocket'],
});

// Track multiple sockets per user
const userSockets = new Map(); // Map<string, Set<Socket>>

// Auth middleware
io.use((socket, next) => {
  const token = socket.handshake.auth.token;
  if (!token) return next(new Error('No token provided'));

  try {
    const payload = jwt.verify(token, JWT_SECRET);
    socket.userID = String(payload.sub);
    next();
  } catch (err) {
    next(new Error('Invalid token'));
  }
});

// Connection handler
io.on('connection', (socket) => {
  const userID = socket.userID;
  console.log(`[Socket] Connected: ${userID} (socket=${socket.id})`);

  if (!userSockets.has(userID)) {
    userSockets.set(userID, new Set());
  }
  userSockets.get(userID).add(socket);

  socket.on('disconnect', () => {
    console.log(`[Socket] Disconnected: ${userID} (socket=${socket.id})`);
    const sockets = userSockets.get(userID);
    if (sockets) {
      sockets.delete(socket);
      if (sockets.size === 0) {
        userSockets.delete(userID);
      }
    }
  });

  socket.on('error', (err) => {
    console.error(`[Socket] Error on ${userID}:`, err.message);
  });
});

// -------------------- Redis Pub/Sub Listener --------------------

redisSub.subscribe('notifications', (err) => {
  if (err) {
    console.error('[RedisSub] Failed to subscribe to "notifications":', err);
  } else {
    console.log('[RedisSub] Subscribed to "notifications" channel');
  }
});

redisSub.on('message', (channel, message) => {
  if (channel !== 'notifications') return;

  console.log('[RedisSub]<', channel, message);

  try {
    const data = JSON.parse(message);
    const { to, type, data: payload } = data;

    if (to) {
      const sockets = userSockets.get(String(to));
      if (sockets && sockets.size > 0) {
        sockets.forEach((s) => s.emit(type, payload));
        console.log(`[Redis -> WS] Sent ${type} to userID=${to}`);
      } else {
        console.log(`[Redis] User ${to} is offline`);
      }
    } else {
      switch (type) {
        case 'video:view':
          io.emit('video:view', {
            videoID: payload.videoID,
            totalViewCount: payload.totalViewCount,
          });
          break;

        case 'video:like':
          io.emit('video:like', {
            videoID:        payload.videoID,
            totalLikeCount: payload.totalLikeCount,
            hasLiked:       payload.hasLiked,
            userID:         payload.userID,
            username:       payload.username,
          });
          break;


        case 'video:comment':
          io.emit('video:comment', {
            videoID: payload.videoID,
            comment: payload.comment,
          });
          break;

        case 'notification':
          io.emit('notification', {
            ...payload,
            createdAt: payload.createdAt || new Date().toISOString(), 
          });
          break;


        default:
          console.warn(`[Redis] Unknown event type: ${type}`);
      }
    }
  } catch (err) {
    console.error('[RedisSub] Failed to parse message:', err.message);
  }
});

// -------------------- Start Server --------------------

httpServer.listen(PORT, () => {
  console.log(`[Startup] WebSocket server listening on port ${PORT}`);
});



// -------------------- Process-level Error Handling --------------------

process.on('uncaughtException', (err) => {
  console.error('[Fatal] Uncaught Exception:', err);
  process.exit(1);
});

process.on('unhandledRejection', (reason, promise) => {
  console.error('[Fatal] Unhandled Rejection at:', promise, 'reason:', reason);
});
