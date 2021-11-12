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
export default class StartStopContainer extends React.Component {
  constructor(props) {
    super(props);
    this.onBtnClick = this.onBtnClick.bind(this);
  }

  getBtnText() {
    return this.props.runningStatus === 'running' ? 'stop' : 'play_arrow';
  }

  onBtnClick() {
    this.props.toggleStatus(this.getNewStatus());
  }

  getNewStatus() {
    return this.props.runningStatus === 'running' ? 'stopped' : 'running';
  }

  render() {
    return (
      <div className="mdl-cell
                      startstop-cell
                      mdl-cell--6-col
                      mdl-cell--1-offset-tablet
                      mdl-cell--3-offset-desktop"
      >
        <button
          id="stopstart"
          className="mdl-button
                    mdl-js-button
                    mdl-button--fab
                    mdl-js-ripple-effect"
          onClick={this.onBtnClick}
        >
          <i className="material-icons">{this.getBtnText()}</i>
        </button>
      </div>
    );
  }
}

StartStopContainer.propTypes = {
  runningStatus: PropTypes.string.isRequired,
  toggleStatus: PropTypes.func.isRequired,
};
