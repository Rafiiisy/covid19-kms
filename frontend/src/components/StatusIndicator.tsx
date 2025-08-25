import React from 'react';
import { CheckCircle, XCircle, Clock, AlertCircle } from 'lucide-react';

interface StatusIndicatorProps {
  status: string;
  message?: string;
  timestamp?: string;
}

export const StatusIndicator: React.FC<StatusIndicatorProps> = ({ status, message, timestamp }) => {
  const getStatusConfig = (status: string) => {
    switch (status.toLowerCase()) {
      case 'success':
        return {
          icon: CheckCircle,
          color: 'text-green-600',
          bgColor: 'bg-green-100',
          borderColor: 'border-green-200',
        };
      case 'error':
        return {
          icon: XCircle,
          color: 'text-red-600',
          bgColor: 'bg-red-100',
          borderColor: 'border-red-200',
        };
      case 'running':
      case 'processing':
        return {
          icon: Clock,
          color: 'text-blue-600',
          bgColor: 'bg-blue-100',
          borderColor: 'border-blue-200',
        };
      default:
        return {
          icon: AlertCircle,
          color: 'text-yellow-600',
          bgColor: 'bg-yellow-100',
          borderColor: 'border-yellow-200',
        };
    }
  };

  const config = getStatusConfig(status);
  const IconComponent = config.icon;

  return (
    <div className={`inline-flex items-center px-3 py-2 rounded-full text-sm font-medium ${config.bgColor} ${config.borderColor} border`}>
      <IconComponent className={`w-4 h-4 mr-2 ${config.color}`} />
      <span className={config.color}>
        {status.charAt(0).toUpperCase() + status.slice(1)}
      </span>
      {message && (
        <span className="ml-2 text-gray-600">- {message}</span>
      )}
      {timestamp && (
        <span className="ml-2 text-gray-500 text-xs">
          {new Date(timestamp).toLocaleString()}
        </span>
      )}
    </div>
  );
};


