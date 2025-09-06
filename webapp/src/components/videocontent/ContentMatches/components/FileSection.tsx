import React from 'react';
import { Track } from '@/api/api';

interface FileSectionProps {
  title: string;
  icon: React.ReactNode;
  files: Track[];
  showName?: boolean;
}

export const FileSection: React.FC<FileSectionProps> = ({
  title,
  icon,
  files,
  showName = false,
}) => (
  <div className="mb-3">
    <div className="d-flex align-items-center gap-2 mb-2">
      {icon}
      <span>{title}</span>
    </div>
    {files.map((file, idx) => (
      <div key={idx} className="ps-4 mb-2 d-flex align-items-center">
        {showName && <div className="text-break">{file.name || 'Unnamed track'}</div>}
        <div className={showName ? 'text-secondary ms-2' : 'text-break'}>{file.relative_path}</div>
      </div>
    ))}
  </div>
);
