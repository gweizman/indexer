import SearchBar from './search/SearchBar';
import { QueryClient, QueryClientProvider, useQuery } from '@tanstack/react-query'
import { useParams, createSearchParams, useNavigate, useSearchParams, Link } from "react-router-dom";
import './container.css'
import './Browse.css'
import FileTree from './browse/FileTree';
import FilePreview from './browse/FilePreview';
import Symbols from './browse/Symbols';

const queryClient = new QueryClient()

function Browse() {
    const [searchParams] = useSearchParams(); 
    const query = searchParams.get('query');
    const { project } = useParams();
    const file_path = useParams()['*'];

    return (
        <div className="container">
            <SearchBar />
            <QueryClientProvider client={queryClient}>
                <div className='contentWrapper'>
                    <div className='columns'>
                        <div className='fileTreeWrapper'>
                            <FileTree project={project} current={file_path} />
                        </div>
                        <div className='filePreviewWrapper'>
                            <FilePreview project={project} file={file_path}/>
                        </div>
                        <div className='symbolsWrapper'>
                            <Symbols />
                        </div>
                    </div>
                </div>
            </QueryClientProvider>
        </div>
    )
}

export default Browse;