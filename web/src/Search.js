import SearchBar from './search/SearchBar';
import { QueryClient, QueryClientProvider, useQuery } from '@tanstack/react-query'
import { createSearchParams, useNavigate, useSearchParams, Link } from "react-router-dom";
import './container.css'
import './Search.css'
import utils from './utils/utils';

const queryClient = new QueryClient()

function Search() {
    const [searchParams] = useSearchParams(); 
    const query = searchParams.get('query');

    return (
        <div className="container">
            <SearchBar />
            <QueryClientProvider client={queryClient}>
                <div className='contentWrapper'>
                    <div className="searchContent">
                        <QueryResults query={query} />
                    </div>
                </div>
            </QueryClientProvider>
        </div>
    )
}

function QueryResults(props) {
    let { isLoading, error, data } = useQuery(['searchResults', props.query], () =>
        fetch('/api/OpenGrok/search/?' + new URLSearchParams({
            query: props.query
        })).then(res =>
            res.json()
        )
    )

    if (isLoading) return 'Loading...'
    if (error) return 'An error has occurred: ' + error.message
    return (
        <div>
        <span>Search results for {props.query}:</span> <br />
        {data.map((object, i) => <QueryResult object={object} query={props.query} />)}
        </div>
    )
}

function QueryResult(props) {
    const filePath = `${props.object.project}/${utils.normalizePath(props.object.file_path)}/${props.object.file_name}`;
    let dataLines = (props.object.data).split('\n')
    dataLines.splice(5)
    return (
        <div className='search-result'>
            <span>In <Link to={utils.buildPath("/browse/", filePath) + "?" + new URLSearchParams({query: props.query})}>{filePath}</Link></span><br />
            <span style={{whiteSpace: "pre-wrap"}}>
            {dataLines.join('\n')}
            </span>
        </div>
    );
}

export default Search;