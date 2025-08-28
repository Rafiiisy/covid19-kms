import React from 'react';

interface WordCloudProps {
  wordFrequency: any;
}

interface WordData {
  word: string;
  count: number;
  positive_count: number;
  negative_count: number;
  neutral_count: number;
  avg_sentiment: number;
}

export const WordCloud: React.FC<WordCloudProps> = ({ wordFrequency }) => {
  if (!wordFrequency || !wordFrequency.words || wordFrequency.words.length === 0) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Word Frequency Analysis</h3>
        <div className="h-64 bg-gray-100 rounded-lg flex items-center justify-center">
          <p className="text-gray-500">No word frequency data available</p>
        </div>
      </div>
    );
  }

  const words = wordFrequency.words as WordData[];
  
  // Sort words by frequency and take top 50 for display
  const topWords = words
    .sort((a, b) => b.count - a.count)
    .slice(0, 50);

  // Calculate size ranges for word scaling
  const maxCount = Math.max(...topWords.map(w => w.count));
  const minCount = Math.min(...topWords.map(w => w.count));

  const getWordSize = (count: number) => {
    const normalized = (count - minCount) / (maxCount - minCount);
    return Math.max(12, Math.min(48, 12 + normalized * 36)); // 12px to 48px
  };

  const getWordColor = (word: WordData) => {
    if (word.avg_sentiment > 0.1) return 'text-green-600';
    if (word.avg_sentiment < -0.1) return 'text-red-600';
    return 'text-gray-600';
  };

  const getWordWeight = (count: number) => {
    const normalized = (count - minCount) / (maxCount - minCount);
    return Math.max(400, Math.min(900, 400 + normalized * 500)); // 400 to 900
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h3 className="text-lg font-semibold text-gray-900 mb-4">Word Frequency Analysis</h3>
      
      {/* Word Cloud Display */}
      <div className="h-80 bg-gray-50 rounded-lg p-4 overflow-hidden relative">
        <div className="flex flex-wrap justify-center items-center h-full gap-2">
          {topWords.map((word, index) => (
            <div
              key={index}
              className={`cursor-pointer hover:scale-110 transition-transform duration-200 ${getWordColor(word)}`}
              style={{
                fontSize: `${getWordSize(word.count)}px`,
                fontWeight: getWordWeight(word.count),
              }}
              title={`${word.word}: ${word.count} occurrences
Sentiment: ${word.avg_sentiment.toFixed(2)}
Positive: ${word.positive_count}, Negative: ${word.negative_count}, Neutral: ${word.neutral_count}`}
            >
              {word.word}
            </div>
          ))}
        </div>
      </div>

      {/* Statistics */}
      <div className="mt-4 grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
        <div className="text-center">
          <div className="text-blue-600 font-semibold">
            {wordFrequency.total_words || 0}
          </div>
          <div className="text-gray-600">Total Words</div>
        </div>
        <div className="text-center">
          <div className="text-green-600 font-semibold">
            {topWords.filter(w => w.avg_sentiment > 0.1).length}
          </div>
          <div className="text-gray-600">Positive Words</div>
        </div>
        <div className="text-center">
          <div className="text-red-600 font-semibold">
            {topWords.filter(w => w.avg_sentiment < -0.1).length}
          </div>
          <div className="text-gray-600">Negative Words</div>
        </div>
        <div className="text-center">
          <div className="text-gray-600 font-semibold">
            {topWords.filter(w => w.avg_sentiment >= -0.1 && w.avg_sentiment <= 0.1).length}
          </div>
          <div className="text-gray-600">Neutral Words</div>
        </div>
      </div>

      {/* Legend */}
      <div className="mt-4 flex justify-center space-x-6 text-xs">
        <div className="flex items-center space-x-2">
          <div className="w-3 h-3 bg-green-600 rounded-full"></div>
          <span>Positive Sentiment</span>
        </div>
        <div className="flex items-center space-x-2">
          <div className="w-3 h-3 bg-red-600 rounded-full"></div>
          <span>Negative Sentiment</span>
        </div>
        <div className="flex items-center space-x-2">
          <div className="w-3 h-3 bg-gray-600 rounded-full"></div>
          <span>Neutral Sentiment</span>
        </div>
      </div>
    </div>
  );
};
