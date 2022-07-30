import SearchBar from './search/SearchBar';
import ProjectBrowser from './browse/ProjectBrowser';

function Home() {
    return (
      <div className="Home">
        <h2>🔍 Monocle Indexer</h2>
        <SearchBar />
        <div>
        <br />
        <span>or browse projects:</span>
        <ProjectBrowser />
        </div>
      </div>
    );
  }

  export default Home;