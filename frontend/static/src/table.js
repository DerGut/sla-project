'use strict';

import React, { useState, useEffect, useRef } from 'react';

const featuredData = "featured-data";
const allData = "all-data";

export default function TableWrapper(props) {
    const [data, setData] = useState([]);

    useEffect(() => {
        fetch(props.type)
            .then((response) => response.json())
            .then((result) => setData(result), (error) => console.error(error));
    }, []);

    return (
        <div className="row">
            <div className="col">
                <h3 id={props.type}>{props.title}</h3>
                {data &&
                    <Table type={props.type} data={data}/>
                }
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
        fetch("upvote/", {method: "POST", body: props.id});
    }, [upvotes]);

    return (
        <>
            <span>{upvotes}</span>
            <i className="material-icons" onClick={() => setUpvotes(upvotes + 1)}>arrow_drop_up</i>
        </>
    );
}
