import React from "react";
export default class StartStopContainer extends React.Component{
    constructor(props){
        super(props);
        this.onBtnClick = this.onBtnClick.bind(this);
    }
    
    getBtnText(){
        return this.props.runningStatus === 'running' ? 'stop' : 'play_arrow';
    }

    onBtnClick(){
        console.log('onBtnClick');
        this.props.toggleStatus(this.getNewStatus());
    }

    getNewStatus(){
        return this.props.runningStatus === 'running' ? 'stopped' : 'running';
    }

    render(){
        return (
            <div className="mdl-cell startstop-cell mdl-cell--6-col mdl-cell--1-offset-tablet mdl-cell--3-offset-desktop">
                <button id="stopstart" className="mdl-button mdl-js-button mdl-button--fab mdl-js-ripple-effect" onClick={this.onBtnClick} >
                    <i className="material-icons">{this.getBtnText()}</i>
                </button>
            </div>
        );
    }
}