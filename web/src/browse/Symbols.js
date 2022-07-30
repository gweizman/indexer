import { QueryClient, QueryClientProvider, useQuery } from '@tanstack/react-query'
import { useState } from 'react';
import { createSearchParams, useNavigate, useSearchParams, Link } from "react-router-dom";
import utils from '../utils/utils';

import './Symbols.css';

function Symbols() {
    const [searchValue, setSearchValue]  = useState('')

    return (
        <div className="SymbolBrowser">
            <strong>Project Symbols</strong>
            <SearchBox value={searchValue} setValue={setSearchValue} />
            <SymbolList project="OpenGrok" query={searchValue} />
        </div>
    )
}

function SearchBox(props) {  
    const handleChange = event => {
        props.setValue(event.target.value)
    };
  
    const handleSubmit = event => {
      event.preventDefault();
    };
  
    return (
      <div className="SearchBar">
        <form onSubmit={handleSubmit}>
            <input type="text" value={props.value} onChange={handleChange} placeholder="Search..."/>
        </form>
      </div>
    );
}

function SymbolList(props) {
  let { isLoading, error, data } = useQuery(['searchResults', props.project, props.query], () =>
      fetch(`/api/${props.project}/definition/?` + new URLSearchParams({
        name: ".*" + props.query + ".*"
      })).then(res =>
          res.json()
      )
  )
 // project, name, language, pattern, signature, file_limited, parent, parent_kind, file_name, file_path, line
  if (isLoading) return 'Loading...'
  if (error) return 'An error has occurred: ' + error.message
  return (
      <div className="symbolResults">
        {data.map((object, i) => <Symbol object={object} />)}
      </div>
  )
}

function Symbol(props) {
  const filePath = utils.normalizePath(utils.buildPath(props.object.file_path, props.object.file_name)) // TODO: Not sure why this is needed. Seems like a bug in db insertions.

  return (
      <div>
          <Link className="reset" to={`/browse/${props.object.project}/${filePath}`}>
          <span>
              {props.object.name} <span className="small">({filePath}:{props.object.line})</span>
          </span>
          </Link>
          <br />
      </div>
  )
}

export default Symbols;