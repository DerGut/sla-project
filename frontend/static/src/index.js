import React from 'react';

import TableWrapper from "./table";

ReactDOM.render(<Root/>, document.getElementById("root"));

function Root() {
    return (
        <>
            <Nav/>
            <Container/>
        </>
    );
}

function Nav() {
    return (
        <nav>
            <div className="nav-wrapper">
                <a href="#" className="brand-logo" id="title-thing">Sample App</a>
                <ul id="nav-mobile" className="right hide-on-med-and-down">
                    <li><a href="#featured-data">Featured Data</a></li>
                    <li><a href="#all-data">All Data</a></li>
                    <li><a href="#process-something">Process Something</a></li>
                </ul>
            </div>
        </nav>
    );
}

function Container() {
    return (
        <div className="container">
            <TableWrapper type="featured-data" title="Featured Data"/>

            <TableWrapper type="all-data" title="All Data"/>
        </div>
    );
}
