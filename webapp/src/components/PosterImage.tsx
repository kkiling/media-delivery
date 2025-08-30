import { useState } from 'react';
import { Image as ImageIcon } from 'react-bootstrap-icons';

interface NoImageFallbackProps {
  text?: string;
}

const NoImageFallback = ({ text = 'No Image Available' }: NoImageFallbackProps) => (
  <div className="d-flex flex-column align-items-center justify-content-center bg-secondary text-white w-100 h-100">
    <ImageIcon size={48} className="mb-2" />
    <p className="mb-0">{text}</p>
  </div>
);

interface PosterImageProps {
  src?: string;
  alt: string;
  minHeight?: number;
}

export const PosterImage = ({ src, alt, minHeight }: PosterImageProps) => {
  const [error, setError] = useState(false);

  if (!src || error) {
    return <NoImageFallback />;
  }

  return (
    <div className="d-flex justify-content-center">
      <div
        style={{
          height: minHeight,
        }}
      >
        <img
          src={src}
          alt={alt}
          onError={() => setError(true)}
          className="w-100 h-100 object-fit-cover rounded"
        />
      </div>
    </div>
  );
};
