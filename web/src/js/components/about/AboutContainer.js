import React from "react";

import twitterIcon from '../../../images/twitter.png';

const NUM_OF_REGIONS_IN_TWEET = 3,
    GLOBAL_REGION_KEY = 'global';
export default class AboutContainer extends React.Component {
    getTweetLink() {
        let tweet = 'My lowest-latency #GCP regions via gcping.com:',
            numRegions = NUM_OF_REGIONS_IN_TWEET;

        for (let i = 0; i < this.props.results.length; i++) {
            if (this.props.results[i] !== GLOBAL_REGION_KEY) {
                tweet += '\n' + this.props.results[i] + ' (' + this.props.regions[this.props.results[i]]['median'] + ' ms)';

                if (--numRegions === 0)
                    break;
            }
        }

        return 'https://twitter.com/share?text=' + encodeURIComponent(tweet);
    }

    render() {
        return (
            <div className="about-content mdl-cell mdl-cell--6-col mdl-cell--1-offset-tablet mdl-cell--3-offset-desktop">
                <center><a id="tweet-link" href={this.getTweetLink()} target="_blank"><img id="tweet-logo" src={twitterIcon} alt="Tweet your latencies"></img></a></center>
                <h2>How does this work?</h2>
                <p>Your browser makes HTTPS requests to <a target="_blank" href="https://cloud.google.com/run">Google Cloud Run</a> services deployed to each region. The median time between request and response is shown.</p>
                <p>The <b>global</b> row uses a <a target="_blank" href="https://cloud.google.com/load-balancing/">Global HTTPS Load Balancer</a> to route requests to the nearest service.</p>
                <p><b>Note:</b> This site is intended to show <i>relative</i> latency
                    to each region, and should not be used to determine the absolute
                    lowest latency possible, or your own network speed or other network
                    conditions. Results are not scientific.</p>
                <p>This is not an official Google project. <a target="_blank" href="https://github.com/GoogleCloudPlatform/gcping">Source available on GitHub</a>.</p>
            </div>
        );
    }
}