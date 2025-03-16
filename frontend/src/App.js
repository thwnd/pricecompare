// App.js - React 메인 컴포넌트
import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Route, Routes, Link, useParams, useNavigate, useLocation } from 'react-router-dom';
import axios from 'axios';
import './App.css';

// 헤더 컴포넌트
const Header = ({ onSearch }) => {
  const [query, setQuery] = useState('');
  const navigate = useNavigate();

  const handleSubmit = (e) => {
    e.preventDefault();
    onSearch(query);
    navigate('/search');
  };

  return (
    <header className="header">
      <div className="logo">
        <Link to="/">가격비교마스터</Link>
      </div>
      <form className="search-form" onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="제품 검색..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
        <button type="submit">검색</button>
      </form>
      <nav className="nav">
        <ul>
          <li><Link to="/category/electronics">전자제품</Link></li>
          <li><Link to="/category/fashion">패션</Link></li>
          <li><Link to="/category/home">홈&리빙</Link></li>
          <li><Link to="/category/beauty">뷰티</Link></li>
        </ul>
      </nav>
    </header>
  );
};

// 제품 카드 컴포넌트
const ProductCard = ({ item }) => {
  return (
    <div className="product-card">
      <div className="product-image">
        <img src={item.image} alt={item.title} />
      </div>
      <div className="product-info">
        <h3><a href={item.url} target="_blank" rel="noopener noreferrer">{item.title}</a></h3>
        <p className="product-price">{item.price}</p>
        <p className="product-source">출처: {item.source}</p>
        <p className="product-description">{item.description.substring(0, 100)}...</p>
      </div>
    </div>
  );
};

// 홈 페이지 컴포넌트
const HomePage = () => {
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);

  useEffect(() => {
    const fetchItems = async () => {
      try {
        setLoading(true);
        const response = await axios.get(`/api/items?page=${page}`);
        setItems(response.data.items);
        setTotalPages(response.data.totalPage);
      } catch (error) {
        console.error('Error fetching items:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchItems();
  }, [page]);

  return (
    <div className="home-page">
      <h1>최신 가격비교 정보</h1>
      
      {/* 광고 슬롯 - 상단 배너 */}
      <div className="ad-banner">
        <div id="ad-slot-top-banner">
          {/* AdSense 코드가 여기에 삽입됩니다 */}
        </div>
      </div>

      {loading ? (
        <div className="loading">로딩 중...</div>
      ) : (
        <>
          <div className="product-grid">
            {items.map((item, index) => (
              <React.Fragment key={index}>
                <ProductCard item={item} />
                {/* 매 8개 아이템마다 인피드 광고 삽입 */}
                {(index + 1) % 8 === 0 && (
                  <div className="ad-infeed">
                    <div id={`ad-slot-infeed-${Math.floor((index + 1) / 8)}`}>
                      {/* AdSense 코드가 여기에 삽입됩니다 */}
                    </div>
                  </div>
                )}
              </React.Fragment>
            ))}
          </div>

          <div className="pagination">
            <button
              onClick={() => setPage(p => Math.max(1, p - 1))}
              disabled={page === 1}
            >
              이전
            </button>
            <span>페이지 {page} / {totalPages}</span>
            <button
              onClick={() => setPage(p => Math.min(totalPages, p + 1))}
              disabled={page === totalPages}
            >
              다음
            </button>
          </div>
        </>
      )}
      
      {/* 광고 슬롯 - 하단 배너 */}
      <div className="ad-banner">
        <div id="ad-slot-bottom-banner">
          {/* AdSense 코드가 여기에 삽입됩니다 */}
        </div>
      </div>
    </div>
  );
};

// 카테고리 페이지 컴포넌트
const CategoryPage = () => {
  const { category } = useParams();
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);

  useEffect(() => {
    const fetchCategoryItems = async () => {
      try {
        setLoading(true);
        const response = await axios.get(`/api/items/${category}?page=${page}`);
        setItems(response.data.items);
        setTotalPages(response.data.totalPage);
      } catch (error) {
        console.error('Error fetching category items:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchCategoryItems();
  }, [category, page]);

  // 카테고리 이름 한글화
  const getCategoryName = () => {
    switch (category) {
      case 'electronics': return '전자제품';
      case 'fashion': return '패션';
      case 'home': return '홈&리빙';
      case 'beauty': return '뷰티';
      default: return category;
    }
  };

  return (
    <div className="category-page">
      <h1>{getCategoryName()} 가격비교</h1>
      
      {/* 광고 슬롯 */}
      <div className="ad-banner">
        <div id="ad-slot-category">
          {/* AdSense 코드가 여기에 삽입됩니다 */}
        </div>
      </div>

      {loading ? (
        <div className="loading">로딩 중...</div>
      ) : (
        <>
          <div className="product-grid">
            {items.map((item, index) => (
              <ProductCard key={index} item={item} />
            ))}
          </div>

          <div className="pagination">
            <button
              onClick={() => setPage(p => Math.max(1, p - 1))}
              disabled={page === 1}
            >
              이전
            </button>
            <span>페이지 {page} / {totalPages}</span>
            <button
              onClick={() => setPage(p => Math.min(totalPages, p + 1))}
              disabled={page === totalPages}
            >
              다음
            </button>
          </div>
        </>
      )}
    </div>
  );
};

// 검색 결과 페이지 컴포넌트
const SearchPage = ({ query }) => {
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);

  useEffect(() => {
    const searchItems = async () => {
      if (!query) return;
      
      try {
        setLoading(true);
        const response = await axios.get(`/api/search?q=${encodeURIComponent(query)}&page=${page}`);
        setItems(response.data.items);
        setTotalPages(response.data.totalPage);
      } catch (error) {
        console.error('Error searching items:', error);
      } finally {
        setLoading(false);
      }
    };

    searchItems();
  }, [query, page]);

  return (
    <div className="search-page">
      <h1>"{query}" 검색 결과</h1>
      
      {/* 광고 슬롯 */}
      <div className="ad-banner">
        <div id="ad-slot-search">
          {/* AdSense 코드가 여기에 삽입됩니다 */}
        </div>
      </div>

      {loading ? (
        <div className="loading">로딩 중...</div>
      ) : (
        <>
          {items.length === 0 ? (
            <div className="no-results">검색 결과가 없습니다.</div>
          ) : (
            <>
              <div className="product-grid">
                {items.map((item, index) => (
                  <ProductCard key={index} item={item} />
                ))}
              </div>

              <div className="pagination">
                <button
                  onClick={() => setPage(p => Math.max(1, p - 1))}
                  disabled={page === 1}
                >
                  이전
                </button>
                <span>페이지 {page} / {totalPages}</span>
                <button
                  onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                  disabled={page === totalPages}
                >
                  다음
                </button>
              </div>
            </>
          )}
        </>
      )}
    </div>
  );
};

// 푸터 컴포넌트
const Footer = () => {
  return (
    <footer className="footer">
      <div className="footer-content">
        <div className="footer-section">
          <h3>가격비교마스터</h3>
          <p>다양한 쇼핑몰의 상품 가격을 한눈에 비교하세요!</p>
        </div>
        <div className="footer-section">
          <h3>카테고리</h3>
          <ul>
            <li><Link to="/category/electronics">전자제품</Link></li>
            <li><Link to="/category/fashion">패션</Link></li>
            <li><Link to="/category/home">홈&리빙</Link></li>
            <li><Link to="/category/beauty">뷰티</Link></li>
          </ul>
        </div>
        <div className="footer-section">
          <h3>회사 정보</h3>
          <p>회사명: 가격비교마스터</p>
          <p>이메일: contact@example.com</p>
          <p>전화번호: 02-123-4567</p>
        </div>
      </div>
      <div className="footer-bottom">
        <p>&copy; 2025 가격비교마스터. All Rights Reserved.</p>
      </div>
    </footer>
  );
};

// 라우터가 필요한 컴포넌트를 감싸는 래퍼
const AppContent = () => {
  const [searchQuery, setSearchQuery] = useState('');

  const handleSearch = (query) => {
    setSearchQuery(query);
  };

  return (
    <div className="app">
      <Header onSearch={handleSearch} />
      <main className="main-content">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/category/:category" element={<CategoryPage />} />
          <Route path="/search" element={<SearchPage query={searchQuery} />} />
        </Routes>
      </main>
      <Footer />
    </div>
  );
};

// 메인 App 컴포넌트
const App = () => {
  // Google AdSense 스크립트 로드
  useEffect(() => {
    // AdSense 스크립트 로드
    const adsScript = document.createElement('script');
    adsScript.src = 'https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js?client=ca-pub-XXXXXXXXXXXXXXXX';
    adsScript.async = true;
    adsScript.crossOrigin = 'anonymous';
    document.head.appendChild(adsScript);

    // 모든 광고 초기화
    window.adsbygoogle = window.adsbygoogle || [];

    return () => {
      document.head.removeChild(adsScript);
    };
  }, []);

  return (
    <Router>
      <AppContent />
    </Router>
  );
};

export default App;