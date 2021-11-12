import React from 'react';
import Ribbon from './components/ribbon/Ribbon';
import Main from './components/main/Main';
import '../css/App.css';

export default class App extends React.Component{
    render(){
        return (
            <div className="mdl-layout mdl-js-layout">
                <Ribbon />
                <Main />
            </div>
        );
    }
}