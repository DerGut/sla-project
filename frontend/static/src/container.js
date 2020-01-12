'use strict';

import React, { useState, useEffect } from 'react';

import TableWrapper from "./table";

export default function Container() {
    // const [upvotes, setUpvotes] = useState({});
    //
    // useEffect(() => {
    //
    // });

    return (
        <div className="container">
            <TableWrapper type="featured-data" title="Featured Data"/>

            <TableWrapper type="all-data" title="All Data"/>
        </div>
    );
}