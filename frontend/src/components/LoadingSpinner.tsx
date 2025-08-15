import React from 'react';
import { RefreshCw } from 'lucide-react';

interface LoadingSpinnerProps {
  message?: string;
  size?: 'sm' | 'md' | 'lg';
}

export const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({ 
  message = 'Processing...', 
  size = 'md' 
}) => {
  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-6 h-6',
    lg: 'w-8 h-8',
  };

  return (
    <div className="flex flex-col items-center justify-center p-8">
      <RefreshCw className={`${sizeClasses[size]} text-blue-600 animate-spin mb-4`} />
      <p className="text-gray-600 font-medium">{message}</p>
      <p className="text-sm text-gray-500 mt-2">This may take a few minutes...</p>
    </div>
  );
};
