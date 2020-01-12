'use strict';

import React, { useState, useEffect, useRef } from 'react';

const featuredData = "featured-data";
const allData = "all-data";

export default function TableWrapper(props) {
    const [data, setData] = useState([]);

    function fetchData() {
        fetch(props.type)
            .then((response) => response.json())
            .then((result) => setData(result), (error) => console.error(error));
    }

    useEffect(() => {
        fetchData();
        const id = setInterval(() => fetchData(), 3000);

        return () => clearInterval(id);
    }, []);

    let table;
    if (data && data.length > 0) {
        table = <Table type={props.type} data={data}/>;
    } else {
        table = <p>Sorry, no data available...</p>
    }

    return (
        <div className="row">
            <div className="col">
                <h3 id={props.type}>{props.title}</h3>
                {table}
            </div>
        </div>
    );
}

function Table(props) {
    return (
        <table>
            <thead>
            <tr>
                <th>Bla</th>
                <th>X</th>
                <th>Y</th>
                <th>Upvotes</th>
            </tr>
            </thead>
            <tbody>
            {props.data.map(doc => (
                <TableRow key={doc._id} type={props.type} doc={doc}/>
            ))}
            </tbody>
        </table>
    );
}

function TableRow(props) {
    let upvote;
    if (props.type === featuredData) {
        upvote = <FixedUpvotesCounter upvotes={props.doc.upvotes} id={props.doc._id}/>
    } else if (props.type === allData) {
        upvote = <InteractiveUpvotesCounter upvotes={props.doc.upvotes} id={props.doc._id}/>
    }
    return (
        <tr id={props.doc._id}>
            <td>{props.doc.val1}</td>
            <td>{props.doc.val2}</td>
            <td>{props.doc.val3}</td>
            <td className="valign-wrapper">
                {upvote}
            </td>
        </tr>
    );
}

function FixedUpvotesCounter(props) {
    return <span>{props.upvotes}</span>;
}

function InteractiveUpvotesCounter(props) {
    const [upvotes, setUpvotes] = useState(props.upvotes);
    const firstUpdate = useRef(true);
    useEffect(() => {
        if (firstUpdate.current) {
            firstUpdate.current = false;
            return;
        }
        fetch("upvote/", {method: "POST", body: props.id})
            .then(response => {
                if (!response.ok) {
                    return response.text();
                }
            })
            .then(result => console.error(result));
    }, [upvotes]);

    return (
        <>
            <span>{upvotes}</span>
            <i className="material-icons" onClick={() => setUpvotes(upvotes + 1)}>arrow_drop_up</i>
        </>
    );
}
