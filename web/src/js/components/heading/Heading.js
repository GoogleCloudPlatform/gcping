import React from "react";
export default class Heading extends React.Component{
    render(){
        return (
            <div className="mdl-cell mdl-cell--6-col mdl-cell--1-offset-tablet mdl-cell--3-offset-desktop">
                <h1>Measure your latency to <a target="_blank" href="https://cloud.google.com/compute/docs/regions-zones/regions-zones">Google Cloud regions</a></h1>
            </div>
        );
    }
}