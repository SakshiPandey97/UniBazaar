import React, { createContext, useState, useContext } from "react";

const SearchContext = createContext();

export const useSearchContext = () => {
    return useContext(SearchContext);
};

export const SearchProvider = ({ children }) => {
    const [searchTerm, setSearchTerm] = useState("");

    return (
        <SearchContext.Provider value={{ searchTerm, setSearchTerm }}>
            {children}
        </SearchContext.Provider>
    );
};
