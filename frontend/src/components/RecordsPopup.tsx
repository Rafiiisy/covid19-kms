import React, { useState, useMemo } from 'react';
import { X, ExternalLink, MessageCircle, ThumbsUp, Calendar, TrendingUp, ChevronLeft, ChevronRight } from 'lucide-react';

interface RecordsPopupProps {
  isOpen: boolean;
  onClose: () => void;
  databaseData: {
    youtube: any[] | null;
    googleNews: any[] | null;
    instagram: any[] | null;
    indonesiaNews: any[] | null;
    summary: any | null;
  } | null;
}

export const RecordsPopup: React.FC<RecordsPopupProps> = ({ isOpen, onClose, databaseData }) => {
  const [activeTab, setActiveTab] = useState('youtube');
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 25;

  // Reset to first page when tab changes
  const handleTabChange = (tabId: string) => {
    setActiveTab(tabId);
    setCurrentPage(1);
  };

  // Get current tab data
  const getCurrentTabData = () => {
    switch (activeTab) {
      case 'youtube':
        return databaseData?.youtube || [];
      case 'google-news':
        return databaseData?.googleNews || [];
      case 'instagram':
        return databaseData?.instagram || [];
      case 'indonesia-news':
        return databaseData?.indonesiaNews || [];
      default:
        return [];
    }
  };

  // Pagination logic
  const currentTabData = getCurrentTabData();
  const totalPages = Math.ceil(currentTabData.length / itemsPerPage);
  const startIndex = (currentPage - 1) * itemsPerPage;
  const endIndex = startIndex + itemsPerPage;
  const currentPageData = currentTabData.slice(startIndex, endIndex);

  // Pagination controls
  const goToPage = (page: number) => {
    setCurrentPage(Math.max(1, Math.min(page, totalPages)));
  };

  const goToPreviousPage = () => {
    setCurrentPage(prev => Math.max(1, prev - 1));
  };

  const goToNextPage = () => {
    setCurrentPage(prev => Math.min(totalPages, prev + 1));
  };

  // Generate page numbers for pagination
  const getPageNumbers = () => {
    const pages = [];
    const maxVisiblePages = 5;
    
    if (totalPages <= maxVisiblePages) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      if (currentPage <= 3) {
        for (let i = 1; i <= 4; i++) {
          pages.push(i);
        }
        pages.push('...');
        pages.push(totalPages);
      } else if (currentPage >= totalPages - 2) {
        pages.push(1);
        pages.push('...');
        for (let i = totalPages - 3; i <= totalPages; i++) {
          pages.push(i);
        }
      } else {
        pages.push(1);
        pages.push('...');
        for (let i = currentPage - 1; i <= currentPage + 1; i++) {
          pages.push(i);
        }
        pages.push('...');
        pages.push(totalPages);
      }
    }
    
    return pages;
  };

  if (!isOpen) return null;

  const tabs = [
    { id: 'youtube', name: 'YouTube Comments', icon: '💬', count: databaseData?.youtube?.length || 0 },
    { id: 'google-news', name: 'Google News', icon: '📰', count: databaseData?.googleNews?.length || 0 },
    { id: 'instagram', name: 'Instagram', icon: '📱', count: databaseData?.instagram?.length || 0 },
    { id: 'indonesia-news', name: 'Indonesia News', icon: '🇮🇩', count: databaseData?.indonesiaNews?.length || 0 },
  ];

  const renderYouTubeContent = () => (
    <div className="space-y-0">
      {currentPageData.length ? (
        currentPageData.map((comment, index) => (
          <div key={startIndex + index} className={`p-4 transition-all duration-200 hover:shadow-md ${
            index % 2 === 0 ? 'bg-red-100' : 'bg-gray-200'
          } ${index < currentPageData.length - 1 ? 'border-b border-gray-300' : ''} ${
            index === 0 ? 'rounded-t-lg' : ''
          } ${
            index === currentPageData.length - 1 ? 'rounded-b-lg' : ''
          }`}>
            <div className="flex items-start space-x-3">
              <div className="flex-shrink-0 w-16 h-12 bg-red-100 rounded flex items-center justify-center">
                <span className="text-red-600 text-xs">💬</span>
              </div>
              <div className="flex-1 min-w-0">
                <h4 className="text-sm font-medium text-gray-900 line-clamp-2 text-left">
                  {comment.description || 'No comment content'}
                </h4>
                <p className="text-xs text-gray-500 mt-1 line-clamp-2 text-left">
                  <span className="font-medium">Video:</span> {comment.title || 'Untitled Video'}
                </p>
                
                {/* Comment-specific metadata */}
                <div className="grid grid-cols-2 gap-2 mt-3 text-xs text-gray-600">
                  {comment.channel_title && (
                    <div className="flex items-center space-x-1">
                      <span className="font-medium">📺</span>
                      <span>{comment.channel_title}</span>
                    </div>
                  )}
                  {comment.word_count && (
                    <div className="flex items-center space-x-1">
                      <span className="font-medium">📝</span>
                      <span>{comment.word_count} words</span>
                    </div>
                  )}
                  {comment.language && (
                    <div className="flex items-center space-x-1">
                      <span className="font-medium">🌐</span>
                      <span>{comment.language}</span>
                    </div>
                  )}
                </div>

                {/* General metadata */}
                <div className="flex items-center space-x-4 mt-3 text-xs text-gray-500">
                  <span className="flex items-center">
                    <Calendar className="w-3 h-3 mr-1" />
                    {comment.published_at ? new Date(comment.published_at).toLocaleDateString() : 'Unknown date'}
                  </span>
                  {comment.covid_relevance_score && (
                    <span className="flex items-center">
                      <TrendingUp className="w-3 h-3 mr-1" />
                      {comment.covid_relevance_score.toFixed(2)}
                    </span>
                  )}
                  <span className="px-2 py-1 rounded text-xs bg-blue-100 text-blue-800">
                    Comment
                  </span>
                </div>
              </div>
            </div>
          </div>
        ))
      ) : (
        <div className="text-center py-8 text-gray-500">
          <p>No YouTube comments available</p>
        </div>
        )}
    </div>
  );

  const renderGoogleNewsContent = () => (
    <div className="space-y-0">
      {currentPageData.length ? (
        currentPageData.map((article, index) => (
          <div key={startIndex + index} className={`p-4 transition-all duration-200 hover:shadow-md ${
            index % 2 === 0 ? 'bg-blue-100' : 'bg-gray-200'
          } ${index < currentPageData.length - 1 ? 'border-b border-gray-300' : ''} ${
            index === 0 ? 'rounded-t-lg' : ''
          } ${
            index === currentPageData.length - 1 ? 'rounded-b-lg' : ''
          }`}>
            <div className="flex items-start space-x-3">
              <div className="flex-shrink-0 w-16 h-12 bg-blue-100 rounded flex items-center justify-center">
                <span className="text-blue-600 text-xs">📰</span>
              </div>
              <div className="flex-1 min-w-0">
                <h4 className="text-sm font-medium text-gray-900 line-clamp-2 text-left">
                  {article.title || 'Untitled Article'}
                </h4>
                
                {/* News-specific metadata */}
                <div className="grid grid-cols-2 gap-2 mt-3 text-xs text-gray-600">
                  {article.author && (
                    <div className="flex items-center space-x-1">
                      <span className="font-medium">✍️</span>
                      <span>{article.author}</span>
                    </div>
                  )}
                  <div className="flex items-center space-x-2">
                    {article.news_source && (
                      <div className="flex items-center space-x-1">
                        <span className="font-medium">🏢</span>
                        <span>{article.news_source}</span>
                      </div>
                    )}
                    {article.language && (
                      <div className="flex items-center space-x-1">
                        <span className="font-medium">🌐</span>
                        <span>{article.language.toUpperCase()}</span>
                      </div>
                    )}
                  </div>
                  {article.category && (
                    <div className="flex items-center space-x-1">
                      <span className="font-medium">📂</span>
                      <span>{article.category}</span>
                    </div>
                  )}
                </div>

                {/* General metadata */}
                <div className="flex items-center space-x-4 mt-3 text-xs text-gray-500">
                  <span className="flex items-center">
                    <Calendar className="w-3 h-3 mr-1" />
                    {article.processed_at ? new Date(article.processed_at).toLocaleDateString() : 'Unknown date'}
                  </span>
                  {article.relevance_score && (
                    <span className="flex items-center">
                      <TrendingUp className="w-3 h-3 mr-1" />
                      {article.relevance_score.toFixed(2)}
                    </span>
                  )}
                  {article.sentiment && (
                    <span className={`px-2 py-1 rounded text-xs ${
                      article.sentiment === 'positive' ? 'bg-green-100 text-green-800' :
                      article.sentiment === 'negative' ? 'bg-red-100 text-red-800' :
                      'bg-gray-100 text-gray-800'
                    }`}>
                      {article.sentiment}
                    </span>
                  )}
                </div>
              </div>
            </div>
          </div>
        ))
      ) : (
        <div className="text-center py-8 text-gray-500">
          <p>No Google News data available</p>
          <p className="mt-2 text-sm text-gray-400">Run the ETL pipeline to fetch news articles</p>
        </div>
      )}
    </div>
  );

  const renderInstagramContent = () => (
    <div className="space-y-0">
      {currentPageData.length ? (
        currentPageData.map((post, index) => (
          <div key={startIndex + index} className={`p-4 transition-all duration-200 hover:shadow-md ${
            index % 2 === 0 ? 'bg-blue-100' : 'bg-gray-200'
          } ${index < currentPageData.length - 1 ? 'border-b border-gray-300' : ''} ${
            index === 0 ? 'rounded-t-lg' : ''
          } ${
            index === currentPageData.length - 1 ? 'rounded-b-lg' : ''
          }`}>
            <div className="flex items-start space-x-3">
              <div className="flex-shrink-0 w-16 h-12 bg-pink-100 rounded flex items-center justify-center">
                <span className="text-pink-600 text-xs">📱</span>
              </div>
              <div className="flex-1 min-w-0">
                <h4 className="text-sm font-medium text-gray-900 line-clamp-2 text-left">
                  {post.title || 'Instagram Post'}
                </h4>
                <p className="text-xs text-gray-500 mt-1 line-clamp-2 text-left">
                  {post.content || 'No caption available'}
                </p>
                
                {/* Instagram-specific metadata */}
                <div className="grid grid-cols-2 gap-2 mt-3 text-xs text-gray-600">
                  {post.likes && (
                    <div className="flex items-center space-x-1">
                      <span className="font-medium">❤️</span>
                      <span>{post.likes.toLocaleString()} likes</span>
                    </div>
                  )}
                  {post.comments && (
                    <div className="flex items-center space-x-1">
                      <span className="font-medium">💬</span>
                      <span>{post.comments} comments</span>
                    </div>
                  )}
                  {post.post_id && (
                    <div className="flex items-center space-x-1">
                      <span className="font-medium">🆔</span>
                      <span className="font-mono text-xs">{post.post_id}</span>
                    </div>
                  )}
                  {post.hashtags && Array.isArray(post.hashtags) && (
                    <div className="flex items-center space-x-1">
                      <span className="font-medium">🏷️</span>
                      <span className="truncate">{post.hashtags.slice(0, 3).join(', ')}</span>
                    </div>
                  )}
                </div>

                {/* General metadata */}
                <div className="flex items-center space-x-4 mt-3 text-xs text-gray-500">
                  <span className="flex items-center">
                    <Calendar className="w-3 h-3 mr-1" />
                    {post.processed_at ? new Date(post.processed_at).toLocaleDateString() : 'Unknown date'}
                  </span>
                  {post.relevance_score && (
                    <span className="flex items-center">
                      <TrendingUp className="w-3 h-3 mr-1" />
                      {post.relevance_score.toFixed(2)}
                    </span>
                  )}
                  {post.sentiment && (
                    <span className={`px-2 py-1 rounded text-xs ${
                      post.sentiment === 'positive' ? 'bg-green-100 text-green-800' :
                      post.sentiment === 'negative' ? 'bg-red-100 text-red-800' :
                      'bg-gray-100 text-gray-800'
                    }`}>
                      {post.sentiment}
                    </span>
                  )}
                </div>
              </div>
            </div>
          </div>
        ))
      ) : (
        <div className="text-center py-8 text-gray-500">
          <p>No Instagram data available</p>
        </div>
      )}
    </div>
  );

  const renderIndonesiaNewsContent = () => (
    <div className="space-y-0">
      {currentPageData.length ? (
        currentPageData.map((news, index) => (
          <div key={startIndex + index} className={`p-4 transition-all duration-200 hover:shadow-md ${
            index % 2 === 0 ? 'bg-blue-100' : 'bg-gray-200'
          } ${index < currentPageData.length - 1 ? 'border-b border-gray-300' : ''} ${
            index === 0 ? 'rounded-t-lg' : ''
          } ${
            index === currentPageData.length - 1 ? 'rounded-b-lg' : ''
          }`}>
            <div className="flex items-start space-x-3">
              <div className="flex-shrink-0 w-16 h-12 bg-red-100 rounded flex items-center justify-center">
                <span className="text-red-600 text-xs">🇮🇩</span>
              </div>
              <div className="flex-1 min-w-0">
                <h4 className="text-sm font-medium text-gray-900 line-clamp-2 text-left">
                  {news.title || 'Untitled News'}
                </h4>
                
                {/* Indonesia News-specific metadata */}
                <div className="grid grid-cols-2 gap-2 mt-3 text-xs text-gray-600">
                  {news.author && (
                    <div className="flex items-center space-x-1">
                      <span className="font-medium">✍️</span>
                      <span>{news.author}</span>
                    </div>
                  )}
                  {news.news_source && (
                    <div className="flex items-center space-x-1">
                      <span className="font-medium">🏢</span>
                      <span>{news.news_source}</span>
                    </div>
                  )}
                  {news.category && (
                    <div className="flex items-center space-x-1">
                      <span className="font-medium">📂</span>
                      <span>{news.category}</span>
                    </div>
                  )}
                  {news.region && (
                    <div className="flex items-center space-x-1">
                      <span className="font-medium">📍</span>
                      <span>{news.region}</span>
                    </div>
                  )}
                </div>

                {/* General metadata */}
                <div className="flex items-center space-x-4 mt-3 text-xs text-gray-500">
                  <span className="flex items-center">
                    <Calendar className="w-3 h-3 mr-1" />
                    {news.processed_at ? new Date(news.processed_at).toLocaleDateString() : 'Unknown date'}
                  </span>
                  {news.relevance_score && (
                    <span className="flex items-center">
                      <TrendingUp className="w-3 h-1 mr-1" />
                      {news.relevance_score.toFixed(2)}
                    </span>
                  )}
                  {news.sentiment && (
                    <span className={`px-2 py-1 rounded text-xs ${
                      news.sentiment === 'positive' ? 'bg-green-100 text-green-800' :
                      news.sentiment === 'negative' ? 'bg-red-100 text-red-800' :
                      'bg-gray-100 text-gray-800'
                    }`}>
                      {news.sentiment}
                    </span>
                  )}
                </div>
              </div>
            </div>
          </div>
        ))
      ) : (
        <div className="text-center py-8 text-gray-500">
          <p>No Indonesia News data available</p>
          <p className="mt-2 text-sm text-gray-400">Run the ETL pipeline to fetch Indonesian news</p>
        </div>
      )}
    </div>
  );

  const renderTabContent = () => {
    switch (activeTab) {
      case 'youtube':
        return renderYouTubeContent();
      case 'google-news':
        return renderGoogleNewsContent();
      case 'instagram':
        return renderInstagramContent();
      case 'indonesia-news':
        return renderIndonesiaNewsContent();
      default:
        return renderYouTubeContent();
    }
  };

  // Pagination component
  const renderPagination = () => {
    if (totalPages <= 1) return null;

    return (
      <div className="flex items-center justify-between px-6 py-4 border-t border-gray-200 bg-gray-50">
        <div className="flex items-center space-x-2 text-sm text-gray-700">
          <span>
            Showing {startIndex + 1} to {Math.min(endIndex, currentTabData.length)} of {currentTabData.length} results
          </span>
        </div>
        
        <div className="flex items-center space-x-2">
          {/* Previous button */}
          <button
            onClick={goToPreviousPage}
            disabled={currentPage === 1}
            className={`p-2 rounded-md transition-colors ${
              currentPage === 1
                ? 'text-gray-400 cursor-not-allowed'
                : 'text-gray-600 hover:text-gray-800 hover:bg-gray-200'
            }`}
          >
            <ChevronLeft className="w-4 h-4" />
          </button>

          {/* Page numbers */}
          <div className="flex items-center space-x-1">
            {getPageNumbers().map((page, index) => (
              <button
                key={index}
                onClick={() => typeof page === 'number' && goToPage(page)}
                disabled={page === '...'}
                className={`px-3 py-1 text-sm rounded-md transition-colors ${
                  page === currentPage
                    ? 'bg-blue-600 text-white'
                    : page === '...'
                    ? 'text-gray-400 cursor-default'
                    : 'text-gray-600 hover:text-gray-800 hover:bg-gray-200'
                }`}
              >
                {page}
              </button>
            ))}
          </div>

          {/* Next button */}
          <button
            onClick={goToNextPage}
            disabled={currentPage === totalPages}
            className={`p-2 rounded-md transition-colors ${
              currentPage === totalPages
                ? 'text-gray-400 cursor-not-allowed'
                : 'text-gray-600 hover:text-gray-800 hover:bg-gray-200'
            }`}
          >
            <ChevronRight className="w-4 h-4" />
          </button>
        </div>
      </div>
    );
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-6xl h-5/6 flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b border-gray-200">
          <div className="text-left">
            <h2 className="text-2xl font-bold text-gray-900">Data Records</h2>
            <p className="text-gray-600">Explore all processed data from different sources</p>
          </div>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <X className="w-6 h-6" />
          </button>
        </div>

        {/* Tabs */}
        <div className="border-b border-gray-200 bg-gray-50">
          <div className="flex space-x-8 px-6">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => handleTabChange(tab.id)}
                className={`py-3 px-1 border-b-2 font-medium text-sm transition-all duration-200 ${
                  activeTab === tab.id
                    ? 'border-blue-500 text-blue-600 bg-white rounded-t-lg px-3'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }`}
              >
                <span className="mr-2">{tab.icon}</span>
                {tab.name}
                <span className={`ml-2 py-1 px-2 rounded-full text-xs ${
                  activeTab === tab.id 
                    ? 'bg-blue-100 text-blue-600' 
                    : 'bg-gray-100 text-gray-600'
                }`}>
                  {tab.count}
                </span>
              </button>
            ))}
          </div>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto">
          <div className="rounded-lg border border-gray-200 shadow-sm overflow-hidden">
            {renderTabContent()}
          </div>
        </div>

        {/* Pagination */}
        {renderPagination()}
      </div>
    </div>
  );
};
