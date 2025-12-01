import React, { useEffect, useRef} from 'react';
import Hls from 'hls.js';

interface HlsPlayerProps {
  playlistUrl: string;
}

const HlsPlayerComponent: React.FC<HlsPlayerProps> = ({ playlistUrl }) => {
  const ref = useRef<HTMLVideoElement>(null);

  useEffect(() => {
    if (ref.current && Hls.isSupported()) {
      const hls = new Hls();
      hls.loadSource(playlistUrl);
      hls.attachMedia(ref.current);

      return () => {
        hls.destroy();
      };
    } else if (ref.current) {
      ref.current.src = playlistUrl;
    }
  }, [playlistUrl]);

  return <video ref={ref} controls autoPlay className="w-full h-80 rounded bg-black" />;
};

export const HlsPlayer = React.memo(HlsPlayerComponent);
