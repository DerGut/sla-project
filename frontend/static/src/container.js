import React from 'react';
import TableWrapper from "./table";

export default function Container() {
    return (
        <>
            <TableWrapper type="featured-data" title="Featured Data"/>

            <TableWrapper type="all-data" title="All Data"/>
        </>
    );
}