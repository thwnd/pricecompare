import React, { useEffect } from 'react';

/**
 * 구글 애드센스 광고를 표시하는 컴포넌트
 * @param {string} adSlot - 애드센스 광고 슬롯 ID
 * @param {string} adFormat - 광고 형식 (자동: 'auto', 인피드: 'fluid')
 * @param {object} style - 인라인 스타일 객체
 * @param {string} className - 추가 CSS 클래스명
 */
const AdSense = ({ adSlot, adFormat = 'auto', style = {}, className = '' }) => {
  useEffect(() => {
    try {
      // 컴포넌트가 마운트된 후 광고 로드
      (window.adsbygoogle = window.adsbygoogle || []).push({});
    } catch (error) {
      console.error('AdSense 광고 로드 오류:', error);
    }
  }, []);

  // 기본 스타일 설정
  const defaultStyle = {
    display: 'block',
    textAlign: 'center',
    ...style
  };

  return (
    <div className={`ad-container ${className}`}>
      <ins
        className="adsbygoogle"
        style={defaultStyle}
        data-ad-client="ca-pub-XXXXXXXXXXXXXXXX" // 실제 AdSense 퍼블리셔 ID로 교체 필요
        data-ad-slot={adSlot}
        data-ad-format={adFormat}
        data-full-width-responsive="true"
      />
    </div>
  );
};

export default AdSense;