import './App.css';
import {
  Route,
  Routes,
  Link,
  BrowserRouter
} from 'react-router-dom'

import Home from "./Home"
import Search from './Search'
import Browse from './Browse'

function App() {
  return (
    <BrowserRouter>
      <div className="App">
        <header>
          <Link to="/">🔍 Monocole</Link>
          <div className="divider"></div>
          Some slogan
        </header>
        <div className="content">
          <Routes>
            <Route exact path="/" element={<Home />} />
            <Route path="/search" element={<Search />} />
            <Route path="/browse/:project/*" element={<Browse />} />
          </Routes>
        </div>
      </div>
    </BrowserRouter>
  );
}

export default App;
