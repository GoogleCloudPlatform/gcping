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
import PropTypes from 'prop-types';
import twitterIcon from '../../../images/twitter.png';

const NUM_OF_REGIONS_IN_TWEET = 3;
const GLOBAL_REGION_KEY = 'global';
export default class AboutContainer extends React.Component {
  getTweetLink() {
    let tweet = 'My lowest-latency #GCP regions via gcping.com:';
    let numRegions = NUM_OF_REGIONS_IN_TWEET;

    for (let i = 0; i < this.props.results.length; i++) {
      if (this.props.results[i] !== GLOBAL_REGION_KEY) {
        const regionKey = this.props.results[i];
        const median = this.props.regions[regionKey]['median'];

        tweet += `\n${regionKey} (${median} ms)`;

        if (--numRegions === 0) {
          break;
        }
      }
    }

    return 'https://twitter.com/share?text=' + encodeURIComponent(tweet);
  }

  render() {
    return (
      <div className="about-content
                      mdl-cell
                      mdl-cell--6-col
                      mdl-cell--1-offset-tablet
                      mdl-cell--3-offset-desktop"
      >
        <center>
          <a
            id="tweet-link"
            href={this.getTweetLink()}
            target="_blank"
            rel="noreferrer"
          >
            <img
              id="tweet-logo"
              src={twitterIcon}
              alt="Tweet your latencies"
            />
          </a>
        </center>
        <h2>How does this work?</h2>
        <p>Your browser makes HTTPS requests to <a target="_blank" href="https://cloud.google.com/run" rel="noreferrer">Google Cloud Run</a> services deployed to each region. The median time between request and response is shown.</p>
        <p>The <b>global</b> row uses a <a target="_blank" href="https://cloud.google.com/load-balancing/" rel="noreferrer">Global HTTPS Load Balancer</a> to route requests to the nearest service.</p>
        <p><b>Note:</b>
          This site is intended to show <i>relative</i> latency
          to each region, and should not be used to determine the absolute
          lowest latency possible,
          or your own network speed or other network
          conditions. Results are not scientific.</p>
        <p>This is not an official Google project. <a target="_blank" href="https://github.com/GoogleCloudPlatform/gcping" rel="noreferrer">Source available on GitHub</a>.</p>
      </div>
    );
  }
}

AboutContainer.propTypes = {
  results: PropTypes.array.isRequired,
  regions: PropTypes.array.isRequired,
};
