/**
 * Copyright 2021 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import React from 'react';
import Heading from '../heading/Heading';
import StartStopContainer from '../controls/StartStopContainer';
import ResultsContainer from '../results/ResultsContainer';
import AboutContainer from '../about/AboutContainer';

const INITIAL_ITERATIONS = 10;
const PING_TEST_RUNNING_STATUS = 'running';
const PING_TEST_STOPPED_STATUS = 'stopped';
const GLOBAL_REGION_KEY = 'global';
export default class Main extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      regions: [],
      results: [],
      loading: true,
      runningStatus: PING_TEST_RUNNING_STATUS,
      globalRegionProxy: '',
      fastestRegionVisible: false,
    };

    this.updateRunningStatus = this.updateRunningStatus.bind(this);
  }

  async componentDidMount() {
    const regions = await this.getEndpoints();
    this.setState({
      regions: regions,
      loading: false,
    });

    // once we're done fetching all endpoints, let's start pinging
    this.pingAllRegions(INITIAL_ITERATIONS);
  }

  /**
     * Fetches the endpoints for different Cloud Run regions.
     * We will later send a request to these endpoints and measure the latency.
     */
  async getEndpoints() {
    return new Promise((resolve, reject)=>{
      const regions = [];

      fetch('/api/endpoints').then((resp)=>{
        return resp.json();
      }).then((endpoints)=>{
        for (const zone of Object.values(endpoints)) {
          const gcpZone = {
            key: zone.Region,
            label: zone.RegionName,
            pingUrl: zone.URL + '/ping',
            latencies: [],
            median: '',
          };

          regions[gcpZone.key] = gcpZone;
        }

        resolve(regions);
      });
    });
  }

  /**
     * Ping all regions to fetch their latency
     */
  async pingAllRegions(iter) {
    const regions = this.state.regions;
    const regionsArr = Object.values(regions);

    // using a for..of instead of a simple for loop
    // to make sure we respect the async operations before proceeding.
    for (const i of new Array(iter)) {
      for (const region of regionsArr) {
        // Takes care of the stopped button
        if (this.state.runningStatus === PING_TEST_STOPPED_STATUS) {
          return;
        }

        const latency = await this.pingSingleRegion(region.key);

        // add the latency to the array of latencies
        // from where we can compute the median and populate the table
        regions[region.key]['latencies'].push(latency);

        const median = this.getMedian(regions[region.key]['latencies']);
        regions[region.key]['median'] = median;

        this.setState({
          regions,
        }, ()=>{
          const results = this.getSortedResults(region.key, median);
          this.setState({
            results,
          });
        });
      }

      // start displaying the fastest region after at least 1 iteration is over.
      // subsequent calls to this won't change anything
      this.setState({
        fastestRegionVisible: true,
      });
    }

    // when all the region latencies have been fetched,
    // let's update our status flag
    this.updateRunningStatus(PING_TEST_STOPPED_STATUS);
  }

  /**
     * Computes the ping time for a single GCP region
     * @param {string} regionKey The key of the GCP region, ex: us-east1
     * @return Promise
     */
  pingSingleRegion(regionKey) {
    return new Promise((resolve) => {
      const gcpZone = this.state.regions[regionKey];
      const start = new Date().getTime();

      fetch(gcpZone.pingUrl, {
        mode: 'no-cors',
        cache: 'no-cache',
      }).then(async (resp) => {
        const latency = new Date().getTime() - start;

        // if we just pinged the global region, the response should contain
        // the region that the Global Load Balancer uses to route the traffic.
        if (regionKey === GLOBAL_REGION_KEY) {
          resp.text().then((val)=>{
            this.setState({
              globalRegionProxy: val.trim(),
            });
          });
        }

        resolve(latency);
      });
    });
  }

  /**
     * Function to update the current status of pinging
     * @param {string} status
     */
  updateRunningStatus(status) {
    this.setState({
      runningStatus: status,
    }, ()=>{
      if (status === PING_TEST_RUNNING_STATUS) {
        this.pingAllRegions(1);
      }
    });
  }

  /**
     * Helper function to return median from a given array
     * @param {*} arr Array of latencies
     * @returns
     */
  getMedian(arr) {
    if (arr.length == 0) {
      return 0;
    }
    const copy = arr.slice(0);
    copy.sort();
    return copy[Math.floor(copy.length/2)];
  }

  /**
     * Helper that adds the regionKey to it's proper position
     * making the results array sorted
     * TODO: Try and use an ordered map here to simply this
     */
  getSortedResults(regionKey, latency) {
    const results = this.state.results;
    const regions = this.state.regions;

    if (!results.length) {
      results.push(regionKey);
      return results;
    }

    // remove any current values with the same regionKey
    for (let i = 0; i < results.length; i++) {
      if (results[i] === regionKey) {
        results.splice(i, 1);
        break;
      }
    }

    // TODO: Probably use Binary search here to merge the following 2 blocks
    if (latency < regions[results[0]].median) {
      results.unshift(regionKey);
      return results;
    } else if (latency > regions[results[results.length - 1]].median) {
      results.push(regionKey);
      return results;
    }

    // add the region to it's proper position
    for (let i = 0; i < results.length - 1; i++) {
      if ( latency >= regions[results[i]].median &&
          latency <= regions[results[i+1]].median) {
        results.splice(i+1, 0, regionKey);
        return results;
      }
    }
  }

  render() {
    return (
      <main className="mdl-layout__content">
        <div className="mdl-grid">
          <Heading />
          <StartStopContainer
            runningStatus={this.state.runningStatus}
            toggleStatus={(status)=>{
              this.updateRunningStatus(status);
            }}
          />
          <ResultsContainer
            loading={this.state.loading}
            regions={this.state.regions}
            results={this.state.results}
            fastestRegionVisible={this.state.fastestRegionVisible}
            globalRegionProxy={this.state.globalRegionProxy}
          />
          <AboutContainer
            results={this.state.results}
            regions={this.state.regions}
          />
        </div>
      </main>
    );
  }
}
