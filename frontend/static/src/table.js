'use strict';

import React, { useState, useEffect } from 'react';

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
                    <Table data={data}/>
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
                <TableRow key={doc._id} doc={doc}/>
            ))}
            </tbody>
        </table>
    );
}

function TableRow(props) {
    return (
        <tr id={props.doc._id}>
            <td>{props.doc.val1}</td>
            <td>{props.doc.val2}</td>
            <td>{props.doc.val3}</td>
            <td className="valign-wrapper">
                <span>{props.doc.upvotes}</span>
                <i className="material-icons" onClick={(e) => voteUp(e, props.doc._id)}>arrow_drop_up</i>
            </td>
        </tr>
    );
}

function Upvote(props) {
    return (
        <>
        </>
    );
}

function voteUp(ev, id) {
    const elem = ev.target.previousElementSibling;
    fetch("upvote/", {
        method: "POST",
        body: id
    })
        .then(response => {
            if (!response.ok) {
                throw new Error("Upvote request failed");
            }
            return response.text();
        })
        .then(number => {
            if (isNaN(number)) {
                throw new Error("No number returned")
            }
            elem.innerHTML = number;
        })
        .catch(error => {
            console.error(error)
        });
}
