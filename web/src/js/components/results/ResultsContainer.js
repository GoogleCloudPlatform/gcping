import React from "react";

const GLOBAL_REGION_KEY = "global";
export default class ResultsContainer extends React.Component{
    getTableBody(){
      let fastestRegion = this.getFastestRegionKey();

      if(this.props.loading){
        return (
          <tr>
            <td colSpan="2">Loading...</td>
          </tr>
        );
      }

      return this.props.results.map((regionKey)=>{
        let cls = (regionKey === fastestRegion && this.props.fastestRegionVisible) ? 'top' : '',
          displayedKey =this.getDisplayedRegionKey(regionKey);

        return (
          <tr className={cls} key={regionKey}>
            <td className="regiondesc">
              {this.props.regions[regionKey]['label']}
              <div className="region">{displayedKey}</div>
            </td>
            <td className="result">
              <div>{this.props.regions[regionKey]['median']} ms</div>
            </td>
          </tr>
        );
      });
    }

    /**
     * Gets the fastest region, excluding the global region
     * @returns string
     */
    getFastestRegionKey(){
      for(let i = 0; i < this.props.results.length; i++){
        if(this.props.results[i] !== GLOBAL_REGION_KEY){
          return this.props.results[i];
        }
      }
    }

    /**
     * Helper function to deduce the region to be displayed in the list
     * @param {string} regionKey
     * @returns
     */
    getDisplayedRegionKey(regionKey){
      // if the region is not global, return it as it is.
      if(regionKey !== GLOBAL_REGION_KEY)
        return regionKey;

      // if the region is global and we have received the region that is used by the Gloabl Load Balancer
      // we display that
      if(this.props.globalRegionProxy.length > 0)
        return (<em>→{this.props.globalRegionProxy}</em>);

      // if the region is global and we don't have the routing region, we show "gloabl"
      return 'global';
    }

    render(){
        return (
            <div className="mdl-cell results-cell mdl-cell--6-col mdl-cell--1-offset-tablet mdl-cell--3-offset-desktop">
          <table className="mdl-data-table mdl-js-data-table mdl-shadow--2dp mdl-color-text--grey-800">
            <thead>
              <tr>
                <th>REGION</th>
                <th>MEDIAN LATENCY</th>
              </tr>
            </thead>
            <tbody>
              {this.getTableBody()}
            </tbody>
          </table>
        </div>
        );
    }
}