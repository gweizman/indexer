import React, { useState } from "react";
import { createSearchParams, useNavigate, useSearchParams } from "react-router-dom";
import './SearchBar.css'

function SearchBar() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate()
  const [value, setValue]  = useState(searchParams.has('query') ? searchParams.get('query') : '')

  const handleChange = event => {
      setValue(event.target.value)
  };

  const handleSubmit = event => {
    navigate({
      pathname:"/search",
      search: createSearchParams({
        query: value
      }).toString(),
    });
    event.preventDefault();
  };

  return (
    <div className="SearchBar">
      <form onSubmit={handleSubmit}>
          <input type="text" value={value} onChange={handleChange} placeholder="Search..."/>
      </form>
    </div>
  );
}
  
export default SearchBar;
  